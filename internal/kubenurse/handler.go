package kubenurse

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/postfinance/kubenurse/internal/servicecheck"
)

func (s *Server) readyHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		if s.ready {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
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

			// checker.Result
			servicecheck.Result

			// kubediscovery
			NeighbourhoodState string                   `json:"neighbourhood_state"`
			Neighbourhood      []servicecheck.Neighbour `json:"neighbourhood"`
		}

		res := s.checker.LastCheckResult
		if res == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Add additional data
		out := Output{
			Result:             *res,
			Headers:            r.Header,
			UserAgent:          r.UserAgent(),
			RequestURI:         r.RequestURI,
			RemoteAddr:         r.RemoteAddr,
			Neighbourhood:      res.Neighbourhood,
			NeighbourhoodState: res.NeighbourhoodState,
		}
		out.Hostname, _ = os.Hostname()

		// Generate output output
		enc := json.NewEncoder(w)
		enc.SetIndent("", " ")
		_ = enc.Encode(out)
	}
}
