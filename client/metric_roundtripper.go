package client

import (
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/tkliner/go-gopay/client/logger"
)

type MetricsRoundTripper struct {
	next http.RoundTripper
	log  logger.Logger
}

func NewMetricsTransport(next http.RoundTripper, log logger.Logger) *MetricsRoundTripper {
	return &MetricsRoundTripper{
		next: next,
		log:  log,
	}
}

// RoundTrip implementuje rozhraní http.RoundTripper.
func (rt *MetricsRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var (
		startTime                 = time.Now()
		dnsStart, dnsDone         time.Time
		connectStart, connectDone time.Time
		firstByteTime             time.Time
	)

	trace := &httptrace.ClientTrace{
		DNSStart: func(info httptrace.DNSStartInfo) {
			dnsStart = time.Now()
			rt.log.Info(req.Context(), "Trace", "phase", "DNS_start", "host", info.Host)
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			dnsDone = time.Now()
			rt.log.Info(req.Context(), "Trace", "phase", "DNS_done", "duration_ms", dnsDone.Sub(dnsStart).Milliseconds())
		},
		ConnectStart: func(network, addr string) {
			connectStart = time.Now()
			rt.log.Info(req.Context(), "Trace", "phase", "connect_start", "address", addr)
		},
		ConnectDone: func(network, addr string, err error) {
			connectDone = time.Now()
			rt.log.Info(req.Context(), "Trace", "phase", "connect_done", "duration_ms", connectDone.Sub(connectStart).Milliseconds())
		},
		GotConn: func(info httptrace.GotConnInfo) {
			rt.log.Info(req.Context(), "Trace", "phase", "got_connection", "reused", info.Reused)
		},
		GotFirstResponseByte: func() {
			firstByteTime = time.Now()
			ttfb := firstByteTime.Sub(startTime)
			rt.log.Info(req.Context(), "Trace", "phase", "first_byte", "ttfb_ms", ttfb.Milliseconds())
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	resp, err := rt.next.RoundTrip(req)

	latency := time.Since(startTime)
	statusCode := 0
	if resp != nil {
		statusCode = resp.StatusCode
	}

	if resp != nil {
		defer func(start time.Time) {
			ttlb := time.Since(start)
			rt.log.Info(req.Context(), "Trace", "phase", "request_complete", "ttlb_ms", ttlb.Milliseconds())
		}(startTime)
	}

	rt.log.Info(req.Context(), "Request metrics",
		"method", req.Method,
		"url", req.URL.Path,
		"status", statusCode,
		"latency_ms", latency.Milliseconds(),
	)

	return resp, err
}
