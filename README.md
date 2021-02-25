# Titan Lightning
**Titan Lightning** is a tool for fast full import of large amounts of data into a [Titan](https://github.com/distributedio/titan) cluster.

it implements parsing and encoding of the [Redis](http://redis.io) [RDB file format](https://github.com/sripathikrishnan/redis-rdb-tools/blob/master/docs/RDB_File_Format.textile).

This tool was heavily inspired by [TiDB Lightning](https://github.com/pingcap/tidb-lightning).

## Installation

### SetUp TiKV cluster

Titan works with 2 TiDB components:

* TiKV
* PD

To setup TiKV and PD, please follow the official [instructions](https://pingcap.com/docs-cn/dev/how-to/deploy/orchestrated/ansible/)

### Run Titan

* Build the binary

```
go get github.com/distributedio/titan
cd $GOPATH/src/github.com/distributedio/titan
make
```

* Edit the configration file

```
pd-addrs="tikv://your-pd-addrs:port"
```

* Run Titan

```
./titan
```

### Run Titan-Lightning

* Build the binary

```
go get github.com/nioshield/titan-lightning
cd $GOPATH/src/github.com/nioshield/titan-lightning
make
```

* Edit the configration file

```
pd-addrs="your-pd-addrs:port"
source-addrs = "./dump.rdb"
```

* Run Titan-Lightning

```
./titan-lightning
```
