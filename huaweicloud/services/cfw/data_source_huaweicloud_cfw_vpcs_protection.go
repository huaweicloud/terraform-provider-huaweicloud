package cfw

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

// @API CFW GET /v1/{project_id}/vpcs/protection
func DataSourceCfwVpcsProtection() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCfwVpcsProtectionRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"object_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the protected object ID.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the enterprise project ID.`,
			},
			"fw_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the firewall instance ID.`,
			},
			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The data of VPC protection information.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"other_protect_vpcs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of other protected VPCs.`,
							Elem:        vpcIdSchema(),
						},
						"other_total": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The total number of other protected VPCs.`,
						},
						"protect_vpcs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of protected VPCs.`,
							Elem:        vpcIdSchema(),
						},
						"self_protect_vpcs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of self-protected VPCs.`,
							Elem:        vpcIdSchema(),
						},
						"self_total": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The total number of self-protected VPCs.`,
						},
						"total": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The total number of protected VPCs.`,
						},
						"total_assets": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The total number of assets.`,
						},
					},
				},
			},
		},
	}
}

func vpcIdSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The VPC ID.`,
			},
		},
	}
}

func buildVpcsProtectionQueryParams(cfg *config.Config, d *schema.ResourceData) string {
	queryParams := fmt.Sprintf("?object_id=%s", d.Get("object_id").(string))

	if v, ok := d.GetOk("fw_instance_id"); ok {
		queryParams = fmt.Sprintf("%s&fw_instance_id=%s", queryParams, v.(string))
	}

	if v := cfg.GetEnterpriseProjectID(d); v != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%s", queryParams, v)
	}

	return queryParams
}

func dataSourceCfwVpcsProtectionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cfw"
		httpUrl = "v1/{project_id}/vpcs/protection"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildVpcsProtectionQueryParams(cfg, d)

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", listPath, &reqOpt)
	if err != nil {
		return diag.Errorf("error retrieving CFW VPCs protection: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenVpcsProtectionData(utils.PathSearch("data", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenVpcsProtectionData(dataResp interface{}) []interface{} {
	if dataResp == nil {
		return nil
	}

	result := map[string]interface{}{
		"other_protect_vpcs": flattenVpcIds(
			utils.PathSearch("other_protect_vpcs", dataResp, make([]interface{}, 0)).([]interface{})),
		"other_total":  utils.PathSearch("other_total", dataResp, nil),
		"protect_vpcs": flattenVpcIds(utils.PathSearch("protect_vpcs", dataResp, make([]interface{}, 0)).([]interface{})),
		"self_protect_vpcs": flattenVpcIds(
			utils.PathSearch("self_protect_vpcs", dataResp, make([]interface{}, 0)).([]interface{})),
		"self_total":   utils.PathSearch("self_total", dataResp, nil),
		"total":        utils.PathSearch("total", dataResp, nil),
		"total_assets": utils.PathSearch("total_assets", dataResp, nil),
	}

	return []interface{}{result}
}

func flattenVpcIds(vpcList []interface{}) []interface{} {
	if len(vpcList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(vpcList))
	for _, v := range vpcList {
		result = append(result, map[string]interface{}{
			"vpc_id": utils.PathSearch("vpc_id", v, ""),
		})
	}

	return result
}
