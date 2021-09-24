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

	// 自动续租
	kv := clientv3.NewKV(client)
	keepRestChan,err := lease.KeepAlive(context.TODO(),leaseid)
	if err != nil{
		fmt.Printf("lease keep-alive error:%v ",err)
		return
	}
	//处理续租应答的协程
	go func() {
		for {
			select {
			case keepresp := <-keepRestChan:
				fmt.Printf("keepresp :%v ",keepresp)
				if keepRestChan == nil{
					fmt.Println("租约已失效了")
					goto END
				}else{// 每秒会续租一次，所以就会收到一次应答
					fmt.Println("收到自动续租的应答")
				}
			}
		}
	END:
	}()

	// put一个kv 让它与租约关联起来 从而实现10秒自动过期
	putResp,err := kv.Put(context.TODO(),"/aaa/keep-alive","v5",clientv3.WithLease(leaseid))
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println("写入成功",putResp.Header.Revision)

	// 定时的看一下key过期了没有
	for{
		getResp,err := kv.Get(context.TODO(),"/aaa/keep-alive")
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