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

// @API IAM POST /v5/virtual-mfa-devices
// @API IAM DELETE /v5/virtual-mfa-devices
// @API IAM GET /v5/mfa-devices
var virtualMFADeviceNonUpdatableParams = []string{"user_id", "name"}

func ResourceIdentityV5VirtualMFADevice() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityV5MFACreate,
		ReadContext:   resourceIdentityV5MFARead,
		UpdateContext: resourceIdentityV5MFAUpdate,
		DeleteContext: resourceIdentityV5MFADelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceMFAImportStateV5,
		},

		CustomizeDiff: config.FlexibleForceNew(virtualMFADeviceNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"base32_string_seed": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceIdentityV5MFACreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	createMFAHttpUrl := "v5/virtual-mfa-devices"
	createMFAPath := iamClient.Endpoint + createMFAHttpUrl
	createMFAOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"virtual_mfa_device_name": d.Get("name").(string),
			"user_id":                 d.Get("user_id").(string),
		},
	}
	createMFAResp, err := iamClient.Request("POST", createMFAPath, &createMFAOpt)
	if err != nil {
		return diag.Errorf("error creating IAM virtual MFA device: %s", err)
	}
	createMFARespBody, err := utils.FlattenResponse(createMFAResp)
	if err != nil {
		return diag.FromErr(err)
	}
	id := utils.PathSearch("virtual_mfa_device.serial_number", createMFARespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating IAM virtual MFA device: serial_number is not found in API response")
	}
	d.SetId(id)
	seed := utils.PathSearch("virtual_mfa_device.base32_string_seed", createMFARespBody, "").(string)
	if seed == "" {
		return diag.Errorf("error creating IAM virtual MFA device: base32_string_seed is not found in API response")
	}
	err = d.Set("base32_string_seed", seed)
	if err != nil {
		return diag.Errorf("error set IAM MFA device field")
	}
	return resourceIdentityV5MFARead(ctx, d, meta)
}

func resourceIdentityV5MFARead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	getMFAHttpUrl := "v5/mfa-devices?user_id=" + d.Get("user_id").(string)
	getPath := iamClient.Endpoint + getMFAHttpUrl
	getMFAOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	response, err := iamClient.Request("GET", getPath, &getMFAOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IAM virtual MFA device")
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error flatten user MFA device")
	}
	mfa := utils.PathSearch("mfa_devices[0]", respBody, nil)
	if mfa == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error")
	}
	id := utils.PathSearch("serial_number", mfa, "").(string)
	if id == "" {
		return diag.Errorf("error found serial_number in the response")
	}
	d.SetId(id)
	if err = d.Set("enabled", utils.PathSearch("enabled", mfa, nil)); err != nil {
		return diag.Errorf("error setting mfa devices: %s", err)
	}
	return nil
}

func resourceIdentityV5MFAUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIdentityV5MFADelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	deleteMFAHttpUrl := "v5/virtual-mfa-devices"
	deleteMFAPath := iamClient.Endpoint + deleteMFAHttpUrl + buildDeleteMFAParamPath(d)
	deleteMFAOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = iamClient.Request("DELETE", deleteMFAPath, &deleteMFAOpt)
	if err != nil {
		return diag.Errorf("error deleting IAM virtual MFA device: %s", err)
	}
	return nil
}

func resourceMFAImportStateV5(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, d.Set("user_id", d.Id())
}

func buildDeleteMFAParamPath(d *schema.ResourceData) string {
	res := fmt.Sprintf("?user_id=%s", d.Get("user_id"))
	res += fmt.Sprintf("&serial_number=%s", d.Id())
	return res
}
