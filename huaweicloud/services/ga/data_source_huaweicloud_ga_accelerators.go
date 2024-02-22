package ga

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GA GET /v1/accelerators
func DataSourceAccelerators() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAcceleratorsRead,
		Schema: map[string]*schema.Schema{
			"accelerator_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the accelerator.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the accelerator.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The current status of the accelerator.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the enterprise project to which the accelerator belongs.",
			},
			"accelerators": {
				Type:        schema.TypeList,
				Elem:        acceleratorsSchema(),
				Computed:    true,
				Description: "The list of the accelerators.",
			},
		},
	}
}

func acceleratorsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the accelerator.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the accelerator.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the accelerator.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the accelerator.",
			},
			"ip_sets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The IP information of the accelerator.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP type of the accelerator.",
						},
						"ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address of the accelerator.",
						},
						"area": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The acceleration zone of the accelerator.",
						},
					},
				},
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the flavor to which the accelerator belongs.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the enterprise project to which the accelerator belongs.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The key/value pairs to associate with the accelerator.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the accelerator.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the accelerator.",
			},
		},
	}
	return &sc
}

func dataSourceAcceleratorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// listAccelerators: Query the list of accelerators
	var (
		listAcceleratorsHttpUrl = "v1/accelerators"
		listAcceleratorsProduct = "ga"
	)
	listAcceleratorsClient, err := cfg.NewServiceClient(listAcceleratorsProduct, region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	listAcceleratorsPath := listAcceleratorsClient.Endpoint + listAcceleratorsHttpUrl

	listAcceleratorsqueryParams := buildListAcceleratorsQueryParams(d, cfg)
	listAcceleratorsPath += listAcceleratorsqueryParams

	listAcceleratorsResp, err := pagination.ListAllItems(
		listAcceleratorsClient,
		"marker",
		listAcceleratorsPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving accelerators")
	}

	listAcceleratorsRespJson, err := json.Marshal(listAcceleratorsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var listAcceleratorsRespBody interface{}
	err = json.Unmarshal(listAcceleratorsRespJson, &listAcceleratorsRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	var mErr *multierror.Error
	mErr = multierror.Append(
		mErr,
		d.Set("accelerators", flattenListAcceleratorsResponseBody(listAcceleratorsRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListAcceleratorsResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("accelerators", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                    utils.PathSearch("id", v, nil),
			"name":                  utils.PathSearch("name", v, nil),
			"status":                utils.PathSearch("status", v, nil),
			"description":           utils.PathSearch("description", v, nil),
			"ip_sets":               flattenIpSets(utils.PathSearch("ip_sets", v, make([]interface{}, 0))),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
			"tags":                  utils.FlattenTagsToMap(utils.PathSearch("tags", v, nil)),
			"flavor_id":             utils.PathSearch("flavor_id", v, nil),
			"created_at":            utils.PathSearch("created_at", v, nil),
			"updated_at":            utils.PathSearch("updated_at", v, nil),
		})
	}
	return rst
}

func flattenIpSets(raw interface{}) []map[string]interface{} {
	curArray := raw.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(curArray))
	for i, ipSets := range curArray {
		result[i] = map[string]interface{}{
			"ip_type":    utils.PathSearch("ip_type", ipSets, nil),
			"ip_address": utils.PathSearch("ip_address", ipSets, nil),
			"area":       utils.PathSearch("area", ipSets, nil),
		}
	}
	return result
}

func buildListAcceleratorsQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	res := ""
	epsID := cfg.GetEnterpriseProjectID(d)

	if v, ok := d.GetOk("accelerator_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if epsID != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%s", res, epsID)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
