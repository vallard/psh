package nr

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/user"
	"regexp"
	"strconv"
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
	// get all the nodes in the config file.
	nodelist, err := nodesFromConfig()
	if err != nil {
		return nil, err
	}

	// first seperate by commas
	elems := strings.Split(nr, ",")
	var telems []string
	// go through each of the comma separated values and further refine.
	for _, e := range elems {
		newRange, err := nodesFromDash(e)
		//fmt.Println(elems)
		if err != nil {
			return []server.Server{}, err
		}
		telems = append(telems, newRange...)
	}

	returnNodes, err := findNodesInList(telems, nodelist)
	if err != nil {
		return nil, err
	}
	return returnNodes, nil
}

// nodesFromDash separates a string like node01-node04 and gives back
// an array of node01,node02,node03,node04
// see line 376 and on of:
// https://sourceforge.net/p/xcat/xcat-core/ci/master/tree/perl-xCAT/xCAT/NodeRange.pm
func nodesFromDash(e string) ([]string, error) {
	var nodes []string

	// match node[01-04]
	//rb := regexp.MustCompile("\[\d+-\d+\]")
	rb := regexp.MustCompile("\\[\\d+[-:]\\d+\\]")
	if rb.MatchString(e) {
		suffix := rb.FindString(e)
		pI := rb.FindStringSubmatchIndex(e)
		prefix := e[:pI[0]]
		//fmt.Printf("%s prefix\n", prefix)
		//fmt.Printf("%s suffix\n", suffix)
		// now strip suffix from [001-003]
		first, last := stripSuffixFromBracketRange(suffix)
		//fmt.Println(first, last)
		nodes, err := makeNodesFromSuffixPoints(prefix, first, last)

		return nodes, err
	}
	// do we support node[01-02]-node[04-06] ?  Legal, but seems crazy.
	// m/[-:]/
	r := regexp.MustCompile("[-:]")
	if r.MatchString(e) {
		nodeParts := r.Split(e, -1)
		if len(nodeParts) != 2 {
			eMsg := fmt.Sprintf("Invalid noderange: %s\n", e)
			return nodes, errors.New(eMsg)
		}
		ne := regexp.MustCompile("[0-9]+")
		fn := ne.FindString(nodeParts[0])
		fIndex := ne.FindStringSubmatchIndex(nodeParts[0])
		sn := ne.FindString(nodeParts[1])
		sIndex := ne.FindStringSubmatchIndex(nodeParts[1])
		// make sure the prefix is the same, so numbers start at same place.
		if fIndex[0] != sIndex[0] {
			eMsg := fmt.Sprintf("Invalid noderange: %s\n", e)
			return nodes, errors.New(eMsg)
		}
		//fmt.Println(fIndex[1])
		fPrefix := nodeParts[0][0:fIndex[0]]
		//fmt.Printf("Prefix: %s\n", fPrefix)
		sPrefix := nodeParts[0][0:sIndex[0]]
		// if they put in: node01-mode03 this doesn't work.
		if sPrefix != fPrefix {
			eMsg := fmt.Sprintf("Invalid noderange: %s\n", e)
			return nodes, errors.New(eMsg)
		}
		nodes, err := makeNodesFromSuffixPoints(sPrefix, fn, sn)
		return nodes, err
	}
	nodes = append(nodes, e)
	return nodes, nil
}

// makeNodesFromSuffixPoints takes in the prefix "node" and the first and last integers
// as strings: "01" and "04" then gives a list of nodes: node01, node02, node03, node04
func makeNodesFromSuffixPoints(prefix string, fn string, sn string) ([]string, error) {
	var nodes []string

	if fn == "" || sn == "" {
		eMsg := fmt.Sprintf("Invalid noderange: %s%s-%s%s\n", prefix, fn, prefix, sn)
		return nodes, errors.New(eMsg)
	}
	// if for some reason they said node1-node1 return one node.
	if fn == sn {
		n := fmt.Sprintf("%s%s", prefix, fn)
		nodes = append(nodes, n)
		return nodes, nil
	}
	// if they put 01-4 this is an error
	if len(fn) != len(sn) {
		eMsg := fmt.Sprintf("Invalid noderange: %s%s-%s%s\n", prefix, fn, prefix, sn)
		return nodes, errors.New(eMsg)
	}

	num1, _ := strconv.Atoi(fn)
	num2, _ := strconv.Atoi(sn)
	if num1 > num2 {
		eMsg := fmt.Sprintf("Invalid noderange: %s%s-%s%s\n", prefix, fn, prefix, sn)
		return nodes, errors.New(eMsg)
	}

	for i := num1; i <= num2; i++ {
		nodes = append(nodes, fmt.Sprintf("%s%0[2]*[3]d", prefix, len(fn), i))
	}

	//nodes = append(nodes, e)
	return nodes, nil
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
			rm := fmt.Sprintf("%s not found in nodelist %s", e, configFiles[0])
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
				return servers, err
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
	if s.IP != "" && !validIP4(s.IP) {
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

func stripSuffixFromBracketRange(s string) (string, string) {
	re := regexp.MustCompile("(?P<first>[0-9]+)-(?P<last>[0-9]+)")
	data := re.FindStringSubmatch(s)
	//fmt.Printf("first: %s, last: %s", data[1], data[2])
	return data[1], data[2]
}
