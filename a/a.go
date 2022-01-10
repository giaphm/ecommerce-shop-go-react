package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	fmt.Println("now.UTC()", now.UTC())

	location, _ := time.LoadLocation("Asia/Ho_Chi_Minh")

	now7 := now.In(location)
	fmt.Println("now7", now7)
}
