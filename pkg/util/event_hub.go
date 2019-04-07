package util

type EventHub struct {
	Send       chan string
	Register   chan *EventClient
	Unregister chan *EventClient
	Clients    map[*EventClient]struct{}
}

type EventClient struct {
	Events chan string
	Filter map[string]struct{}
}

func NewEventHub() *EventHub {
	e := &EventHub{
		Send:       make(chan string),
		Register:   make(chan *EventClient),
		Unregister: make(chan *EventClient),
		Clients:    make(map[*EventClient]struct{}),
	}
	e.run()
	return e
}

func (e *EventHub) NewClient(filter []string) *EventClient {
	c := &EventClient{
		Events: make(chan string),
		Filter: make(map[string]struct{}),
	}
	for _, e := range filter {
		c.Filter[e] = struct{}{}
	}
	e.Register <- c
	return c
}

func (e *EventHub) run() {
	for {
		select {
		// Add new client
		case c := <-e.Register:
			e.Clients[c] = struct{}{}
			// Remove client
		case c := <-e.Unregister:
			if _, ok := e.Clients[c]; ok {
				close(c.Events)
				delete(e.Clients, c)
			}
			// Send to all clients
		case event := <-e.Send:
			for c, _ := range e.Clients {
				if _, ok := c.Filter[event]; ok {
					select {
					case c.Events <- event:
					default:
						close(c.Events)
						delete(e.Clients, c)
					}
				}
			}
		}
	}
}

func (c *EventClient) WaitEvent(eventMatch string) {
	for {
		select {
		case event := <-c.Events:
			if event == eventMatch {
				return
			}
		}
	}
}

func (c *EventClient) Drain() {
	for {
		select {
		case <-c.Events:
		default:
			return
		}
	}
}
