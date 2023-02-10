package test

import "testing"

func AssertErrorString(t *testing.T, err error, msg string) {
	if err.Error() != msg {
		t.Errorf("got error %s, expected error %s", err.Error(), msg)
	}
}

func AssertNil(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Error was not nil, got: %s", err.Error())
	}
}
