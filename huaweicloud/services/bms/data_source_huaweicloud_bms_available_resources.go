package bms

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

// @API BMS GET /v1/{project_id}/baremetalservers/available_resource
func DataSourceAvailableResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAvailableResourcesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"available_resource": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     availableResourcesAvailableResourceSchema(),
			},
		},
	}
}

func availableResourcesAvailableResourceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     availableResourcesAvailableResourceFlavorsSchema(),
			},
		},
	}
	return &sc
}

func availableResourcesAvailableResourceFlavorsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceAvailableResourcesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v1/{project_id}/baremetalservers/available_resource"
		product = "bms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating BMS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getQueryParams := buildGetAvailableResourcesQueryParams(d)
	getPath += getQueryParams

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving BMS available resources: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("available_resource", flattenAvailableResourcesAvailableResource(getRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetAvailableResourcesQueryParams(d *schema.ResourceData) string {
	res := ""
	for _, v := range d.Get("availability_zone").([]interface{}) {
		res = fmt.Sprintf("%s&availability_zone=%v", res, v)
	}
	res = "?" + res[1:]
	return res
}

func flattenAvailableResourcesAvailableResource(resp interface{}) []interface{} {
	curJson := utils.PathSearch("available_resource", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"availability_zone": utils.PathSearch("availability_zone", v, nil),
			"flavors":           flattenAvailableResourcesAvailableFlavorResource(v),
		})
	}
	return rst
}

func flattenAvailableResourcesAvailableFlavorResource(resp interface{}) []interface{} {
	curJson := utils.PathSearch("flavors", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"flavor_id": utils.PathSearch("flavor_id", v, nil),
			"status":    utils.PathSearch("status", v, nil),
		})
	}
	return rst
}
