package settings

import (
	"encoding/json"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type Settings struct {
	URL string `json:"url"`
}

func Load(settings backend.DataSourceInstanceSettings) *Settings {
	s := &Settings{}
	if settings.JSONData != nil && len(settings.JSONData) > 1 {
		_ = json.Unmarshal(settings.JSONData, s)
	}
	if s.URL == "" {
		s.URL = "localhost:10000"
	}
	return s
}
