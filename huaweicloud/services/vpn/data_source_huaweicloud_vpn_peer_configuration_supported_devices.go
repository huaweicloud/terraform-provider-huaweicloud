package vpn

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPN GET /v5/{project_id}/peer-configuration/supported-devices
func DataSourceVpnPeerConfigurationSupportedDevices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpnPeerConfigurationSupportedDevicesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"supported_devices": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     supportedDeviceSchema(),
			},
		},
	}
}

func supportedDeviceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vendor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"model": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVpnPeerConfigurationSupportedDevicesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v5/{project_id}/peer-configuration/supported-devices"
		product = "vpn"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving VPN peer configuration supported devices: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("supported_devices", flattenGetSupportedDevicesBody(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetSupportedDevicesBody(resp interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("supported_devices", resp, make([]any, 0))
	curArray := curJson.([]any)

	result := make([]map[string]interface{}, 0, len(curArray))
	for _, device := range curArray {
		result = append(result, map[string]interface{}{
			"vendor":  utils.PathSearch("vendor", device, nil),
			"type":    utils.PathSearch("type", device, nil),
			"model":   utils.PathSearch("model", device, nil),
			"version": utils.PathSearch("version", device, nil),
		})
	}

	return result
}
