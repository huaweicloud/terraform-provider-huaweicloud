package iam

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM POST /v3.0/OS-MFA/virtual-mfa-devices
// @API IAM DELETE /v3.0/OS-MFA/virtual-mfa-devices
// @API IAM GET /v3.0/OS-MFA/users/{user_id}/virtual-mfa-device
// @API IAM GET /v3.0/OS-USER/users/{user_id}
// @API IAM GET /v3/auth/domains
func ResourceIdentityVirtualMFADevice() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMFACreate,
		ReadContext:   resourceMFARead,
		DeleteContext: resourceMFADelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceMFAImportState,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The virtual MFA device name.`,
			},
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The user ID which the virtual MFA device belongs to.`,
			},

			// Attribute
			"base32_string_seed": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The base32 seed, which a third-patry system can use to generate a CAPTCHA code.`,
			},
			"qr_code_png": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The QR code PNG image.`,
			},
		},
	}
}

func resourceMFACreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	product := "iam"
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating IAM Client: %s", err)
	}

	createMFAHttpUrl := "v3.0/OS-MFA/virtual-mfa-devices"
	createMFAPath := client.Endpoint + createMFAHttpUrl
	createMFAOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"virtual_mfa_device": map[string]string{
				"name":    d.Get("name").(string),
				"user_id": d.Get("user_id").(string),
			},
		},
	}
	createMFAResp, err := client.Request("POST", createMFAPath, &createMFAOpt)
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

	// Base32StringSeed and QRCodePNG must be set here, because they are not available via ShowUserMFADevice
	seed := utils.PathSearch("virtual_mfa_device.base32_string_seed", createMFARespBody, "").(string)
	if seed == "" {
		return diag.Errorf("error creating IAM virtual MFA device: base32_string_seed is not found in API response")
	}
	d.Set("base32_string_seed", seed)

	// Get domain name and user name combine with `base32_string_seed` to form `qr_code_png`.
	domainName, err := getDomainName(client)
	if err != nil {
		return diag.FromErr(err)
	}
	userName, err := getUserName(client, d.Get("user_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	qRCodePNG := fmt.Sprintf("optauth://totp/huawei:%s@%s?secret=%s", domainName, userName, seed)
	d.Set("qr_code_png", qRCodePNG)

	return nil
}

func getDomainName(client *golangsdk.ServiceClient) (string, error) {
	getDomainNameHttpUrl := "v3/auth/domains"
	getDomainNamePath := client.Endpoint + getDomainNameHttpUrl
	getDomainNameOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getDomainNameResp, err := client.Request("GET", getDomainNamePath, &getDomainNameOpt)
	if err != nil {
		return "", fmt.Errorf("error getting IAM domain name: %s", err)
	}
	getDomainNameRespBody, err := utils.FlattenResponse(getDomainNameResp)
	if err != nil {
		return "", err
	}
	domainName := utils.PathSearch("domains[0].name", getDomainNameRespBody, "").(string)
	if domainName == "" {
		return "", fmt.Errorf("error getting IAM domain name: name is not found in API response")
	}
	return domainName, nil
}

func getUserName(client *golangsdk.ServiceClient, userID string) (string, error) {
	getUserNameHttpUrl := "v3.0/OS-USER/users/{user_id}"
	getUserNamePath := client.Endpoint + getUserNameHttpUrl
	getUserNamePath = strings.ReplaceAll(getUserNamePath, "{user_id}", userID)
	getUserNameOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getUserNameResp, err := client.Request("GET", getUserNamePath, &getUserNameOpt)
	if err != nil {
		return "", fmt.Errorf("error getting IAM user name: %s", err)
	}
	getUserNameRespBody, err := utils.FlattenResponse(getUserNameResp)
	if err != nil {
		return "", err
	}
	userName := utils.PathSearch("user.name", getUserNameRespBody, "").(string)
	if userName == "nil" {
		return "", fmt.Errorf("error getting IAM user name: name is not found in API response")
	}
	return userName, nil
}

func resourceMFARead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	getMFAProduct := "iam"
	getMFAClient, err := cfg.NewServiceClient(getMFAProduct, region)
	if err != nil {
		return diag.Errorf("error creating IAM Client: %s", err)
	}

	getMFAHttpUrl := "v3.0/OS-MFA/users/{user_id}/virtual-mfa-device"
	getMFAPath := getMFAClient.Endpoint + getMFAHttpUrl
	getMFAPath = strings.ReplaceAll(getMFAPath, "{user_id}", d.Get("user_id").(string))
	getMFAOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getMFAResp, err := getMFAClient.Request("GET", getMFAPath, &getMFAOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IAM virtual MFA device")
	}
	getMFARespBody, err := utils.FlattenResponse(getMFAResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("virtual_mfa_device.serial_number", getMFARespBody, "").(string)
	if id == "nil" {
		return diag.Errorf("error getting IAM virtual MFA device: serial_number is not found in API response")
	}
	d.SetId(id)

	index := strings.LastIndex(id, ":mfa/") + len(":mfa/")
	d.Set("name", id[index:])

	return nil
}

func resourceMFADelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	dproduct := "iam"
	client, err := cfg.NewServiceClient(dproduct, region)
	if err != nil {
		return diag.Errorf("error creating IAM Client: %s", err)
	}

	deleteMFAHttpUrl := "v3.0/OS-MFA/virtual-mfa-devices?user_id={user_id}&serial_number={id}"
	deleteMFAPath := client.Endpoint + deleteMFAHttpUrl
	deleteMFAPath = strings.ReplaceAll(deleteMFAPath, "{user_id}", d.Get("user_id").(string))
	deleteMFAPath = strings.ReplaceAll(deleteMFAPath, "{id}", d.Id())
	deleteMFAOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = client.Request("DELETE", deleteMFAPath, &deleteMFAOpt)
	if err != nil {
		return diag.Errorf("error deleting IAM virtual MFA device: %s", err)
	}

	// When input any serial number, API always return 200 and no more information, so need to check it truly deleted
	getMFAHttpUrl := "v3.0/OS-MFA/users/{user_id}/virtual-mfa-device"
	getMFAPath := client.Endpoint + getMFAHttpUrl
	getMFAPath = strings.ReplaceAll(getMFAPath, "{user_id}", d.Get("user_id").(string))
	getMFAOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("GET", getMFAPath, &getMFAOpt)
	if err == nil {
		return diag.Errorf("error deleting IAM virtual MFA device: the virtual MFA device still exists")
	}

	return nil
}

func resourceMFAImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, d.Set("user_id", d.Id())
}
