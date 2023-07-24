package json

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// Decode interface to given struct using json tag.
func Decode(input, output interface{}) error {
	// Only allow pointer output.
	if reflect.ValueOf(output).Kind() != reflect.Ptr {
		return errors.New("given output is not a pointer")
	}

	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		// Tell decoder to use json tag.
		TagName: "json",
		Result:  output,
		// Make sure to add `squash` option, to correctly restore anonymous embedded struct.
		// https://github.com/mitchellh/mapstructure/issues/113
		//
		// Also see decoder_test.go for more information.
		Squash: true,
	})
	if err != nil {
		return errors.Wrap(err, "Decode: failed create new decoder")
	}

	if err := d.Decode(input); err != nil {
		return errors.Wrap(err, "Decode: failed decode")
	}

	return nil
}
