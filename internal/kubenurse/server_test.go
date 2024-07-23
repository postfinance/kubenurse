package kubenurse

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestCombined(t *testing.T) {
	r := require.New(t)

	os.Setenv("KUBENURSE_EXTRA_CHECKS", "cloudy_endpoint:http://cloudy.enpdoint:1234/test|ep_number_two:http://interesting.endpoint:8080/abcd")
	fakeClient := fake.NewFakeClient()
	kubenurse, err := New(fakeClient)
	r.NoError(err)
	r.NotNil(kubenurse)

	r.Equal(map[string]string{
		"ep_number_two":   "http://interesting.endpoint:8080/abcd",
		"cloudy_endpoint": "http://cloudy.enpdoint:1234/test",
	}, kubenurse.checker.ExtraChecks)

	t.Run("start/stop", func(t *testing.T) {
		r := require.New(t)
		errc := make(chan error, 1)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		go func() {
			// blocks until shutdown is called
			err := kubenurse.Run(ctx)

			errc <- err
			close(errc)
			cancel()
		}()

		// Shutdown, Run() should stop after function completes
		err := kubenurse.Shutdown()
		r.NoError(err)

		err = <-errc // blocks until kubenurse.Run() finishes and eventually returns an error
		r.NoError(err)
	})
}
