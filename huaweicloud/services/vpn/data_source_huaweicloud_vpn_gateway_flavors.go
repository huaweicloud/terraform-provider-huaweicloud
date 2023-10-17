// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product VPN
// ---------------------------------------------------------------

package vpn

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceVpnGatewayFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceVpnGatewayFlavorsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the availability zone to get the flavors.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the flavor name.`,
			},
			"attachment_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "vpc",
				Description: `Specifies the attachment type.`,
			},
			"flavors": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `The list of message templates.`,
			},
		},
	}
}

func resourceVpnGatewayFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getGatewayFlavors: Query VPN gateway flavors
	var (
		getGatewayFlavorsHttpUrl = "v5/{project_id}/vpn-gateways/availability-zones"
		getGatewayFlavorsProduct = "vpn"
	)
	getGatewayFlavorsClient, err := cfg.NewServiceClient(getGatewayFlavorsProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	getGatewayFlavorsPath := getGatewayFlavorsClient.Endpoint + getGatewayFlavorsHttpUrl
	getGatewayFlavorsPath = strings.ReplaceAll(getGatewayFlavorsPath, "{project_id}", getGatewayFlavorsClient.ProjectID)

	getGatewayFlavorsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getGatewayFlavorsResp, err := getGatewayFlavorsClient.Request("GET", getGatewayFlavorsPath, &getGatewayFlavorsOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving VpnGatewayFlavors")
	}

	getGatewayFlavorsRespBody, err := utils.FlattenResponse(getGatewayFlavorsResp)
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
		d.Set("flavors", filterFlavors(utils.PathSearch("availability_zones", getGatewayFlavorsRespBody, nil), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func filterFlavors(all interface{}, d *schema.ResourceData) []interface{} {
	availabilityZone := d.Get("availability_zone").(string)
	attachmentType := d.Get("attachment_type").(string)
	allMap := all.(map[string]interface{})
	rst := make([]interface{}, 0, len(allMap))

	for k, v := range allMap {
		if flavorName, ok := d.GetOk("name"); ok && flavorName.(string) != k {
			continue
		}
		vMap := v.(map[string]interface{})
		if attachmentType == "er" && utils.StrSliceContains(utils.ExpandToStringList(vMap["er"].([]interface{})),
			availabilityZone) {
			rst = append(rst, k)
		}
		if attachmentType == "vpc" && utils.StrSliceContains(utils.ExpandToStringList(vMap["vpc"].([]interface{})),
			availabilityZone) {
			rst = append(rst, k)
		}
	}
	return rst
}
