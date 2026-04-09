// Copyright (C) 2022 Specter Ops, Inc.
//
// This file is part of AzureHound.
//
// AzureHound is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// AzureHound is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Name       string
	Shorthand  string
	Usage      string
	Required   bool
	Persistent bool
	Default    interface{}
	MinValue   int
	MaxValue   int
}

func (s Config) Value() interface{} {
	switch reflect.ValueOf(s.Default).Kind() {
	case reflect.Slice:
		return viper.GetStringSlice(s.Name)
	case reflect.Int:
		return viper.GetInt(s.Name)
	case reflect.Bool:
		return viper.GetBool(s.Name)
	default:
		return viper.Get(s.Name)
	}
}

func (s Config) Set(value interface{}) {
	viper.Set(s.Name, value)
}

type Options struct {
	ConfigFile  string
	ConfigName  string
	ConfigType  string
	ConfigPaths []string
	EnvPrefix   string
}

func Init(cmd *cobra.Command, configs []Config) {
	for _, config := range configs {
		viper.SetDefault(config.Name, config.Default)
		if cmd != nil {
			if config.Persistent {
				setFlag(config, cmd.PersistentFlags(), cmd.MarkPersistentFlagRequired)
			} else {
				setFlag(config, cmd.LocalFlags(), cmd.MarkFlagRequired)
			}
		}
	}
}

func setFlag(config Config, flagSet *pflag.FlagSet, markRequired func(string) error) error {
	switch config.Default.(type) {
	case int:
		flagSet.IntP(config.Name, config.Shorthand, 0, config.Usage)
	case bool:
		flagSet.BoolP(config.Name, config.Shorthand, false, config.Usage)
	case []string:
		flagSet.StringSliceP(config.Name, config.Shorthand, []string{}, config.Usage)
	default:
		flagSet.StringP(config.Name, config.Shorthand, "", config.Usage)
	}

	if config.Required {
		return markRequired(config.Name)
	} else {
		return nil
	}
}

func LoadValues(cmd *cobra.Command, options Options) {
	if cmd != nil {
		viper.BindPFlags(cmd.Flags())
	}
	viper.SetEnvPrefix(options.EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	if options.ConfigFile != "" {
		// If set, ConfigFile gets the highest priority by prepending to ConfigPaths; ConfigPaths are searched
		// in priority order from ConfigPaths[0] (highest priority) to ConfigPaths[len(ConfigPaths)-1] (lowest priority)
		options.ConfigPaths = append([]string{filepath.Dir(options.ConfigFile)}, options.ConfigPaths...)
		basename := filepath.Base(options.ConfigFile)
		ext := filepath.Ext(basename)
		if ext != "" {
			options.ConfigType = ext[1:]
		}
		options.ConfigName = strings.TrimSuffix(basename, ext)
	}

	setConfigSearchPaths(options.ConfigName, options.ConfigType, options.ConfigPaths)

	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError, *os.PathError:
			fmt.Fprintf(os.Stderr, "No configuration file located at %s\n", options.ConfigFile)
		default:
			fmt.Fprintf(os.Stderr, "Unable to read config file: %s\n", err)
		}
	}

	if cmd != nil {
		// Ensure all required values that actually have been set don't return an error. (See https://github.com/spf13/viper/issues/397)
		cmd.Flags().VisitAll(func(flag *pflag.Flag) {
			if viper.IsSet(flag.Name) && viper.Get(flag.Name) != nil {
				switch reflect.ValueOf(viper.Get(flag.Name)).Kind() {
				case reflect.Slice:
					value := strings.Join(viper.GetStringSlice(flag.Name), ",")
					cmd.Flags().Set(flag.Name, value)
				default:
					cmd.Flags().Set(flag.Name, fmt.Sprintf("%v", viper.Get(flag.Name)))
				}
			}
		})
	}
}

func setConfigSearchPaths(name string, extension string, paths []string) {
	viper.SetConfigName(name)
	if extension != "" {
		viper.SetConfigType(extension)
	}
	for _, path := range paths {
		viper.AddConfigPath(path)
	}
}

func ConfigFileUsed() string {
	return viper.ConfigFileUsed()
}
