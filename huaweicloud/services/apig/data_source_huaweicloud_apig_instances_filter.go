package apig

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

// @API APIG POST /v2/{project_id}/apigw/resource-instances/filter
func DataSourceInstancesFilter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstancesFilterRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"without_any_tag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to query resources without tags. Defaults to **false**.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `The list of the tags to be queried.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The key of the tag.`,
						},
						"values": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of values of the tag.`,
						},
					},
				},
			},
			"matches": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `The fields to be queried.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The key to be matched.`,
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The value of the matching field.`,
						},
					},
				},
			},
			"instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All dedicated instances that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the instance.`,
						},
						"resource_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the instance.`,
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The tag list associated with the instance.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The key of the instance tag.`,
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The value of the instance tag.`,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildGetInstancesBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"without_any_tag": d.Get("without_any_tag"),
		"tags":            buildGetInstanceTagsBodyParams(d.Get("tags").([]interface{})),
		"matches":         buildGetInstanceMatchesBodyParams(d.Get("matches").([]interface{})),
	}
	return bodyParams
}

func buildGetInstanceTagsBodyParams(tags []interface{}) []map[string]interface{} {
	if len(tags) == 0 {
		return nil
	}

	res := make([]map[string]interface{}, len(tags))
	for i, tag := range tags {
		res[i] = map[string]interface{}{
			// For interface, `key` and `values` is required, so they cannot be ignored.
			"key":    utils.PathSearch("key", tag, nil),
			"values": utils.PathSearch("values", tag, nil),
		}
	}
	return res
}

func buildGetInstanceMatchesBodyParams(matches []interface{}) []map[string]interface{} {
	if len(matches) == 0 {
		return nil
	}

	res := make([]map[string]interface{}, len(matches))
	for i, v := range matches {
		// For interface, `key` and `value` is required, so they cannot be ignored.
		res[i] = map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
		}
	}
	return res
}

func dataSourceInstancesFilterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/apigw/resource-instances/filter"
	)
	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildGetInstancesBodyParams(d)),
	}

	// Paging is not effective, the default value of limit is 1000.
	resp, err := client.Request("POST", listPath, &listOpt)
	if err != nil {
		return diag.Errorf("error retrieving APIG instances: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instances", flattenInstancesResp(
			utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenInstancesResp(instances []interface{}) []interface{} {
	if len(instances) == 0 {
		return nil
	}

	rst := make([]interface{}, len(instances))
	for i, v := range instances {
		rst[i] = map[string]interface{}{
			"resource_id":   utils.PathSearch("resource_id", v, nil),
			"resource_name": utils.PathSearch("resource_name", v, nil),
			"tags":          flattenResourceInstanceTagsResp(utils.PathSearch("tags", v, make([]interface{}, 0)).([]interface{})),
		}
	}
	return rst
}

func flattenResourceInstanceTagsResp(tags []interface{}) []interface{} {
	if len(tags) == 0 {
		return nil
	}

	rst := make([]interface{}, len(tags))
	for i, tag := range tags {
		rst[i] = map[string]interface{}{
			"key":   utils.PathSearch("key", tag, nil),
			"value": utils.PathSearch("value", tag, nil),
		}
	}
	return rst
}
