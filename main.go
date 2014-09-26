package main

import (
	"fmt"

	"github.com/lucasweiblen/pushbulletclient/client"
)

func main() {
	cli := client.NewClient("teste")
	resp, code, _ := cli.GetContacts()
	fmt.Println(resp)
	fmt.Println(code)
}
