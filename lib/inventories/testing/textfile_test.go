package testing

import (
	"testing"

	"github.com/jtopjian/bagel/lib/inventories"

	"github.com/stretchr/testify/assert"
)

func TestTextFile(t *testing.T) {
	options := map[string]interface{}{
		"file": "fixtures/hosts.txt",
	}

	textfile, err := inventories.New("textfile", options)
	if err != nil {
		t.Fatal(err)
	}

	expected := []inventories.Target{
		inventories.Target{Address: "host1.example.com", Name: "host1.example.com"},
		inventories.Target{Address: "host2.example.com", Name: "host2.example.com"},
		inventories.Target{Address: "192.168.100.1", Name: "192.168.100.1"},
		inventories.Target{Address: "fe80::f816:3eff:fe8c:c73a", Name: "fe80::f816:3eff:fe8c:c73a"},
		inventories.Target{Address: "[fe80::f816:3eff:fe8c:c73a]", Name: "[fe80::f816:3eff:fe8c:c73a]"},
	}

	actual, err := textfile.Discover()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, actual)
}
