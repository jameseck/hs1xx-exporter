package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/mattn/go-colorable"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ripienaar/hs1xxplug"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"time"
)

var debug bool

type Config struct {
	MetricsPort     int
	MetricsInterval string
	PlugIPs         []string
}

func initConfig(homedir string, configPath string) (config *Config, err error) {
	config = &Config{}

	viper.SetDefault("MetricsInterval", "15s")
	viper.SetDefault("MetricsPort", 9115)
	viper.SetDefault("PlugIPs", []string{})

	viper.SetConfigName("hs1xx-exporter")
	viper.AddConfigPath(configPath)
	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Infof("initConfig: Config file not found - using defaults")
		} else {
			// Config file was found but another error was produced
			return &Config{}, errors.Wrap(err, "initConfig: Error reading the config file")
		}
	}

	if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		log.Infof("Initializing config file watcher for %s", viper.ConfigFileUsed())
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			if e.Op == fsnotify.Write {
				log.Infof("Config file %s changed", e.Name)
				viper.Unmarshal(&config)
				if err != nil {
					log.Errorf("viperOnConfigChange: Error unmarshaling config: %+v", err)
				}
			}
		})
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return &Config{}, errors.Wrap(err, "initConfig: Error unmarshaling config")
	}

	return config, err
}

func initLogger(logLevel log.Level) {

	log.SetLevel(logLevel)
	log.SetOutput(colorable.NewColorableStdout())
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})
}

func main() {

	configPath := flag.StringP("config-path", "c", ".", "The path to search for config files (hs1xx.{yaml,toml,json})")
	flag.BoolVarP(&debug, "debug", "d", false, "Debug logging")
	flag.Parse()

	var logLevel log.Level
	if debug {
		logLevel = log.DebugLevel
	} else {
		logLevel = log.InfoLevel
	}
	initLogger(logLevel)

	homedir := os.Getenv("HOME")

	config, err := initConfig(homedir, *configPath)
	if err != nil {
		log.Errorf("%v", err)
		os.Exit(1)
	}

	go initRunMetrics(config.MetricsPort)

	for {
		for i := range config.PlugIPs {
			log.Debugf("Connecting to %s", config.PlugIPs[i])
			plug := hs1xxplug.NewPlug(config.PlugIPs[i])

			log.Debugf("Calling Energy func on %s", config.PlugIPs[i])
			energyData, err := plug.Energy()
			if err != nil {
				log.Error(err)
			}

			log.Debugf("Calling Info func on %s", config.PlugIPs[i])
			infoData, err := plug.Info()
			if err != nil {
				log.Error(err)
			}

			hs1xxInfoRelayState.With(prometheus.Labels{"ip": infoData.Address, "alias": infoData.Alias}).Set(float64(infoData.RelayState))
			hs1xxInfoOnTimeSeconds.With(prometheus.Labels{"ip": infoData.Address, "alias": infoData.Alias}).Set(float64(infoData.OnTimeSeconds))
			hs1xxInfoSignalStrength.With(prometheus.Labels{"ip": infoData.Address, "alias": infoData.Alias}).Set(float64(infoData.SignalStrength))

			hs1xxEnergyVoltageMv.With(prometheus.Labels{"ip": infoData.Address, "alias": infoData.Alias}).Set(float64(energyData.MilliVolt))
			hs1xxEnergyVolt.With(prometheus.Labels{"ip": infoData.Address, "alias": infoData.Alias}).Set(float64(energyData.Volt))
			hs1xxEnergyCurrentMa.With(prometheus.Labels{"ip": infoData.Address, "alias": infoData.Alias}).Set(float64(energyData.MilliAmp))
			hs1xxEnergyCurrentAmp.With(prometheus.Labels{"ip": infoData.Address, "alias": infoData.Alias}).Set(float64(energyData.Amp))
			hs1xxEnergyPowerMw.With(prometheus.Labels{"ip": infoData.Address, "alias": infoData.Alias}).Set(float64(energyData.PowerUseMilliWatt))
			hs1xxEnergyPowerW.With(prometheus.Labels{"ip": infoData.Address, "alias": infoData.Alias}).Set(float64(energyData.PowerUseWatt))
			hs1xxEnergyTotalWh.With(prometheus.Labels{"ip": infoData.Address, "alias": infoData.Alias}).Set(float64(energyData.TotalMilliWatt))
			hs1xxEnergyTotalW.With(prometheus.Labels{"ip": infoData.Address, "alias": infoData.Alias}).Set(float64(energyData.TotalWatt))
		}
		interval, _ := time.ParseDuration(config.MetricsInterval)
		log.Debugf("Sleeping %v", interval)
		time.Sleep(interval)
	}
}
