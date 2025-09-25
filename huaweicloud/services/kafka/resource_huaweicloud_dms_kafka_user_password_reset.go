package kafka

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

var userPasswordResetNonUpdatableParams = []string{"instance_id", "user_name", "new_password"}

// @API Kafka PUT /v2/{project_id}/instances/{instance_id}/users/{user_name}
func ResourceUserPasswordReset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserPasswordResetCreate,
		ReadContext:   resourceUserPasswordResetRead,
		UpdateContext: resourceUserPasswordResetUpdate,
		DeleteContext: resourceUserPasswordResetDelete,

		CustomizeDiff: config.FlexibleForceNew(userPasswordResetNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The region where the user whose password is to be reset is located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the Kafka instance to which the user belongs.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the user to reset the password.`,
			},
			"new_password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: `The new password of the user.`,
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

func resourceUserPasswordResetCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		httpUrl    = "v2/{project_id}/instances/{instance_id}/users/{user_name}"
		instanceId = d.Get("instance_id").(string)
		userName   = d.Get("user_name").(string)
	)
	client, err := cfg.NewServiceClient("dmsv2", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	httpUrl = strings.ReplaceAll(httpUrl, "{project_id}", client.ProjectID)
	httpUrl = strings.ReplaceAll(httpUrl, "{instance_id}", instanceId)
	httpUrl = strings.ReplaceAll(httpUrl, "{user_name}", userName)
	resetPath := client.Endpoint + httpUrl

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"new_password": d.Get("new_password"),
		},
		OkCodes: []int{204},
	}

	_, err = client.Request("PUT", resetPath, &opt)
	if err != nil {
		return diag.Errorf("error resetting password for user (%s) under instance (%s): %s", userName, instanceId, err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate resource ID: %s", err)
	}
	d.SetId(randUUID)

	return nil
}

func resourceUserPasswordResetRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceUserPasswordResetUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceUserPasswordResetDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to reset user's password. Deleting
this resource will not clear the password reset operation record, but will only remove the resource information from the
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
