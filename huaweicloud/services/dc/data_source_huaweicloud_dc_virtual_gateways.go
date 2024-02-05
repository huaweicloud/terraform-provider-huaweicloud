package dc

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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DC GET /v3/{project_id}/dcaas/virtual-gateways
func DataSourceDCVirtualGateways() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDCVirtualGatewaysRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"virtual_gateway_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the ID of the virtual gateway.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the name of the virtual gateway.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the ID of the VPC connected to the virtual gateway.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the enterprise project ID.",
			},
			"virtual_gateways": {
				Type:        schema.TypeList,
				Elem:        virtualGatewaySchema(),
				Computed:    true,
				Description: `Indicates the virtual gateway list.`,
			},
		},
	}
}

func virtualGatewaySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the virtual gateway ID.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the ID of the VPC connected by the virtual gateway.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the virtual gateway name.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the virtual gateway type.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the virtual gateway status..",
			},
			"asn": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Indicates the local BGP ASN of the virtual gateway.",
			},
			"local_ep_group": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Indicates the IPv4 subnets connected by the virtual gateway.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the virtual gateway description.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the ID of the enterprise project that the virtual gateway belongs to.",
			},
		},
	}
	return &sc
}

func resourceDCVirtualGatewaysRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// listVirtualGateways: Query the List of DC virtual gateways.
	var (
		listVirtualGatewaysHttpUrl = "v3/{project_id}/dcaas/virtual-gateways"
		listVirtualGatewaysProduct = "dc"
	)
	listVirtualGatewaysClient, err := cfg.NewServiceClient(listVirtualGatewaysProduct, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	listVirtualGatewaysPath := listVirtualGatewaysClient.Endpoint + listVirtualGatewaysHttpUrl
	listVirtualGatewaysPath = strings.ReplaceAll(listVirtualGatewaysPath, "{project_id}", listVirtualGatewaysClient.ProjectID)

	listVirtualGatewaysQueryParams := buildListVirtualGatewaysQueryParams(d, cfg.GetEnterpriseProjectID(d))
	listVirtualGatewaysPath += listVirtualGatewaysQueryParams

	listVirtualGatewaysResp, err := pagination.ListAllItems(
		listVirtualGatewaysClient,
		"marker",
		listVirtualGatewaysPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving DC virtual gateways: %s", err)
	}

	listVirtualGatewaysRespJson, err := json.Marshal(listVirtualGatewaysResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listVirtualGatewaysRespBody interface{}
	err = json.Unmarshal(listVirtualGatewaysRespJson, &listVirtualGatewaysRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("virtual_gateways", filterListVirtualGatewaysBody(
			flattenListVirtualGatewaysBody(listVirtualGatewaysRespBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListVirtualGatewaysBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("virtual_gateways", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                    utils.PathSearch("id", v, nil),
			"vpc_id":                utils.PathSearch("vpc_id", v, nil),
			"name":                  utils.PathSearch("name", v, nil),
			"type":                  utils.PathSearch("type", v, nil),
			"status":                utils.PathSearch("status", v, nil),
			"asn":                   utils.PathSearch("bgp_asn", v, nil),
			"local_ep_group":        utils.PathSearch("local_ep_group", v, nil),
			"description":           utils.PathSearch("description", v, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
		})
	}
	return rst
}

func filterListVirtualGatewaysBody(all []interface{}, d *schema.ResourceData) []interface{} {
	name := d.Get("name").(string)
	if name == "" {
		return all
	}

	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if name != fmt.Sprint(utils.PathSearch("name", v, nil)) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func buildListVirtualGatewaysQueryParams(d *schema.ResourceData, enterpriseProjectId string) string {
	res := ""
	if v, ok := d.GetOk("virtual_gateway_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		res = fmt.Sprintf("%s&vpc_id=%v", res, v)
	}

	if enterpriseProjectId != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, enterpriseProjectId)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
