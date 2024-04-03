package util

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
)

func Load(jsonStr string, makeInstance func(typ string) any) (any, error) {
	// json to map
	m := make(map[string]any)
	e := json.Unmarshal([]byte(jsonStr), &m)
	if e != nil {
		return nil, e
	}

	data := makeInstance(m["type"].(string))

	// decoder to copy map values to my struct using json tags
	cfg := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &data,
		TagName:  "json",
		Squash:   true,
	}
	decoder, e := mapstructure.NewDecoder(cfg)
	if e != nil {
		return nil, e
	}
	// copy map to struct
	e = decoder.Decode(m)
	return data, e
}
