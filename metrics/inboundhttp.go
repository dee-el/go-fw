package metrics

import "time"

// InboundHTTPRecorder is metric recorder for every inbound traffics to HTTP API.
type InboundHTTPRecorder interface {
	Record(path, method, status string, duration time.Duration)
}
