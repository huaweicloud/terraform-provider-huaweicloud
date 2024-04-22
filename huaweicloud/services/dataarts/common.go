package dataarts

import (
	"encoding/json"
	"errors"

	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"
)

// ParseQueryError400 is a method used to parse whether a 404 error message means the resources not found.
// For the DataArts Studio service, there are some known 404 error codes:
// + Workspace:
//   - DLM.4001: Instance or workspace does not exist.
//   - DLS.6036: Workspace does not exist.
//
// + Data Service:
//   - DLM.4205: Catalog does not found.
//   - DLM.4018: API does not exist (during API detail query).
//
// + Security:
//   - DLM.3027: Permission set does not found.
//   - DLS.1000: Data secrecy level does not found.
func ParseQueryError400(err error, specErrors []string) error {
	var err400 golangsdk.ErrDefault400
	if errors.As(err, &err400) {
		var apiError interface{}
		if jsonErr := json.Unmarshal(err400.Body, &apiError); jsonErr != nil {
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
