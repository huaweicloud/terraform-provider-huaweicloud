package fgs

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

// @API FunctionGraph GET /v2/{project_id}/{resource_type}/tags
func DataSourceResourceTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the resource tags are located.`,
			},

			// Attribute(s).
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of resource tags.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The key of the tag.`,
						},
						"values": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The values of the tag.`,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			// Internal attribute(s).
			"sys_tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The key of the system tag.`,
						},
						"values": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The values of the system tag.`,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
				Description: utils.SchemaDesc(
					`The list of system tags.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func getResourceTags(client *golangsdk.ServiceClient) (interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/functions/tags"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func flattenResourceTags(tags []interface{}) []map[string]interface{} {
	if len(tags) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(tags))
	for _, tag := range tags {
		result = append(result, map[string]interface{}{
			"key":    utils.PathSearch("key", tag, nil),
			"values": utils.PathSearch("values", tag, make([]interface{}, 0)),
		})
	}

	return result
}

func dataSourceResourceTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	respBody, err := getResourceTags(client)
	if err != nil {
		return diag.Errorf("error querying resource tags: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("tags", flattenResourceTags(utils.PathSearch("tags", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("sys_tags", flattenResourceTags(utils.PathSearch("sys_tags", respBody, make([]interface{}, 0)).([]interface{}))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
