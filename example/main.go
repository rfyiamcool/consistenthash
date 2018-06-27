package main

import (
	"fmt"
	ch "github.com/rfyiamcool/consistenthash"
)

func main() {
	hash := ch.New(3, nil)

	hash.Add("node1", "node2", "node3")

	fmt.Println(hash.Get("key1"))
	fmt.Println(hash.Get("key2"))
	fmt.Println(hash.Get("key3"))
	fmt.Println(hash.Get("key4"))
	fmt.Println(hash.Get("key5"))
	fmt.Println(hash.Get("key6"))
	fmt.Println(hash.Get("key7"))
}
