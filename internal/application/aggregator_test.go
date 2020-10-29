package application_test

import (
	"testing"

	"github.com/brycedarling/go-practical-microservices/internal/application"
)

type testAggregator struct {
	StartCalled bool
	StopCalled  bool
}

func (a *testAggregator) Start() {
	a.StartCalled = true
}

func (a *testAggregator) Stop() {
	a.StopCalled = true
}

func TestAggregatorsStart(t *testing.T) {
	agg := &testAggregator{}
	aggs := application.Aggregators{agg}

	aggs.Start()

	if !agg.StartCalled {
		t.Errorf("want %v, got %v", true, agg.StartCalled)
	}
}

func TestAggregatorsStop(t *testing.T) {
	agg := &testAggregator{}
	aggs := application.Aggregators{agg}

	aggs.Stop()

	if !agg.StopCalled {
		t.Errorf("want %v, got %v", true, agg.StopCalled)
	}
}
