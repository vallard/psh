package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/vallard/psh/server"
)

var configFiles = []string{
	"~/.psh",
}

var usr, _ = user.Current()
var dir = usr.HomeDir

func errorFunction(errMessage string) {
	fmt.Printf("%s\n", errMessage)
	panic(1)
}

func main() {

	// the first argument is the range or group of nodes
	nodes := getNodeRange(os.Args[1])

	// open the file for ssh stuff.
	cmd := getCommand(os.Args[2:])
	fmt.Printf("Running command: %s", cmd)
	servers, err := openConfig()
	if err != nil {
		errorFunction(err.Error())
	}

	for _, server := range servers {
		fmt.Println(server.Host)
	}

}

// get the command from the command line
func getCommand(arr []string) string {
	return strings.Join(arr, " ")
}

// substitute params for files
func expandShell(f string) string {
	f = strings.Replace(f, "~", dir, 1)
	return f
}

// open the config file and return all the servers from it.
func openConfig() ([]server.Server, error) {
	// look for config file
	var servers []server.Server
	for _, filename := range configFiles {
		// expand shell params
		filename = expandShell(filename)
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}

		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			server, err := parseLine(scanner.Text())
			if err != nil {
				continue
			}
			servers = append(servers, server)
		}
	}
	return servers, nil
}

func parseLine(str string) (server.Server, error) {
	params := strings.Split(str, ",")
	s := server.Server{Host: params[0], User: params[1], Key: params[2]}
	return s, nil
}
