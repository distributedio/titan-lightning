package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/distributedio/configo"
	"github.com/nioshield/titan-lightning/conf"
)

func main() {
	var confPath string
	flag.StringVar(&confPath, "c", "conf/import.toml", "conf file path")
	flag.Parse()

	cfg := &conf.Import{}
	if err := configo.Load(confPath, cfg); err != nil {
		fmt.Printf("unmarshal config file failed, %s\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	back, err := NewBackend(ctx, &cfg.Backend)
	if err != nil {
		fmt.Printf("init backend failed, %s\n", err)
		os.Exit(1)

	}
	_ = back

}
