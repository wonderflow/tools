package main

import (
	//"errors"
	"encoding/json"
	//"fmt"
	"github.com/coreos/etcd/client"
	"github.com/qiniu/log.v1"
	"golang.org/x/net/context"
	//	"strings"
	"time"
)

type Con struct {
	XX string `json:"x"`
	YY string `json:"y"`
}

/*
func GetIndex(client *etcd.Client, jobname string) (int, string, error) {
	jobdir := "/cloudagent/" + jobname
	response, err := client.AddChild(jobdir, jobname, 0)
	if err != nil {
		fmt.Printf("use etcd to get index error: %v\n", err)
		return 0, "", err
	}
	mykey := response.Node.Key
	response, err = client.Get(jobdir, true, true)
	if err != nil {
		fmt.Printf("get etcd jobdir error: %v\n", err)
		return 0, "", err
	}
	for i := 0; i < response.Node.Nodes.Len(); i++ {
		if response.Node.Nodes[i].Key == mykey {
			return i, mykey, nil
		}
	}
	// this line would never reach.
	return 0, "", errors.New("etcd add child error!")
}
*/

func main() {

	// new etcd client
	cfg := client.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// get etcd api from client
	kapi := client.NewKeysAPI(c)

	// generate an etcd value for use,  json code
	etcdValue := &Con{
		XX: "hh",
		YY: "xxxxx",
	}
	log.Print(etcdValue)

	// set "/foo" key with "bar" value
	bar, err := json.Marshal(etcdValue)
	if err != nil {
		log.Error(err)
	}
	log.Print(bar)
	str := string(bar)
	log.Printf("Setting '/foo' key with %v value", str)
	resp, err := kapi.Set(context.Background(), "/xx", str, nil)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Set is done. Metadata is %q\n", resp)
	}

	// get an etcd node value
	// eg. get "/foo" key's value
	log.Print("Getting '/foo' key value")

	var test Con
	resp, err = kapi.Get(context.Background(), "/xx", nil)
	if err != nil {
		log.Fatal(err)
	} else {
		// print common key info
		log.Printf("Get is done. Metadata is %q\n", resp)
		// print value
		log.Printf("%q key has %q value\n", resp.Node.Key, resp.Node.Value)
	}
	json.Unmarshal([]byte(resp.Node.Value), &test)
	log.Printf("%v", test)

	// etcd watcher
	wat := kapi.Watcher("/test", &c.WatcherOptions{AfterIndex: uint64(0)})
	var m map[string]string
	ch := make(chan string)
	go func() {
		for {
			resp, err := wat.Next(context.Background())
			if err != nil {
				log.Println(err)
				break
			} else {
				log.Println(resp.Node.Key, resp.Node.Value)
				m["test"] = resp.Node.Value
				ch <- resp.Node.Value
			}
		}
		close(ch)
	}()
	for v := range ch {
		fmt.Println(v)
	}

	// more :  https://github.com/coreos/etcd/blob/master/client/keys.go

	//machines := []string{"http://localhost:4001"}
	//client := etcd.NewClient(machines)
	/*
		resp, err := client.Get("/jobs", true, true)
		if err != nil {
			fmt.Println("Get etcd job list error: ", err)
		}

		for _, val := range resp.Node.Nodes {
			//fmt.Println(val.Key, "\t", val.Value)
			subres, _ := client.Get(val.Key, true, true)

			fmt.Println(subres.Node.Nodes)
			//fmt.Println(strings.Split(val.Key, "/")[2])
		}

			for {
				if _, err := client.Set("cloud_agent"+"/"+"10.10.101.170", "test", 0); err != nil {
					fmt.Printf("etcd set error: %v.\n", err)
				}
				fmt.Printf("%v beat sent.\n", time.Now())
				time.Sleep(time.Second * 30)
			}
	*/
	//response, err := client.AddChild("/remove1", "10.10.101.101", 30)
	//response, err = client.Get("/remove1", true, true)
	//client.Delete("jobs", true)
	//fmt.Printf("%v %v\n", response.Node.Nodes[0].Value, err)
}
