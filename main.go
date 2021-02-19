package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"

	"os"
	ospath "path"
	"time"

	"github.com/distributedio/configo"
	"github.com/nioshield/titan-lightning/conf"
	kvlog "github.com/pingcap/tidb-lightning/lightning/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	rolling "github.com/arthurkiller/rollingwriter"
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
	if err := ConfigureZap(cfg.Logger.Name, cfg.Logger.Path, cfg.Logger.Level,
		cfg.Logger.TimeRotate, cfg.Logger.Compress); err != nil {
		fmt.Printf("create logger failed, %s\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	l, err := NewLightning(ctx, cfg)
	if err != nil {
		zap.L().Error("new lightning err", zap.Error(err))
		return
	}
	if err := l.Run(); err != nil {
		fmt.Println("import err", err)
		return
	}
}

func ConfigureZap(name, path, level, pattern string, compress bool) error {
	writer, err := Writer(path, pattern, compress)
	if err != nil {
		return err
	}

	var lv = zap.NewAtomicLevel()
	switch level {
	case "debug":
		lv.SetLevel(zap.DebugLevel)
	case "info":
		lv.SetLevel(zap.InfoLevel)
	case "warn":
		lv.SetLevel(zap.WarnLevel)
	case "error":
		lv.SetLevel(zap.ErrorLevel)
	case "panic":
		lv.SetLevel(zap.PanicLevel)
	case "fatal":
		lv.SetLevel(zap.FatalLevel)
	default:
		return fmt.Errorf("unknown log level(%s)", level)
	}
	timeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Local().Format("2006-01-02 15:04:05.999999999"))
	}

	encoderCfg := zapcore.EncoderConfig{
		NameKey:        "Name",
		StacktraceKey:  "Stack",
		MessageKey:     "Message",
		LevelKey:       "Level",
		TimeKey:        "TimeStamp",
		CallerKey:      "Caller",
		EncodeTime:     timeEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	output := zapcore.AddSync(writer)
	var zapOpts []zap.Option
	zapOpts = append(zapOpts, zap.AddCaller())

	logger := zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), output, lv), zapOpts...)
	logger.Named(name)
	log := logger.With(zap.Int("PID", os.Getpid()))
	zap.ReplaceGlobals(log)
	kvlog.SetAppLogger(log)
	//http change log level
	http.Handle("/titan-ligh/log/level", lv)
	return nil
}

//Writer generate the rollingWriter
func Writer(path, pattern string, compress bool) (io.Writer, error) {
	if path == "stdout" {
		return os.Stdout, nil
	} else if path == "stderr" {
		return os.Stderr, nil
	}
	var opts []rolling.Option
	opts = append(opts, rolling.WithRollingTimePattern(pattern))
	if compress {
		opts = append(opts, rolling.WithCompress())
	}
	dir, filename := ospath.Split(path)
	opts = append(opts, rolling.WithLogPath(dir), rolling.WithFileName(filename), rolling.WithLock())
	writer, err := rolling.NewWriter(opts...)
	if err != nil {
		return nil, fmt.Errorf("create IOWriter failed, %s", err)
	}
	return writer, nil
}
