package metrics

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"strings"

	"github.com/drand/drand/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// PrivateMetrics about the internal world (go process, private stuff)
	PrivateMetrics = prometheus.NewRegistry()
	// HTTPMetrics about the public surface area (http requests, cdn stuff)
	HTTPMetrics = prometheus.NewRegistry()
	// GroupMetrics about the group surface (grp, group-member stuff)
	GroupMetrics = prometheus.NewRegistry()
	// ClientMetrics about the drand client requests to servers
	ClientMetrics = prometheus.NewRegistry()

	// APICallCounter (Group) how many grpc calls
	APICallCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_call_counter",
		Help: "Number of API calls that we have received",
	}, []string{"api_method"})
	// GroupDialFailures (Group) how manuy failures connecting outbound
	GroupDialFailures = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "dial_failures",
		Help: "Number of times there have been network connection issues",
	}, []string{"peer_address"})
	// GroupConnections (Group) how many GrpcClient connections are present
	GroupConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "group_connections",
		Help: "Number of peers with current GrpcClient connections",
	})

	// HTTPCallCounter (HTTP) how many http requests
	HTTPCallCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_call_counter",
		Help: "Number of HTTP calls received",
	}, []string{"code", "method"})
	// HTTPLatency (HTTP) how long http request handling takes
	HTTPLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "http_response_duration",
		Help:        "histogram of request latencies",
		Buckets:     prometheus.DefBuckets,
		ConstLabels: prometheus.Labels{"handler": "http"},
	}, []string{"method"})
	// HTTPInFlight (HTTP) how many http requests exist
	HTTPInFlight = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "http_in_flight",
		Help: "A gauge of requests currently being served.",
	})

	// Client observation metrics

	// ClientWatchLatency measures the latency of the watch channel from the client's perspective.
	ClientWatchLatency = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "client_watch_latency",
		Help: "Duration between time round received and time round expected.",
	})

	// ClientHTTPHeartbeatLatency measures the latency of HTTP hearbeat andomness requests.
	ClientHTTPHeartbeatLatency = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "client_http_heartbeat_latency",
		Help: "Duration between time round received and time round expected, measured on heartbeats.",
	})

	// ClientHTTPHeartbeatFailure measures the number of times HTTP heartbeats fail.
	ClientHTTPHeartbeatFailure = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "client_http_heartbeat_failure",
		Help: "The number of times HTTP heartbeat requests fail.",
	})

	// ClientInFlight measures how many active requests have been made
	ClientInFlight = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "client_in_flight",
		Help: "A gauge of in-flight drand client http requests.",
	},
		[]string{"url"},
	)

	// ClientRequests measures how many total requests have been made
	ClientRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "client_api_requests_total",
			Help: "A counter for requests from the drand client.",
		},
		[]string{"code", "method", "url"},
	)

	// ClientDNSLatencyVec tracks the observed DNS resolution times
	ClientDNSLatencyVec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "client_dns_duration_seconds",
			Help:    "Client drand dns latency histogram.",
			Buckets: []float64{.005, .01, .025, .05},
		},
		[]string{"event", "url"},
	)

	// ClientTLSLatencyVec tracks observed TLS connection times
	ClientTLSLatencyVec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "client_tls_duration_seconds",
			Help:    "Client drand tls latency histogram.",
			Buckets: []float64{.05, .1, .25, .5},
		},
		[]string{"event", "url"},
	)

	// ClientLatencyVec tracks raw http request latencies
	ClientLatencyVec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "client_request_duration_seconds",
			Help:    "A histogram of client request latencies.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"url"},
	)

	metricsBound = false
)

func bindMetrics() {
	if metricsBound {
		return
	}
	metricsBound = true

	// The private go-level metrics live in private.
	PrivateMetrics.Register(prometheus.NewGoCollector())
	PrivateMetrics.Register(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))

	// Group metrics
	group := []prometheus.Collector{
		APICallCounter,
		GroupDialFailures,
		GroupConnections,
	}
	for _, c := range group {
		GroupMetrics.Register(c)
		PrivateMetrics.Register(c)
	}

	// HTTP metrics
	http := []prometheus.Collector{
		HTTPCallCounter,
		HTTPLatency,
		HTTPInFlight,
	}
	for _, c := range http {
		HTTPMetrics.Register(c)
		PrivateMetrics.Register(c)
	}

	// Client metrics
	RegisterClientMetrics(ClientMetrics)
	RegisterClientMetrics(PrivateMetrics)
}

// RegisterClientMetrics registers drand client metrics with the given registry
func RegisterClientMetrics(r prometheus.Registerer) {
	// Client metrics
	client := []prometheus.Collector{
		ClientDNSLatencyVec,
		ClientInFlight,
		ClientLatencyVec,
		ClientRequests,
		ClientTLSLatencyVec,
		ClientWatchLatency,
		ClientHTTPHeartbeatLatency,
		ClientHTTPHeartbeatFailure,
	}
	for _, c := range client {
		r.Register(c)
	}
}

// PeerHandler abstracts a helper for relaying http requests to a group peer
type PeerHandler func(ctx context.Context) (map[string]http.Handler, error)

// Start starts a prometheus metrics server with debug endpoints.
func Start(metricsBind string, pprof http.Handler, peerHandler PeerHandler) net.Listener {
	log.DefaultLogger.Debug("metrics", "private listener started", "at", metricsBind)
	bindMetrics()

	l, err := net.Listen("tcp", metricsBind)
	if err != nil {
		log.DefaultLogger.Warn("metrics", "listen failed", "err", err)
		return nil
	}
	s := http.Server{Addr: l.Addr().String()}
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(PrivateMetrics, promhttp.HandlerOpts{Registry: PrivateMetrics}))

	if peerHandler != nil {
		mux.Handle("/peer/", &lazyPeerHandler{peerHandler})
	}

	if pprof != nil {
		mux.Handle("/debug/pprof/", pprof)
	}

	mux.HandleFunc("/debug/gc", func(w http.ResponseWriter, req *http.Request) {
		runtime.GC()
		fmt.Fprintf(w, "GC run complete")
	})
	s.Handler = mux
	go func() {
		log.DefaultLogger.Warn("metrics", "listen finished", "err", s.Serve(l))
	}()
	return l
}

// GroupHandler provides metrics shared to other group members
// This HTTP handler, which would typically be mounted at `/metrics` exposes `GroupMetrics`
func GroupHandler() http.Handler {
	return promhttp.HandlerFor(GroupMetrics, promhttp.HandlerOpts{Registry: GroupMetrics})
}

// lazyPeerHandler is a structure that defers learning who current
// group members are until an http request is received.
type lazyPeerHandler struct {
	peerHandler PeerHandler
}

func (l *lazyPeerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	addr := strings.Replace(r.URL.Path, "/peer/", "", 1)
	if strings.Contains(addr, "/") {
		addr = addr[:strings.Index(addr, "/")]
	}

	handlers, err := l.peerHandler(r.Context())
	if err != nil {
		log.DefaultLogger.Warn("metrics", "failed to get peer handlers", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	handler, ok := handlers[addr]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// The request to make to the peer is for its "/metrics" endpoint
	// Note that at present this shouldn't matter, since the only handler
	// mounted for the other end of these requests is `GroupHandler()` above,
	// so all paths / requests should see group metrics as a response.
	r.URL.Path = "/metrics"
	handler.ServeHTTP(w, r)
}
