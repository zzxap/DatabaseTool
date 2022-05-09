package dal

import (
	//"DataBaseManage/public"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	//"github.com/etcd-io/etcd/api/v3rpc/rpctypes"
	"github.com/etcd-io/etcd/clientv3"
)

//ETCDCTL_API=3 etcdctl get "" --prefix=true
//ETCDCTL_API=3 etcdctl get "" --from-key

var cli *clientv3.Client
var err error

func InitETCD() bool {
	//log.Println("inidb")
	etcdservers := Dbhost + ":" + Dbport
	array := strings.Split(etcdservers, ",")
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   array,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Println("init etcd fail")
		fmt.Println(err)
		return false
	}
	log.Println("init etcd success")
	//defer cli.Close()
	return true
}
func GetETCD() *clientv3.Client {
	if cli == nil {
		InitETCD()
	}
	if cli != nil {
		//log.Println("client is not nil")
		return cli
	}
	log.Println("client is nil")
	return nil
}
func PutETCD(key, value string) bool {
	//fmt.Println("put key=" + key + " value=" + value)
	if len(value) == 0 || len(key) == 0 {
		return false
	}

	GetETCD()
	if cli == nil {
		fmt.Println("db init error")
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(10))
	_, err := cli.Put(ctx, key, value)
	cancel()
	if err != nil {
		switch err {
		case context.Canceled:
			fmt.Println("ctx is canceled by another routine: %v", err)
		case context.DeadlineExceeded:
			fmt.Println("ctx is attached with a deadline is exceeded: %v", err)
		//case rpctypes.ErrEmptyKey:
		//	fmt.Println("client-side error: %v", err)
		default:
			fmt.Println("bad cluster endpoints, which are not etcd servers: %v", err)
		}
		return false
	}

	//fmt.Println("success put key=" + key + " value=" + value)
	return true
}
func GetMap(key string) map[string]string {
	//fmt.Println("getmap key=" + key)

	GetETCD()

	kv := make(map[string]string)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(5))
	resp, err := cli.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	cancel()
	if err != nil {
		log.Println("err %v", err)
		return kv
	}

	for _, ev := range resp.Kvs {

		kv[string(ev.Key)] = string(ev.Value)

	}
	return kv

}

func GetMapArray(key string) []map[string]string {
	GetETCD()
	//fmt.Println("get key=" + key)
	//kv := make(map[string]string)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(5))
	resp, err := cli.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	cancel()
	if err != nil {
		log.Println("err %v", err)
		return nil
	}
	final_result := make([]map[string]string, 0)

	//log.Println(resp.Kvs)
	for _, ev := range resp.Kvs {
		m := make(map[string]string)
		m["key"] = string(ev.Key) // strings.Replace(string(ev.Key), key, "", -1)
		m["value"] = string(ev.Value)

		final_result = append(final_result, m)

	}
	return final_result

}
func DeletePrefix(key string) bool {
	GetETCD()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(5))
	_, err := cli.Delete(ctx, key, clientv3.WithPrefix()) //
	//withPrefix()是未了获取该key为前缀的所有key-value
	cancel()

	if err != nil {
		return false
	}

	return true

}
func Delete(key string) bool {
	GetETCD()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(5))
	_, err := cli.Delete(ctx, key) //, clientv3.WithPrefix()
	//withPrefix()是未了获取该key为前缀的所有key-value
	cancel()

	if err != nil {
		return false
	}

	return true

}

//WithPrevKV
func Watch(key string) {
	wc := cli.Watch(context.Background(), key, clientv3.WithPrefix(), clientv3.WithPrefix())
	for v := range wc {
		if v.Err() != nil {
			//panic(err)
		}
		for _, e := range v.Events {
			fmt.Printf("type:%v\n kv:%v  prevKey:%v  ", e.Type, e.Kv, e.Kv.Value)
		}
	}
}
