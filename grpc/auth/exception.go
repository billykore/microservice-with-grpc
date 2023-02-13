package main

import (
	"fmt"
)

func invalidGrantTypeException(req *request) error {
	return fmt.Errorf("invalid grant type %q", req.grantType)
}
