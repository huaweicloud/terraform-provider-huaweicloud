package dew

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

// @API DEW POST /v1/{project_id}/csms/{resource_instances}/action
func DataSourceCSMSSecretsByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCSMSSecretsByTagsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the resource.",
			},
			"resource_instances": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the resource instances.",
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the operation type.",
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the tag key.",
						},
						"values": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Specifies the set of tag values.",
						},
					},
				},
				Description: "Specifies the list of tags.",
			},
			"sequence": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the `36` byte sequence number of a request message.",
			},
			"matches": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the search field.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the field for fuzzy match.",
						},
					},
				},
				Description: "Specifies the key-value pair to be matched.",
			},
			"resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        resourceSchema(),
				Description: "The list of the filtered secrets.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of the filtered secrets.",
			},
		},
	}
}

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secret ID.",
			},
			"resource_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secret name.",
			},
			"resource_detail": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the secret.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The secret name.",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The secret status.",
						},
						"kms_key_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of KMS key used to encrypt secret.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the secret.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The creation time of the secret, the value is a timestamp.",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The update time of the secret, the value is a timestamp.",
						},
						"scheduled_delete_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The time of the secret to be scheduled deleted, the value is a timestamp.",
						},
						"secret_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The secret type.",
						},
						"auto_rotation": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Automatic rotation.",
						},
						"rotation_period": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The secret rotation period.",
						},
						"rotation_config": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The secret rotation config.",
						},
						"rotation_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The rotation time of the secret",
						},
						"next_rotation_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The next rotation time of the secret.",
						},
						"event_subscriptions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The list of events subscribed to by secret.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The enterprise project ID.",
						},
					},
				},
				Description: "The secret detail.",
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The tag name.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The tag value.",
						},
					},
				},
				Description: "The tag list.",
			},
			"sys_tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The system tag key.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The system tag value.",
						},
					},
				},
				Description: "The system tag list.",
			},
		},
	}
}

func dataSourceCSMSSecretsByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "kms"
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	requestPath := client.Endpoint + "v1/{project_id}/csms/{resource_instances}/action"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{resource_instances}", d.Get("resource_instances").(string))

	allResources := make([]interface{}, 0)
	allCount := 0
	offset := 0

	listCSMSOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	bodyParams := buildCSMSParams(d)
	for {
		if d.Get("action").(string) == "filter" {
			bodyParams["limit"] = 1
			bodyParams["offset"] = offset
		}
		listCSMSOpt.JSONBody = utils.RemoveNil(bodyParams)
		resp, err := client.Request("POST", requestPath, &listCSMSOpt)

		if err != nil {
			return diag.Errorf("error retrieving CSMS secrets by tags: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		resources := flattenCSMSResources(utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{}))

		// When the action is count, resources is null, but total_count is not null.
		if d.Get("action").(string) == "count" {
			totalCount := utils.PathSearch("total_count", respBody, float64(0)).(float64)
			allCount = int(totalCount)
			break
		}

		if len(resources) == 0 {
			break
		}
		allResources = append(allResources, resources...)
		allCount += len(resources)
		offset += len(resources)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("resources", allResources),
		d.Set("total_count", allCount),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildCSMSParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"action":   d.Get("action"),
		"tags":     buildTags(d),
		"matches":  buildMatches(d),
		"sequence": utils.ValueIgnoreEmpty(d.Get("sequence")),
	}
}

func buildTags(d *schema.ResourceData) []map[string]interface{} {
	v, ok := d.GetOk("tags")
	if !ok {
		return nil
	}

	bodyParams := make([]map[string]interface{}, len(v.([]interface{})))

	for i, tag := range v.([]interface{}) {
		bodyParams[i] = map[string]interface{}{
			"key":    utils.PathSearch("key", tag, nil),
			"values": utils.PathSearch("values", tag, nil),
		}
	}

	return bodyParams
}

func buildMatches(d *schema.ResourceData) []map[string]interface{} {
	v, ok := d.GetOk("matches")
	if !ok {
		return nil
	}

	bodyParams := make([]map[string]interface{}, len(v.([]interface{})))

	for i, match := range v.([]interface{}) {
		bodyParams[i] = map[string]interface{}{
			"key":   utils.PathSearch("key", match, nil),
			"value": utils.PathSearch("value", match, nil),
		}
	}

	return bodyParams
}

func flattenCSMSResources(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))

	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"resource_id":     utils.PathSearch("resource_id", v, nil),
			"resource_name":   utils.PathSearch("resource_name", v, nil),
			"resource_detail": flattenResourceDetail(utils.PathSearch("resource_detail", v, nil)),
			"tags":            flattenTags(utils.PathSearch("tags", v, make([]interface{}, 0)).([]interface{})),
			"sys_tags":        flattenSysTags(utils.PathSearch("sys_tags", v, make([]interface{}, 0)).([]interface{})),
		})
	}
	return rst
}

func flattenResourceDetail(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	resourceDetail := map[string]interface{}{
		"id":                    utils.PathSearch("id", resp, nil),
		"name":                  utils.PathSearch("name", resp, nil),
		"state":                 utils.PathSearch("state", resp, nil),
		"kms_key_id":            utils.PathSearch("kms_key_id", resp, nil),
		"description":           utils.PathSearch("description", resp, nil),
		"create_time":           utils.PathSearch("create_time", resp, nil),
		"update_time":           utils.PathSearch("update_time", resp, nil),
		"scheduled_delete_time": utils.PathSearch("scheduled_delete_time", resp, nil),
		"secret_type":           utils.PathSearch("secret_type", resp, nil),
		"auto_rotation":         utils.PathSearch("auto_rotation", resp, nil),
		"rotation_period":       utils.PathSearch("rotation_period", resp, nil),
		"rotation_config":       utils.PathSearch("rotation_config", resp, nil),
		"rotation_time":         utils.PathSearch("rotation_time", resp, nil),
		"next_rotation_time":    utils.PathSearch("next_rotation_time", resp, nil),
		"event_subscriptions":   utils.PathSearch("event_subscriptions", resp, nil),
		"enterprise_project_id": utils.PathSearch("enterprise_project_id", resp, nil),
	}
	return []interface{}{resourceDetail}
}

func flattenTags(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))

	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}
	return rst
}

func flattenSysTags(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))

	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}
	return rst
}
