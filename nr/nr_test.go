package nr

import (
	"testing"
)

func TestExpandShell(t *testing.T) {
	ex := dir + "/foo"
	str := ExpandShell("~/foo")
	if !(ex == str) {
		t.Error(`ExpandShell("~/foo")`)
	}
}

func TestNoExpandShell(t *testing.T) {
	ex := "/home/newton/blast"
	if !(ExpandShell(ex) == ex) {
		t.Error(`ExpandShell("/home/newton/blast")`)
	}
}

func TestNodesFromConfig(t *testing.T) {
	configFiles = []string{"test_data/goodlist"}
	nodes, err := nodesFromConfig()
	if err != nil {
		t.Errorf("Error from test_data/goodlist nodes: %v", err)
	}

	s := nodes[0]
	if s.Host != "node01" {
		t.Errorf("first node in test_data/goodlist is %s instead of node01.", s.Host)
	}

	configFiles = []string{"test_data/good01"}
	nodes, err = nodesFromConfig()
	if err != nil {
		t.Errorf("good list test_data/good01 failed and generated an error")
	}
}

func TestNodesFromConfigWithDash(t *testing.T) {
	configFiles = []string{"test_data/goodlist"}

	var tests = []struct {
		input string
		want  int
	}{
		{"node01", 1},
		{"node01-node04", 4},
		{"node01-node02,node04", 3},
		{"node[01-04]", 4},
	}
	for _, p := range tests {
		nodes, err := GetNodeRange(p.input)
		if err != nil {
			t.Errorf("Error getting node range: %v", err)
		}
		if len(nodes) != p.want {
			t.Errorf("GetNodeRange(%s) should return %d nodes but returned: %d", p.input, p.want, len(nodes))
		}
	}
}

func TestBadConfig(t *testing.T) {
	configFiles = []string{"test_data/doesnotexist"}
	_, err := nodesFromConfig()
	if err == nil {
		t.Errorf("bad list does not generate an error")
	}
}

func TestBadIPaddress(t *testing.T) {
	configFiles = []string{"test_data/bad01"}
	_, err := nodesFromConfig()
	if err == nil {
		t.Errorf("Bad IP address not detected")
	}
}

// Examples
//func ExampleNodeRange
