package identitycenter

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var identityCenterApplicationInstanceProfileNonUpdateParams = []string{"instance_id", "application_instance_id", "profile_id"}

// @API IdentityCenter DELETE /v1/instances/{instance_id}/application-instances/{application_instance_id}/profiles/{profile_id}
func ResourceIdentityCenterApplicationInstanceProfileDel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterApplicationInstanceProfileDeleteCreate,
		UpdateContext: resourceIdentityCenterApplicationInstanceProfileDeleteUpdate,
		ReadContext:   resourceIdentityCenterApplicationInstanceProfileDeleteRead,
		DeleteContext: resourceIdentityCenterApplicationInstanceProfileDeleteDelete,

		CustomizeDiff: config.FlexibleForceNew(identityCenterApplicationInstanceProfileNonUpdateParams),

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"application_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"profile_id": {
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

func resourceIdentityCenterApplicationInstanceProfileDeleteCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/instances/{instance_id}/application-instances/{application_instance_id}/profiles/{profile_id}"
		product = "identitycenter"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	deleteProfilePath := client.Endpoint + httpUrl
	deleteProfilePath = strings.ReplaceAll(deleteProfilePath, "{instance_id}", d.Get("instance_id").(string))
	deleteProfilePath = strings.ReplaceAll(deleteProfilePath, "{application_instance_id}", d.Get("application_instance_id").(string))
	deleteProfilePath = strings.ReplaceAll(deleteProfilePath, "{profile_id}", d.Get("profile_id").(string))

	deleteProfileOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deleteProfilePath, &deleteProfileOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting IdentityCenter application profile")
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	return nil
}

func resourceIdentityCenterApplicationInstanceProfileDeleteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIdentityCenterApplicationInstanceProfileDeleteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIdentityCenterApplicationInstanceProfileDeleteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for operating the component. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
