package iam

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IAM GET /v5/{resource_type}/{resource_id}/tags
func DataSourceV5ResourceTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV5TagsRead,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resource ID to be queried.`,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resource type to be queried.`,
			},
			"tags": common.TagsComputedSchema(`The key/value pairs associated with the resource.`),
		},
	}
}

func dataSourceV5TagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM Client: %s", err)
	}

	resourceType := d.Get("resource_type").(string)
	resourceId := d.Get("resource_id").(string)
	tags, err := getV5ResourceTags(iamClient, resourceType, resourceId)
	if err != nil {
		return diag.Errorf("error retrieving tags of %s (%s): %s", resourceType, resourceId, err)
	}

	randomId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomId)

	mErr := multierror.Append(
		d.Set("tags", flattenTagsToMap(tags)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting tags fields: %s", err)
	}

	return nil
}
