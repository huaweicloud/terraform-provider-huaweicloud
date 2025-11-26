package workspace

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

var appGroupAuthorizationNotificationResendNonUpdatableParams = []string{
	"records",
	"records.*.id",
	"records.*.account",
	"records.*.account_auth_type",
	"records.*.account_auth_name",
	"records.*.app_group_id",
	"records.*.app_group_name",
	"records.*.mail_send_type",
	"is_notification_record",
}

// @API Workspace POST /v1/{project_id}/mails/actions/send
// @API Workspace POST /v1/{project_id}/mails/actions/send-by-authorization
func ResourceAppGroupAuthorizationNotificationResend() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppGroupAuthorizationNotificationResendCreate,
		ReadContext:   resourceAppGroupAuthorizationNotificationResendRead,
		UpdateContext: resourceAppGroupAuthorizationNotificationResendUpdate,
		DeleteContext: resourceAppGroupAuthorizationNotificationResendDelete,

		CustomizeDiff: config.FlexibleForceNew(appGroupAuthorizationNotificationResendNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the application group authorization notifications to be resent are located.`,
			},
			"records": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the authorization notification record or authorization record.",
						},
						"account": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the authorized account.",
						},
						"account_auth_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The type of the authorized object.",
						},
						"account_auth_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the authorized object.",
						},
						"app_group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the application group.",
						},
						"app_group_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the application group.",
						},
						"mail_send_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The type of authorization notification.",
						},
					},
				},
				Description: `The list of record IDs to resend authorization notification.`,
			},
			"is_notification_record": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to resend according to the authorization notification records.`,
			},
			// Internal parameter(s).
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceAppGroupAuthorizationNotificationResendCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                  = meta.(*config.Config)
		region               = cfg.GetRegion(d)
		httpUrl              = "v1/{project_id}/mails/actions/send-by-authorization"
		isNotificationRecord = d.Get("is_notification_record").(bool)
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	if isNotificationRecord {
		httpUrl = "v1/{project_id}/mails/actions/send"
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: utils.RemoveNil(buildAppGroupAuthorizationNotificationResendBodyParams(d.Get("records").([]interface{}),
			isNotificationRecord)),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error resending application group authorization notification: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	return nil
}

func buildAppGroupAuthorizationNotificationResendBodyParams(records []interface{}, isNotificationRecord bool) map[string]interface{} {
	if !isNotificationRecord {
		return map[string]interface{}{
			"records": utils.PathSearch("[*].id", records, nil),
		}
	}

	result := make([]map[string]interface{}, len(records))
	for i, v := range records {
		result[i] = map[string]interface{}{
			"id":                utils.PathSearch("id", v, nil),
			"account":           utils.ValueIgnoreEmpty(utils.PathSearch("account", v, nil)),
			"account_auth_type": utils.ValueIgnoreEmpty(utils.PathSearch("account_auth_type", v, nil)),
			"account_auth_name": utils.ValueIgnoreEmpty(utils.PathSearch("account_auth_name", v, nil)),
			"app_group_id":      utils.ValueIgnoreEmpty(utils.PathSearch("app_group_id", v, nil)),
			"app_group_name":    utils.ValueIgnoreEmpty(utils.PathSearch("app_group_name", v, nil)),
			"mail_send_type":    utils.ValueIgnoreEmpty(utils.PathSearch("mail_send_type", v, nil)),
		}
	}

	return map[string]interface{}{
		"records": result,
	}
}

func resourceAppGroupAuthorizationNotificationResendRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppGroupAuthorizationNotificationResendUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppGroupAuthorizationNotificationResendDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for resending application group authorization notifications. Deleting
    this resource will not clear the corresponding request record, but will only remove the resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
