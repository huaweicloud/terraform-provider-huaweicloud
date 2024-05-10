package cbh

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

// @API CBH GET /v2/{project_id}/cbs/available-zone
func DataSourceAvailabilityZones() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceAvailabilityZonesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func datasourceAvailabilityZonesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                          = meta.(*config.Config)
		region                       = cfg.GetRegion(d)
		httpUrl                      = "v2/{project_id}/cbs/available-zone"
		listAvailabilityZonesProduct = "cbh"
	)
	listAvailabilityZonesClient, err := cfg.NewServiceClient(listAvailabilityZonesProduct, region)
	if err != nil {
		return diag.Errorf("error creating CBH client: %s", err)
	}

	listAvailabilityZonesPath := listAvailabilityZonesClient.Endpoint + httpUrl
	listAvailabilityZonesPath = strings.ReplaceAll(listAvailabilityZonesPath, "{project_id}",
		listAvailabilityZonesClient.ProjectID)
	listAvailabilityZonesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listResp, err := listAvailabilityZonesClient.Request("GET", listAvailabilityZonesPath, &listAvailabilityZonesOpt)
	if err != nil {
		return diag.Errorf("error retrieving CBH availability zones: %s", err)
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("availability_zones", filterAvailabilityZones(flattenAvailabilityZones(listRespBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func filterAvailabilityZones(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("name"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("name", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("display_name"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("display_name", v, nil)) {
			continue
		}
		rst = append(rst, v)
	}

	return rst
}

func flattenAvailabilityZones(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("availability_zone", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name":         utils.PathSearch("id", v, nil),
			"region_id":    utils.PathSearch("region_id", v, nil),
			"display_name": utils.PathSearch("display_name", v, nil),
			"type":         utils.PathSearch("type", v, nil),
			"status":       utils.PathSearch("status", v, nil),
		})
	}

	return rst
}
