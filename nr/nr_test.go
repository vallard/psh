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

/*
func TestNodeRange(t *testing.T) {
	var tests = []struct {
		input string
		want  []server.Server
	}{
		{"node01-node02",
			[]server.Server{
				{"node01", "10.93.234.10", "bob", "~/.ssh/was"},
				{"node02", "10.93.234.11", "bob", "~/.ssh/was"},
			},
		},
	}

	for _, test := range tests {
		got, err := GetNodeRange(test);
		if got != test.want {
			t.Errorf("GetNodeRange(%q) = %v", test.input, got)
		}
	}
}
*/

// Examples
//func ExampleNodeRange
