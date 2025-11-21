package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM GET /v5/mfa-devices
func DataSourceIdentityV5VirtualMfaDevices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityV5VirtualMfaDevicesRead,

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the IAM user.",
			},
			"devices": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"serial_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityV5VirtualMfaDevicesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	var allDevices []interface{}
	var marker string
	var path string
	for {
		path = fmt.Sprintf("%sv5/mfa-devices", client.Endpoint) + buildListMfaDevicesV5Params(d, marker)
		reqOpt := &golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		r, err := client.Request("GET", path, reqOpt)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving virtual MFA devices")
		}
		resp, err := utils.FlattenResponse(r)
		if err != nil {
			return diag.FromErr(err)
		}
		devices := flattenListMfaDevicesV5Response(resp)
		allDevices = append(allDevices, devices...)

		marker = utils.PathSearch("page_info.next_marker", resp, "").(string)
		if marker == "" {
			break
		}
	}
	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	mErr := multierror.Append(nil,
		d.Set("devices", allDevices),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting virtual MFA devices fields: %s", err)
	}
	return nil
}

func buildListMfaDevicesV5Params(d *schema.ResourceData, marker string) string {
	res := "?limit=100"
	if v, ok := d.GetOk("user_id"); ok {
		res = fmt.Sprintf("%s&user_id=%v", res, v)
	}
	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}
	return res
}

func flattenListMfaDevicesV5Response(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	devices := utils.PathSearch("mfa_devices", resp, make([]interface{}, 0)).([]interface{})
	result := make([]interface{}, len(devices))
	for i, device := range devices {
		result[i] = map[string]interface{}{
			"enabled":       utils.PathSearch("enabled", device, nil),
			"serial_number": utils.PathSearch("serial_number", device, nil),
			"user_id":       utils.PathSearch("user_id", device, nil),
		}
	}
	return result
}
