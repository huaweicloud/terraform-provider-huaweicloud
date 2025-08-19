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

// @API FunctionGraph GET /v2/{project_id}/{resource_type}/{resource_id}/tags
func DataSourceFunctionTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFunctionTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the function is located and tags to be queried.`,
			},

			// Required parameter(s).
			"function_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the function to which the tags belong.`,
			},

			// Attribute(s).
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The key of the tag.`,
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The value of the tag.`,
						},
					},
				},
				Description: `The list of the function tags.`,
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
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The value of the system tag.`,
						},
					},
				},
				Description: utils.SchemaDesc(
					`The list of the system tags.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func getFunctionTags(client *golangsdk.ServiceClient, functionId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/functions/{function_id}/tags"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{function_id}", functionId)

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

func flattenFunctionTags(tags []interface{}) []map[string]interface{} {
	if len(tags) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(tags))
	for _, tag := range tags {
		result = append(result, map[string]interface{}{
			"key":   utils.PathSearch("key", tag, nil),
			"value": utils.PathSearch("value", tag, nil),
		})
	}

	return result
}

func dataSourceFunctionTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	functionId := d.Get("function_id").(string)
	resp, err := getFunctionTags(client, functionId)
	if err != nil {
		return diag.Errorf("error querying FunctionGraph function tags: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("tags", flattenFunctionTags(utils.PathSearch("tags", resp, make([]interface{}, 0)).([]interface{}))),
		d.Set("sys_tags", flattenFunctionTags(utils.PathSearch("sys_tags", resp, make([]interface{}, 0)).([]interface{}))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
