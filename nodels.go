package main

import (
	"fmt"
	"os"

	"github.com/vallard/psh/nr"
)

// prints the usage of the tool and exits with the value past in.
func useage(rc int) {
	fmt.Printf("nodels [noderange]\nReturns list of nodes in the .psh config file\n")
	os.Exit(rc)
}

func main() {
	if len(os.Args) > 2 {
		useage(1)
	}
	n := ""

	if len(os.Args) > 1 {
		n = os.Args[1]
	}
	// the first argument is the range or group of nodes
	nodes, err := nr.GetNodeRange(n)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	for _, node := range nodes {
		fmt.Println(node.Host)
	}
}
