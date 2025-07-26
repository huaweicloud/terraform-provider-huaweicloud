package config

import (
	"context"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetRawConfigTags(rawConfig cty.Value) map[string]interface{} {
	tagsRaw := rawConfig.GetAttr("tags")

	if tagsRaw.IsNull() || !tagsRaw.IsKnown() || !tagsRaw.Type().IsMapType() {
		return nil
	}

	result := make(map[string]interface{})
	for k, v := range tagsRaw.AsValueMap() {
		result[k] = v.AsString()
	}
	return result
}

func MergeDefaultTags() schema.CustomizeDiffFunc {
	return func(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
		cfg := meta.(*Config)
		defaultTags := cfg.DefaultTags

		rawConfig := d.GetRawConfig()
		resourceTags := GetRawConfigTags(rawConfig)

		mergedTags := defaultTags
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
