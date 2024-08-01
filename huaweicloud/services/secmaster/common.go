package secmaster

import (
	"encoding/json"
	"errors"

	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"
)

// ParseQueryError403 is a method used to parse whether a 404 error message means the resources not found.
// + SecMaster.20010001: workspace does not exist.
func ParseQueryError403(err error, specError string) error {
	var err403 golangsdk.ErrDefault403
	if errors.As(err, &err403) {
		var apiError interface{}
		if jsonErr := json.Unmarshal(err403.Body, &apiError); jsonErr != nil {
			return err
		}

		errCode, searchErr := jmespath.Search("code", apiError)
		if searchErr != nil {
			return err
		}

		if errCode == specError {
			return golangsdk.ErrDefault404{}
		}
	}

	return err
}
