package settings

import (
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type SheeshSettings struct {
	Token    string `yaml:"token"`
	Status   string `yaml:"status"`
	GameName string `yaml:"game-name"`
}

func InitSettings(logger *zap.Logger, f string) *SheeshSettings {
	var sheeshSettings SheeshSettings

	fileSettings, err := os.ReadFile(f)
	if err != nil {
		logger.Fatal("Error loading settings.yaml file")
	}

	err = yaml.Unmarshal(fileSettings, &sheeshSettings)
	if err != nil {
		logger.Fatal("YAML failed to unmarshal to SheeshSettings", zap.Error(err))
	}

	return &sheeshSettings
}
