package tui

import (
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestCalcPercent(t *testing.T) {
	is := is.New(t)

	is.Equal(calcPercent(5*time.Second, 100*time.Second), 0.05)
}
