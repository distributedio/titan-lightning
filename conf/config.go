package conf

import "time"

type Import struct {
	Backend           Backend       `cfg:"backend"`
	Security          Security      `cfg:"security"`
	SwitchModInterval time.Duration `cfg:"switch-mod-interval;20m;;switch mod tick interval"`
	PdAddrs           string        `cfg:"pd-addrs; mocktikv://; ;pd address in tidb"`
	Logger            Logger        `cfg:"logger"`
	PIDFileName       string        `cfg:"pid-filename; titan.pid; ; the file name to record connd PID"`
}

type Backend struct {
	MaxOpenFile    uint64 `cfg:"max-open-file;6;;max opened file num"`
	ReginSplitSize string `cfg:"regin-split-size; 96M; ; regin split size"`
	SortedDir      string `cfg:"sorted-dir; ./data; ; sorted sstable file path"`
	Concurrency    int    `cfg:"concurrency;16;;concurrency num"`
	SendKVPairs    int    `cfg:"send-kv-pairs;32768;;send kv paris"`
}

type Security struct {
	CAPath   string `toml:"ca-path" json:"ca-path"`
	CertPath string `toml:"cert-path" json:"cert-path"`
	KeyPath  string `toml:"key-path" json:"key-path"`
}

type Logger struct {
	Name       string `cfg:"name; titan; ; the default logger name"`
	Path       string `cfg:"path; logs/titan; ; the default log path (or stdout/stderr)"`
	Level      string `cfg:"level; info; ; log level(debug, info, warn, error, panic, fatal)"`
	Compress   bool   `cfg:"compress; false; boolean; true for enabling log compress"`
	TimeRotate string `cfg:"time-rotate; 0 0 0 * * *; ; log time rotate pattern(s m h D M W)"`
}
