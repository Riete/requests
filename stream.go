package requests

import (
	"bufio"
	"encoding/json"
	"io"
	"reflect"
)

func readStringStream[T any](r io.Reader, event chan T) error {
	var err error
	br := bufio.NewReader(r)
	var line string
	for {
		if line, err = br.ReadString('\n'); err != nil {
			break
		}
		event <- reflect.ValueOf(line).Interface().(T)
	}
	if err == io.EOF {
		return nil
	}
	return err
}

func readJsonStream[T any](r io.Reader, event chan T) error {
	var err error
	decoder := json.NewDecoder(r)
	for {
		var data T
		if err = decoder.Decode(&data); err != nil {
			break
		}
		event <- data
	}
	if err == io.EOF {
		return nil
	}
	return err
}

func ReadStream[T any](r io.ReadCloser, event chan T) error {
	defer close(event)
	defer r.Close()
	if reflect.TypeFor[T]().Kind() == reflect.String {
		return readStringStream(r, event)
	}
	return readJsonStream(r, event)
}
