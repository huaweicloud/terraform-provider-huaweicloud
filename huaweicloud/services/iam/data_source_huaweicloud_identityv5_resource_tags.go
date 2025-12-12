package iam

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM GET /v5/{resource_type}/{resource_id}/tags
func DataSourceIdentityv5ResourceTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIAMIdentityv5TagsRead,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Required: true,
			},

			"tags": common.TagsComputedSchema(),
		},
	}
}

func dataSourceIAMIdentityv5TagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM Client: %s", err)
	}

	resourceType := d.Get("resource_type").(string)
	resourceId := d.Get("resource_id").(string)
	getResourceTagHttpUrl := "v5/{resource_type}/{resource_id}/tags"
	getResourceTagPath := iamClient.Endpoint + getResourceTagHttpUrl
	getResourceTagPath = strings.ReplaceAll(getResourceTagPath, "{resource_type}", resourceType)
	getResourceTagPath = strings.ReplaceAll(getResourceTagPath, "{resource_id}", resourceId)
	getResourceTagOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	getResourceTagResp, err := iamClient.Request("GET", getResourceTagPath, &getResourceTagOpt)
	if err != nil {
		return diag.Errorf("error retrieving IAM resource tag: %s", err)
	}

	getResourceTagRespBody, err := utils.FlattenResponse(getResourceTagResp)
	if err != nil {
		return diag.FromErr(err)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)
	mErr := multierror.Append(nil,
		d.Set("tags", flattenTagsToMap(utils.PathSearch("tags", getResourceTagRespBody, nil))),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting tags fields: %s", err)
	}
	return nil
}
