package cerebra

import (
	"fmt"
	"strings"
	"time"

	"github.com/reef-pi/hal"
)

type cerebraOutlet struct {
	index int
	state bool
}

type cerebra struct {
	meta    hal.Metadata
	outlets []*cerebraOutlet
	hid     *cerebraHid
}

func NewCerebra(m hal.Metadata, sn string) *cerebra {

	c := &cerebra{
		meta:    m,
		outlets: make([]*cerebraOutlet, 6),
		hid:     &cerebraHid{serial: strings.ToUpper(sn)},
	}
	for i := range c.outlets {
		c.outlets[i] = &cerebraOutlet{index: i + 1}
	}
	go c.poll()
	return c
}

func (c *cerebra) Metadata() hal.Metadata {
	return c.meta
}

func (c *cerebra) Close() error {
	return nil
}

func (c *cerebra) Pins(cap hal.Capability) ([]hal.Pin, error) {
	switch cap {
	case hal.DigitalOutput:
		pins := make([]hal.Pin, len(c.outlets))
		for i, outlet := range c.outlets {
			pins[i] = outlet
		}
		return pins, nil
	default:
		return nil, fmt.Errorf("capability not supported")
	}
}

func (c *cerebra) DigitalOutputPins() []hal.DigitalOutputPin {
	pins := make([]hal.DigitalOutputPin, len(c.outlets))
	for i, outlet := range c.outlets {
		pins[i] = outlet
	}
	return pins
}

func (c *cerebra) DigitalOutputPin(pin int) (hal.DigitalOutputPin, error) {
	if pin > 0 && pin <= len(c.outlets) {
		return c.outlets[pin-1], nil
	}
	return nil, fmt.Errorf("unknown pin:%d", pin)
}

func (c *cerebra) poll() {
	for range time.Tick(time.Millisecond * 250) {
		if c.hid.ensureConnection() {
			for _, outlet := range c.outlets {
				if outlet.state {
					c.hid.powerOn(outlet.index)
				} else {
					c.hid.powerOff(outlet.index)
				}
			}
		}
	}
}

func (c *cerebraOutlet) Close() error {
	return nil
}

func (c *cerebraOutlet) Name() string {
	return fmt.Sprint("Outlet ", c.index)
}

func (c *cerebraOutlet) Number() int {
	return c.index
}

func (c *cerebraOutlet) LastState() bool {
	return c.state
}

func (c *cerebraOutlet) Write(b bool) error {
	c.state = b
	return nil
}
