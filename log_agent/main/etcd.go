package main

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	etcd_client "github.com/coreos/etcd/clientv3"
	"logagent/log_agent/tailf"
	"strings"
	"time"
)

type EtcdClient struct {
	client *etcd_client.Client
	keys    []string
}

func initEtcd(addr string, key string) (err error) {
	client, err := etcd_client.New(etcd_client.Config{
		Endpoints:            []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout:          5 * time.Second,
	})
	if err != nil {
		logs.Error("connect to etcd failed ", err)
		return
	}

	etcdClient := &EtcdClient{
		client: client,
	}

	if strings.HasSuffix(key, "/") == false {
		key = key + "/"
	}

	for _, ip := range localIPArray{
		etcdKey := fmt.Sprintf("%s%s", key, ip)
		etcdClient.keys = append(etcdClient.keys, etcdKey)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := client.Get(ctx, etcdKey)
		if err != nil {
			logs.Error("client get from etcd faild", err)
			continue
		}
		cancel()

		fmt.Println(resp)
	}

}
