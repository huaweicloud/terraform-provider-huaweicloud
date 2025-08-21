package cts

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CTS GET /v3/{project_id}/{resource_type}/{resource_id}/tags
func DataSourceCtsResourcesTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCtsResourcesTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to query the resource tags.`,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resource type to be queried.`,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resource ID to be queried.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of tags that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The tag key.`,
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The tag value.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCtsResourcesTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cts", region)
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	resourceId := d.Get("resource_id").(string)
	httpUrl := "v3/{project_id}/{resource_type}/{resource_id}/tags"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{resource_type}", d.Get("resource_type").(string))
	path = strings.ReplaceAll(path, "{resource_id}", resourceId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	response, err := client.Request("GET", path, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CTS resource tags: %s", err)
	}

	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("tags", flattenResourceTags(utils.PathSearch("tags", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenResourceTags(tags []interface{}) []interface{} {
	if len(tags) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(tags))
	for _, tag := range tags {
		result = append(result, map[string]interface{}{
			"key":   utils.PathSearch("key", tag, nil),
			"value": utils.PathSearch("value", tag, nil),
		})
	}

	return result
}
