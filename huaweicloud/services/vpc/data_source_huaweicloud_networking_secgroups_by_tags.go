package vpc

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

// @API VPC POST /v2.0/{project_id}/security-groups/resource_instances/action
func DataSourceNetworkingSecGroupsByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingSecGroupsByTagsRead,

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
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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

func dataSourceNetworkingSecGroupsByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		listNetworkingSecGroupsByTagsHttpUrl = "v2.0/{project_id}/security-groups/resource_instances/action"
		listNetworkingSecGroupsByTagsProduct = "vpc"
	)
	client, err := cfg.NewServiceClient(listNetworkingSecGroupsByTagsProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	listNetworkingSecGroupsByTagsPath := client.Endpoint + listNetworkingSecGroupsByTagsHttpUrl
	listNetworkingSecGroupsByTagsPath = strings.ReplaceAll(listNetworkingSecGroupsByTagsPath, "{project_id}", client.ProjectID)

	listNetworkingSecGroupsByTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	bodyParams := buildNetworkingSecGroupsByTagsBodyParams(d)
	action := d.Get("action").(string)
	resources := make([]interface{}, 0)
	offset := 0
	totalCount := float64(0)
	for {
		if action == "filter" {
			bodyParams["offset"] = offset
		}

		listNetworkingSecGroupsByTagsOpt.JSONBody = utils.RemoveNil(bodyParams)
		listNetworkingSecGroupsByTagsResp, err := client.Request("POST", listNetworkingSecGroupsByTagsPath, &listNetworkingSecGroupsByTagsOpt)
		if err != nil {
			return diag.FromErr(err)
		}
		listNetworkingSecGroupsByTagsRespBody, err := utils.FlattenResponse(listNetworkingSecGroupsByTagsResp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalCount = utils.PathSearch("total_count", listNetworkingSecGroupsByTagsRespBody, float64(0)).(float64)
		if action == "count" {
			break
		}

		data := utils.PathSearch("resources", listNetworkingSecGroupsByTagsRespBody, make([]interface{}, 0)).([]interface{})
		if len(data) == 0 {
			break
		}

		resources = append(resources, data...)

		offset += len(data)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("resources", flattenNetworkingSecGroupsByTagsResponseBody(resources)),
		d.Set("total_count", totalCount),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildNetworkingSecGroupsTagsBodyParams(d *schema.ResourceData) []map[string]interface{} {
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

func buildNetworkingSecGroupsMatchesBodyParams(d *schema.ResourceData) []map[string]interface{} {
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

func buildNetworkingSecGroupsByTagsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"action":  d.Get("action"),
		"tags":    buildNetworkingSecGroupsTagsBodyParams(d),
		"matches": buildNetworkingSecGroupsMatchesBodyParams(d),
	}

	return bodyParams
}

func flattenNetworkingSecGroupsByTagsResponseBody(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	resources := make([]interface{}, len(resp))
	for i, v := range resp {
		resources[i] = map[string]interface{}{
			"resource_name":   utils.PathSearch("resource_name", v, nil),
			"resource_id":     utils.PathSearch("resource_id", v, nil),
			"resource_detail": utils.JsonToString(utils.PathSearch("resource_detail", v, nil)),
			"tags":            utils.FlattenTagsToMap(utils.PathSearch("tags", v, nil)),
		}
	}
	return resources
}
