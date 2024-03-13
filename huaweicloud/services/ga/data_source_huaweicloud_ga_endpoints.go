package ga

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GA GET /v1/endpoint-groups/{endpoint_group_id}/endpoints
func DataSourceEndpoints() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEndpointsRead,
		Schema: map[string]*schema.Schema{
			"endpoint_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the endpoint group to which the endpoint belongs.",
			},
			"endpoint_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the endpoint.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The status of the endpoint.",
			},
			"resource_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the backend resource corresponding to the endpoint.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of the backend resource corresponding to the endpoint.",
			},
			"health_state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The health status of the endpoint.",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IP address of the backend resource corresponding to the endpoint.",
			},
			"endpoints": {
				Type:        schema.TypeList,
				Elem:        endpointsSchema(),
				Computed:    true,
				Description: "The list of the endpoints.",
			},
		},
	}
}

func endpointsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the endpoint.",
			},
			"endpoint_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the endpoint group to which the endpoint belongs.",
			},
			"resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the backend resource corresponding to the endpoint.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the backend resource corresponding to the endpoint.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the endpoint.",
			},
			"weight": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The weight of traffic distribution to the endpoint.",
			},
			"health_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The health status of the endpoint.",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address of the backend resource corresponding to the endpoint.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the endpoint.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the endpoint.",
			},
		},
	}
	return &sc
}

func dataSourceEndpointsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// listEndpoints: Query the list of endpoints
	var (
		listEndpointsHttpUrl = "v1/endpoint-groups/{endpoint_group_id}/endpoints"
		listEndpointsProduct = "ga"
	)
	listEndpointsClient, err := cfg.NewServiceClient(listEndpointsProduct, region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	listEndpointsPath := listEndpointsClient.Endpoint + listEndpointsHttpUrl
	listEndpointsPath = strings.ReplaceAll(listEndpointsPath, "{endpoint_group_id}", d.Get("endpoint_group_id").(string))

	listEndpointsqueryParams := buildListEndpointsQueryParams(d)
	listEndpointsPath += listEndpointsqueryParams

	listEndpointsResp, err := pagination.ListAllItems(
		listEndpointsClient,
		"marker",
		listEndpointsPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving endpoints")
	}

	listEndpointsRespJson, err := json.Marshal(listEndpointsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var listEndpointsRespBody interface{}
	err = json.Unmarshal(listEndpointsRespJson, &listEndpointsRespBody)
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
		d.Set("endpoints", filterListEndpointsResponseBody(flattenListEndpointsResponseBody(listEndpointsRespBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListEndpointsResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("endpoints", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                utils.PathSearch("id", v, nil),
			"endpoint_group_id": utils.PathSearch("endpoint_group_id", v, nil),
			"resource_id":       utils.PathSearch("resource_id", v, nil),
			"resource_type":     utils.PathSearch("resource_type", v, nil),
			"status":            utils.PathSearch("status", v, nil),
			"weight":            utils.PathSearch("weight", v, nil),
			"health_state":      utils.PathSearch("health_state", v, nil),
			"ip_address":        utils.PathSearch("ip_address", v, nil),
			"created_at":        utils.PathSearch("created_at", v, nil),
			"updated_at":        utils.PathSearch("updated_at", v, nil),
		})
	}
	return rst
}

func filterListEndpointsResponseBody(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("resource_id"); ok &&
			fmt.Sprint(param) != utils.PathSearch("resource_id", v, nil) {
			continue
		}

		if param, ok := d.GetOk("resource_type"); ok &&
			fmt.Sprint(param) != utils.PathSearch("resource_type", v, nil) {
			continue
		}

		if param, ok := d.GetOk("health_state"); ok &&
			fmt.Sprint(param) != utils.PathSearch("health_state", v, nil) {
			continue
		}

		if param, ok := d.GetOk("ip_address"); ok &&
			fmt.Sprint(param) != utils.PathSearch("ip_address", v, nil) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func buildListEndpointsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("endpoint_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
