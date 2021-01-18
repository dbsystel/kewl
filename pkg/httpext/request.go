package httpext

import (
	"encoding/json"
	"io/ioutil"
	"reflect"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
)

// ReadBodyAsJSON will try to read the body as JSON into the provided object
func (d *Request) ReadBodyAsJSON(obj interface{}) error {
	bytes, err := d.ReadBody()
	if err != nil {
		return errors.Wrap(err, "unable to read body")
	}
	if len(bytes) == 0 {
		return errors.New("body is empty")
	}

	if err := json.Unmarshal(bytes, obj); err != nil {
		return errors.Wrapf(err, "could not unmarshal body to: %v", reflect.TypeOf(obj))
	}

	return nil
}

// ReadBody will read the whole body, close it and return the contents
func (d *Request) ReadBody() ([]byte, error) {
	defer func() {
		if err := d.Body.Close(); err != nil {
			d.logger.Error(err, "could not close body")
		}
	}()
	bytes, err := ioutil.ReadAll(d.Body)
	err = errors.Wrap(err, "could not read request body")
	return bytes, err
}

// Logger returns the logger for the request
func (d *Request) Logger() logr.Logger {
	return d.logger
}
