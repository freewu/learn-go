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
		client *clientv3.Client
		err error
		kv clientv3.KV
		putOp clientv3.Op
		getOp clientv3.Op
		opResp clientv3.OpResponse
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

	// 创建Op: operation
	putOp = clientv3.OpPut("/aaa/bbb/op", "opopop")
	// 执行OP
	if opResp, err = kv.Do(context.TODO(), putOp); err != nil {
		fmt.Println(err)
		return
	}

	// kv.Do(op)
	// kv.Put
	// kv.Get
	// kv.Delete

	fmt.Println("写入Revision:", opResp.Put().Header.Revision)

	// 创建Op
	getOp = clientv3.OpGet("/aaa/bbb/op")

	// 执行OP
	if opResp, err = kv.Do(context.TODO(), getOp); err != nil {
		fmt.Println(err)
		return
	}

	// 打印
	fmt.Println("数据Revision:", opResp.Get().Kvs[0].ModRevision)    // create rev == mod rev
	fmt.Println("数据value:", string(opResp.Get().Kvs[0].Value))
}