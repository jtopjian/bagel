package testing

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jtopjian/bagel/lib/site"
)

var nilMap = make(map[string]interface{})

var expectedSite = &site.Site{
	Roles: map[string]site.Role{
		"memcached": site.Role{
			Inventories: []string{
				"static1", "static2",
			},
		},
		"mysql": site.Role{
			Inventories: []string{
				"mysql_nodes",
			},
		},
	},

	Inventories: map[string]site.Inventory{
		"static1": site.Inventory{
			Type:       "textfile",
			Connection: "ssh",
			Options: map[string]interface{}{
				"file": "/my/file.txt",
			},
		},
		"static2": site.Inventory{
			Type:       "textfile",
			Connection: "ssh",
			Options: map[string]interface{}{
				"file": "/my/other/file.txt",
			},
		},
		"mysql_nodes": site.Inventory{
			Type:       "textfile",
			Connection: "ssh",
			Options: map[string]interface{}{
				"file": "/my/mysql/nodes.txt",
			},
		},
	},

	Connections: map[string]site.Connection{
		"ssh": site.Connection{
			Type: "ssh",
			Options: map[string]interface{}{
				"private_key": "/path/to/id_rsa",
				"port":        22,
			},
		},
	},
}

func TestSite(t *testing.T) {
	actualSite, err := site.New("fixtures/site.yaml")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedSite, actualSite)
}
