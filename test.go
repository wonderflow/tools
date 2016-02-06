package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	url = "http://101.71.89.40:8787/v4/apps/suna/repos/sunarepo/points"
	//bodyType = "Content-Type: application/json"
)

func main() {
	cnt := 0
	client := &http.Client{}
	r := make(chan string)
	for i := 0; i < 8; i++ {
		go func() {
			for {
				body := "newserie,hello=" + strconv.Itoa(2) + " value=" + strconv.Itoa(1) + " " + strconv.FormatInt(time.Now().UnixNano(), 10)
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

	}
}
