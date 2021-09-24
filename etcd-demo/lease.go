package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	var (
		config clientv3.Config
		err error
		client *clientv3.Client
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

	// 申请一个10秒的租约
	lease := clientv3.NewLease(client)
	leaseGrantResp, err := lease.Grant(context.TODO(),10)
	if err != nil{
		fmt.Printf("lease grant error:%v ",err)
		return
	}

	//拿到租约id
	leaseid := leaseGrantResp.ID

	// put一个kv 让它与租约关联起来 从而实现10秒自动过期
	kv := clientv3.NewKV(client)
	putResp,err := kv.Put(context.TODO(),"/aaa/bbb/ccc","111",clientv3.WithLease(leaseid))
	if err != nil{
		fmt.Printf("kv put error: %v",err)
		return
	}
	fmt.Println("写入成功",putResp.Header.Revision)

	// 定时的查看key 是否过期
	for{
		getResp,err := kv.Get(context.TODO(),"/aaa/bbb/ccc")
		if err != nil{
			fmt.Printf("kv get error: %v",err)
			return
		}
		if getResp.Count == 0{
			fmt.Println("kv过期了")
			break
		}
		fmt.Println("还没过期：",getResp.Kvs)
		time.Sleep(time.Second * 2)
	}
}
