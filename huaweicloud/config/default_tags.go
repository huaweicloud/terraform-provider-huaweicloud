package config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func MergeDefaultTags() schema.CustomizeDiffFunc {
	return func(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
		var (
			cfg          = meta.(*Config)
			resourceTags = utils.TryMapValueAnalysis(utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "tags"))
			mergedTags   = make(map[string]interface{})
		)

		for k, v := range cfg.DefaultTags {
			mergedTags[k] = v
		}

		for k, v := range resourceTags {
			mergedTags[k] = v
		}

		err := d.SetNew("tags", mergedTags)
		if err != nil {
			return err
		}

		return nil
	}
}
