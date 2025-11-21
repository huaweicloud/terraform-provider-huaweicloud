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

var identityCenterEmailVerifyNonUpdateParams = []string{"identity_store_id", "user_id"}

// @API IdentityCenter POST /v1/identity-stores/{identity_store_id}/users/{user_id}/verify-email
func ResourceIdentityCenterEmailVerify() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterEmailVerifyCreate,
		UpdateContext: resourceIdentityCenterEmailVerifyUpdate,
		ReadContext:   resourceIdentityCenterEmailVerifyRead,
		DeleteContext: resourceIdentityCenterEmailVeirfyDelete,

		CustomizeDiff: config.FlexibleForceNew(identityCenterEmailVerifyNonUpdateParams),

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
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceIdentityCenterEmailVerifyCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/identity-stores/{identity_store_id}/users/{user_id}/verify-email"
		product = "identitystore"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	verifyEmailPath := client.Endpoint + httpUrl
	verifyEmailPath = strings.ReplaceAll(verifyEmailPath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))
	verifyEmailPath = strings.ReplaceAll(verifyEmailPath, "{user_id}", fmt.Sprintf("%v", d.Get("user_id")))

	verifyEmailOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("POST", verifyEmailPath, &verifyEmailOpt)
	if err != nil {
		return diag.Errorf("error verifying Identity Center user email: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	return nil
}

func resourceIdentityCenterEmailVerifyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIdentityCenterEmailVerifyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIdentityCenterEmailVeirfyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for operating the component. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
