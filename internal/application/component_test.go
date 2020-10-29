package application_test

import (
	"testing"

	"github.com/brycedarling/go-practical-microservices/internal/application"
)

type testComponent struct {
	StartCalled bool
	StopCalled  bool
}

func (c *testComponent) Start() {
	c.StartCalled = true
}

func (c *testComponent) Stop() {
	c.StopCalled = true
}

func TestComponentsStart(t *testing.T) {
	comp := &testComponent{}
	comps := application.Components{comp}

	comps.Start()

	if !comp.StartCalled {
		t.Errorf("want %v, got %v", true, comp.StartCalled)
	}
}

func TestComponentsStop(t *testing.T) {
	comp := &testComponent{}
	comps := application.Components{comp}

	comps.Stop()

	if !comp.StopCalled {
		t.Errorf("want %v, got %v", true, comp.StopCalled)
	}
}
