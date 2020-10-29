package application

// Aggregator is a
type Aggregator interface {
	Start()
	Stop()
}

// Aggregators ...
type Aggregators []Aggregator

// Start ...
func (a Aggregators) Start() {
	for _, aggregator := range a {
		aggregator.Start()
	}
}

// Stop ...
func (a Aggregators) Stop() {
	for _, aggregator := range a {
		aggregator.Stop()
	}
}
