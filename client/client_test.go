package client

import (
	"fmt"
	"testing"
)

func TestNewClient(t *testing.T) {
	cli := NewClient("teste")
	fmt.Println(cli)
}
