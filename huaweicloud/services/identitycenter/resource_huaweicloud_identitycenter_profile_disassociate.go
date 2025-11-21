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

var identityCenterProfileDisassociateNonUpdateParams = []string{"instance_id", "identity_store_id", "accessor_type", "accessor_id"}

// @API IdentityCenter POST /v1/instances/{instance_id}/disassociate-profile
func ResourceIdentityCenterProfileDisassociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterProfileDisassociateCreate,
		UpdateContext: resourceIdentityCenterProfileDisassociateUpdate,
		ReadContext:   resourceIdentityCenterProfileDisassociateRead,
		DeleteContext: resourceIdentityCenterProfileDisassociateDelete,

		CustomizeDiff: config.FlexibleForceNew(identityCenterProfileDisassociateNonUpdateParams),

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"identity_store_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"accessor_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"accessor_id": {
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

func resourceIdentityCenterProfileDisassociateCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/instances/{instance_id}/disassociate-profile"
		product = "identitycenter"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	disassociateProfilePath := client.Endpoint + httpUrl
	disassociateProfilePath = strings.ReplaceAll(disassociateProfilePath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))

	disassociateProfileOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildDisassociateProfileUserBodyParams(d)),
	}

	_, err = client.Request("POST", disassociateProfilePath, &disassociateProfileOpt)
	if err != nil {
		return diag.Errorf("error disassociating Identity Center profile: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	return nil
}

func buildDisassociateProfileUserBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"id":                utils.ValueIgnoreEmpty(d.Get("accessor_id")),
		"accessor_type":     utils.ValueIgnoreEmpty(d.Get("accessor_type")),
		"identity_store_id": utils.ValueIgnoreEmpty(d.Get("identity_store_id")),
	}
	return bodyParams
}

func resourceIdentityCenterProfileDisassociateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIdentityCenterProfileDisassociateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIdentityCenterProfileDisassociateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for operating the component. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
