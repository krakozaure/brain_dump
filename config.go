package main

import (
	"fmt"
	"log"
	"os"

	"github.com/imdario/mergo"
)

type Formats map[string]string

type Profile struct {
	File         string            `json:"file"`
	FileMode     string            `json:"file_mode"`
	TemplateFile string            `json:"template_file"`
	TemplateVars map[string]string `json:"template_vars"`
	Formats      Formats           `json:"formats"`
	KeysCase     string            `json:"keys_case"`
	Editor       string            `json:"editor"`
	EditorArgs   []string          `json:"editor_args"`
}

type Config map[string]Profile

func (c Config) String() string {
	return marshal(c)
}

func (c Config) hasProfile(name string) bool {
	if name == "" {
		return false
	}
	for currentName := range c {
		if currentName == name {
			return true
		}
	}
	return false
}

func (c Config) getProfile(name string) Profile {
	var err error

	profile := getDefaultProfile()

	if c.hasProfile(name) {
		err = mergo.Merge(&profile, c[name], mergo.WithOverride)
	} else if c.hasProfile("default") {
		err = mergo.Merge(&profile, c["default"], mergo.WithOverride)
	}
	if err != nil {
		log.Println("Unable to merge default and user profile. Default profile is used.")
		profile = getDefaultProfile()
	}

	return profile
}

func getDefaultConfig() Config {
	return Config{"default": getDefaultProfile()}
}

func getDefaultProfile() Profile {
	return Profile{
		File:     "~/Documents/brain_dump.md",
		FileMode: "append",
		Formats: Formats{
			"date":     "2006-01-02",
			"time":     "15:04:05",
			"datetime": "2006-01-02 15:04:05",
		},
		KeysCase: "snake_case",
		Editor:   "$EDITOR",
	}
}

func getUserConfig() (Config, error) {
	var configFile string

	cwdConfigFile := fmt.Sprintf("%s.json", APP_NAME)
	envConfigFile := os.Getenv("BRAINDUMP_CONFIG_FILE")

	if fileExists(cwdConfigFile) {
		configFile = cwdConfigFile
	} else if envConfigFile != "" {
		configFile = envConfigFile
	} else {
		configFile = APP_CONFIG_FILE
	}

	if fileExists(configFile) {
		return loadConfigFile(configFile)
	}
	return Config{}, fmt.Errorf("Unable to find the configuration file")
}

func loadConfigFile(path string) (Config, error) {
	config := Config{}

	path = expandText(path)
	data, err := readFile(path)
	if err != nil {
		return config, err
	}

	err = unmarshalString(data, &config)
	if err != nil {
		return config, err
	}

	return config, err
}
