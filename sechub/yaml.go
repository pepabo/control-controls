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
	// 	k := hub.region
	// 	v := hub
	// 	regions = append(regions, yaml.MapItem{
	// 		Key:   k,
	// 		Value: v,
	// 	})
	// }

	regions := map[string]*SecHub{}
	for _, hub := range s.Regions {
		k := hub.region
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

type SecHubForUnmarshal struct {
	AutoEnable *bool                          `yaml:"autoEnable,omitempty"`
	Standards  map[string]*Standard           `yaml:"standards,omitempty"`
	Regions    map[string]*SecHubForUnmarshal `yaml:"regions,omitempty"`
}

func (s *SecHub) UnmarshalYAML(b []byte) error {
	tmp := &SecHubForUnmarshal{}
	if err := yaml.Unmarshal(b, tmp); err != nil {
		return err
	}
	s.AutoEnable = tmp.AutoEnable
	for k, std := range tmp.Standards {
		std.Key = k
		s.Standards = append(s.Standards, std)
	}
	for r, tmphub := range tmp.Regions {
		hub := New(r)
		hub.AutoEnable = tmphub.AutoEnable
		for k, std := range tmphub.Standards {
			std.Key = k
			hub.Standards = append(hub.Standards, std)
		}
		s.Regions = append(s.Regions, hub)
	}
	return nil
}
