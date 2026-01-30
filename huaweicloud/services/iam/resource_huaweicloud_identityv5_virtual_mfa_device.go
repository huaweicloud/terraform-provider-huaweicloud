package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var v5VirtualMFADeviceNonUpdatableParams = []string{"user_id", "name"}

// @API IAM POST /v5/virtual-mfa-devices
// @API IAM DELETE /v5/virtual-mfa-devices
// @API IAM GET /v5/mfa-devices
func ResourceV5VirtualMfaDevice() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV5MfaDeviceCreate,
		ReadContext:   resourceV5MfaDeviceRead,
		UpdateContext: resourceV5MfaDeviceUpdate,
		DeleteContext: resourceV5MfaDeviceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceV5MfaDeviceImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(v5VirtualMFADeviceNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the MFA device`,
			},
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the user`,
			},
			// Attributes.
			"base32_string_seed": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The key information used for third-party generation of image verification codes`,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the MFA device is enabled.`,
			},
			// Internal parameter(s).
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceV5MfaDeviceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		userId = d.Get("user_id").(string)
	)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	createMFAHttpUrl := "v5/virtual-mfa-devices"
	createMFAPath := client.Endpoint + createMFAHttpUrl
	createMFAOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"virtual_mfa_device_name": d.Get("name").(string),
			"user_id":                 userId,
		},
	}
	createMFAResp, err := client.Request("POST", createMFAPath, &createMFAOpt)
	if err != nil {
		return diag.Errorf("error creating virtual MFA device for user (%s): %s", userId, err)
	}

	createMFARespBody, err := utils.FlattenResponse(createMFAResp)
	if err != nil {
		return diag.FromErr(err)
	}

	serialNumber := utils.PathSearch("virtual_mfa_device.serial_number", createMFARespBody, "").(string)
	if serialNumber == "" {
		return diag.Errorf("unable to find serial number from API response")
	}

	d.SetId(serialNumber)

	seed := utils.PathSearch("virtual_mfa_device.base32_string_seed", createMFARespBody, "").(string)
	if seed == "" {
		return diag.Errorf("unable to find key information from API response")
	}

	if err := d.Set("base32_string_seed", seed); err != nil {
		return diag.Errorf("error setting base32_string_seed: %s", err)
	}

	return resourceV5MfaDeviceRead(ctx, d, meta)
}

func GetV5VirtualMfaDevice(client *golangsdk.ServiceClient, userId string) (interface{}, error) {
	getHttpUrl := fmt.Sprintf("v5/mfa-devices?user_id=%s", userId)
	getPath := client.Endpoint + getHttpUrl
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	response, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return nil, err
	}

	device := utils.PathSearch("mfa_devices[0]", respBody, nil)
	if device == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v5/mfa-devices",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("there are no MFA devices under the user (%s)", userId)),
			},
		}
	}

	return device, nil
}

func resourceV5MfaDeviceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		userId = d.Get("user_id").(string)
	)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	device, err := GetV5VirtualMfaDevice(client, userId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error getting virtual MFA device of user (%s)", userId))
	}

	serialNumber := utils.PathSearch("serial_number", device, "").(string)
	if serialNumber == "" {
		return diag.Errorf("unable to find serial number from API response")
	}

	d.SetId(serialNumber)
	return diag.FromErr(d.Set("enabled", utils.PathSearch("enabled", device, nil)))
}

func resourceV5MfaDeviceUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV5MfaDeviceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	deleteMFAHttpUrl := fmt.Sprintf("v5/virtual-mfa-devices?user_id=%s&serial_number=%s", d.Get("user_id"), d.Id())
	deleteMFAPath := client.Endpoint + deleteMFAHttpUrl
	deleteMFAOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deleteMFAPath, &deleteMFAOpt)
	if err != nil {
		return diag.Errorf("error deleting virtual MFA device of user (%s): %s", d.Get("user_id").(string), err)
	}
	return nil
}

func resourceV5MfaDeviceImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, d.Set("user_id", d.Id())
}
