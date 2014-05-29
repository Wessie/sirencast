package icecast

import "testing"

func TestContainerSingleEntry(t *testing.T) {
	t.Parallel()
	var (
		c = NewContainer()
		s = &Source{Name: "test"}
	)

	c.Add(s)
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
		highestSource = &Source{Name: "test"}
		prios         = []int{10, 20, 30, 40, highest, 5, 15, 25, 35, 45}
	)

	for _, prio := range prios {
		if prio != highest {
			s = &Source{Name: "test"}
		} else {
			s = highestSource
		}

		t.Logf("added %p for prio %d", s, prio)
		c.AddPriority(s, prio)
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

	c.RemovePriority(highestSource, highest)

	if res := c.Top(); res == highestSource {
		t.Logf("got: %p want something else", res)
		t.Error("container did not remove highest source")
	}
}

func BenchmarkContainerAdd(b *testing.B) {
	var (
		s = &Source{Name: "test"}
		c = NewContainer()
	)

	for i := 0; i < b.N; i++ {
		c.Add(s)
	}
}

func BenchmarkContainerAddPriority(b *testing.B) {
	var (
		s = &Source{Name: "test"}
		c = NewContainer()
	)

	for i := 0; i < b.N; i++ {
		c.AddPriority(s, i%5)
	}
}

func BenchmarkContainerAddRemove(b *testing.B) {
	var (
		s = &Source{Name: "test"}
		c = NewContainer()
	)

	for i := 0; i < b.N; i++ {
		c.Add(s)
		c.Remove(s)
	}
}

func BenchmarkContainerAddPriorityRemove(b *testing.B) {
	var (
		s = &Source{Name: "test"}
		c = NewContainer()
	)
	for i := 0; i < b.N; i++ {
		c.AddPriority(s, i%10)
		c.RemovePriority(s, i%10)
	}
}

func BenchmarkContainerTop(b *testing.B) {
	var (
		s = &Source{Name: "test"}
		c = NewContainer()
	)

	for i := 0; i < 10; i++ {
		c.Add(s)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = c.Top()
	}
}
