package kubenurse

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestCombined(t *testing.T) {
	r := require.New(t)

	fakeClient := fake.NewFakeClient()
	kubenurse, err := New(fakeClient)
	r.NoError(err)
	r.NotNil(kubenurse)

	t.Run("start/stop", func(t *testing.T) {
		r := require.New(t)
		errc := make(chan error, 1)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
