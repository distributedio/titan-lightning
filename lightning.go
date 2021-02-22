package main

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/nioshield/titan-lightning/conf"
	sstpb "github.com/pingcap/kvproto/pkg/import_sstpb"
	kv "github.com/pingcap/tidb-lightning/lightning/backend"
	"github.com/pingcap/tidb-lightning/lightning/common"
	"github.com/tent/rdb"
	"go.uber.org/zap"
)

const (
	DefaultTable = "defatuTable"
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
	go l.tickerWork(ctx)
	l.switchMode(ctx, sstpb.SwitchMode_Import)
	if err := l.process(ctx); err != nil {
		cancel()
		return err
	}
	cancel()
	l.switchMode(ctx, sstpb.SwitchMode_Normal)
	return nil
}

func (l *Lightning) process(ctx context.Context) error {
	openEngin, err := l.bk.OpenEngine(l.ctx, DefaultTable, 0)
	if err != nil {
		zap.L().Error("open engin failed", zap.Error(err))
		return err
	}
	w, err := openEngin.LocalWriter(ctx)
	if err != nil {
		zap.L().Error("get local writer failed", zap.Error(err))
		return err
	}
	callbak := NewRdbDecode(ctx, w, l.cfg.NameSpace)
	f, err := l.reader()
	if err != nil {
		zap.L().Error("get reader failed", zap.Error(err))
		return err
	}
	if err := rdb.Decode(f, callbak); err != nil {
		zap.L().Error("decode failed", zap.Error(err))
		return err
	}
	if err := w.Close(); err != nil {
		zap.L().Error("close writer failed", zap.Error(err))
		return err
	}
	close, err := openEngin.Close(ctx)
	if err != nil {
		zap.L().Error("get close engin failed", zap.Error(err))
		return err
	}
	if err := close.Import(ctx); err != nil {
		zap.L().Error("close engin import failed", zap.Error(err))
		return err
	}
	if err := close.Cleanup(ctx); err != nil {
		zap.L().Error("close engin cleanup failed", zap.Error(err))
		return err
	}
	return nil
}

func (l *Lightning) reader() (io.Reader, error) {
	return os.Open(l.cfg.SourceAddrs)
}
