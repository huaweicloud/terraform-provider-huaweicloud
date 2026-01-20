package dds

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

// @API DDS POST /v3/{project_id}/instances/action
func DataSourceInstancesByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstancesByTagsRead,

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
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
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

func buildInstancesTagsBodyParams(d *schema.ResourceData, action string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"action":  action,
		"matches": buildFilterMatchesBodyParams(d.Get("matches").([]interface{})),
		"tags":    buildFilterTagsBodyParams(d.Get("tags").([]interface{})),
	}

	return bodyParams
}

func buildFilterMatchesBodyParams(matches []interface{}) []map[string]interface{} {
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

func buildFilterTagsBodyParams(tags []interface{}) []map[string]interface{} {
	if len(tags) == 0 {
		return nil
	}

	res := make([]map[string]interface{}, len(tags))
	for i, tagRaw := range tags {
		if tag, ok := tagRaw.(map[string]interface{}); ok {
			res[i] = map[string]interface{}{
				"key":    utils.PathSearch("key", tag, nil),
				"values": utils.PathSearch("values", tag, nil),
			}
		}
	}

	return res
}

func dataSourceInstancesByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v3/{project_id}/instances/action"
		tagsAction = d.Get("action").(string)
		offset     = 0
		result     = make([]interface{}, 0)
		totalCount float64
	)

	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	reqBodyParams := buildInstancesTagsBodyParams(d, tagsAction)
	if tagsAction == "filter" {
		reqBodyParams["limit"] = 100
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		if tagsAction == "filter" {
			reqBodyParams["offset"] = offset
		}

		listOpt.JSONBody = utils.RemoveNil(reqBodyParams)
		resp, err := client.Request("POST", listPath, &listOpt)
		if err != nil {
			return diag.Errorf("error retrieving instances: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalCount = utils.PathSearch("total_count", respBody, float64(0)).(float64)
		if tagsAction == "count" {
			break
		}

		instances := utils.PathSearch("instances", respBody, make([]interface{}, 0)).([]interface{})
		if len(instances) == 0 {
			break
		}

		result = append(result, instances...)

		offset += len(instances)
	}

	datasourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(datasourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instances", flattenInstancesByTags(result)),
		d.Set("total_count", totalCount),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenInstancesByTags(instances []interface{}) []interface{} {
	if len(instances) == 0 {
		return nil
	}

	rst := make([]interface{}, len(instances))
	for i, v := range instances {
		rst[i] = map[string]interface{}{
			"instance_id":   utils.PathSearch("instance_id", v, nil),
			"instance_name": utils.PathSearch("instance_name", v, nil),
			"tags":          flattenTagsResp(utils.PathSearch("tags", v, make([]interface{}, 0)).([]interface{})),
		}
	}
	return rst
}

func flattenTagsResp(tags []interface{}) []interface{} {
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
