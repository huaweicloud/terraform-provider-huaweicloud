package er

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ER GET /v3/{project_id}/enterprise-router/availability-zones
func DataSourceAvailabilityZones() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAvailabilityZonesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// Attributes
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAvailabilityZonesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		listAvailabilityZoneHttpUrl = "v3/{project_id}/enterprise-router/availability-zones"
		listAvailabilityZoneProduct = "er"
	)
	listAvailabilityZoneClient, err := cfg.NewServiceClient(listAvailabilityZoneProduct, region)
	if err != nil {
		return diag.Errorf("error creating ER client: %s", err)
	}

	listAvailabilityZonePath := listAvailabilityZoneClient.Endpoint + listAvailabilityZoneHttpUrl
	listAvailabilityZonePath = strings.ReplaceAll(listAvailabilityZonePath, "{project_id}",
		listAvailabilityZoneClient.ProjectID)

	listAvailabilityZoneOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	listAvailabilityZoneResp, err := listAvailabilityZoneClient.Request("GET", listAvailabilityZonePath,
		&listAvailabilityZoneOpt)
	if err != nil {
		return diag.FromErr(err)
	}
	listAvailabilityZoneRespBody, err := utils.FlattenResponse(listAvailabilityZoneResp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("names", flattenListAvailabilityZone(listAvailabilityZoneRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListAvailabilityZone(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("availability_zones[?state == 'available'].code | sort(@)", resp, make([]string, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, v.(string))
	}
	return rst
}
