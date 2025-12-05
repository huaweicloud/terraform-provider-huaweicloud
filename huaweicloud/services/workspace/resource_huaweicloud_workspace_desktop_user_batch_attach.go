package workspace

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var desktopUserBatchAttachNonUpdatableParams = []string{
	"desktops",
}

// @API Workspace POST /v2/{project_id}/desktops/attach
func ResourceDesktopUserBatchAttach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDesktopUserBatchAttachCreate,
		ReadContext:   resourceDesktopUserBatchAttachRead,
		UpdateContext: resourceDesktopUserBatchAttachUpdate,
		DeleteContext: resourceDesktopUserBatchAttachDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(desktopUserBatchAttachNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the desktops are located.`,
			},

			// Required parameters.
			"desktops": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        desktopUserBatchAttachDesktopSchema(),
				Description: `The list of desktop information to be assigned.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func desktopUserBatchAttachDesktopSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"desktop_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the desktop to be assigned.`,
			},
			"computer_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The desktop name.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The user to whom the desktop belongs.`,
			},
			"user_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Schema: required.`,
			},
			"user_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The email of user.`,
			},
			"is_clear_data": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Whether to clean up desktop data when binding.`,
			},
			"attach_user_infos": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        desktopUserBatchAttachUserInfoSchema(),
				Description: `The list of user information to be assigned.`,
			},
		},
	}
}

func desktopUserBatchAttachUserInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "USER",
				Description: `The object type.`,
			},
			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The id of the user.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the desktop assignment object.`,
			},
			"user_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The user group to which the desktop user belongs.`,
			},
		},
	}
}

func buildDesktopUserBatchAttachUserInfosBodyParams(userInfos []interface{}) []interface{} {
	if len(userInfos) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(userInfos))
	for _, userInfo := range userInfos {
		result = append(result, utils.RemoveNil(map[string]interface{}{
			"type":       utils.ValueIgnoreEmpty(utils.PathSearch("type", userInfo, "USER")),
			"user_id":    utils.ValueIgnoreEmpty(utils.PathSearch("user_id", userInfo, "")),
			"user_name":  utils.ValueIgnoreEmpty(utils.PathSearch("user_name", userInfo, "")),
			"user_group": utils.ValueIgnoreEmpty(utils.PathSearch("user_group", userInfo, "")),
		}))
	}

	return result
}

func buildDesktopUserBatchAttachDesktopsBodyParams(desktops []interface{}) []interface{} {
	if len(desktops) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(desktops))
	for _, desktop := range desktops {
		desktopMap := map[string]interface{}{
			"desktop_id":    utils.ValueIgnoreEmpty(utils.PathSearch("desktop_id", desktop, "")),
			"computer_name": utils.ValueIgnoreEmpty(utils.PathSearch("computer_name", desktop, "")),
			"user_name":     utils.ValueIgnoreEmpty(utils.PathSearch("user_name", desktop, "")),
			"user_group":    utils.ValueIgnoreEmpty(utils.PathSearch("user_group", desktop, "")),
			"user_email":    utils.ValueIgnoreEmpty(utils.PathSearch("user_email", desktop, "")),
			"is_clear_data": utils.PathSearch("is_clear_data", desktop, true),
		}

		if v, ok := utils.PathSearch("attach_user_infos", desktop, []interface{}{}).([]interface{}); ok && len(v) > 0 {
			desktopMap["attach_user_infos"] = buildDesktopUserBatchAttachUserInfosBodyParams(v)
		}

		result = append(result, utils.RemoveNil(desktopMap))
	}

	return result
}

func buildDesktopUserBatchAttachBodyParams(d *schema.ResourceData) map[string]interface{} {
	body := map[string]interface{}{}

	if v, ok := d.GetOk("desktops"); ok {
		body["desktops"] = buildDesktopUserBatchAttachDesktopsBodyParams(v.([]interface{}))
	}

	return body
}

func resourceDesktopUserBatchAttachCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/desktops/attach"
	)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildDesktopUserBatchAttachBodyParams(d),
	}

	resp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating desktop user batch attach: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", respBody, "")
	if jobId.(string) != "" {
		_, err = waitForWorkspaceJobCompleted(ctx, client, jobId.(string), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceDesktopUserBatchAttachRead(ctx, d, meta)
}

func resourceDesktopUserBatchAttachRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDesktopUserBatchAttachUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDesktopUserBatchAttachDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for batch attaching users to a desktop. Deleting this
    resource will not clear the corresponding request record, but will only remove the resource information
    from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
