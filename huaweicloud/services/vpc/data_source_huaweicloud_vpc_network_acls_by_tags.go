package vpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPC POST /v3/{project_id}/firewalls/resource-instances/filter
func DataSourceNetworkAclsByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkAclsByTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
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

func dataSourceNetworkAclsByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		listNetworkAclsByTagsHttpUrl = "v3/{project_id}/firewalls/resource-instances/filter"
		listNetworkAclsByTagsProduct = "vpc"
	)
	client, err := cfg.NewServiceClient(listNetworkAclsByTagsProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	listNetworkAclsByTagsPath := client.Endpoint + listNetworkAclsByTagsHttpUrl
	listNetworkAclsByTagsPath = strings.ReplaceAll(listNetworkAclsByTagsPath, "{project_id}", client.ProjectID)

	listNetworkAclsByTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	bodyParams := buildNetworkAclsByTagsBodyParams(d)
	resources := make([]interface{}, 0)
	offset := 0
	totalCount := float64(0)
	for {
		listNetworkAclsByTagsReqPath := listNetworkAclsByTagsPath + fmt.Sprintf("?offset=%v", offset)
		listNetworkAclsByTagsOpt.JSONBody = utils.RemoveNil(bodyParams)
		listNetworkAclsByTagsResp, err := client.Request("POST", listNetworkAclsByTagsReqPath, &listNetworkAclsByTagsOpt)
		if err != nil {
			return diag.FromErr(err)
		}
		listNetworkAclsByTagsRespBody, err := utils.FlattenResponse(listNetworkAclsByTagsResp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalCount = utils.PathSearch("total_count", listNetworkAclsByTagsRespBody, float64(0)).(float64)
		data := utils.PathSearch("resources", listNetworkAclsByTagsRespBody, make([]interface{}, 0)).([]interface{})
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
		d.Set("resources", flattenNetworkAclsByTagsResponseBody(resources)),
		d.Set("total_count", totalCount),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildAclsByTagsTagsBodyParams(d *schema.ResourceData) []map[string]interface{} {
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

func buildAclsByTagsMatchesBodyParams(d *schema.ResourceData) []map[string]interface{} {
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

func buildNetworkAclsByTagsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"action":  d.Get("action"),
		"tags":    buildAclsByTagsTagsBodyParams(d),
		"matches": buildAclsByTagsMatchesBodyParams(d),
	}

	return bodyParams
}

func flattenNetworkAclsByTagsResponseBody(resp []interface{}) []interface{} {
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
