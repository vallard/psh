package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/vallard/psh/nr"
	"github.com/vallard/psh/server"
)

func errorFunction(errMessage string) {
	fmt.Printf("%s\n", errMessage)
	os.Exit(1)
}

// prints the usage of the tool and exits with the value past in.
func useage(rc int) {
	fmt.Printf("psh <noderange> command args args ...\n")
	os.Exit(rc)
}

func main() {
	if len(os.Args) < 2 {
		useage(1)
	}

	// the first argument is the range or group of nodes
	nodes, err := nr.GetNodeRange(os.Args[1])
	if err != nil {
		errorFunction(err.Error())
	}

	// open the file for ssh stuff.
	cmd := getCommand(os.Args[2:])
	//fmt.Printf("Running command: %s\n", cmd)

	ch := make(chan string)
	for _, server := range nodes {
		//fmt.Println(server.Host)
		go runcmd(server, cmd, ch)
	}

	for range nodes {
		fmt.Println(<-ch)
	}
}

func runcmd(server server.Server, cmd string, ch chan<- string) {
	sshConfig := &ssh.ClientConfig{
		User: server.User,
		Auth: []ssh.AuthMethod{
			PublicKeyFile(server.Key),
		},
	}

	// create connection
	s := getServerAndPort(server)
	connection, err := ssh.Dial("tcp", s, sshConfig)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	// create session
	session, err := connection.NewSession()
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	defer session.Close()
	// create terminal
	/*
		modes := ssh.TerminalModes{
			ssh.ECHO:          0,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}
			if err = session.RequestPty("xterm", 80, 40, modes); err != nil {
				session.Close()
				ch <- fmt.Sprintf("request for pseudo terminal failed: %s", err)
				return
			}
	*/
	var stdoutBuff bytes.Buffer
	//stdout, err := session.StdoutPipe()
	session.Stdout = &stdoutBuff

	// copy pipe stuff
	/*stdin, err := session.StdinPipe()
	if err != nil {
		//ch <- fmt.Errorf("Unable to setup stdin for session: %v", err)
		ch <- fmt.Sprintf("Unable to setup stdin for session: %v", err)
		return
	}
	go io.Copy(stdin, os.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		//ch <- fmt.Errorf("Unable to setup stdout for session: %v", err)
		ch <- fmt.Sprintf("Unable to setup stdout for session: %v", err)
		return
	}
	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		//ch <- fmt.Errorf("Unable to setup stderr for session: %v", err)
		ch <- fmt.Sprintf("Unable to setup stderr for session: %v", err)
		return
	}
	go io.Copy(os.Stderr, stderr)
	*/
	//log.Printf("Running command on %s", server.IP)
	err = session.Run(cmd)
	if err != nil {
		ch <- fmt.Sprintf("Unable to setup stderr for session: %v", err)
	} else {
		ch <- server.Host + ": " + stdoutBuff.String()
	}
	return
}

func getServerAndPort(server server.Server) string {
	if strings.Contains(server.IP, ":") {
		return server.IP
	}
	s := []string{server.IP, "22"}
	return strings.Join(s, ":")
}

// get the key file. cred: http://blog.ralch.com/tutorial/golang-ssh-connection/
func PublicKeyFile(file string) ssh.AuthMethod {
	file = nr.ExpandShell(file)
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

// get the command from the command line
func getCommand(arr []string) string {
	return strings.Join(arr, " ")
}
