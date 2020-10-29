package application

// Component ...
type Component interface {
	Start()
	Stop()
}

// Components ...
type Components []Component

// Start ...
func (c Components) Start() {
	for _, component := range c {
		component.Start()
	}
}

// Stop ...
func (c Components) Stop() {
	for _, component := range c {
		component.Stop()
	}
}
