package identitycenter

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var identityCenterPasswordResetNonUpdateParams = []string{"identity_store_id", "user_id", "mode"}

// @API IdentityCenter POST /v1/identity-stores/{identity_store_id}/users/{user_id}/reset-password
func ResourceIdentityCenterPasswordReset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterPasswordResetCreate,
		UpdateContext: resourceIdentityCenterPasswordResetUpdate,
		ReadContext:   resourceIdentityCenterPasswordResetRead,
		DeleteContext: resourceIdentityCenterPasswordResetDelete,

		CustomizeDiff: config.FlexibleForceNew(identityCenterPasswordResetNonUpdateParams),

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
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mode": {
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

func resourceIdentityCenterPasswordResetCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/identity-stores/{identity_store_id}/users/{user_id}/reset-password"
		product = "identitystore"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	resetPwdModePath := client.Endpoint + httpUrl
	resetPwdModePath = strings.ReplaceAll(resetPwdModePath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))
	resetPwdModePath = strings.ReplaceAll(resetPwdModePath, "{user_id}", fmt.Sprintf("%v", d.Get("user_id")))

	resetPwdModeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildResetPwdModeBodyParams(d)),
	}

	_, err = client.Request("POST", resetPwdModePath, &resetPwdModeOpt)
	if err != nil {
		return diag.Errorf("error resetting Identity Center user password: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	return nil
}

func buildResetPwdModeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"mode": utils.ValueIgnoreEmpty(d.Get("mode").(string)),
	}
	return bodyParams
}

func resourceIdentityCenterPasswordResetRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIdentityCenterPasswordResetUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIdentityCenterPasswordResetDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for operating the component. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
