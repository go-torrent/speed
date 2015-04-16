// Package speed provides a gauge to calculate a average of bytes transfered per second
// It's borrowed from some other code, and I believe this is some king of
// http://en.wikipedia.org/wiki/Moving_average or other sampling-based strategy
package speed

import (
	"math"
	"sync"
	"time"
)

const (
	// DefaultResolution is the frequency the clock tick
	DefaultResolution = 4
	// MaxTick is the max value for the tick
	MaxTick = math.MaxUint16
)

var (
	globalClock = newLockedClock(DefaultResolution)
)

func init() {
  globalClock.Start()
}

// A Clock represents a periodic ticker
type Clock interface {
	Start()
	Tick() int
	Resolution() int
}

type lockedClock struct {
	lk    sync.Mutex
	clock Clock
}

func (r *lockedClock) Start() {
	r.lk.Lock()
	r.clock.Start()
	r.lk.Unlock()
}

func (r *lockedClock) Tick() (t int) {
	r.lk.Lock()
	t = r.clock.Tick()
	r.lk.Unlock()
	return
}

func (r *lockedClock) Resolution() int {
	return r.clock.Resolution()
}

func newLockedClock(r int) Clock {
	return &lockedClock{clock: NewClock(r)}
}

type clock struct {
	t      int
	r      int
	ticker *time.Ticker
}

// Start starts the clock tick if not started
func (c *clock) Start() {
	if c.ticker != nil {
		return
	}

	c.ticker = newTickerFor(c.r)
	go c.tickLoop()
}

func (c *clock) Tick() int {
	return c.t
}

func (c *clock) Resolution() int {
	return c.r
}

func newTickerFor(r int) *time.Ticker {
	return time.NewTicker(time.Duration(1000/r) * time.Millisecond)
}

func (c *clock) tickLoop() {
	for range c.ticker.C {
		c.doTick()
	}
}

func (c *clock) doTick() int {
	c.t = (c.t + 1) & MaxTick
	return c.t
}

// NewClock returns a new Clock with a resolution of r ticks per second
func NewClock(r int) Clock {
	return &clock{
		t: 1,
		r: r,
	}
}

type Gauge struct {
	clock      Clock
	size       int
	buffer     []int
	bufferSize int
	pointer    int
	last       int
}

func NewGauge() Gauge {
	return NewGaugeWithClock(globalClock)
}

func NewGaugeWithClock(c Clock) Gauge {
	b := 5 // buffer
	res := c.Resolution()
	size := res * b

	return Gauge{
		clock:      c,
		size:       size,
		buffer:     make([]int, size),
		bufferSize: 1,
		pointer:    1,
		last:       (c.Tick() - 1) & MaxTick,
	}
}

func (g *Gauge) Progress(delta int) float32 {
	tick := g.clock.Tick()
	dist := (tick - g.last) & MaxTick

	if dist > g.size {
		dist = g.size
	}
	g.last = tick

	for ; dist > 0; dist-- {
		if g.pointer == g.size {
			g.pointer = 0
		}

		var v int
		if g.pointer == 0 {
			v = g.size
		} else {
			v = g.pointer
		}

		g.buffer[g.pointer] = g.buffer[v-1]

		if g.bufferSize-1 < g.pointer {
			g.bufferSize++
		}

		g.pointer++
	}

	if delta != 0 {
		g.buffer[g.pointer-1] += delta
	}

	var top = float32(g.buffer[g.pointer-1])
	var btm float32
	if g.bufferSize < g.size {
		btm = 0
	} else {
		var v int
		if g.pointer == g.size {
			v = 0
		} else {
			v = g.pointer
		}
		btm = float32(g.buffer[v])
	}

	if g.bufferSize < g.clock.Resolution() {
		return float32(top)
	}

	return (top - btm) * float32(g.clock.Resolution()) / float32(g.bufferSize)
}

func (g *Gauge) Read() float32 {
	return g.Progress(0)
}
