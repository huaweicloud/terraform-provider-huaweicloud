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

// @API CFW GET /v1/{project_id}/firewall/east-west/enterprise-router
func DataSourceCfwEastWestFirewallEr() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCfwEastWestFirewallErRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"fw_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the firewall instance ID.`,
			},
			"er_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of enterprise routers associated with the east-west firewall.`,
				Elem:        eastWestFirewallErSchema(),
			},
		},
	}
}

func eastWestFirewallErSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"er_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The enterprise router ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The enterprise router name.`,
			},
		},
	}
}

func buildEastWestFirewallErQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?fw_instance_id=%s", d.Get("fw_instance_id").(string))
}

func dataSourceCfwEastWestFirewallErRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cfw"
		httpUrl = "v1/{project_id}/firewall/east-west/enterprise-router"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildEastWestFirewallErQueryParams(d)

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", listPath, &reqOpt)
	if err != nil {
		return diag.Errorf("error retrieving CFW east-west firewall enterprise router: %s", err)
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
		d.Set("er_list", flattenEastWestFirewallErList(utils.PathSearch("data.er_list", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenEastWestFirewallErList(erList []interface{}) []interface{} {
	if erList == nil {
		return nil
	}

	result := make([]interface{}, 0, len(erList))
	for _, er := range erList {
		result = append(result, map[string]interface{}{
			"er_id": utils.PathSearch("er_id", er, nil),
			"name":  utils.PathSearch("name", er, nil),
		})
	}

	return result
}
