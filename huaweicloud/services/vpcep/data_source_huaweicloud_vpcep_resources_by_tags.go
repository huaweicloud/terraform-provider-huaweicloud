package vpcep

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPCEP POST /v1/{project_id}/{resource_type}/resource_instances/action
func DataSourceVpcepResourcesByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcepResourcesByTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
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
			"tags_any": {
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
			"not_tags": {
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
			"not_tags_any": {
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
			"resource_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the resource type.`,
			},
			"without_any_tag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies if the resource has no tags.`,
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
						"resource_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_detail": {
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

func flattenVpcepTags(tags []interface{}) []map[string]interface{} {
	rst := make([]map[string]interface{}, 0, len(tags))
	for _, tag := range tags {
		rst = append(rst, map[string]interface{}{
			"key":   utils.PathSearch("key", tag, nil),
			"value": utils.PathSearch("value", tag, nil),
		})
	}
	return rst
}

func flattenVpcepResourcesByTagsResponseBody(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	resources := make([]interface{}, len(resp))
	for i, v := range resp {
		resources[i] = map[string]interface{}{
			"resource_name": utils.PathSearch("resource_name", v, nil),
			"resource_id":   utils.PathSearch("resource_id", v, nil),
			"tags":          flattenVpcepTags(utils.PathSearch("tags", v, make([]interface{}, 0)).([]interface{})),
		}
	}
	return resources
}

func expandTags(rawTags []interface{}) []map[string]interface{} {
	tags := make([]map[string]interface{}, len(rawTags))
	for i, raw := range rawTags {
		tag := raw.(map[string]interface{})
		tags[i] = map[string]interface{}{
			"key":    tag["key"].(string),
			"values": tag["values"].([]interface{}),
		}
	}
	return tags
}

func expandMatches(rawTags []interface{}) []map[string]interface{} {
	tags := make([]map[string]interface{}, len(rawTags))
	for i, raw := range rawTags {
		tag := raw.(map[string]interface{})
		tags[i] = map[string]interface{}{
			"key":   tag["key"].(string),
			"value": tag["value"].(string),
		}
	}
	return tags
}

func buildVpcepResourcesByTagsBodyParams(d *schema.ResourceData) map[string]interface{} {
	body := map[string]interface{}{}
	if v, ok := d.GetOk("action"); ok {
		body["action"] = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		body["tags"] = expandTags(v.([]interface{}))
	}
	if v, ok := d.GetOk("tags_any"); ok {
		body["tags_any"] = expandTags(v.([]interface{}))
	}
	if v, ok := d.GetOk("not_tags"); ok {
		body["not_tags"] = expandTags(v.([]interface{}))
	}
	if v, ok := d.GetOk("not_tags_any"); ok {
		body["not_tags_any"] = expandTags(v.([]interface{}))
	}
	if v, ok := d.GetOk("matches"); ok {
		body["matches"] = expandMatches(v.([]interface{}))
	}
	return body
}

func dataSourceVpcepResourcesByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v1/{project_id}/{resource_type}/resource_instances/action"
	)

	vpcepClient, err := cfg.VPCEPClient(region)
	if err != nil {
		return diag.Errorf("error creating VPC endpoint client: %s", err)
	}

	requestPath := vpcepClient.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", vpcepClient.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{resource_type}", d.Get("resource_type").(string))

	var allVpcepResources []interface{}
	offset := 0
	totalCount := 0

	for {
		requestBody := buildVpcepResourcesByTagsBodyParams(d)
		if requestBody["action"] == "filter" {
			requestBody["limit"] = "1000"
			requestBody["offset"] = strconv.Itoa(offset)
		}

		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         requestBody,
		}
		resp, err := vpcepClient.Request("POST", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error querying Vpcep resources by tags: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		vpcepResources := utils.PathSearch("resources", respBody, []interface{}{}).([]interface{})
		totalCount = int(utils.PathSearch("total_count", respBody, float64(0)).(float64))
		if len(vpcepResources) == 0 {
			break
		}
		allVpcepResources = append(allVpcepResources, vpcepResources...)
		offset += len(vpcepResources)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("resources", flattenVpcepResourcesByTagsResponseBody(allVpcepResources)),
		d.Set("total_count", totalCount),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
