package nr

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/vallard/psh/server"
)

var usr, _ = user.Current()
var dir = usr.HomeDir

var configFiles = []string{
	"~/.psh",
}

// public function to get the noderange from the config file for a bunch of nodes.
func GetNodeRange(nr string) ([]server.Server, error) {
	elems := strings.Split(nr, ",")
	nodelist, err := nodesFromConfig()
	if err != nil {
		return nil, err
	}

	returnNodes, err := findNodesInList(elems, nodelist)
	if err != nil {
		return nil, err
	}
	return returnNodes, nil
}

// given a bunch of nodes get them from
func findNodesInList(elem []string, nodelist []server.Server) ([]server.Server, error) {
	var returnNodes []server.Server
	for _, e := range elem {
		found := false
		for _, allNode := range nodelist {
			if e == allNode.Host {
				found = true
				returnNodes = append(returnNodes, allNode)
				continue
			}
		}
		if !found {
			rm := fmt.Sprintf("findNodesInList: %s not found in nodelist", e)
			return nil, errors.New(rm)
		}
	}
	return returnNodes, nil
}

// substitute params for files
func ExpandShell(f string) string {
	f = strings.Replace(f, "~", dir, 1)
	return f
}

// open the config file and return all the servers from it.
func nodesFromConfig() ([]server.Server, error) {
	// look for config file
	var servers []server.Server
	for _, filename := range configFiles {
		// expand shell params
		filename = ExpandShell(filename)
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
	s := server.Server{Host: params[0], IP: params[1], User: params[2], Key: params[3]}
	return s, nil
}
