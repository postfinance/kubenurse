package kubenurse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTTLCache(t *testing.T) {

	c := TTLCache[string]{}
	c.Init(10 * time.Millisecond)

	c.Insert("node-a")
	time.Sleep(5 * time.Millisecond)

	require.Equal(t, 1, c.ActiveEntries(), "after 5ms and with a 10ms TTL, there should have been 1 entry in the cache")
	time.Sleep(5 * time.Millisecond)
	require.Equal(t, 0, c.ActiveEntries(), "after 10ms and with a 10ms TTL, there should be no entry in the cache")

}
