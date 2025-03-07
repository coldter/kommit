package settings

import (
	"errors"
	"fmt"
	"github.com/kirsle/configdir"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
	"sync"
)

type Settings struct {
	changed bool
}

var (
	settings *Settings
	mu       sync.Mutex
)

func ReadSettings() (*Settings, error) {
	mu.Lock()
	defer mu.Unlock()

	if settings != nil {
		return settings, nil
	}

	settings = &Settings{}
	configPath := configdir.LocalConfig("kommit")
	_ = viper.BindEnv("baseURL", "KOMMIT_API_BASEURL")

	err := configdir.MakePath(configPath)
	if err != nil {
		return nil, err
	}

	viper.SetConfigName("settings")
	viper.SetConfigType("json")
	viper.AddConfigPath(configPath)
	configFile := path.Join(configPath, "settings.json")
	if abs, err := filepath.Abs(configFile); err == nil {
		configFile = abs
	}

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		var configParseError viper.ConfigParseError
		switch {
		case errors.As(err, &configFileNotFoundError):
			// Force config creation
			if err := viper.SafeWriteConfig(); err != nil {
				return nil, err
			}
		case errors.As(err, &configParseError):
			fmt.Println("Warning")
			fmt.Printf("could not parse JSON config from file %s\n", configFile)
			return nil, err
		default:
			return nil, err
		}
	}

	return settings, nil
}

func (s *Settings) SetToken(token string) {
	viper.Set("token", token)
	s.changed = true
}

func (s *Settings) SetBaseUrl(baseUrl string) {
	viper.Set("baseUrl", baseUrl)
	s.changed = true
}

func (s *Settings) GetToken() string {
	return viper.GetString("token")
}

func PersistChanges() {
	if settings == nil || !settings.changed {
		return
	}

	if err := TryToPersistChanges(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	}
}

func TryToPersistChanges() error {
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to persist turso settings file: %w", err)
	}
	return nil
}

func (s *Settings) GetBaseURL() string {
	return viper.GetString("baseURL")
}

func (s *Settings) GetConfigPath() string {
	return path.Join(configdir.LocalConfig("kommit"), "settings.json")
}

func (s *Settings) GetAlLSettings() map[string]interface{} {
	return viper.AllSettings()
}
