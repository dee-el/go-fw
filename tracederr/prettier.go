package tracederr

type Prettier struct {
	Error  string    `json:"error"`
	Traces []string  `json:"traces,omitempty"`
	Cause  *Prettier `json:"cause,omitempty"`
}

func PrintErrors(err error, pf *Prettier) *Prettier {
	if pf == nil {
		pf = &Prettier{}
	}

	msg := err.Error()
	e, ok := err.(*Error)
	if ok {
		if e.message != "" {
			msg = e.message
		}

		var traces []string
		for _, l := range e.traces {
			traces = append(traces, l.String())
		}

		pf.Traces = traces
	}

	pf.Error = msg

	err = Unwrap(err)
	if err != nil {
		pf.Cause = PrintErrors(err, pf.Cause)
	}

	return pf
}
