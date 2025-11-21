package identitycenter

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var identityCenterMfaManagementSettingNonUpdateParams = []string{"instance_id"}

// @API IdentityCenter POST /v1/instances/{instance_id}/mfa-devices/management-settings
// @API IdentityCenter GET /v1/instances/{instance_id}/mfa-devices/management-settings
func ResourceIdentityCenterMfaManagementSetting() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterMfaManagementSettingCreateOrUpdate,
		UpdateContext: resourceIdentityCenterMfaManagementSettingCreateOrUpdate,
		ReadContext:   resourceIdentityCenterMfaManagementSettingRead,
		DeleteContext: resourceIdentityCenterMfaManagementSettingDelete,

		CustomizeDiff: config.FlexibleForceNew(identityCenterMfaManagementSettingNonUpdateParams),

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"identity_store_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_permission": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceIdentityCenterMfaManagementSettingCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	instanceId := d.Get("instance_id").(string)

	var (
		httpUrl = "v1/instances/{instance_id}/mfa-devices/management-settings"
		product = "identitycenter"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	putPath := client.Endpoint + httpUrl
	putPath = strings.ReplaceAll(putPath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))

	putOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildPutMfaDeviceManagementForIdentityStoreBodyParams(d)),
	}
	_, err = client.Request("POST", putPath, &putOpt)
	if err != nil {
		return diag.Errorf("error updating Identity Center sso configuration: %s", err)
	}
	if d.IsNewResource() {
		d.SetId(instanceId)
	}

	return resourceIdentityCenterMfaManagementSettingRead(ctx, d, meta)
}

func buildPutMfaDeviceManagementForIdentityStoreBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"identity_store_id": utils.ValueIgnoreEmpty(d.Get("identity_store_id")),
		"user_permission":   utils.ValueIgnoreEmpty(d.Get("user_permission")),
	}
	return bodyParams
}

func resourceIdentityCenterMfaManagementSettingRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v1/instances/{instance_id}/mfa-devices/management-settings"
		product = "identitycenter"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center mfa management settings.")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("identity_store_id", utils.PathSearch("identity_store_id", getRespBody, nil)),
		d.Set("user_permission", utils.PathSearch("user_permission", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildDefaultMfaManagementForIdentityStoreBodyParams(identityStoreId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"identity_store_id": identityStoreId,
		"user_permission":   "ALL_ACTIONS",
	}
	return bodyParams
}

func resourceIdentityCenterMfaManagementSettingDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	identityStoreId := d.Get("identity_store_id").(string)

	var (
		httpUrl = "v1/instances/{instance_id}/mfa-devices/management-settings"
		product = "identitycenter"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	putPath := client.Endpoint + httpUrl
	putPath = strings.ReplaceAll(putPath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))

	putOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildDefaultMfaManagementForIdentityStoreBodyParams(identityStoreId)),
	}
	_, err = client.Request("POST", putPath, &putOpt)
	if err != nil {
		return diag.Errorf("error updating Identity Center sso configuration: %s", err)
	}

	errorMsg := "Deleting mfa management setting is unsupported. The resource is removed from the state," +
		" and the mfa management setting is reset to default setting."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
