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

var identityCenterPasswordPolicyNonUpdateParams = []string{"identity_store_id"}

// @API IdentityCenter PUT /v1/identity-stores/{identity_store_id}/password-policy
// @API IdentityCenter GET /v1/identity-stores/{identity_store_id}/password-policy
func ResourceIdentityCenterPasswordPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterPasswordPolicyCreateOrUpdate,
		UpdateContext: resourceIdentityCenterPasswordPolicyCreateOrUpdate,
		ReadContext:   resourceIdentityCenterPasswordPolicyRead,
		DeleteContext: resourceIdentityCenterPasswordPolicyDelete,

		CustomizeDiff: config.FlexibleForceNew(identityCenterPasswordPolicyNonUpdateParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"identity_store_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"minimum_password_length": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"require_lowercase_characters": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"require_numbers": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"require_symbols": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"require_uppercase_characters": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"max_password_age": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"password_reuse_prevention": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceIdentityCenterPasswordPolicyCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	identityStoreId := d.Get("identity_store_id").(string)

	var (
		httpUrl = "v1/identity-stores/{identity_store_id}/password-policy"
		product = "identitystore"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdatePasswordPolicyBodyParams(d)),
	}
	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating Identity Center password policy: %s", err)
	}
	if d.IsNewResource() {
		d.SetId(identityStoreId)
	}

	return resourceIdentityCenterPasswordPolicyRead(ctx, d, meta)
}

func buildUpdatePasswordPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"password_policy": map[string]interface{}{
			"minimum_password_length":      d.Get("minimum_password_length").(int),
			"require_lowercase_characters": d.Get("require_lowercase_characters").(bool),
			"require_numbers":              d.Get("require_numbers").(bool),
			"require_symbols":              d.Get("require_symbols").(bool),
			"require_uppercase_characters": d.Get("require_uppercase_characters").(bool),
			"max_password_age":             d.Get("max_password_age").(int),
			"password_reuse_prevention":    d.Get("password_reuse_prevention").(int),
		},
	}
	return bodyParams
}

func buildDefaultPasswordPolicyBodyParams() map[string]interface{} {
	bodyParams := map[string]interface{}{
		"password_policy": map[string]interface{}{
			"minimum_password_length":      8,
			"require_lowercase_characters": true,
			"require_numbers":              true,
			"require_symbols":              true,
			"require_uppercase_characters": true,
			"password_reuse_prevention":    1,
		},
	}
	return bodyParams
}

func resourceIdentityCenterPasswordPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v1/identity-stores/{identity_store_id}/password-policy"
		product = "identitystore"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center password policy.")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("minimum_password_length", utils.PathSearch("password_policy.minimum_password_length", getRespBody, nil)),
		d.Set("require_lowercase_characters", utils.PathSearch("password_policy.require_lowercase_characters", getRespBody, nil)),
		d.Set("require_numbers", utils.PathSearch("password_policy.require_numbers", getRespBody, nil)),
		d.Set("require_symbols", utils.PathSearch("password_policy.require_symbols", getRespBody, nil)),
		d.Set("require_uppercase_characters", utils.PathSearch("password_policy.require_uppercase_characters", getRespBody, nil)),
		d.Set("max_password_age", utils.PathSearch("password_policy.max_password_age", getRespBody, nil)),
		d.Set("password_reuse_prevention", utils.PathSearch("password_policy.password_reuse_prevention", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceIdentityCenterPasswordPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/identity-stores/{identity_store_id}/password-policy"
		product = "identitystore"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildDefaultPasswordPolicyBodyParams()),
	}
	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating Identity Center password policy: %s", err)
	}

	errorMsg := "Deleting password policy is unsupported. The resource is removed from the state," +
		" and the password policy is reset to default setting."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
