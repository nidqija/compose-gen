package tui

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ServiceItem represents a selectable service template with parsed metadata
type ServiceItem struct {
	Key         string
	Label       string
	Services    map[string]interface{}
	Volumes     map[string]interface{}
	Image       string
	Ports       []string
	ExtraFiles  []string
	EnvVars     []string
	Selected    bool
}

// GenerationLog represents the result of generating compose and scaffolding files
type GenerationLog struct {
	ServicesGenerated []string
	FilesScaffolded   []string
	ExistingPreserved []string
	Error             error
	Mode              string // "merged" or "overwritten" or "created"
}

// LoadServices parses all templates from the embedded filesystem
func LoadServices(fs embed.FS) ([]ServiceItem, error) {
	entries, err := fs.ReadDir("templates")
	if err != nil {
		return nil, err
	}

	var items []ServiceItem

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yml") {
			continue
		}

		data, err := fs.ReadFile("templates/" + entry.Name())
		if err != nil {
			continue
		}

		var raw struct {
			Name     string                 `yaml:"name"`
			Services map[string]interface{} `yaml:"services"`
			Volumes  map[string]interface{} `yaml:"volumes"`
		}

		if err := yaml.Unmarshal(data, &raw); err != nil {
			continue
		}

		key := strings.TrimSuffix(entry.Name(), ".yml")
		label := raw.Name
		if label == "" {
			label = key
		}

		// Inspect image, ports, env from raw.Services
		var image string
		var ports []string
		var envVars []string

		for _, v := range raw.Services {
			if vMap, ok := v.(map[string]interface{}); ok {
				if img, ok := vMap["image"].(string); ok {
					image = img
				}
				if pList, ok := vMap["ports"].([]interface{}); ok {
					for _, p := range pList {
						ports = append(ports, fmt.Sprintf("%v", p))
					}
				}
				if envList, ok := vMap["environment"].([]interface{}); ok {
					for _, e := range envList {
						envVars = append(envVars, fmt.Sprintf("%v", e))
					}
				} else if envMap, ok := vMap["environment"].(map[string]interface{}); ok {
					for k := range envMap {
						envVars = append(envVars, k)
					}
				}
			}
		}

		// Map known extra files scaffolded
		var extraFiles []string
		switch key {
		case "grafana-k6":
			extraFiles = append(extraFiles, "tests/k6/load-test.js")
		case "nginx":
			extraFiles = append(extraFiles, "configs/nginx.conf", "html/")
		case "otel":
			extraFiles = append(extraFiles, "configs/otel-collector-config.yaml")
		case "cloudflared":
			extraFiles = append(extraFiles, "configs/.env")
		}

		items = append(items, ServiceItem{
			Key:        key,
			Label:      label,
			Services:   raw.Services,
			Volumes:    raw.Volumes,
			Image:      image,
			Ports:      ports,
			ExtraFiles: extraFiles,
			EnvVars:    envVars,
			Selected:   false,
		})
	}

	return items, nil
}

// GenerateCompose writes the docker-compose.yml and scaffolds extra files
func GenerateCompose(fs embed.FS, selected []ServiceItem, mergeWithExisting bool) GenerationLog {
	log := GenerationLog{
		Mode: "created",
	}

	targetFile := "docker-compose.yml"
	composeMap := map[string]interface{}{
		"services": make(map[string]interface{}),
	}
	volumeMap := make(map[string]interface{})

	// Check if existing file should be merged
	if mergeWithExisting {
		if existingData, err := os.ReadFile(targetFile); err == nil {
			var existingConfig map[string]interface{}
			if err := yaml.Unmarshal(existingData, &existingConfig); err == nil {
				if existingServices, ok := existingConfig["services"].(map[string]interface{}); ok {
					for k, v := range existingServices {
						composeMap["services"].(map[string]interface{})[k] = v
					}
				}
				if existingVolumes, ok := existingConfig["volumes"].(map[string]interface{}); ok {
					for k, v := range existingVolumes {
						volumeMap[k] = v
					}
				}
				log.Mode = "merged"
			}
		}
	} else if _, err := os.Stat(targetFile); err == nil {
		log.Mode = "overwritten"
	}

	// Add selected services
	for _, s := range selected {
		composeMap["services"].(map[string]interface{})[s.Key] = s.Services
		log.ServicesGenerated = append(log.ServicesGenerated, s.Key)

		if s.Volumes != nil {
			for volKey, volVal := range s.Volumes {
				volumeMap[volKey] = volVal
			}
		}
	}

	if len(volumeMap) > 0 {
		composeMap["volumes"] = volumeMap
	}

	// Marshal YAML
	outData, err := yaml.Marshal(composeMap)
	if err != nil {
		log.Error = fmt.Errorf("failed to encode docker-compose.yml: %w", err)
		return log
	}

	if err := os.WriteFile(targetFile, outData, 0644); err != nil {
		log.Error = fmt.Errorf("failed to write docker-compose.yml: %w", err)
		return log
	}

	// Scaffold extra files
	for _, s := range selected {
		switch s.Key {
		case "grafana-k6":
			targetScript := "tests/k6/load-test.js"
			if _, err := os.Stat(targetScript); err == nil {
				log.ExistingPreserved = append(log.ExistingPreserved, targetScript)
			} else {
				data, err := fs.ReadFile("templates/configs/load-test.js")
				if err == nil {
					_ = os.MkdirAll(filepath.Dir(targetScript), 0755)
					if err := os.WriteFile(targetScript, data, 0644); err == nil {
						log.FilesScaffolded = append(log.FilesScaffolded, targetScript)
					}
				}
			}

		case "nginx":
			targetConf := "configs/nginx.conf"
			if _, err := os.Stat(targetConf); err == nil {
				log.ExistingPreserved = append(log.ExistingPreserved, targetConf)
			} else {
				data, err := fs.ReadFile("templates/configs/nginx.conf")
				if err == nil {
					_ = os.MkdirAll(filepath.Dir(targetConf), 0755)
					if err := os.WriteFile(targetConf, data, 0644); err == nil {
						log.FilesScaffolded = append(log.FilesScaffolded, targetConf)
					}
				}
			}
			_ = os.MkdirAll("html", 0755)

		case "otel":
			targetConf := "configs/otel-collector-config.yaml"
			if _, err := os.Stat(targetConf); err == nil {
				log.ExistingPreserved = append(log.ExistingPreserved, targetConf)
			} else {
				data, err := fs.ReadFile("templates/configs/otel-collector-config.yaml")
				if err == nil {
					_ = os.MkdirAll(filepath.Dir(targetConf), 0755)
					if err := os.WriteFile(targetConf, data, 0644); err == nil {
						log.FilesScaffolded = append(log.FilesScaffolded, targetConf)
					}
				}
			}

		case "cloudflared":
			targetEnv := "configs/.env"
			if _, err := os.Stat(targetEnv); err == nil {
				log.ExistingPreserved = append(log.ExistingPreserved, targetEnv)
			} else {
				defaultEnv := "# Base URL for Mobile & Web Services\nAPP_BASE_URL=http://localhost:8025\n"
				_ = os.MkdirAll(filepath.Dir(targetEnv), 0755)
				if err := os.WriteFile(targetEnv, []byte(defaultEnv), 0644); err == nil {
					log.FilesScaffolded = append(log.FilesScaffolded, targetEnv)
				}
			}
		}
	}

	return log
}

