package identitycenter

import (
	"context"
	"fmt"
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

// @API IdentityCenter POST /v1/identity-stores/{identity_store_id}/users/retrieve-mfa-devices
func DataSourceIdentityCenterMfaDevices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityCenterMfaDevicesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"identity_store_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mfa_devices": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"device_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mfa_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"registered_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityCenterMfaDevicesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v1/identity-stores/{identity_store_id}/users/retrieve-mfa-devices"
		product = "identitystore"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(batchListMfaDevicesForUserBodyParams(d)),
	}
	listResp, err := client.Request("POST", listPath, &opt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center user mfa devices.")
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

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("mfa_devices", flattenMfaDevices(utils.PathSearch("user_mfa_devices_entry_list|[0].mfa_devices", listRespBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenMfaDevices(data interface{}) []interface{} {
	devices := data.([]interface{})
	if len(devices) < 1 {
		return nil
	}
	result := make([]interface{}, 0, len(devices))
	for _, device := range devices {
		result = append(result, map[string]interface{}{
			"registered_date": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("registered_date", device, float64(0)).(float64))/1000, false),
			"device_name":     utils.PathSearch("device_name", device, nil),
			"display_name":    utils.PathSearch("display_name", device, nil),
			"mfa_type":        utils.PathSearch("mfa_type", device, nil),
			"device_id":       utils.PathSearch("device_id", device, nil),
		})
	}
	return result
}

func batchListMfaDevicesForUserBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"user_list": []map[string]interface{}{
			{
				"user_id":           utils.ValueIgnoreEmpty(d.Get("user_id")),
				"identity_store_id": utils.ValueIgnoreEmpty(d.Get("identity_store_id")),
			},
		},
	}
	return bodyParams
}
