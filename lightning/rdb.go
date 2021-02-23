package lightning

import (
	"context"

	"github.com/distributedio/titan/db"
	kv "github.com/pingcap/tidb-lightning/lightning/backend"
	"github.com/pingcap/tidb-lightning/lightning/common"
	"go.uber.org/zap"
)

func NewRdbDecode(ctx context.Context, w *kv.LocalEngineWriter, ns string) *RdbDecode {
	return &RdbDecode{
		ctx: ctx,
		ns:  ns,
		w:   w,
	}
}

type RdbDecode struct {
	ctx   context.Context
	w     *kv.LocalEngineWriter
	meta  Meta
	db    *db.DB
	ns    string
	nowTs int64
}

func (r *RdbDecode) IsExpired(expire int64) bool {
	if expire > r.nowTs {
		return false
	}
	return true
}

func (r *RdbDecode) write(rows kv.Rows) error {
	return r.w.WriteRows(r.ctx, nil, rows)
}

// StartRDB is called when parsing of a valid RDB file starts.
func (r *RdbDecode) StartRDB() {}

// StartDatabase is called when database n starts.
// Once a database starts, another database will not start until EndDatabase is called.
func (r *RdbDecode) StartDatabase(n int) {
	if r.db == nil {
		r.db = &db.DB{
			Namespace: r.ns,
			ID:        db.DBID(n),
		}
	} else {
		r.db.ID = db.DBID(n)
	}
}

// AUX field
func (r *RdbDecode) Aux(key, value []byte) {}

// ResizeDB hint
func (r *RdbDecode) ResizeDatabase(dbSize, expiresSize uint32) {}

// Set is called once for each string key.
func (r *RdbDecode) Set(key, value []byte, expiry int64) {
	if r.IsExpired(expiry) {
		return
	}
	meta := NewStringMeta()
	meta.Value = value
	metaKey := db.MetaKey(r.db, key)
	var kvs []common.KvPair
	if expiry > 0 {
		meta.Object.ExpireAt = expiry
		if ekey, err := ExpireKey(metaKey, expiry); err == nil {
			kvs = append(kvs, common.KvPair{Key: ekey, Val: meta.ID})
		}
	}

	kvs = append(kvs, common.KvPair{Key: metaKey, Val: meta.Encode()})
	if err := r.write(kv.MakeRowsFromKvPairs(kvs)); err != nil {
		zap.L().Error("write string err", zap.String("key", string(key)), zap.Error(err))
	}
}

// StartHash is called at the beginning of a hash.
// Hset will be called exactly length times before EndHash.
func (r *RdbDecode) StartHash(key []byte, length, expiry int64) {
	if r.IsExpired(expiry) {
		zap.L().Info("hash key expired", zap.String("key", string(key)))
		return
	}
	meta := NewHashMeta()
	if expiry > 0 {
		meta.Object.ExpireAt = expiry
	}
	r.meta = meta
}

// Hset is called once for each field=value pair in a hash.
func (r *RdbDecode) Hset(key, field, value []byte) {
	if r.meta == nil {
		return
	}
	meta := r.meta.(*HashMeta)
	var kvs []common.KvPair
	kvs = append(kvs, common.KvPair{Key: HashItemKey(r.db, meta, field), Val: value})

	if err := r.write(kv.MakeRowsFromKvPairs(kvs)); err != nil {
		zap.L().Error("write hash iter err", zap.String("key", string(key)), zap.Error(err))
	}
}

// EndHash is called when there are no more fields in a hash.
func (r *RdbDecode) EndHash(key []byte) {
	if r.meta == nil {
		return
	}
	meta := r.meta.(*HashMeta)
	metaKey := db.MetaKey(r.db, key)
	var kvs []common.KvPair
	kvs = append(kvs, common.KvPair{Key: metaKey, Val: meta.Encode()})
	if meta.ExpireAt > 0 {
		if ekey, err := ExpireKey(metaKey, meta.ExpireAt); err == nil {
			kvs = append(kvs, common.KvPair{Key: ekey, Val: meta.ID})
		}
	}

	if err := r.write(kv.MakeRowsFromKvPairs(kvs)); err != nil {
		zap.L().Error("write hash meta err", zap.String("key", string(key)), zap.Error(err))
	}
}

// StartSet is called at the beginning of a set.
// Sadd will be called exactly cardinality times before EndSet.
func (r *RdbDecode) StartSet(key []byte, cardinality, expiry int64) {
	if r.IsExpired(expiry) {
		return
	}
	meta := NewSetMeta()
	if expiry > 0 {
		meta.Object.ExpireAt = expiry
	}
	r.meta = meta
}

// Sadd is called once for each member of a set.
func (r *RdbDecode) Sadd(key, member []byte) {
	if r.meta == nil {
		return
	}
	meta := r.meta.(*SetMeta)
	var kvs []common.KvPair
	kvs = append(kvs, common.KvPair{Key: SetItemKey(r.db, meta, member), Val: db.SetNilValue})
	if err := r.write(kv.MakeRowsFromKvPairs(kvs)); err != nil {
		zap.L().Error("write set iter err", zap.String("key", string(key)), zap.Error(err))
	}
	meta.Len += 1

}

// EndSet is called when there are no more fields in a set.
func (r *RdbDecode) EndSet(key []byte) {
	if r.meta == nil {
		return
	}
	meta := r.meta.(*SetMeta)
	metaKey := db.MetaKey(r.db, key)
	var kvs []common.KvPair
	kvs = append(kvs, common.KvPair{Key: metaKey, Val: meta.Encode()})
	if meta.ExpireAt > 0 {
		if ekey, err := ExpireKey(metaKey, meta.ExpireAt); err == nil {
			kvs = append(kvs, common.KvPair{Key: ekey, Val: meta.ID})
		}
	}
	if err := r.write(kv.MakeRowsFromKvPairs(kvs)); err != nil {
		zap.L().Error("write set meta err", zap.String("key", string(key)), zap.Error(err))
	}

}

// StartList is called at the beginning of a list.
// Rpush will be called exactly length times before EndList.
// If length of the list is not known, then length is -1
func (r *RdbDecode) StartList(key []byte, length, expiry int64) {
	if r.IsExpired(expiry) {
		return
	}
	meta := NewLListMeta()
	if expiry > 0 {
		meta.Object.ExpireAt = expiry
	}
	r.meta = meta
}

// Rpush is called once for each value in a list.
func (r *RdbDecode) Rpush(key, value []byte) {
	if r.meta == nil {
		return
	}
	meta := r.meta.(*LListMeta)
	meta.Rindex++
	iterKey, err := LListItemKey(r.db, meta, meta.Rindex)
	if err != nil {
		zap.L().Error("get llist iter key err", zap.String("key", string(key)), zap.Error(err))
	}
	meta.Len++
	if meta.Len == 1 {
		meta.Lindex = meta.Rindex
	}
	var kvs []common.KvPair
	kvs = append(kvs, common.KvPair{Key: iterKey, Val: value})
	if err := r.write(kv.MakeRowsFromKvPairs(kvs)); err != nil {
		zap.L().Error("write set iter err", zap.String("key", string(key)), zap.Error(err))
	}
	meta.Len += 1
}

// EndList is called when there are no more values in a list.
func (r *RdbDecode) EndList(key []byte) {
	if r.meta == nil {
		return
	}
	meta := r.meta.(*LListMeta)
	metaKey := db.MetaKey(r.db, key)
	var kvs []common.KvPair
	kvs = append(kvs, common.KvPair{Key: metaKey, Val: meta.Encode()})
	if meta.ExpireAt > 0 {
		if ekey, err := ExpireKey(key, meta.ExpireAt); err == nil {
			kvs = append(kvs, common.KvPair{Key: ekey, Val: meta.ID})
		}
	}
	if err := r.write(kv.MakeRowsFromKvPairs(kvs)); err != nil {
		zap.L().Error("write llist meta err", zap.String("key", string(key)), zap.Error(err))
	}
}

// StartZSet is called at the beginning of a sorted set.
// Zadd will be called exactly cardinality times before EndZSet.
func (r *RdbDecode) StartZSet(key []byte, cardinality, expiry int64) {
	if r.IsExpired(expiry) {
		return
	}
	meta := NewZSetMeta()
	if expiry > 0 {
		meta.Object.ExpireAt = expiry
	}
	r.meta = meta

}

// Zadd is called once for each member of a sorted set.
func (r *RdbDecode) Zadd(key []byte, score float64, member []byte) {
	if r.meta == nil {
		return
	}
	meta := r.meta.(*ZSetMeta)
	dkey := db.DataKey(r.db, meta.ID)
	iterKey := ZsetItemKey(dkey, member)
	bytesScore, err := db.EncodeFloat64(score)
	if err != nil {
		zap.L().Info("write zset encode score err", zap.Error(err))
		return
	}
	var kvs []common.KvPair
	kvs = append(kvs, common.KvPair{Key: iterKey, Val: bytesScore})
	scoreKey := ZsetScoreKey(dkey, bytesScore, member)
	kvs = append(kvs, common.KvPair{Key: scoreKey, Val: db.NilValue})

	if err := r.write(kv.MakeRowsFromKvPairs(kvs)); err != nil {
		zap.L().Error("write zset iter err", zap.String("key", string(key)), zap.Error(err))
	}

	meta.Len += 1
}

// EndZSet is called when there are no more members in a sorted set.
func (r *RdbDecode) EndZSet(key []byte) {
	if r.meta == nil {
		return
	}
	meta := r.meta.(*ZSetMeta)
	metaKey := db.MetaKey(r.db, key)
	var kvs []common.KvPair
	kvs = append(kvs, common.KvPair{Key: metaKey, Val: meta.Encode()})
	if meta.ExpireAt > 0 {
		if ekey, err := ExpireKey(metaKey, meta.ExpireAt); err == nil {
			kvs = append(kvs, common.KvPair{Key: ekey, Val: meta.ID})
		}
	}
	if err := r.write(kv.MakeRowsFromKvPairs(kvs)); err != nil {
		zap.L().Error("write zset meta err", zap.String("key", string(key)), zap.Error(err))
	}
}

// EndDatabase is called at the end of a database.
func (r *RdbDecode) EndDatabase(n int) {}

// EndRDB is called when parsing of the RDB file is complete.
func (r *RdbDecode) EndRDB() {
}
