package iam

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"base32_string_seed": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"qr_code_png": {
				Type:     schema.TypeString,
				Computed: true,
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

	id, err := jmespath.Search("virtual_mfa_device.serial_number", createMFARespBody)
	if err != nil {
		return diag.Errorf("error creating IAM virtual MFA device: serial_number is not found in API response")
	}
	d.SetId(id.(string))

	// Base32StringSeed and QRCodePNG must be set here, because they are not available via ShowUserMFADevice
	seed, err := jmespath.Search("virtual_mfa_device.base32_string_seed", createMFARespBody)
	if err != nil {
		return diag.Errorf("error creating IAM virtual MFA device: base32_string_seed is not found in API response")
	}
	d.Set("base32_string_seed", seed.(string))

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
	domainName, err := jmespath.Search("domains[0].name", getDomainNameRespBody)
	if err != nil {
		return "", fmt.Errorf("error getting IAM domain name: name is not found in API response")
	}
	return domainName.(string), nil
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
	userName, err := jmespath.Search("user.name", getUserNameRespBody)
	if err != nil {
		return "", fmt.Errorf("error getting IAM user name: name is not found in API response")
	}
	return userName.(string), nil
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

	id, err := jmespath.Search("virtual_mfa_device.serial_number", getMFARespBody)
	if err != nil {
		return diag.Errorf("error getting IAM virtual MFA device: serial_number is not found in API response")
	}
	d.SetId(id.(string))

	index := strings.LastIndex(id.(string), ":mfa/") + len(":mfa/")
	d.Set("name", id.(string)[index:])

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
