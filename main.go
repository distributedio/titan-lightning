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
	l, err := NewLightning(ctx, cfg)
	if err != nil {
		fmt.Println("new lightning err", err)
		return
	}
	if err := l.Run(); err != nil {
		fmt.Println("import err", err)
		return
	}
}
