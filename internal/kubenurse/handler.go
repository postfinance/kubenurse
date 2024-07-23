package kubenurse

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/postfinance/kubenurse/internal/servicecheck"
)

func (s *Server) readyHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		if s.ready.Load() {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}
}

func (s *Server) aliveHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type Output struct {
			Hostname   string              `json:"hostname"`
			Headers    map[string][]string `json:"headers"`
			UserAgent  string              `json:"user_agent"`
			RequestURI string              `json:"request_uri"`
			RemoteAddr string              `json:"remote_addr"`
			Result     map[string]any      `json:"last_check_result"`
		}

		res := s.checker.LastCheckResult
		if res == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Add additional data
		out := Output{
			Result:     res,
			Headers:    r.Header,
			UserAgent:  r.UserAgent(),
			RequestURI: r.RequestURI,
			RemoteAddr: r.RemoteAddr,
		}
		out.Hostname, _ = os.Hostname()

		// Generate output output
		enc := json.NewEncoder(w)
		enc.SetIndent("", " ")
		_ = enc.Encode(out)
	}
}

func (s *Server) alwaysHappyHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(_ http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get(servicecheck.NeighbourOriginHeader)
		if origin != "" {
			s.neighboursTTLCache.Insert(origin)
		}
	}
}
