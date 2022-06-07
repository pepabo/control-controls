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
		if std.Controls == nil {
			std.Controls = &Controls{}
		}
		std.Key = k
		s.Standards = append(s.Standards, std)
	}
	for r, tmphub := range tmp.Regions {
		hub := New(r)
		hub.AutoEnable = tmphub.AutoEnable
		for k, std := range tmphub.Standards {
			if std.Controls == nil {
				std.Controls = &Controls{}
			}
			std.Key = k
			hub.Standards = append(hub.Standards, std)
		}
		s.Regions = append(s.Regions, hub)
	}
	return nil
}

func (c *Controls) UnmarshalYAML(b []byte) error {
	s := struct {
		Enable  []string      `yaml:"enable,flow,omitempty"`
		Disable yaml.MapSlice `yaml:"disable,omitempty"`
	}{}
	if err := yaml.Unmarshal(b, &s); err == nil {
		c.Enable = s.Enable
		c.Disable = s.Disable
		return nil
	}

	// fallback as slice
	s2 := struct {
		Enable  []string `yaml:"enable,flow,omitempty"`
		Disable []string `yaml:"disable,flow,omitempty"`
	}{}
	if err := yaml.Unmarshal(b, &s2); err != nil {
		return err
	}

	c.Enable = s2.Enable
	for _, d := range s2.Disable {
		c.Disable = append(c.Disable, yaml.MapItem{
			Key:   d,
			Value: "",
		})
	}
	return nil
}
