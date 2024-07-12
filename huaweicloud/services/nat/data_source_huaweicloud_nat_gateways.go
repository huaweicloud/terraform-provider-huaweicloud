// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product NAT
// ---------------------------------------------------------------

package nat

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

// @API NAT GET /v2/{project_id}/nat_gateways
func DataSourcePublicGateways() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePublicGatewaysRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the NAT gateways are located.",
			},
			"gateway_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the NAT gateway.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the NAT gateway.",
			},
			"spec": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The specification of the NAT gateways.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The current status of the NAT gateways.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the VPC to which the NAT gateways belong.",
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The network ID of the downstream interface (the next hop of the DVR)" +
					"of the NAT gateways.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the enterprise project to which the NAT gateways belong.",
			},
			"gateways": {
				Type:        schema.TypeList,
				Elem:        gatewayPublicGatewaysSchema(),
				Computed:    true,
				Description: "The list of the NAT gateway.",
			},
		},
	}
}

func gatewayPublicGatewaysSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the NAT gateway.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the NAT gateway.",
			},
			"spec": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The specification of the NAT gateway.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the NAT gateway.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the NAT gateway.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the NAT gateway.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the VPC to which the NAT gateway belongs.",
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "The network ID of the downstream interface (the next hop of the DVR)" +
					"of the NAT gateway.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the enterprise project to which the NAT gateway belongs.",
			},
		},
	}
	return &sc
}

func dataSourcePublicGatewaysRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// listGateways: Query the list of NAT gateways
	var (
		listGatewaysHttpUrl = "v2/{project_id}/nat_gateways"
		listGatewaysProduct = "nat"
	)
	listGatewaysClient, err := cfg.NewServiceClient(listGatewaysProduct, region)
	if err != nil {
		return diag.Errorf("error creating NAT client: %s", err)
	}

	listGatewaysPath := listGatewaysClient.Endpoint + listGatewaysHttpUrl
	listGatewaysPath = strings.ReplaceAll(listGatewaysPath, "{project_id}", listGatewaysClient.ProjectID)

	listGatewaysqueryParams := buildListPublicGatewaysQueryParams(d, cfg)
	listGatewaysPath += listGatewaysqueryParams

	listGatewaysOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	listGatewaysResp, err := listGatewaysClient.Request("GET", listGatewaysPath, &listGatewaysOpt)

	if err != nil {
		return diag.Errorf("error retrieving NAT gateways %s", err)
	}

	listGatewaysRespBody, err := utils.FlattenResponse(listGatewaysResp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("gateways", flattenListGatewaysResponseBodyPublicGateways(listGatewaysRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListGatewaysResponseBodyPublicGateways(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("nat_gateways", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                    utils.PathSearch("id", v, nil),
			"name":                  utils.PathSearch("name", v, nil),
			"spec":                  utils.PathSearch("spec", v, nil),
			"status":                utils.PathSearch("status", v, nil),
			"description":           utils.PathSearch("description", v, nil),
			"created_at":            utils.PathSearch("created_at", v, nil),
			"vpc_id":                utils.PathSearch("router_id", v, nil),
			"subnet_id":             utils.PathSearch("internal_network_id", v, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
		})
	}
	return rst
}

func buildListPublicGatewaysQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	res := ""
	enterpriseProjectID := cfg.GetEnterpriseProjectID(d)

	if v, ok := d.GetOk("gateway_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("spec"); ok {
		res = fmt.Sprintf("%s&spec=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		res = fmt.Sprintf("%s&router_id=%v", res, v)
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		res = fmt.Sprintf("%s&internal_network_id=%v", res, v)
	}
	if enterpriseProjectID != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%s", res, enterpriseProjectID)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
