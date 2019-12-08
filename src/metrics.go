package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	hs1xxInfoRelayState = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hs1xx_info_relay_state",
		Help: "The state of the relay (0=off, 1=on)",
	},
		[]string{
			"alias",
			"ip",
		})
	hs1xxInfoOnTimeSeconds = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hs1xx_info_on_time_seconds",
		Help: "The time in seconds the device has been on",
	},
		[]string{
			"alias",
			"ip",
		})
	hs1xxInfoSignalStrength = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hs1xx_info_signal_strength",
		Help: "The Wifi signal strength in dB",
	},
		[]string{
			"alias",
			"ip",
		})
	hs1xxEnergyVoltageMv = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hs1xx_energy_voltage_mv",
		Help: "The realtime voltage in Millivolts",
	},
		[]string{
			"alias",
			"ip",
		})
	hs1xxEnergyVolt = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hs1xx_energy_volt",
		Help: "The realtime voltage in Volts",
	},
		[]string{
			"alias",
			"ip",
		})
	hs1xxEnergyCurrentMa = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hs1xx_energy_current_ma",
		Help: "The realtime current in Milliamps",
	},
		[]string{
			"alias",
			"ip",
		})
	hs1xxEnergyCurrentAmp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hs1xx_energy_current_amp",
		Help: "The realtime current in Amps",
	},
		[]string{
			"alias",
			"ip",
		})
	hs1xxEnergyPowerMw = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hs1xx_energy_power_mw",
		Help: "The realtime power usage in Milliwatts",
	},
		[]string{
			"alias",
			"ip",
		})
	hs1xxEnergyPowerW = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hs1xx_energy_power_watts",
		Help: "The realtime power usage in Watts",
	},
		[]string{
			"alias",
			"ip",
		})
	hs1xxEnergyTotalWh = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hs1xx_energy_power_total_wh",
		Help: "The total power usage in WattHours",
	},
		[]string{
			"alias",
			"ip",
		})
	hs1xxEnergyTotalW = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hs1xx_energy_power_total_w",
		Help: "The total power usage in Watts",
	},
		[]string{
			"alias",
			"ip",
		})
)

func initRunMetrics(port int) {

	prometheus.MustRegister(hs1xxInfoRelayState)
	prometheus.MustRegister(hs1xxInfoOnTimeSeconds)
	prometheus.MustRegister(hs1xxInfoSignalStrength)

	prometheus.MustRegister(hs1xxEnergyVoltageMv)
	prometheus.MustRegister(hs1xxEnergyVolt)
	prometheus.MustRegister(hs1xxEnergyCurrentMa)
	prometheus.MustRegister(hs1xxEnergyCurrentAmp)
	prometheus.MustRegister(hs1xxEnergyPowerMw)
	prometheus.MustRegister(hs1xxEnergyPowerW)
	prometheus.MustRegister(hs1xxEnergyTotalWh)
	prometheus.MustRegister(hs1xxEnergyTotalW)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
             <head><title>hs1xx Exporter</title></head>
             <body>
             <h1>hs1xx Exporter</h1>
             <p><a href='/metrics'>Metrics</a></p>
             </body>
             </html>`))
	})

	http.Handle("/metrics", promhttp.Handler())
	log.Error(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil))
}
