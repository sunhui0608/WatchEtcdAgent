// Copyright 2016 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"log"
	"path"

	"go.etcd.io/etcd/clientv3"
	"time"
	"flag"
	"os"
	"io/ioutil"
	"strings"
)


var (
	dialTimeout    = 5 * time.Second
	requestTimeout = 10 * time.Second
	endpoints      = []string{"localhost:2379"}
	help bool
)


func watch() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	rch := cli.Watch(context.Background(), "uin")
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
	// PUT "foo" : "bar"
}

func getAllKeyValue(endpoint [] string,  prefix string, path string) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoint,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	resp, err := cli.Get(context.Background(),prefix, clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
		exist, _ := DirExists(path)
		//length := len(prefix)
		key := string(ev.Key[:])
		value := string (ev.Value[:])
		if !exist {
			os.Mkdir(path, os.ModePerm)
		}
		updateSecret(path, key, value)
	}

}
func watchWithPrefix(endpoint []string, prefix string, path string) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoint,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	rch := cli.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {

			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			exist, _ := DirExists(path)
			//length := len(prefix)
			key := string(ev.Kv.Key[:])
			value := string (ev.Kv.Value[:])
			//if err != nil {
			//	fmt.Printf("get dir error![%v]\n", err)
			//	return
			//}
			if !exist {
				os.Mkdir(path, os.ModePerm)
			}
			if ev.Type.String()  == "PUT" {
				fmt.Printf("Put process\n")
				updateSecret(path, key, value)
			} else if ev.Type.String()  == "DELETE" {
				fmt.Printf("Delete process\n")
				deleteSecret(path, key, value)
			} else {
				fmt.Printf("Do nothing\n")
			}
		}
	}
	// PUT "foo1" : "bar"
}


func DirExists(path string) (bool, error) {
	exist, err := os.Stat(path)
	if err == nil && exist.IsDir() {
		return true, nil
	}

	return false, err
}


func FileExists(path string) (bool, error) {
	exist, err := os.Stat(path)
	if err == nil && !exist.IsDir() {
		return true, nil
	}

	return false, err
}
func deleteSecret(rootPath string, key string, value string) {
	secretPath := path.Join(rootPath, key)
	exist, _ := DirExists(secretPath)

	//if err != nil {
	//	fmt.Printf("get dir error![%v]\n", err)
	//	return
	//}
	if exist {
		fmt.Printf("removing %s\n", secretPath)
		os.RemoveAll(secretPath)
	}

}
func updateSecret(rootPath string, key string, value string) {
	secretPath := path.Join(rootPath, key)
	exist, _ := DirExists(secretPath)

	//if err != nil {
	//	fmt.Printf("get dir error![%v]\n", err)
	//	return
	//}
	if !exist {
		os.Mkdir(secretPath,os.ModePerm)
	}

	writeSecretFile(secretPath, value)

}

/*
COSAccessKeyId=AKID5yc1B6BEwRikX8gaQ1NTIrAE2ay92mFS
COSSecretKey=80e1A3h5FBtNbcPxGPL3ZqthFYdU6TbY
COSAccessToken=109dbb14ca0c30ef4b7e2fc9612f26788cadbfac3
COSAccessTokenExpire=2017-08-29T20:30:00
 */
func writeSecretFile(secretPath string, value string) {
	tempPass := path.Join(secretPath, "/tempPassword")
	//exist, _ := FileExists(tempPass)
	////if err != nil {
	////	fmt.Printf("get file error![%v]\n", err)
	////	return
	////}
	//var f *os.File
	//var err1 error
	//if exist {
	//	fmt.Println("文件存在")
	//	f, err1 = os.OpenFile(tempPass, os.O_WRONLY, 0666) //打开文件
	//} else {
	//	fmt.Println("文件不存在")
	//	f, err1 = os.Create(tempPass) //创建文件
	//}
	//check(err1)
	//n, err1 := io.WriteString(f, value) //写入文件(字符串)
	//fmt.Printf("写入 %d 个字节\n", n)
	//check(err1)
	value = strings.Replace(value, "@@", "\n", -1);
	err := ioutil.WriteFile(tempPass, []byte(value), 0666)
	if err != nil {
		log.Fatal(err)
	}

	//var d1 = []byte(value)
	//err2 := ioutil.WriteFile(t, d1, 0666) //写入文件(字节数组)
	//check(err2)
	//f.Close();

}

//func check(e error) {
//	if e != nil {
//		panic(e)
//	}
//}
func watchWithRange() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// watches within ['foo1', 'foo4'), in lexicographical order
	rch := cli.Watch(context.Background(), "uin1", clientv3.WithRange("uin4"))
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
	// PUT "foo1" : "bar"
	// PUT "foo2" : "bar"
	// PUT "foo3" : "bar"
}

func watchWithProgressNotify() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}

	rch := cli.Watch(context.Background(), "uin", clientv3.WithProgressNotify())
	wresp := <-rch
	fmt.Printf("wresp.Header.Revision: %d\n", wresp.Header.Revision)
	fmt.Println("wresp.IsProgressNotify:", wresp.IsProgressNotify())
	// wresp.Header.Revision: 0
	// wresp.IsProgressNotify: true
}

func usage() {
	fmt.Fprintf(os.Stderr, `Etcd watch version: etcdWatch/1.0.0
Usage: WatchFunction [-h] [-endpoints endpoint] [-prefix prefix] [-path storepath]

Options:
`)
	flag.PrintDefaults()
}


func main() {
	endPoint := flag.String("endpoints", "http://localhost:2379", "EndPoint of etcd server")
	prefix := flag.String("prefix", "uin", "Prefix of watch keys")
	storePath := flag.String("path", "/Users/sunhui/data_etcd", "Pull secret into store path")
	flag.BoolVar(&help, "h", false, "this help")
	flag.Usage = usage
	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	//fmt.Printf("endpoint %s : prefix %s : path %s\n", *endPoint, *prefix, *storePath)
	getAllKeyValue([]string{*endPoint}, *prefix, *storePath)
	watchWithPrefix([]string{*endPoint}, *prefix, *storePath)

}