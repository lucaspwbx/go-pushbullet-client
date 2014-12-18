package main

import (
	"fmt"
	"log"

	"github.com/lucasweiblen/pushbulletclient/client"
)

func main() {
	//200
	cli := client.NewClient("swbpcaTIjyV5eAYZnjfL2GZqFiiqrBHH")
	user, err := cli.GetMe()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(user)

	//400
	cli2 := client.NewClient("foo")
	_, err = cli2.GetMe()
	if err != nil {
		fmt.Println(err)
	}
}
