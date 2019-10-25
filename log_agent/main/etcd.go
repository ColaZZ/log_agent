package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	etcd_client "go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"logagent/log_agent/tailf"
	"strings"
	"time"
)

type EtcdClient struct {
	client *etcd_client.Client
	keys   []string
}

var (
	etcdClient *EtcdClient
)

func initEtcd(addr string, key string) (collectConf []tailf.CollectConf, err error) {
	client, err := etcd_client.New(etcd_client.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
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

	for _, ip := range localIPArray {
		etcdKey := fmt.Sprintf("%s%s", key, ip)
		etcdClient.keys = append(etcdClient.keys, etcdKey)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := client.Get(ctx, etcdKey)
		if err != nil {
			logs.Error("client get from etcd faild", err)
			continue
		}
		cancel()

		for _, v := range resp.Kvs {
			if string(v.Key) == etcdKey {
				err = json.Unmarshal(v.Value, &collectConf)
				if err != nil {
					logs.Error("unmarshal failed ", err)
					continue
				}
				logs.Debug("log is %v ", collectConf)
			}
		}
		//fmt.Println(resp)
		initEtcdWatcher()
	}
	return
}

func initEtcdWatcher() {
	for _, key := range etcdClient.keys {
		go watchKey(key)
	}
}

func watchKey(key string) {
	cli, err := etcd_client.New(etcd_client.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		logs.Error("etcd connect failed ", err)
		return
	}

	logs.Debug("begin watch key:%s", key)
	for {
		rch := cli.Watch(context.Background(), key)
		var collectConf []tailf.CollectConf
		var getConfSucc = true
		for wresp := range rch {
			for _, ev := range wresp.Events {
				if ev.Type == mvccpb.DELETE {
					logs.Warn("key[%s]'s config is deleted ", key)
					continue
				}
				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					err = json.Unmarshal(ev.Kv.Value, &collectConf)
					if err != nil {
						logs.Error("key[%s], unmarshal failed %s ", key, err)
						getConfSucc = false
						continue
					}
				}
				logs.Debug("get config from etcd, %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
			if getConfSucc {
				logs.Debug("get config from etcd succ, %v", collectConf)
				tailf.UpdateConfig(collectConf)
			}
		}
	}

}
