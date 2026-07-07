package taurusdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API TaurusDB GET /v3/{project_id}/proxy/flavors
func DataSourceTaurusDBProxyAzFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBProxyAzFlavorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"az_codes": {
				Type:     schema.TypeString,
				Required: true,
			},
			"proxy_engine_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"proxy_flavor_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     proxyAzFlavorGroupsSchema(),
			},
		},
	}
}

func proxyAzFlavorGroupsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"group_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"proxy_flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     proxyAzFlavorGroupFlavorsSchema(),
			},
		},
	}
}

func proxyAzFlavorGroupFlavorsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vcpus": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ram": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"spec_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"az_status": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"supported_ipv6": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceTaurusDBProxyAzFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/proxy/flavors"
	)
	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildGetProxyAzFlavorsQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving TaurusDB proxy AZ flavors: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("proxy_flavor_groups", flattenGetProxyAzFlavorGroupsBody(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetProxyAzFlavorsQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("?az_codes=%v", d.Get("az_codes").(string))
	res = fmt.Sprintf("%s&proxy_engine_name=%v", res, d.Get("proxy_engine_name"))
	return res
}

func flattenGetProxyAzFlavorGroupsBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("proxy_flavor_groups", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"group_type":    utils.PathSearch("group_type", v, nil),
			"proxy_flavors": flattenGetProxyAzFlavorGroupFlavorsBody(v),
		})
	}
	return res
}

func flattenGetProxyAzFlavorGroupFlavorsBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("proxy_flavors", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"id":             utils.PathSearch("id", v, nil),
			"db_type":        utils.PathSearch("db_type", v, nil),
			"vcpus":          utils.PathSearch("vcpus", v, nil),
			"ram":            utils.PathSearch("ram", v, nil),
			"spec_code":      utils.PathSearch("spec_code", v, nil),
			"az_status":      utils.PathSearch("az_status", v, nil),
			"supported_ipv6": utils.PathSearch("supported_ipv6", v, nil),
		})
	}
	return res
}
