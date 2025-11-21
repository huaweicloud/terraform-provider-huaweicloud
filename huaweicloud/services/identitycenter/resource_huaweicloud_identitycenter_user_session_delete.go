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

var identityCenterUserSessionDeleteNonUpdateParams = []string{"identity_store_id", "user_id", "session_ids"}

// @API IdentityCenter POST /v1/identity-stores/{identity_store_id}/users/{user_id}/sessions/batch-delete
func ResourceIdentityCenterUserSessionDelete() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterUserSessionDeleteCreate,
		UpdateContext: resourceIdentityCenterUserSessionDeleteUpdate,
		ReadContext:   resourceIdentityCenterUserSessionDeleteRead,
		DeleteContext: resourceIdentityCenterUserSessionDeleteDelete,

		CustomizeDiff: config.FlexibleForceNew(identityCenterUserSessionDeleteNonUpdateParams),

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
			"session_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

func resourceIdentityCenterUserSessionDeleteCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/identity-stores/{identity_store_id}/users/{user_id}/sessions/batch-delete"
		product = "identitystore"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	batchDeletePath := client.Endpoint + httpUrl
	batchDeletePath = strings.ReplaceAll(batchDeletePath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))
	batchDeletePath = strings.ReplaceAll(batchDeletePath, "{user_id}", fmt.Sprintf("%v", d.Get("user_id")))

	batchDeleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildBatchDeleteSessionsBodyParams(d)),
	}
	_, err = client.Request("POST", batchDeletePath, &batchDeleteOpt)
	if err != nil {
		return diag.Errorf("error deleting Identity Center user sessions: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	return nil
}

func buildBatchDeleteSessionsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"session_ids": utils.ValueIgnoreEmpty(d.Get("session_ids")),
	}
	return bodyParams
}

func resourceIdentityCenterUserSessionDeleteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIdentityCenterUserSessionDeleteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIdentityCenterUserSessionDeleteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for operating the component. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
