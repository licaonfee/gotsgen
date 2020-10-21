package gotsgen_test

import (
	"testing"
	"time"

	"github.com/intercloud/gotsgen"
)

func TestQuery(t *testing.T) {
	duration, _ := time.ParseDuration("24h")
	end := time.Now()
	start := end.Add(-duration)

	ts, err := gotsgen.Query(start, end, 200, "rand")

	if err != nil {
		t.Errorf("Unexpected error %s\n", err.Error())
	}
	if len(ts.XValues) != 200 {
		t.Errorf("Expected time series to have 200 values got %d\n", len(ts.XValues))
	}

	ts, err = gotsgen.Query(start, end, 200, "fake")
	if err == nil || err.Error() != "Unknown generator type" {
		t.Errorf("Expected error Unknown generator type but got %v\n", err)
	}

	ts, err = gotsgen.Query(start, start, 200, "norm")
	if err == nil || err.Error() != "Bad time range" {
		t.Errorf("Expected error Unknown generator type but got %v\n", err)
	}

	ts, err = gotsgen.Query(end, start, 200, "deriv")
	if err == nil || err.Error() != "Bad time range" {
		t.Errorf("Expected error Unknown generator type but got %v\n", err)
	}
}

func BenchmarkQuery(b *testing.B) {
	duration := time.Hour * 24
	end := time.Now()
	start := end.Add(-duration)
	samples := uint(duration.Seconds())
	type args struct {
		gentype string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "norm",
			args: args{gentype: "norm"},
		},
		{
			name: "deriv",
			args: args{gentype: "deriv"},
		},
		{
			name: "rand",
			args: args{gentype: "rand"},
		},
	}
	for _, bb := range tests {
		b.Run(bb.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				gotsgen.Query(start, end, samples, bb.name)
			}
		})
	}
}
