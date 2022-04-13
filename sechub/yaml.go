package sechub

import (
	"github.com/goccy/go-yaml"
)

func (s *SecHub) MarshalYAML() ([]byte, error) {
	stds := yaml.MapSlice{}
	for _, std := range s.Standards {
		k := std.Key
		v := &Standard{
			Enable:   std.Enable,
			Controls: std.Controls,
		}
		stds = append(stds, yaml.MapItem{
			Key:   k,
			Value: v,
		})
	}

	// regions := yaml.MapSlice{}
	// for _, hub := range s.Regions {
	// 	k := hub.Region
	// 	v := hub
	// 	regions = append(regions, yaml.MapItem{
	// 		Key:   k,
	// 		Value: v,
	// 	})
	// }

	regions := map[string]*SecHub{}
	for _, hub := range s.Regions {
		k := hub.Region
		regions[k] = hub
	}

	s2 := struct {
		AutoEnable *bool              `yaml:"autoEnable,omitempty"`
		Standards  yaml.MapSlice      `yaml:"standards,omitempty"`
		Regions    map[string]*SecHub `yaml:"regions,omitempty"`
	}{
		AutoEnable: s.AutoEnable,
		Standards:  stds,
		Regions:    regions,
	}
	return yaml.Marshal(s2)
}
