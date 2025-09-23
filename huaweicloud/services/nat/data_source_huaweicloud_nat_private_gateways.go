package nat

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

// @API NAT GET /v3/{project_id}/private-nat/gateways
func DataSourcePrivateGateways() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePrivateGatewaysRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the private NAT gateways are located.",
			},
			"gateway_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the private NAT gateway.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the private NAT gateway.",
			},
			"spec": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The specification of the private NAT gateways.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The current status of the private NAT gateways.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the VPC to which the private NAT gateways belong.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the subnet to which the private NAT gateways belong.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the enterprise project to which the private NAT gateways belong.",
			},
			"description": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The description of the private NAT gateway.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The key/value pairs to associate the private NAT gateways.",
			},
			"gateways": {
				Type:        schema.TypeList,
				Elem:        gatewayGatewaysSchema(),
				Computed:    true,
				Description: "The list of the private NAT gateways.",
			},
		},
	}
}

func gatewayGatewaysSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the private NAT gateway.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the private NAT gateway.",
			},
			"spec": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The specification of the private NAT gateway.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the private NAT gateway.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the private NAT gateway.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the private NAT gateway.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the private NAT gateway.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the VPC to which the private NAT gateway belongs.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the subnet to which the private NAT gateway belongs.",
			},
			"ngport_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address of the NG port of the private NAT gateway.",
			},
			"rule_max": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of rules of the private NAT gateway.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the enterprise project to which the private NAT gateway belongs.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The key/value pairs to associate the private NAT gateway.",
			},
		},
	}
	return &sc
}

func dataSourcePrivateGatewaysRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/private-nat/gateways"
		product = "nat"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating NAT client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildListGatewaysQueryParams(d, cfg)
	resp, err := pagination.ListAllItems(
		client,
		"marker",
		requestPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving private NAT gateways %s", err)
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		return diag.FromErr(err)
	}
	var respBody interface{}
	err = json.Unmarshal(respJson, &respBody)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	curJson := utils.PathSearch("gateways", respBody, make([]interface{}, 0))
	curArray := curJson.([]interface{})

	var mErr *multierror.Error
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("gateways", flattenListGatewaysResponseBodyGateways(filterListGatewaysResponseByTags(curArray, d))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListGatewaysResponseBodyGateways(resp []interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"id":                    utils.PathSearch("id", v, nil),
			"name":                  utils.PathSearch("name", v, nil),
			"spec":                  utils.PathSearch("spec", v, nil),
			"description":           utils.PathSearch("description", v, nil),
			"status":                utils.PathSearch("status", v, nil),
			"created_at":            utils.PathSearch("created_at", v, nil),
			"updated_at":            utils.PathSearch("updated_at", v, nil),
			"vpc_id":                utils.PathSearch("downlink_vpcs[0].vpc_id", v, nil),
			"subnet_id":             utils.PathSearch("downlink_vpcs[0].virsubnet_id", v, nil),
			"ngport_ip_address":     utils.PathSearch("downlink_vpcs[0].ngport_ip_address", v, nil),
			"rule_max":              utils.PathSearch("rule_max", v, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
			"tags":                  utils.FlattenTagsToMap(utils.PathSearch("tags", v, nil)),
		})
	}
	return rst
}

func filterListGatewaysResponseByTags(all []interface{}, d *schema.ResourceData) []interface{} {
	tagFilter := d.Get("tags").(map[string]interface{})
	if len(tagFilter) == 0 {
		return all
	}

	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		tags := utils.FlattenTagsToMap(utils.PathSearch("tags", v, nil))
		tagMap := utils.ExpandToStringMap(tags)
		if !utils.HasMapContains(tagMap, tagFilter) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func buildListGatewaysQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	res := ""
	epsID := cfg.GetEnterpriseProjectID(d)
	descriptionList := d.Get("description").([]interface{})

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
		res = fmt.Sprintf("%s&vpc_id=%v", res, v)
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		res = fmt.Sprintf("%s&virsubnet_id=%v", res, v)
	}
	if len(descriptionList) > 0 {
		for _, v := range descriptionList {
			res = fmt.Sprintf("%s&description=%v", res, v)
		}
	}
	if epsID != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%s", res, epsID)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
