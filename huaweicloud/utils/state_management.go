package utils

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func RefreshObjectParamOriginValues(d *schema.ResourceData, objectParamKeys []string) error {
	var mErr *multierror.Error

	for _, objectParamKey := range objectParamKeys {
		originVal, ok := d.GetOk(objectParamKey)
		if !ok {
			mErr = multierror.Append(mErr, fmt.Errorf("invalid param key: '%s'", objectParamKey))
		}

		// Store the origin value
		originParamKey := fmt.Sprintf("%s_origin", objectParamKey)
		// lintignore:R001
		if err := d.Set(originParamKey, originVal); err != nil {
			mErr = multierror.Append(mErr, fmt.Errorf("failed to set origin value for attribute storage '%s': %v",
				originParamKey, err))
		}
	}

	return mErr.ErrorOrNil()
}
