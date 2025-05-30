package eps

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var actionNonUpdatableParams = []string{
	"enterprise_project_id",
	"action",
}

// @API ER POST /v1.0/enterprise-projects/{enterprise_project_id}/action
func ResourceAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceActionCreate,
		UpdateContext: resourceActionUpdate,
		ReadContext:   resourceActionRead,
		DeleteContext: resourceActionDelete,

		CustomizeDiff: config.FlexibleForceNew(actionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of enterprise project to be operated.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The action type.`,
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

func resourceActionCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("eps", cfg.Region)
	if err != nil {
		return diag.Errorf("error creating EPS client: %s", err)
	}

	var (
		httpUrl = "v1.0/enterprise-projects/{enterprise_project_id}/action"
		epsId   = d.Get("enterprise_project_id").(string)
		action  = d.Get("action").(string)
	)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{enterprise_project_id}", epsId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		OkCodes: []int{204},
		JSONBody: map[string]interface{}{
			"action": action,
		},
	}
	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("unable to %s the enterprise project (%s): %s", action, epsId, err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	return nil
}

func resourceActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for operating the enterprise project. Deleting this
resource will not clear the corresponding request record, but will only remove the resource information from the tfstate
file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
