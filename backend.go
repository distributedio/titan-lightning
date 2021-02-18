package main

import (
	"context"
	"errors"
	"math"

	"github.com/docker/go-units"
	"github.com/nioshield/titan-lightning/conf"
	kv "github.com/pingcap/tidb-lightning/lightning/backend"
	"github.com/pingcap/tidb-lightning/lightning/common"
)

type Backend struct {
	b   kv.Backend
	cfg *conf.Backend
}

func NewBackend(ctx context.Context, cfg *conf.Backend) (*Backend, error) {
	tls, err := common.NewTLS(cfg.Security.CAPath, cfg.Security.CertPath, cfg.Security.KeyPath, "")
	var rLimit uint64
	rLimit, err = kv.GetSystemRLimit()
	if err != nil {
		return nil, err
	}
	maxOpenFiles := int(rLimit / uint64(cfg.MaxOpenFile))
	if maxOpenFiles < 0 {
		maxOpenFiles = math.MaxInt32
	}
	reginSplitSize, err := units.RAMInBytes(string(cfg.ReginSplitSize))
	if err != nil {
		return nil, err
	}
	bk, err := kv.NewLocalBackend(ctx, tls, cfg.PdAddrs, reginSplitSize,
		cfg.SortedDir, cfg.Concurrency, cfg.SendKVPairs,
		false, nil, maxOpenFiles)
	if err != nil {
		return nil, errors.New("build local backend failed")
	}
	return &Backend{
		cfg: cfg,
		b:   bk,
	}, nil
}

func (bk *Backend) OpenEngine(ctx context.Context, preKey string, enginID int32) (*kv.OpenedEngine, error) {
	return bk.b.OpenEngine(ctx, preKey, enginID)
}

//OpenEngine
// localwrite-> WriteRows
// close

//import
//clean
//compatch
//switchmode
