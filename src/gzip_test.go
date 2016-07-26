package common

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"testing"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func getSeries() string {
	x := ""
	for i := 0; i < 5; i++ {
		x += `stream_charge,hub=test,method=publish,nodeId=cs8,type=rtmp,uid=1380528809 audio="672",audioPerSecond="[1 62 47 46 47 47 47 47 47 47 47 46 47 47 47]",bytes=1978677.0,bytesPerSecond="[359 112914 177883 128467 121710 136394 177293 129999 130281 126589 174938 129498 132758 128851 170743]",data="1",dataPerSecond="[1 0 0 0 0 0 0 0 0 0 0 0 0 0 0]",domain="10.200.20.28",duration="14.554434268000001",durationPerSecond="[0.554532272 0.9999290780000001 0.9999652840000001 1.000078423 0.9999213600000001 1.000022566 1.000040598 0.999971373 0.9999826100000001 1.000068104 0.9998950710000001 1.000004277 1.000106385 0.9999182070000001 0.99999866]",isV1=true,reqId="1ZmfF79ATM2C6Fj7",status="connected",streamId="test:test/fsf22",version="v1",video="344",videoPerSecond="[1 17 26 25 25 24 26 25 25 25 25 24 25 26 25]" 
        `
	}
	return x
}

func Benchmark_Raw(b *testing.B) {
	x := getSeries()
	for n := 0; n < b.N; n++ {
		_, err := httpRequest("http://192.168.201.81:8086/write?db=foo_test", "POST", x, "text/plain", "", "")
		if err != nil {
			b.Error(err)
		}
		//fmt.Println(ret)
	}
}

func Benchmark_Gzip(b *testing.B) {
	var gpool GzipPool
	x := getSeries()
	for n := 0; n < b.N; n++ {
		gz, _ := gpool.Gzip(x)
		_, err := httpRequest("http://192.168.201.81:8086/write?db=foo_test", "POST", string(gz), "text/plain", "Content-Encoding", "gzip")
		if err != nil {
			b.Error(err)
		}
		//fmt.Println(ret)
	}
}

func Benchmark_UnGzip(b *testing.B) {
	var gpool GzipPool
	var bgzip bytes.Buffer
	w := gzip.NewWriter(&bgzip)
	x := getSeries()
	w.Write([]byte(x))
	w.Close()
	s := bgzip.String()
	for n := 0; n < b.N; n++ {
		gpool.UnGzip(s)
	}
}

func httpRequest(url, method, bodyValue, contentType, headKey, headValue string) (ret string, err error) {
	body := strings.NewReader(bodyValue)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}
	req.Header.Add(headKey, headValue)
	req.Header.Set("Content-Type", contentType)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	ret = string(b)
	return
}
