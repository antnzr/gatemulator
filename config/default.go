package config

import (
	"bufio"
	"os"
	"path"
	"strconv"
	"strings"
)

type Config struct {
	Env  string
	Port int
}

func LoadConfig(configFilePath string) (*Config, error) {
	envFromFile := make(map[string]string)

	if configFilePath != "" {
		envFromFileLoaded, err := LoadEnvFromFile(configFilePath)
		if err != nil {
			return nil, err
		}
		envFromFile = envFromFileLoaded
	}

	env, err := getEnvOrDefault("ENV", "development", envFromFile)
	if err != nil {
		return nil, err
	}

	port, err := getEnvAsIntOrDefault("PORT", 34000, envFromFile)
	if err != nil {
		return nil, err
	}

	config := &Config{
		Env:  env,
		Port: port,
	}

	return config, err
}

func getEnvOrDefault(key, fallback string, envFromFile map[string]string) (string, error) {
	if value, ok := os.LookupEnv(key); ok {
		return value, nil
	}
	if envFromFile[key] != "" {
		return envFromFile[key], nil
	}
	return fallback, nil
}

func getEnvAsIntOrDefault(key string, fallback int, envFromFile map[string]string) (int, error) {
	if value, ok := os.LookupEnv(key); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue, nil
		}
	}
	if value, ok := envFromFile[key]; ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue, nil
		}
	}
	return fallback, nil
}

func LoadEnvFromFile(filename string) (map[string]string, error) {
	envVars := make(map[string]string)

	root, _ := os.Getwd()
	filePath := path.Join(root, filename)

	if _, err := os.Stat(filePath); err != nil {
		return nil, nil
	}

	file, err := os.Open(path.Join(root, filename))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "=") && !strings.HasPrefix(line, "#") {
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			envVars[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return envVars, nil
}
