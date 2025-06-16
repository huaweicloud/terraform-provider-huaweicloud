package dew

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

var scheduledDeleteSecretTaskNonUpdatableParams = []string{
	"secret_name",
	"action",
	"recovery_window_in_days",
}

// @API DEW POST /v1/{project_id}/secrets/{secret_name}/scheduled-deleted-tasks/create
// @API DEW POST /v1/{project_id}/secrets/{secret_name}/scheduled-deleted-tasks/cancel
func ResourceCsmsScheduledDeleteSecretTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCsmsScheduledDeleteSecretTaskCreate,
		ReadContext:   resourceCsmsScheduledDeleteSecretTaskRead,
		UpdateContext: resourceCsmsScheduledDeleteSecretTaskUpdate,
		DeleteContext: resourceCsmsScheduledDeleteSecretTaskDelete,

		CustomizeDiff: config.FlexibleForceNew(scheduledDeleteSecretTaskNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the region in which to create the resource.`,
			},
			"secret_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the secret name.`,
			},
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"create", "cancel"}, false),
				Description:  `Specifies the secret name.`,
			},
			"recovery_window_in_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
				Description: `Specifies the recovery window in days.`,
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

func buildReqsuestOpts(action string, recoveryDays int) golangsdk.RequestOpts {
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	if action == "create" {
		requestOpt.JSONBody = map[string]interface{}{
			"recovery_window_in_days": recoveryDays,
		}
	}

	return requestOpt
}

func resourceCsmsScheduledDeleteSecretTaskCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		product    = "kms"
		httpUrl    = "v1/{project_id}/secrets/{secret_name}/scheduled-deleted-tasks/{action}"
		requestOpt golangsdk.RequestOpts
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CSMS client: %s", err)
	}

	secretName := d.Get("secret_name").(string)
	action := d.Get("action").(string)
	recoveryDays := d.Get("recovery_window_in_days").(int)

	requestOpt = buildReqsuestOpts(action, recoveryDays)

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{secret_name}", secretName)
	requestPath = strings.ReplaceAll(requestPath, "{action}", action)

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error managing (%s) CSMS scheduled delete secret task: %s", action, err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	return nil
}

func resourceCsmsScheduledDeleteSecretTaskRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceCsmsScheduledDeleteSecretTaskUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceCsmsScheduledDeleteSecretTaskDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to manage scheduled delete secret task.
Deleting this resource will not recover the scheduled delete task, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
