package es_handler

type eventHub struct {
	send       chan string
	register   chan *client
	unregister chan *client
	clients    map[*client]struct{}
}

type client struct {
	events chan string
	filter map[string]struct{}
}

func newEventHub() *eventHub {
	return &eventHub{
		send:       make(chan string),
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[*client]struct{}),
	}
}

func (e *eventHub) newClient(filter []string) *client {
	c := &client{
		events: make(chan string),
		filter: make(map[string]struct{}),
	}
	for _, e := range filter {
		c.filter[e] = struct{}{}
	}
	e.register <- c
	return c
}

func (e *eventHub) run() {
	for {
		select {
		// Add new client
		case c := <-e.register:
			e.clients[c] = struct{}{}
			// Remove client
		case c := <-e.unregister:
			if _, ok := e.clients[c]; ok {
				close(c.events)
				delete(e.clients, c)
			}
			// Send to all clients
		case event := <-e.send:
			for c, _ := range e.clients {
				if _, ok := c.filter[event]; ok {
					select {
					case c.events <- event:
					default:
						close(c.events)
						delete(e.clients, c)
					}
				}
			}
		}
	}
}

func (c *client) waitEvent(eventMatch string) {
	for {
		select {
		case event := <-c.events:
			if event == eventMatch {
				return
			}
		}
	}
}

func (c *client) drain() {
	for {
		select {
		case <-c.events:
		default:
			return
		}
	}
}
