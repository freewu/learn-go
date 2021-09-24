package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err error
		kv clientv3.KV
		watcher clientv3.Watcher
		getResp *clientv3.GetResponse
		watchStartRevision int64
		watchRespChan <-chan clientv3.WatchResponse
		watchResp clientv3.WatchResponse
		event *clientv3.Event
	)

	// 配置
	config = clientv3.Config{
		Endpoints:[]string{"192.168.150.129:2379","192.168.150.129:2378","192.168.150.129:2377"},
		DialTimeout:time.Second * 5,
	}
	// 连接
	if client,err = clientv3.New(config);err != nil{
		fmt.Printf("connect error: %v",err)
		return
	}

	kv = clientv3.NewKV(client)
	// 模拟etcd中KV的变化
	go func() {
		for {
			_, _ = kv.Put(context.TODO(), "/aaa/bbb/watch", strconv.Itoa(rand.Int()))
			_, _ = kv.Delete(context.TODO(), "/aaa/bbb/watch")
			time.Sleep(5 * time.Second)
		}
	}()

	// 先GET到当前的值，并监听后续变化
	if getResp, err = kv.Get(context.TODO(), "/aaa/bbb/watch"); err != nil {
		fmt.Println(err)
		return
	}

	// 现在key是存在的
	if len(getResp.Kvs) != 0 {
		fmt.Println("当前值:", string(getResp.Kvs[0].Value))
	}

	// 当前etcd集群事务ID, 单调递增的
	watchStartRevision = getResp.Header.Revision + 1

	// 创建一个watcher
	watcher = clientv3.NewWatcher(client)
	// 启动监听
	fmt.Println("从该版本向后监听:", watchStartRevision)
	ctx, _ := context.WithCancel(context.TODO())
	//ctx, cancelFunc := context.WithCancel(context.TODO())
	//time.AfterFunc(5 * time.Second, func() {
	//	cancelFunc()
	//})

	watchRespChan = watcher.Watch(ctx, "/aaa/bbb/watch", clientv3.WithRev(watchStartRevision))

	// 处理kv变化事件
	for watchResp = range watchRespChan {
		for _, event = range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println(string(event.Kv.Key)," 修改为:", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println(string(event.Kv.Key)," 被删除了", "Revision:", event.Kv.ModRevision)
			}
		}
	}
}