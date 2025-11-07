package config

import (
	"sync"
	"time"

	"github.com/spf13/viper"
)

type Loader interface {
	Get(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetDuration(key string) time.Duration
}

type reader struct {
	*viper.Viper
}

var once sync.Once
var onceReader *reader

func Reader() Loader {
	once.Do(func() {
		viper.AutomaticEnv()
		onceReader = &reader{viper.GetViper()}
	})
	return onceReader
}
func (r *reader) Get(key string) string {
	return r.Viper.GetString(key)
}

func (r *reader) GetBool(key string) bool {
	return r.Viper.GetBool(key)
}

func (r *reader) GetInt(key string) int {
	return r.Viper.GetInt(key)
}

func (r *reader) GetDuration(key string) time.Duration {
	return r.Viper.GetDuration(key)
}
