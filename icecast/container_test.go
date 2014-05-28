package icecast

import "testing"

func TestContainerSingleEntry(t *testing.T) {
	t.Parallel()
	var (
		c = NewContainer()
		s = new(Source)
	)

	c.Add("test", s)
	if res := c.Top(); res != s {
		t.Log("got: %p want: %p", res, s)
		t.Error("container did not return expected top source")
	}

	if res := c.Get(DefaultPriority); res != s {
		t.Logf("got: %p want: %p", res, s)
		t.Error("container did not return expected source")
	}

	if res := c.GetByName("test"); res != s {
		t.Logf("got: %p want: %p", res, s)
		t.Error("container did not return expected (named) source")
	}
}

func TestContainerPriorities(t *testing.T) {
	t.Parallel()
	var (
		s *Source

		c             = NewContainer()
		highest       = 50
		highestSource = new(Source)
		prios         = []int{10, 20, 30, 40, highest, 5, 15, 25, 35, 45}
	)

	for _, prio := range prios {
		if prio != highest {
			s = new(Source)
		} else {
			s = highestSource
		}

		t.Logf("added %p for prio %d", s, prio)
		c.AddPriority("test", s, prio)
	}

	if res := c.Top(); res != highestSource {
		t.Logf("got: %p want: %p", res, highestSource)
		t.Error("container did not return expected top source")
	}

	if res := c.Get(highest); res != highestSource {
		t.Logf("got: %p want: %p", res, highestSource)
		t.Error("container did not return expected source")
	}

	// This one should equal the last thing we added in the loop, so `s`
	if res := c.GetByName("test"); res != s {
		t.Logf("got: %p want: %p", res, s)
		t.Error("container did not return expected (named) source")
	}
}
