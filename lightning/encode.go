package lightning

import (
	"encoding/binary"
	"math"

	"github.com/distributedio/titan/db"
)

var (
	expireKeyPrefix = []byte("$sys:0:at:")
)

func ExpireKey(key []byte, ts int64) ([]byte, error) {
	var buf []byte
	buf = append(buf, expireKeyPrefix...)
	encode, err := db.EncodeInt64(ts)
	if err != nil {
		return nil, err
	}
	buf = append(buf, encode...)
	buf = append(buf, ':')
	buf = append(buf, key...)
	return buf, nil
}

func HashItemKey(dbInfo *db.DB, meta *HashMeta, field []byte) []byte {
	b := db.DataKey(dbInfo, meta.ID)
	b = append(b, ':')
	return append(b, field...)
}

func SetItemKey(dbInfo *db.DB, meta *SetMeta, member []byte) []byte {
	key := db.DataKey(dbInfo, meta.ID)
	key = append(key, ':')
	key = append(key, member...)
	return key
}

func LListItemKey(dbInfo *db.DB, meta *LListMeta, index float64) ([]byte, error) {
	key := db.DataKey(dbInfo, meta.ID)
	key = append(key, ':')
	encode, err := db.EncodeFloat64(index)
	if err != nil {
		return nil, err
	}
	key = append(key, encode...)
	return key, nil
}

func ZsetItemKey(dkey []byte, member []byte) []byte {
	var memberKey []byte
	memberKey = append(memberKey, dkey...)
	memberKey = append(memberKey, ':', 'M', ':')
	memberKey = append(memberKey, member...)
	return memberKey
}

func ZsetScoreKey(dkey []byte, score []byte, member []byte) []byte {
	var scoreKey []byte
	scoreKey = append(scoreKey, dkey...)
	scoreKey = append(scoreKey, ':', 'S', ':')
	scoreKey = append(scoreKey, score...)
	scoreKey = append(scoreKey, ':')
	scoreKey = append(scoreKey, member...)
	return scoreKey
}

type Meta interface {
	Encode() []byte
	Decode([]byte) error
}

type StringMeta struct {
	db.Object
	Value []byte
}

func NewStringMeta() *StringMeta {
	now := db.Now()
	return &StringMeta{
		Object: db.Object{
			ID:        db.UUID(),
			Type:      db.ObjectString,
			Encoding:  db.ObjectEncodingRaw,
			CreatedAt: now,
			UpdatedAt: now,
			ExpireAt:  0,
		},
	}
}

func (meta *StringMeta) Encode() []byte {
	b := db.EncodeObject(&meta.Object)
	b = append(b, meta.Value...)
	return b
}

func (meta *StringMeta) Decode(b []byte) error {
	obj, err := db.DecodeObject(b)
	if err != nil {
		return err
	}

	if obj.Type != db.ObjectString {
		return db.ErrTypeMismatch
	}

	if obj.Encoding != db.ObjectEncodingRaw {
		return db.ErrTypeMismatch
	}
	meta.Object = *obj
	if len(b) >= db.ObjectEncodingLength {
		meta.Value = b[db.ObjectEncodingLength:]
	}
	return nil
}

func NewHashMeta() *HashMeta {
	now := db.Now()
	return &HashMeta{
		Object: db.Object{
			ID:        db.UUID(),
			Type:      db.ObjectHash,
			Encoding:  db.ObjectEncodingHT,
			CreatedAt: now,
			UpdatedAt: now,
			ExpireAt:  0,
		},
	}
}

type HashMeta struct {
	db.Object
}

func (meta *HashMeta) Encode() []byte {
	return db.EncodeObject(&meta.Object)
}

//DecodeHashMeta decode meta data into meta field
func (meta *HashMeta) Decode(b []byte) error {
	obj, err := db.DecodeObject(b)
	if err != nil {
		return err
	}
	meta.Object = *obj
	return nil
}

func NewLListMeta() *LListMeta {
	now := db.Now()
	return &LListMeta{
		Object: db.Object{
			ID:        db.UUID(),
			Type:      db.ObjectList,
			Encoding:  db.ObjectEncodingLinkedlist,
			CreatedAt: now,
			UpdatedAt: now,
			ExpireAt:  0,
		},
		Len:    0,
		Lindex: 0,
		Rindex: 0,
	}
}

type LListMeta struct {
	db.Object
	Len    int64
	Lindex float64
	Rindex float64
}

// Marshal encodes meta data into byte slice
func (meta *LListMeta) Encode() []byte {
	b := db.EncodeObject(&meta.Object)
	m := make([]byte, 24)
	binary.BigEndian.PutUint64(m[0:8], uint64(meta.Len))
	binary.BigEndian.PutUint64(m[8:16], math.Float64bits(meta.Lindex))
	binary.BigEndian.PutUint64(m[16:24], math.Float64bits(meta.Rindex))
	return append(b, m...)
}

// Unmarshal parses meta data into meta field
func (meta *LListMeta) Decode(b []byte) (err error) {
	if len(b[db.ObjectEncodingLength:]) != 24 {
		return db.ErrInvalidLength
	}
	m := b[db.ObjectEncodingLength:]
	meta.Len = int64(binary.BigEndian.Uint64(m[:8]))
	meta.Lindex = math.Float64frombits(binary.BigEndian.Uint64(m[8:16]))
	meta.Rindex = math.Float64frombits(binary.BigEndian.Uint64(m[16:24]))
	return nil
}

func NewSetMeta() *SetMeta {
	now := db.Now()
	return &SetMeta{
		Object: db.Object{
			ID:        db.UUID(),
			Type:      db.ObjectSet,
			Encoding:  db.ObjectEncodingHT,
			CreatedAt: now,
			UpdatedAt: now,
			ExpireAt:  0,
		},
		Len: 0,
	}
}

type SetMeta struct {
	db.Object
	Len int64
}

func (meta *SetMeta) Encode() []byte {
	b := db.EncodeObject(&meta.Object)
	m := make([]byte, 8)
	binary.BigEndian.PutUint64(m[:8], uint64(meta.Len))
	return append(b, m...)

}

func (meta *SetMeta) Decode(b []byte) error {
	obj, err := db.DecodeObject(b)
	if err != nil {
		return err
	}
	if obj.Type != db.ObjectSet {
		return db.ErrTypeMismatch
	}

	m := b[db.ObjectEncodingLength:]
	if len(m) != 8 {
		return db.ErrInvalidLength
	}
	meta.Object = *obj
	meta.Len = int64(binary.BigEndian.Uint64(m[:8]))
	return nil
}

func NewZSetMeta() *ZSetMeta {
	now := db.Now()
	return &ZSetMeta{
		Object: db.Object{
			ID:        db.UUID(),
			Type:      db.ObjectZSet,
			Encoding:  db.ObjectEncodingHT,
			CreatedAt: now,
			UpdatedAt: now,
			ExpireAt:  0,
		},
		Len: 0,
	}
}

type ZSetMeta struct {
	db.Object
	Len int64
}

func (meta *ZSetMeta) Encode() []byte {
	b := db.EncodeObject(&meta.Object)
	m := make([]byte, 8)
	binary.BigEndian.PutUint64(m[:8], uint64(meta.Len))
	return append(b, m...)

}

func (meta *ZSetMeta) Decode(b []byte) error {
	obj, err := db.DecodeObject(b)
	if err != nil {
		return err
	}
	if obj.Type != db.ObjectZSet {
		return db.ErrTypeMismatch
	}

	m := b[db.ObjectEncodingLength:]
	if len(m) != 8 {
		return db.ErrInvalidLength
	}
	meta.Object = *obj
	meta.Len = int64(binary.BigEndian.Uint64(m[:8]))
	return nil
}
