package site

import (
	"fmt"
	"sync"

	"github.com/jtopjian/bagel/lib/inventories"
	"github.com/jtopjian/bagel/lib/utils"
)

// Inventory represents an inventory of targets.
type Inventory struct {
	Auth       string                 `yaml:"auth"`
	Type       string                 `yaml:"type" required:"true"`
	Options    map[string]interface{} `yaml:"options"`
	Connection string                 `yaml:"connection" required:"true"`

	Targets []inventories.Target `yaml:"-"`
	mux     sync.Mutex
}

// UnmarshalYAML is a custom unmarshaler to help initialize and
// validate an inventory.
func (r *Inventory) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type tmp Inventory
	var s struct {
		tmp `yaml:",inline"`
	}

	err := unmarshal(&s)
	if err != nil {
		return fmt.Errorf("unable to parse YAML: %s", err)
	}

	*r = Inventory(s.tmp)

	if err := utils.ValidateTags(r); err != nil {
		return err
	}

	// If options weren't specified, create an empty map.
	if r.Options == nil {
		r.Options = make(map[string]interface{})
	}

	return nil
}

// DiscoverTargets will run Discover and return the targets.
func (r *Inventory) DiscoverTargets() error {
	r.mux.Lock()
	defer r.mux.Unlock()

	t, err := inventories.New(r.Type, r.Options)
	if err != nil {
		return err
	}

	discoveredTargets, err := t.Discover()
	if err != nil {
		return err
	}

	for _, target := range discoveredTargets {
		r.Targets = append(r.Targets, inventories.Target{
			Name:    target.Name,
			Address: target.Address,
		})
	}

	return nil
}
