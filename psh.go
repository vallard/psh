package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"

	"github.com/vallard/psh/nr"
	"github.com/vallard/psh/server"
)

const maxSSHSessions int = 20

var wg sync.WaitGroup

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
	ch := make(chan string, maxSSHSessions)

	wg.Add(len(nodes))
	for _, server := range nodes {
		go runcmd(server, cmd, ch)
	}
	/*for range nodes {
		fmt.Print(<-ch)
	}*/
	wg.Wait()
}

func runcmd(server server.Server, cmd string, ch chan<- string) {
	defer wg.Done()
	auth := PublicKeyFile(server.Key)
	if auth == nil {
		fmt.Println("Invalid server key")
		return
	}
	sshConfig := &ssh.ClientConfig{
		User: server.User,
		Auth: []ssh.AuthMethod{
			auth,
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
			ch <- fmt.Sprintf("request for pseudo terminal failed: %s", err)
			return
		}
	*/

	stdout, err := session.StdoutPipe()
	if err != nil {
		//ch <- fmt.Errorf("Unable to setup stdout for session: %v", err)
		ch <- fmt.Sprintf("Unable to setup stdout for session: %v", err)
		return
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		//ch <- fmt.Errorf("Unable to setup stderr for session: %v", err)
		ch <- fmt.Sprintf("Unable to setup stderr for session: %v", err)
		return
	}

	err = session.Run(cmd)
	if err != nil {
		// might need to do something here...
		//fmt.Printf("%s: %v\n", err)
		//return
	}

	// check the stdout for data and display.
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Printf("%s: %s\n", server.Host, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		fmt.Printf("%s: %v\n", server.Host, err)
	}

	// check the error output for errors.
	errScanner := bufio.NewScanner(stderr)
	for errScanner.Scan() {
		fmt.Printf("%s: %s\n", server.Host, errScanner.Text())
	}
	if err = errScanner.Err(); err != nil {
		fmt.Printf("%s: %v\n", server.Host, err)
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
	fmt.Println("File", file)
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("couldn't read file: %s\n", file)
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		fmt.Printf("couldn't parse file: %s\n", file)
		return nil
	}

	return ssh.PublicKeys(key)
}

// get the command from the command line
func getCommand(arr []string) string {
	return strings.Join(arr, " ")
}
