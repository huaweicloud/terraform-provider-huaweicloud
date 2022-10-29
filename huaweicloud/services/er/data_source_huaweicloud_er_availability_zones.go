// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product ER
// ---------------------------------------------------------------

package er

import (
	"context"
	"log"
	"strings"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceAvailabilityZones() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAvailabilityZonesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"names": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `The name list of the availability zones.`,
			},
		},
	}
}

func dataSourceAvailabilityZonesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	var mErr *multierror.Error

	// getAvailabilityZones: Query the availability zone list of Enterprise router
	var (
		getAvailabilityZonesHttpUrl = "v3/{project_id}/enterprise-router/availability-zones"
		getAvailabilityZonesProduct = "er"
	)
	getAvailabilityZonesClient, err := config.NewServiceClient(getAvailabilityZonesProduct, region)
	if err != nil {
		return diag.Errorf("error creating Instance Client: %s", err)
	}

	getAvailabilityZonesPath := getAvailabilityZonesClient.Endpoint + getAvailabilityZonesHttpUrl
	getAvailabilityZonesPath = strings.Replace(getAvailabilityZonesPath, "{project_id}", getAvailabilityZonesClient.ProjectID, -1)
	getAvailabilityZonesPath = strings.Replace(getAvailabilityZonesPath, "{id}", d.Id(), -1)

	getAvailabilityZonesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getAvailabilityZonesResp, err := getAvailabilityZonesClient.Request("GET", getAvailabilityZonesPath, &getAvailabilityZonesOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Instance")
	}

	getAvailabilityZonesRespBody, err := utils.FlattenResponse(getAvailabilityZonesResp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)
	log.Printf("[Lance] The data-source ID is: %s", uuid)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("names", parseAvailabilityZones(utils.PathSearch("availability_zones", getAvailabilityZonesRespBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func parseAvailabilityZones(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curArray := resp.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		if state := utils.PathSearch("state", v, nil); state == "available" {
			rst = append(rst, utils.PathSearch("code", v, nil))
		}
	}
	return rst
}
