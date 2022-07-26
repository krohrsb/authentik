package web

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"goauthentik.io/internal/config"
	"goauthentik.io/internal/utils/sentry"
)

var (
	Requests = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "authentik_main_requests",
		Help: "The total number of configured providers",
	}, []string{"dest"})
)

func RunMetricsServer() {
	m := mux.NewRouter()
	l := log.WithField("logger", "authentik.router.metrics")
	m.Use(sentry.SentryNoSampleMiddleware)
	m.Path("/metrics").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		promhttp.InstrumentMetricHandler(
			prometheus.DefaultRegisterer, promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{
				DisableCompression: true,
			}),
		).ServeHTTP(rw, r)

		// Get upstream metrics
		re, err := http.NewRequest("GET", "http://localhost:8000/-/metrics/", nil)
		if err != nil {
			l.WithError(err).Warning("failed to get upstream metrics")
			return
		}
		re.SetBasicAuth("monitor", config.Get().SecretKey)
		res, err := http.DefaultClient.Do(re)
		if err != nil {
			l.WithError(err).Warning("failed to get upstream metrics")
			return
		}
		bm, err := ioutil.ReadAll(res.Body)
		if err != nil {
			l.WithError(err).Warning("failed to get upstream metrics")
			return
		}
		_, err = rw.Write(bm)
		if err != nil {
			l.WithError(err).Warning("failed to get upstream metrics")
			return
		}
	})
	l.WithField("listen", config.Get().Web.ListenMetrics).Info("Starting Metrics server")
	err := http.ListenAndServe(config.Get().Web.ListenMetrics, m)
	if err != nil {
		l.WithError(err).Warning("Failed to start metrics server")
	}
	l.WithField("listen", config.Get().Web.ListenMetrics).Info("Stopping Metrics server")
}
