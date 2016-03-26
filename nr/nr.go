package nr

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/user"
	"regexp"
	"strings"

	"github.com/vallard/psh/server"
)

var usr, _ = user.Current()
var dir = usr.HomeDir

var configFiles = []string{
	"~/.psh",
}

// GetNodeRange gets the list of servers from an encoded string of nodes.
// The node input is a string like node01-node99.  GetNodeRange will return
// all of the nodes in Server objects.  node01, node02, ... node99
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
	s := server.Server{}
	l := len(params)
	switch {
	case l < 1:
		return s, errors.New("Invalid line in config file")
	case l < 2:
		s.Host = params[0]
	case l < 3:
		s.Host = params[0]
		s.IP = params[1]
	case l < 4:
		s.Host = params[0]
		s.IP = params[1]
		s.User = params[2]
	case l < 5:
		s.Host = params[0]
		s.IP = params[1]
		s.User = params[2]
		s.Key = params[3]
	case l > 4:
		s.Host = params[0]
		s.IP = params[1]
		s.User = params[2]
		s.Key = params[3]
	}
	if s.IP == "" {
		return s, nil
	}
	if s.IP != "" && validIP4(s.IP) {
		return s, nil
	} else {
		e := fmt.Sprintf("IP address is not valid for %s", s.Host)
		return s, errors.New(e)
	}
	return s, nil
}

// from: https://www.socketloop.com/tutorials/golang-validate-ip-address
func validIP4(ipAddress string) bool {
	ipAddress = strings.Trim(ipAddress, " ")
	re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	if re.MatchString(ipAddress) {
		return true
	}
	return false
}
