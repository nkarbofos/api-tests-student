package wait

import (
	"testing"
	"time"
)

const pollInterval = 200 * time.Millisecond

// Until повторяет проверку до первого успеха или истечения таймаута и заменяет фиксированный sleep в асинхронных сценариях
func Until(t *testing.T, timeout time.Duration, description string, check func() bool) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if check() {
			return
		}
		time.Sleep(pollInterval)
	}

	t.Fatalf("%s за %s", description, timeout)
}
