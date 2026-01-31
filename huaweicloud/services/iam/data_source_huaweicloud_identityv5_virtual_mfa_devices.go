package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM GET /v5/mfa-devices
func DataSourceV5VirtualMfaDevices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV5VirtualMfaDevicesRead,

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the IAM user.`,
			},
			"devices": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the MFA device is enabled.`,
						},
						"serial_number": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The serial number of the MFA device.`,
						},
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the user.`,
						},
					},
				},
				Description: `The list of virtual MFA devices that matched filter parameters.`,
			},
		},
	}
}

func listV5VirtualMfaDevices(client *golangsdk.ServiceClient, d *schema.ResourceData, marker string) ([]interface{}, error) {
	var (
		httpUrl = "v5/mfa-devices"
		limit   = 100
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = fmt.Sprintf("%s?limit=%v", listPath, limit)
	listPath += buildListMfaDevicesV5Params(d)
	for {
		listPathWithMarker := listPath
		if marker != "" {
			listPathWithMarker = fmt.Sprintf("%s&marker=%s", listPathWithMarker, marker)
		}

		reqOpt := &golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		resp, err := client.Request("GET", listPathWithMarker, reqOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		devices := utils.PathSearch("mfa_devices", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, devices...)
		if len(devices) < limit {
			break
		}

		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return result, nil
}

func dataSourceV5VirtualMfaDevicesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	devices, err := listV5VirtualMfaDevices(client, d, "")
	if err != nil {
		return diag.Errorf("error querying virtual MFA devices: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("devices", flattenV5MfaDevices(devices)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting virtual MFA devices fields: %s", err)
	}

	return nil
}

func buildListMfaDevicesV5Params(d *schema.ResourceData) string {
	if v, ok := d.GetOk("user_id"); ok {
		return fmt.Sprintf("&user_id=%v", v)
	}

	return ""
}

func flattenV5MfaDevices(devices []interface{}) []interface{} {
	if len(devices) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(devices))
	for _, device := range devices {
		result = append(result, map[string]interface{}{
			"enabled":       utils.PathSearch("enabled", device, nil),
			"serial_number": utils.PathSearch("serial_number", device, nil),
			"user_id":       utils.PathSearch("user_id", device, nil),
		})
	}
	return result
}
