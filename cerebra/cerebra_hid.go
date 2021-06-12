package cerebra

import (
	"fmt"
	"strings"

	"github.com/karalabe/hid"
)

type cerebraHid struct {
	serial  string
	device  *hid.Device
	counter byte
}

func (c *cerebraHid) ensureConnection() bool {
	//If we're already connected, we don't need to do anything
	if c.device != nil {
		return true
	}

	//Look for a usb device matching the requested serial number
	for _, d := range hid.Enumerate(0x275a, 0x8044) {
		if c.serial == "00000000" || c.serial == strings.ToUpper(d.Serial) {
			dev, e := d.Open()
			if e != nil {
				fmt.Printf("cerebra: could not open USB: %s\n", e.Error())
				return false
			}
			fmt.Printf("cerebra: found device with serial number: %s\n", d.Serial)
			c.device = dev
			return true
		}
	}
	return false
}

func (c *cerebraHid) write(data []byte) {
	if c.device != nil {
		_, e := c.device.Write(data)
		if e != nil {
			fmt.Printf("cerebra: could not write USB: %s\n", e.Error())
			c.device = nil
		}
	}
}

/*
The below messages are hard-coded functions to turn on and off an outlet.  These should be able to be abstracted
into a general-purpose function to get or set any property of any component on the device.  For a full list of
function IDs and property codes, refer to com.scitronix.star.hardware.android/databases/star.db on the Cerebra
head unit SD card.
*/

func (c *cerebraHid) powerOn(index int) {
	data := [...]byte{
		//Common header
		0,
		254,
		0,
		0,
		8,
		//Write header
		59,
		0,
		0,
		c.counter,
		6,           //Function type (5=get, 6=set)
		byte(index), //Function ID
		0,
		4, //value length
		7, //Property code
		0,
		'o',
		'n',
		0,
	}
	c.counter++
	c.write(data[:])
}

func (c *cerebraHid) powerOff(index int) {
	data := [...]byte{
		//Common header
		0,
		254,
		0,
		0,
		8,
		//Write header
		59,
		0,
		0,
		c.counter,
		6,           //Function type (5=get, 6=set)
		byte(index), //Function ID
		0,
		5, //value length
		7, //Property code
		0,
		'o',
		'f',
		'f',
		0,
	}
	c.counter++
	c.write(data[:])
}
