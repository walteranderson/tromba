package config

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type ProjectConfig struct {
	Name   string
	Input  string
	Output string
}

type Config struct {
	ProjectPath   string
	ProjectConfig ProjectConfig
}

func Load() *Config {
	path := getProjectPath()
	projectConfig := getProjectConfig(path)
	return &Config{
		ProjectPath:   path,
		ProjectConfig: projectConfig,
	}
}

func getProjectPath() string {
	var projectPath string
	if len(os.Args) > 1 {
		path, err := filepath.Abs(os.Args[1])
		if err != nil {
			log.Fatal("ERROR: could not get project path", err)
		}
		projectPath = path
	}

	if projectPath == "" {
		defaultPath, err := os.Getwd()
		if err != nil {
			log.Fatal("ERROR: could not get default project path", err)
		}
		projectPath = defaultPath
	}

	validateProjectPath(projectPath)
	return projectPath
}

func validateProjectPath(path string) {
	file, err := os.Open(filepath.Join(path, "tromba.toml"))
	if err != nil {
		log.Fatal("ERROR: could not find tromba.toml in project directory")
	}
	defer file.Close()
}

func getProjectConfig(path string) ProjectConfig {
	conf := ProjectConfig{
		Input:  "src",
		Output: "dist",
	}

	file, err := os.Open(filepath.Join(path, "tromba.toml"))
	if err != nil {
		log.Fatal("ERROR: could not find tromba.toml in project directory")
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("ERROR: could not read tromba.toml", err)
	}

	err = toml.Unmarshal(data, &conf)
	if err != nil {
		log.Fatal("ERROR: could not unmarshal tromba.toml", err)
	}

	return conf
}
