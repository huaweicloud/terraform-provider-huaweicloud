package das

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var emailTemplatesBatchActionNonUpdatableParams = []string{
	"subscribe",
	"email_template_ids",
}

// @API DAS POST /v3/{project_id}/batch-inspection/batch-subscribe
func ResourceEmailTemplatesBatchAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEmailTemplatesBatchActionCreate,
		ReadContext:   resourceEmailTemplatesBatchActionRead,
		UpdateContext: resourceEmailTemplatesBatchActionUpdate,
		DeleteContext: resourceEmailTemplatesBatchActionDelete,

		CustomizeDiff: config.FlexibleForceNew(emailTemplatesBatchActionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the email templates are located.`,
			},

			// Required parameters.
			"subscribe": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Whether to subscribe to the email templates.`,
			},
			"email_template_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: `The list of email template IDs.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{Internal: true},
				),
			},
		},
	}
}

func buildEmailTemplatesBatchActionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"subscribe":              d.Get("subscribe"),
		"email_template_id_list": d.Get("email_template_ids"),
	}
}

func resourceEmailTemplatesBatchActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	httpUrl := "v3/{project_id}/batch-inspection/batch-subscribe"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildEmailTemplatesBatchActionBodyParams(d),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error subscribing to DAS email templates: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	return resourceEmailTemplatesBatchActionRead(ctx, d, meta)
}

func resourceEmailTemplatesBatchActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceEmailTemplatesBatchActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceEmailTemplatesBatchActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for subscribing to email templates. Deleting this resource
will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
