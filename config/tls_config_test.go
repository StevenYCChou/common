// Copyright 2016 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"encoding/json"

	"gopkg.in/yaml.v2"
)

// LoadTLSConfig parses the given file into a tls.Config.
func LoadTLSConfig(filename string) (*tls.Config, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg := TLSConfig{}
	switch filepath.Ext(filename) {
	case ".yml":
		if err = yaml.UnmarshalStrict(content, &cfg); err != nil {
			return nil, err
		}
	case ".json":
		decoder := json.NewDecoder(bytes.NewReader(content))
		decoder.DisallowUnknownFields()
		if err = decoder.Decode(&cfg); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Unknown extension: %s", filepath.Ext(filename))
	}
	return NewTLSConfig(&cfg)
}

var expectedTLSConfigs = []struct {
	filename string
	config   *tls.Config
}{
	{
		filename: "tls_config.empty.good.json",
		config:   &tls.Config{},
	}, {
		filename: "tls_config.insecure.good.json",
		config:   &tls.Config{InsecureSkipVerify: true},
	}, {
		filename: "tls_config.tlsversion.good.json",
		config:   &tls.Config{MinVersion: tls.VersionTLS11},
	},
	{
		filename: "tls_config.empty.good.yml",
		config:   &tls.Config{},
	}, {
		filename: "tls_config.insecure.good.yml",
		config:   &tls.Config{InsecureSkipVerify: true},
	}, {
		filename: "tls_config.tlsversion.good.yml",
		config:   &tls.Config{MinVersion: tls.VersionTLS11},
	},
}

func TestValidTLSConfig(t *testing.T) {
	for _, cfg := range expectedTLSConfigs {
		got, err := LoadTLSConfig("testdata/" + cfg.filename)
		if err != nil {
			t.Fatalf("Error parsing %s: %s", cfg.filename, err)
		}
		// non-nil functions are never equal.
		got.GetClientCertificate = nil
		if !reflect.DeepEqual(got, cfg.config) {
			t.Fatalf("%v: unexpected config result: \n\n%v\n expected\n\n%v", cfg.filename, got, cfg.config)
		}
	}
}
