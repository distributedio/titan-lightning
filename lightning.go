package main

import (
	"context"
	"time"

	"github.com/distributedio/titan/db"
	"github.com/nioshield/titan-lightning/conf"
	sstpb "github.com/pingcap/kvproto/pkg/import_sstpb"
	kv "github.com/pingcap/tidb-lightning/lightning/backend"
	"github.com/pingcap/tidb-lightning/lightning/common"
	"go.uber.org/zap"
)

type Lightning struct {
	ctx context.Context
	cfg *conf.Import
	bk  *Backend
	tls *common.TLS
}

func NewLightning(ctx context.Context, cfg *conf.Import) (*Lightning, error) {
	l := &Lightning{
		ctx: ctx,
		cfg: cfg,
	}
	var err error

	if l.tls, err = common.NewTLS(cfg.Security.CAPath, cfg.Security.CertPath, cfg.Security.KeyPath, cfg.PdAddrs); err != nil {
		zap.L().Error("tlserr", zap.Error(err))
		return nil, err
	}

	if l.bk, err = NewBackend(ctx, &cfg.Backend, l.tls, cfg.PdAddrs); err != nil {
		zap.L().Error("new backerr", zap.Error(err))
		return nil, err
	}
	return l, nil
}

func (l *Lightning) tickerWork(ctx context.Context) {
	modeTicker := time.NewTicker(l.cfg.SwitchModInterval)
	defer modeTicker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-modeTicker.C:
			l.switchMode(ctx, sstpb.SwitchMode_Import)
		}
	}
}

func (l *Lightning) switchMode(ctx context.Context, mode sstpb.SwitchMode) {
	var minState kv.StoreState
	if mode == sstpb.SwitchMode_Import {
		minState = kv.StoreStateOffline
	} else {
		minState = kv.StoreStateDisconnected
	}
	// we ignore switch mode failure since it is not fatal.
	// no need log the error, it is done in kv.SwitchMode already.
	_ = kv.ForAllStores(
		ctx,
		l.tls,
		minState,
		func(c context.Context, store *kv.Store) error {
			return kv.SwitchMode(c, l.tls, store.Address, mode)
		},
	)
}

func (l *Lightning) Run() error {
	ctx, cancel := context.WithCancel(l.ctx)
	defer cancel()
	go l.tickerWork(ctx)
	l.switchMode(ctx, sstpb.SwitchMode_Import)
	openEngin, err := l.bk.OpenEngine(l.ctx, "test", 0)
	if err != nil {
		return err
	}
	w, err := openEngin.LocalWriter(ctx)
	if err != nil {
		return err
	}
	rows := l.kvPars()
	if err := w.WriteRows(ctx, nil, rows); err != nil {
		return err
	}
	close, err := openEngin.Close(ctx)
	if err != nil {
		return err
	}
	if err := close.Import(ctx); err != nil {
		return err
	}
	if err := close.Cleanup(ctx); err != nil {
		return err
	}
	l.switchMode(ctx, sstpb.SwitchMode_Normal)
	return nil
}

func (l *Lightning) kvPars() kv.Rows {
	d := &db.DB{Namespace: "default", ID: db.DBID(0)}
	key := db.MetaKey(d, []byte("strkey"))
	now := db.Now()
	obj := &db.Object{
		CreatedAt: now,
		UpdatedAt: now,
		ExpireAt:  0,
		ID:        db.UUID(),
		Type:      db.ObjectString,
		Encoding:  db.ObjectEncodingRaw,
	}
	val := db.EncodeObject(obj)
	val = append(val, []byte("testval")...)
	result := make([]common.KvPair, 0)
	result = append(result, common.KvPair{Key: key, Val: val})
	zap.L().Info("add kv", zap.String("key", string(key)), zap.String("val", string(val)))
	return kv.MakeRowsFromKvPairs(result)
}
