package lts

import (
	"encoding/json"
	"errors"

	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"
)

// parseQueryError500 is a method used to parse whether a 500 error message means the resources not found.
// For the LTS service, there are some known 404 error codes:
// + LTS.2504: the member does not found.
func parseQueryError500(err error, specErrors []string) error {
	var err500 golangsdk.ErrDefault500
	if errors.As(err, &err500) {
		var apiError interface{}
		if jsonErr := json.Unmarshal(err500.Body, &apiError); jsonErr != nil {
			return err
		}

		errCode, searchErr := jmespath.Search("error_code", apiError)
		if searchErr != nil {
			return err
		}

		for _, v := range specErrors {
			if errCode == v {
				return golangsdk.ErrDefault404{}
			}
		}
	}
	return err
}
