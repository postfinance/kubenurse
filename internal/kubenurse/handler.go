package kubenurse

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/postfinance/kubenurse/pkg/kubediscovery"
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
			APIServerDirect string `json:"api_server_direct"`
			APIServerDNS    string `json:"api_server_dns"`
			MeIngress       string `json:"me_ingress"`
			MeService       string `json:"me_service"`

			// kubediscovery
			NeighbourhoodState string                    `json:"neighbourhood_state"`
			Neighbourhood      []kubediscovery.Neighbour `json:"neighbourhood"`
		}

		// Run checks now
		res, haserr := s.checker.Run()
		if haserr {
			w.WriteHeader(http.StatusInternalServerError)
		}

		// Add additional data
		out := Output{
			APIServerDNS:       res.APIServerDNS,
			APIServerDirect:    res.APIServerDirect,
			MeIngress:          res.MeIngress,
			MeService:          res.MeService,
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
