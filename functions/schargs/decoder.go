package schargs

import (
	"fmt"
)

// decoder for extract variable name from string.
// variable like: ${var1}
type decoder struct {
	data string
	off  int // next read offset in data
}

func newDecoder(data string) *decoder {
	return &decoder{
		data: data,
		off:  0,
	}
}

func (d *decoder) next() {
	d.off += 1
}

// return (value, hasMore)
func (d *decoder) read() (byte, bool) {
	i, data := d.off, d.data
	if i < len(data) {
		c := data[i]
		d.off += 1
		return c, true
	}
	return 0, false
}

// return (value, hasMore)
func (d *decoder) peek() (byte, bool) {
	i, data := d.off, d.data
	if i < len(data) {
		c := data[i]
		return c, true
	}
	return 0, false
}

// isVariableBegin for check the next word whether is a valid variable begging.
func (d *decoder) isVariableBegin(b byte) bool {
	// b == '$' and next byte is '{' means next string is a variable.
	if b == '$' {
		b1, ok := d.peek()
		if !ok {
			// No more data.
			return false
		}
		if b1 == '{' {
			d.next()
			return true
		}
	}
	return false
}

func (d *decoder) extractVariable() (string, error) {
	var (
		v     []byte
		found bool
	)

LOOP:
	for {
		b, ok := d.read()
		if !ok {
			break LOOP
		}
		if d.isVariableBegin(b) {
			return "", ErrUnsupportedVariableNested
		}
		if b == '}' {
			found = true
			break LOOP
		}
		v = append(v, b)
	}

	if !found {
		return "", ErrInvalidVariableFormat
	}

	r := string(v)
	if r == "" {
		return "", ErrVariableIsEmpty
	}
	return r, nil
}

// ExtractVariables for extract all valid variables in giving data.
func (d *decoder) extractVariables() ([]string, error) {
	var keys []string
LOOP:
	for {
		b, ok := d.read()
		if !ok {
			break LOOP
		}
		if d.isVariableBegin(b) {
			key, err := d.extractVariable()
			if err != nil {
				return nil, err
			}
			keys = append(keys, key)
		}
	}
	return keys, nil
}

func (d *decoder) encode(valueMap map[string]string) (string, error) {
	var (
		value string
	)

	if valueMap == nil {
		// To avoids nil-pointer panic.
		valueMap = make(map[string]string)
	}
	// store the new data.
	nd := make([]byte, 0, len(d.data))

LOOP:
	for {
		b, ok := d.read()
		if !ok {
			break LOOP
		}
		if d.isVariableBegin(b) {
			key, err := d.extractVariable()
			if err != nil {
				return "", err
			}
			value, ok = valueMap[key]
			if !ok {
				return "", fmt.Errorf("varialbe [%s] not defined", key)
			}
			nd = append(nd, value...)
		} else {
			nd = append(nd, b)
		}
	}
	return string(nd), nil
}
