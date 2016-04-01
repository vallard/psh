package main

import (
	"testing"

	"github.com/vallard/psh/nr"
)

func TestBadCommand(t *testing.T) {
	nodes, err := nr.GetNodeRange("cosa")
	if err != nil {
		t.Errorf("%v\n", err)
	}
	// open the file for ssh stuff.
	cmd := "this-is-not-a-real-command-i-hope.sh"
	//fmt.Printf("Running command: %s\n", cmd)
	ch := make(chan string, maxSSHSessions)

	wg.Add(len(nodes))
	for _, server := range nodes {
		go runcmd(server, cmd, ch)
	}

	wg.Wait()
}
