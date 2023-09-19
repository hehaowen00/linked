package main

import (
	"bufio"
	"os"
	"strings"
)

type Config struct {
	DBPath         string `env:"DATABASE_PATH"`
	MigrationsPath string `env:"MIGRATIONS_PATH"`
}

func LoadEnvConfig(config map[string]string) error {
	f, err := os.Open(".env")
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			config[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
