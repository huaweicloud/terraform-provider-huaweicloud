package coc

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API COC GET /v1/group-resource-relations
func DataSourceCocGroupResourceRelations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCocGroupResourceRelationsRead,

		Schema: map[string]*schema.Schema{
			"cloud_service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vendor": {
				Type:     schema.TypeString,
				Required: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"component_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_id_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"az_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"agent_state": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"charging_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavor_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"is_collected": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cmdb_resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_service_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ep_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ep_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_id": {
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
						"agent_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"agent_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"inner_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"properties": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ingest_properties": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operable": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCocGroupResourceRelationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("coc", region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	groupResourceRelations, err := queryGroupResourceRelations(client, d)
	if err != nil {
		return diag.Errorf("error querying group resource relations: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		nil,
		d.Set("data", flattenCocGetGroupResourceRelations(groupResourceRelations)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func queryGroupResourceRelations(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/group-resource-relations"
		offset  = 0
		result  = make([]interface{}, 0)
	)
	listPath := client.Endpoint + httpUrl
	listPath += buildGetGroupResourceRelationsRequiredParams(d)
	listPath += buildGetGroupResourceRelationsParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		groupResourceRelations := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
		if len(groupResourceRelations) < 1 {
			break
		}
		result = append(result, groupResourceRelations...)
		offset += len(groupResourceRelations)
	}

	return result, nil
}

func buildGetGroupResourceRelationsRequiredParams(d *schema.ResourceData) string {
	res := "?limit=100"
	if v, ok := d.GetOk("cloud_service_name"); ok {
		res = fmt.Sprintf("%s&provider=%v", res, v)
	}
	if v, ok := d.GetOk("vendor"); ok {
		res = fmt.Sprintf("%s&vendor=%v", res, v)
	}
	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, v)
	}
	return res
}

func buildGetGroupResourceRelationsParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("application_id"); ok {
		res = fmt.Sprintf("%s&application_id=%v", res, v)
	}
	if v, ok := d.GetOk("component_id"); ok {
		res = fmt.Sprintf("%s&component_id=%v", res, v)
	}
	if v, ok := d.GetOk("group_id"); ok {
		res = fmt.Sprintf("%s&group_id=%v", res, v)
	}
	if v, ok := d.GetOk("resource_id_list"); ok {
		res += buildQueryStringParams("resource_id_list", v.(*schema.Set).List())
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("region_id"); ok {
		res = fmt.Sprintf("%s&region_id=%v", res, v)
	}
	if v, ok := d.GetOk("az_id"); ok {
		res = fmt.Sprintf("%s&az_id=%v", res, v)
	}
	if v, ok := d.GetOk("ip_type"); ok {
		res = fmt.Sprintf("%s&ip_type=%v", res, v)
	}
	if v, ok := d.GetOk("ip"); ok {
		res = fmt.Sprintf("%s&ip=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if v, ok := d.GetOk("agent_state"); ok {
		res = fmt.Sprintf("%s&agent_state=%v", res, v)
	}
	if v, ok := d.GetOk("image_name"); ok {
		res = fmt.Sprintf("%s&image_name=%v", res, v)
	}
	if v, ok := d.GetOk("os_type"); ok {
		res = fmt.Sprintf("%s&os_type=%v", res, v)
	}
	if v, ok := d.GetOk("tag"); ok {
		res = fmt.Sprintf("%s&tag=%v", res, v)
	}
	if v, ok := d.GetOk("charging_mode"); ok {
		res = fmt.Sprintf("%s&charging_mode=%v", res, v)
	}
	if v, ok := d.GetOk("flavor_name"); ok {
		res = fmt.Sprintf("%s&flavor_name=%v", res, v)
	}
	if v, ok := d.GetOk("ip_list"); ok {
		res += buildQueryStringParams("ip_list", v.(*schema.Set).List())
	}
	if v, ok := d.GetOk("is_collected"); ok {
		res = fmt.Sprintf("%s&is_collected=%v", res, v)
	}

	return res
}

func flattenCocGetGroupResourceRelations(groupResourceRelations []interface{}) []interface{} {
	result := make([]interface{}, 0, len(groupResourceRelations))

	for _, groupResourceRelation := range groupResourceRelations {
		result = append(result, map[string]interface{}{
			"id":                 utils.PathSearch("id", groupResourceRelation, nil),
			"cmdb_resource_id":   utils.PathSearch("cmdb_resource_id", groupResourceRelation, nil),
			"group_id":           utils.PathSearch("group_id", groupResourceRelation, nil),
			"group_name":         utils.PathSearch("group_name", groupResourceRelation, nil),
			"resource_id":        utils.PathSearch("resource_id", groupResourceRelation, nil),
			"name":               utils.PathSearch("name", groupResourceRelation, nil),
			"cloud_service_name": utils.PathSearch("provider", groupResourceRelation, nil),
			"type":               utils.PathSearch("type", groupResourceRelation, nil),
			"region_id":          utils.PathSearch("region_id", groupResourceRelation, nil),
			"ep_id":              utils.PathSearch("ep_id", groupResourceRelation, nil),
			"ep_name":            utils.PathSearch("ep_name", groupResourceRelation, nil),
			"project_id":         utils.PathSearch("project_id", groupResourceRelation, nil),
			"domain_id":          utils.PathSearch("domain_id", groupResourceRelation, nil),
			"tags": flattenCocGetGroupResourceRelationsTags(
				utils.PathSearch("tags", groupResourceRelation, nil)),
			"agent_id":          utils.PathSearch("agent_id", groupResourceRelation, nil),
			"agent_state":       utils.PathSearch("agent_state", groupResourceRelation, nil),
			"inner_ip":          utils.PathSearch("inner_ip", groupResourceRelation, nil),
			"properties":        utils.JsonToString(utils.PathSearch("properties", groupResourceRelation, nil)),
			"ingest_properties": utils.JsonToString(utils.PathSearch("ingest_properties", groupResourceRelation, nil)),
			"operable":          utils.PathSearch("operable", groupResourceRelation, nil),
			"create_time":       utils.PathSearch("create_time", groupResourceRelation, nil),
		})
	}

	return result
}

func flattenCocGetGroupResourceRelationsTags(rawParams interface{}) []interface{} {
	if paramsList, ok := rawParams.([]interface{}); ok {
		if len(paramsList) == 0 {
			return nil
		}
		rst := make([]interface{}, 0, len(paramsList))
		for _, params := range paramsList {
			raw := params.(map[string]interface{})
			m := map[string]interface{}{
				"key":   utils.PathSearch("key", raw, nil),
				"value": utils.PathSearch("value", raw, nil),
			}
			rst = append(rst, m)
		}

		return rst
	}

	return nil
}
