package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

/*
tips:
在用clientv3连接etcd时 undefined: resolver.BuildOption

google.golang.org/grpc 1.26后的版本是不支持clientv3的。
要把这个改成1.26版本的就可以了。
具体操作方法是在go.mod里加上：
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

http://t.zoukankan.com/anmutu-p-etcd.html
 */

func main(){
	var (
		config clientv3.Config
		err error
		client *clientv3.Client
		kv clientv3.KV
		putResp *clientv3.PutResponse
		getResp *clientv3.GetResponse
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
	//fmt.Printf("%v",client)

	// 写etcd的键值对
	kv = clientv3.NewKV(client)
	putResp, err = kv.Put(context.TODO(),"/aaa/bbb/ccc","bye",clientv3.WithPrevKV())
	if err != nil{
		fmt.Printf("kv put error: %v",err)
	} else {
		//获取版本信息
		fmt.Println("Revision:",putResp.Header.Revision)
		if putResp.PrevKv != nil{
			fmt.Println("key:",string(putResp.PrevKv.Key)) // key: /aaa/bbb/ccc
			fmt.Println("Value:",string(putResp.PrevKv.Value)) // Value: bye
			fmt.Println("Version:",string(putResp.PrevKv.Version))
		}
	}

	// 读取etcd的键值对
	kv = clientv3.NewKV(client)
	getResp,err = kv.Get(context.TODO(),"/aaa/bbb/ccc")
	if err != nil {
		fmt.Printf("kv get error: %v",err)
		return
	}
	fmt.Printf("kvs: %v \n",getResp.Kvs) // [key:"/aaa/bbb/ccc" create_revision:2 mod_revision:3 version:2 value:"bye" ]
	fmt.Printf("kvs.length: %v\n", len(getResp.Kvs))

	// WithCountOnly
	kv = clientv3.NewKV(client)
	getResp,err = kv.Get(context.TODO(),"/aaa/bbb/ccc",clientv3.WithCountOnly())
	if err != nil {
		fmt.Printf("WithCountOnly err: %v",err)
		return
	}
	fmt.Println(getResp.Kvs,getResp.Count) // [],1

	// 读取前缀
	kv = clientv3.NewKV(client)
	getResp,err = kv.Get(context.TODO(),"/aaa/bbb",clientv3.WithPrefix())
	if err != nil {
		fmt.Printf("WithPrefix err: %v",err)
		return
	}
	fmt.Println(getResp.Kvs) // [key:"/aaa/bbb/ccc" create_revision:6 mod_revision:7 version:2 value:"bye" ]

	// 删除
	//用于读写etcd的键值对
	kv = clientv3.NewKV(client)
	delResp,err := kv.Delete(context.TODO(),"/aaa/bbb/ccc",clientv3.WithPrevKV())
	if err != nil{
		fmt.Println(err)
		return
	}else{
		if len(delResp.PrevKvs) > 0 {
			for idx,kvpair := range delResp.PrevKvs{
				fmt.Printf("idx: %v key: %v value: %v",idx,string(kvpair.Key),string(kvpair.Value)) // idx: 0 key: /aaa/bbb/ccc value: bye
			}
		}
	}
}
