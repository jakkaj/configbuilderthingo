package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

type UserConfig struct {
	Name  string `yaml:"name"`
	Tools []Tool `yaml:"tools"`
}

type Tool struct {
	Type   string `yaml:"type"`
	Params []struct {
		Name  string `yaml:"name"`
		Value string `yaml:"value"`
	} `yaml:"params"`
}

type SystemConfig []struct {
	Name       string `yaml:"name"`
	Type       string `yaml:"type"`
	Parameters []struct {
		Name       string `yaml:"name"`
		IsRequired bool   `yaml:"is_required"`
	} `yaml:"parameters"`
	Action struct {
		Type         string `yaml:"type"`
		Path         string `yaml:"path"`
		ArgsTemplate string `yaml:"args_template"`
	} `yaml:"action"`
}

func main() {
	userConfig := parseUserConfig("artefact_config.yaml")
	systemConfig := parseSystemConfig("tool_config.yaml")

	for _, userTool := range userConfig.Tools {
		for _, systemTool := range systemConfig {
			if userTool.Type == systemTool.Type {
				command, err := generateCommand(systemTool.Action.Path, systemTool.Action.ArgsTemplate, userTool.Params)
				if err != nil {
					log.Fatalf("Error generating command: %v\n", err)
				}
				fmt.Println(command)
			}
		}
	}
}

func parseUserConfig(filePath string) UserConfig {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading user-facing config file: %v\n", err)
	}

	var config UserConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling user-facing config: %v\n", err)
	}
	return config
}

func parseSystemConfig(filePath string) SystemConfig {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading system-facing config file: %v\n", err)
	}

	var config SystemConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling system-facing config: %v\n", err)
	}
	return config
}

func generateCommand(path, argsTemplate string, params []struct {
	Name  string
	Value string
}) (string, error) {
	paramMap := make(map[string]string)
	for _, param := range params {
		paramMap[param.Name] = param.Value
	}

	tmpl, err := template.New("args").Parse(argsTemplate)
	if err != nil {
		return "", err
	}

	var argsBuilder strings.Builder
	err = tmpl.Execute(&argsBuilder, paramMap)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s", path, argsBuilder.String()), nil
}
