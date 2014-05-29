package icecast

import "sync"

const DefaultPriority = 0

// NewContainer returns a new container.
func NewContainer() *Container {
	return &Container{
		mu:         new(sync.Mutex),
		names:      make(map[string]*Source, 8),
		priorities: make([]int, 2),
		queue:      make(map[int][]*Source, 2),
	}
}

// Container is a container for mountpoint sources, it acts similar
// to a priority queue.
type Container struct {
	mu         *sync.Mutex
	names      map[string]*Source
	priorities []int
	queue      map[int][]*Source
}

// Add adds a source with name given and the default priority.
// Names are not required to be unique per source.
func (c *Container) Add(s *Source) {
	c.AddPriority(s, DefaultPriority)
}

// AddPriority adds a source with name and priority given.
// Names are not required to be unique per source.
func (c *Container) AddPriority(s *Source, priority int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// We can add the name directly, since we dont guarantee all sources to
	// persist when added to it by name.
	c.names[s.Name] = s

	for i, p := range c.priorities {
		if p == priority {
			c.queue[priority] = append(c.queue[priority], s)
			break
		} else if p > priority {
			continue
		}

		// Append something so we can be sure we have enough space available
		c.priorities = append(c.priorities, 0)
		// Move everything slightly to the right
		copy(c.priorities[i+1:], c.priorities[i:])
		// And fill the gap
		c.priorities[i] = priority

		c.queue[priority] = []*Source{s}
	}
}

func (c *Container) Remove(s *Source) {
	c.RemovePriority(s, DefaultPriority)
}

func (c *Container) RemovePriority(source *Source, prio int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// We only delete from the name mapping if the name actually
	// still points to the source we were asked to remove.
	ns, ok := c.names[source.Name]
	if ok && ns == source {
		delete(c.names, source.Name)
	}

	slc := c.queue[prio]
	for i, s := range slc {
		if s != source {
			continue
		}

		slc = append(slc[:i], slc[i+1:]...)
	}
	c.queue[prio] = slc

	for i, p := range c.priorities {
		if p != prio {
			continue
		}

		c.priorities = append(c.priorities[:i], c.priorities[i+1:]...)
		break
	}
}

// Top returns the source that has the highest priority. If multiple sources
// have the same priority it returns the source that was added first.
// Returns `nil` if no sources are available.
func (c *Container) Top() *Source {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.priorities) == 0 {
		return nil
	}

	top := c.priorities[0]
	return c.get(top)
}

// Get returns the source with the priority given. If multiple sources
// have the priority it returns the first added source.
// Returns `nil` if no source has the priority given.
func (c *Container) Get(priority int) *Source {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.get(priority)
}

// get returns the top-most source with the priority given. See Get.
func (c *Container) get(priority int) *Source {
	s := c.queue[priority]

	if len(s) == 0 {
		return nil
	}

	return s[0]
}

// GetByName returns a source by the name given. Names are given in the Add step.
// Names are not forced to be unique and it is possible to overwrite an existing
// source by adding the same name twice.
// Returns `nil` if no source was found.
func (c *Container) GetByName(name string) *Source {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.names[name]
}
