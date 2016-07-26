package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	dir := "/root/.influx/data"
	fmt.Println(filepath.Dir(dir))
}
