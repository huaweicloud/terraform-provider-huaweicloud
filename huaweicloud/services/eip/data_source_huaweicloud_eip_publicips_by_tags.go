package eip

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EIP POST /v2.0/{project_id}/publicips/resource_instances/action
func DataSourcePublicipsByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePublicipsByTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"matches": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// In the API documentation, it is of type `Object`,
						// but here it has been changed to type `string`.
						"resource_detail": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildPublicipsByTagsBodyParams(d *schema.ResourceData) map[string]interface{} {
	action := d.Get("action").(string)
	bodyParams := map[string]interface{}{
		"action":  action,
		"tags":    buildPublicipsFilterTagsBodyParams(d.Get("tags").([]interface{})),
		"matches": buildPublicipsFilterMatchesBodyParams(d.Get("matches").([]interface{})),
	}

	return bodyParams
}

func buildPublicipsFilterTagsBodyParams(tags []interface{}) []map[string]interface{} {
	if len(tags) == 0 {
		return nil
	}

	res := make([]map[string]interface{}, len(tags))
	for i, tagRaw := range tags {
		if tag, ok := tagRaw.(map[string]interface{}); ok {
			res[i] = map[string]interface{}{
				"key": utils.PathSearch("key", tag, nil),
				"values": utils.ExpandToStringList(
					utils.PathSearch("values", tag, make([]interface{}, 0)).([]interface{})),
			}
		}
	}

	return res
}

func buildPublicipsFilterMatchesBodyParams(matches []interface{}) []map[string]interface{} {
	if len(matches) == 0 {
		return nil
	}

	res := make([]map[string]interface{}, len(matches))
	for i, matchRaw := range matches {
		if match, ok := matchRaw.(map[string]interface{}); ok {
			res[i] = map[string]interface{}{
				"key":   utils.PathSearch("key", match, nil),
				"value": utils.PathSearch("value", match, nil),
			}
		}
	}

	return res
}

func dataSourcePublicipsByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		product    = "vpc"
		httpUrl    = "v2.0/{project_id}/publicips/resource_instances/action"
		tagsAction = d.Get("action").(string)
		offset     = 0
		result     = make([]interface{}, 0)
		totalCount float64
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VPC EIP client: %s", err)
	}

	reqBodyParams := buildPublicipsByTagsBodyParams(d)
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		if tagsAction == "filter" {
			reqBodyParams["offset"] = offset
		}

		requestOpt.JSONBody = utils.RemoveNil(reqBodyParams)
		resp, err := client.Request("POST", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving EIP public IPs by tags: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalCount = utils.PathSearch("total_count", respBody, float64(0)).(float64)
		if tagsAction == "count" {
			break
		}

		resourcesResp := utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{})
		if len(resourcesResp) == 0 {
			break
		}

		result = append(result, resourcesResp...)

		offset += len(resourcesResp)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("resources", flattenPublicipsByTags(result)),
		d.Set("total_count", totalCount),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPublicipsByTags(resourcesResp []interface{}) []interface{} {
	if len(resourcesResp) == 0 {
		return nil
	}

	rst := make([]interface{}, len(resourcesResp))
	for i, v := range resourcesResp {
		rst[i] = map[string]interface{}{
			"resource_detail": utils.JsonToString(utils.PathSearch("resource_detail", v, nil)),
			"resource_id":     utils.PathSearch("resource_id", v, nil),
			"resource_name":   utils.PathSearch("resource_name", v, nil),
			"tags": flattenPublicipsTagsResp(
				utils.PathSearch("tags", v, make([]interface{}, 0)).([]interface{})),
		}
	}

	return rst
}

func flattenPublicipsTagsResp(tagsResp []interface{}) []interface{} {
	if len(tagsResp) == 0 {
		return nil
	}

	rst := make([]interface{}, len(tagsResp))
	for i, tag := range tagsResp {
		rst[i] = map[string]interface{}{
			"key":   utils.PathSearch("key", tag, nil),
			"value": utils.PathSearch("value", tag, nil),
		}
	}

	return rst
}
