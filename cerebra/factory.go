package cerebra

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/reef-pi/hal"
)

const snParam = "Serial Number"

type cerebraFactory struct {
	meta       hal.Metadata
	parameters []hal.ConfigParameter
}

var factory *cerebraFactory
var once sync.Once

func CerebraFactory() hal.DriverFactory {
	once.Do(func() {
		factory = &cerebraFactory{
			meta: hal.Metadata{
				Name:        "cerebra",
				Description: "Supports Cerebra MultiBar",
				Capabilities: []hal.Capability{
					hal.DigitalOutput,
				},
			},
			parameters: []hal.ConfigParameter{
				{
					Name:    snParam,
					Type:    hal.String,
					Order:   0,
					Default: "00000000",
				},
			},
		}
	})

	return factory
}

func (f *cerebraFactory) Metadata() hal.Metadata {
	return f.meta
}

func (f *cerebraFactory) GetParameters() []hal.ConfigParameter {
	return f.parameters
}

//We set hal.ConfigParameter.Type to be a string, but numeric values still get passed in as int.
//Use this function to ensure the incoming value is always a string.
func getFormattedSerial(v interface{}) string {
	switch v := v.(type) {
	case int:
		return fmt.Sprintf("%08d", v)
	case string:
		return v
	default:
		return ""
	}
}

func (f *cerebraFactory) ValidateParameters(parameters map[string]interface{}) (bool, map[string][]string) {
	var failures = make(map[string][]string)
	var v interface{}
	var ok bool

	if v, ok = parameters[snParam]; ok {
		s := getFormattedSerial(v)

		if len(s) != 8 {
			failure := fmt.Sprint(snParam, " must be 8 digits long.")
			failures[snParam] = append(failures[snParam], failure)
		}
		_, err := strconv.ParseInt(s, 16, 64)
		if err != nil {
			failure := fmt.Sprint(snParam, " must only contain digits 0-9, A-F.")
			failures[snParam] = append(failures[snParam], failure)
		}
	} else {
		failure := fmt.Sprint(snParam, " is missing.")
		failures[snParam] = append(failures[snParam], failure)
	}

	return len(failures) == 0, failures
}

func (f *cerebraFactory) NewDriver(parameters map[string]interface{}, hardwareResources interface{}) (hal.Driver, error) {
	if valid, failures := f.ValidateParameters(parameters); !valid {
		return nil, errors.New(hal.ToErrorString(failures))
	}

	driver := NewCerebra(f.meta, getFormattedSerial(parameters[snParam]))
	return driver, nil
}
