package zk_recipes

import (
	"github.com/go-zookeeper/zk"
	"log"
	"time"
)

type ZkLogger struct {
}

func (*ZkLogger) Printf(msg string, args ...interface{}) {
	log.Printf(msg, args)
}

type Config struct {
	Servers        []string
	SessionTimeout time.Duration
	Stop           chan bool
}

type treeNode struct {
	data     *[]byte
	children map[string]*treeNode
}

type TreeCache struct {
	config *Config
	conn   *zk.Conn
	root   *treeNode
}

func New(config *Config) *TreeCache {
	return &TreeCache{config: config}
}
