package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	//	url = "http://101.71.89.40:8787/v4/apps/suna/repos/sunarepo/points"
	url = "http://127.0.0.1:8086/write?db=foo_test"
	//bodyType = "Content-Type: application/json"
)

func main() {
	cnt := 0
	client := &http.Client{}
	batch_size := 100
	r := make(chan string)
	for i := 0; i < 20; i++ {
		go func() {
			for {
				body := ""
				for i := 0; i < batch_size; i++ {
					body += "newserie,hello=" + strconv.Itoa(2) + " value=" + strconv.Itoa(i) + " " + strconv.FormatInt(time.Now().UnixNano(), 10) + "\n"

				}
				buffer := bytes.NewBuffer([]byte(body))

				req1, err := http.NewRequest("POST", url, buffer)
				if err != nil {
					fmt.Println(err)
				}

				req1.Header.Set("Content-Type", "text/plain")
				resp, err := client.Do(req1)

				if err != nil {
					fmt.Println("xx", err)
				}

				if resp != nil {
					resp.Body.Close()
				}

				r <- "."
			}
		}()
	}

	for s := range r {
		if s == "." {
			cnt++
		}
		if cnt%1000 == 0 {
			fmt.Printf("%v\n", cnt)
		}
		if cnt == 100000 {
			return
		}

	}
}
