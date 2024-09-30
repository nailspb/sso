package errorHelper

import "fmt"

func WrapError(operation string, message string, err error) error {
	return fmt.Errorf("%s in %s -> `%w`", message, operation, err)
}
