package management_test

import (
	"testing"

	mgmtapp "github.com/ai-finops/ai-finops/services/management-service/internal/application/management"
)

func TestSoftThresholdBreached(t *testing.T) {
	tests := []struct {
		name     string
		spend    float64
		budget   float64
		pct      int32
		want     bool
	}{
		{"under", 50, 100, 80, false},
		{"exact soft", 80, 100, 80, true},
		{"over", 95, 100, 80, true},
		{"zero budget", 10, 0, 80, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := mgmtapp.SoftThresholdBreached(tc.spend, tc.budget, tc.pct); got != tc.want {
				t.Fatalf("got %v want %v", got, tc.want)
			}
		})
	}
}

func TestUsagePercent(t *testing.T) {
	if got := mgmtapp.UsagePercent(25, 100); got != 25 {
		t.Fatalf("got %v", got)
	}
	if got := mgmtapp.UsagePercent(10, 0); got != 0 {
		t.Fatalf("zero budget => 0, got %v", got)
	}
}
