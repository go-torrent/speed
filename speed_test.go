package speed

import (
	"testing"
	"time"
)

func TestManualTick(t *testing.T) {
	c := &clock{t: 1, r: 4}
	g := NewGaugeWithClock(c)

	if s := g.Progress(10); s != 10 {
		t.Errorf("Expected %v, got %v", 10, s)
	}

	if s := g.Progress(20); s != 30 {
		t.Errorf("Expected %v, got %v", 30, s)
	}

	for i := 0; i < 2; i++ {
		c.doTick()

		if s := g.Progress(0); s != 30 {
			t.Errorf("Expected %v, got %v", 30, s)
		}
	}

	c.doTick() // 1 second has passed

	// 24 = 30 / 1.25 (1 + 1/r seconds)
	if s := g.Progress(0); s != 24 {
		t.Errorf("Expected %v, got %v", 24, s)
	}

	c.doTick()

	// 20 = 30 / 1.5 (1 + 2/r seconds)
	if s := g.Progress(0); s != 20 {
		t.Errorf("Expected %v, got %v", 20, s)
	}
}

func TestSpeedometer(t *testing.T) {
	g := NewGauge()

	if s := g.Progress(100); s != 100 {
		t.Errorf("Expected %v, got %v", 100, s)
	}

	oneTick := time.Duration(1000/DefaultResolution) * time.Millisecond
	<-time.After(time.Second - oneTick)

	if s := g.Read(); s != 80 {
		t.Errorf("Expected %v, got %v", 80, s)
	}
}
