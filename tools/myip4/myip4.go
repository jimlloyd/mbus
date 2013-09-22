package main

import (
	"fmt"
	"github.com/jimlloyd/mbus/utils"
)

func main() {

	host, err := utils.MyIp4()
	if err != nil {
		fmt.Println("Error getting IP:", err)
		return
	}

	fmt.Println(host)
}