package workspace

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var desktopUserBatchDetachNonUpdatableParams = []string{
	"desktops",
}

// @API Workspace POST /v2/{project_id}/desktops/batch-detach
func ResourceDesktopUserBatchDetach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDesktopUserBatchDetachCreate,
		ReadContext:   resourceDesktopUserBatchDetachRead,
		UpdateContext: resourceDesktopUserBatchDetachUpdate,
		DeleteContext: resourceDesktopUserBatchDetachDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(desktopUserBatchDetachNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the desktops and users are located.`,
			},
			"desktops": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 100,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"desktop_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The ID of the desktop to be detached.`,
						},
						"is_detach_all_users": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Whether to detach all users.`,
						},
						"detach_user_infos": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The ID of the user.`,
									},
									"user_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The name of the user or user group.`,
									},
									"user_group": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The user group which the user belongs to.`,
									},
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The type of the object.`,
									},
								},
							},
							Description: `The list of users to be detached.`,
						},
					},
				},
				Description: `The list of desktop user detach information.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildDesktopUserBatchDetachBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"desktops": buildDesktopUserBatchDetachDesktops(d.Get("desktops").([]interface{})),
	}
	return bodyParams
}

func buildDesktopUserBatchDetachDesktops(rawParams []interface{}) []map[string]interface{} {
	if len(rawParams) == 0 {
		return nil
	}

	desktops := make([]map[string]interface{}, len(rawParams))
	for i, rawParam := range rawParams {
		desktop := rawParam.(map[string]interface{})
		desktops[i] = map[string]interface{}{
			"desktop_id":          desktop["desktop_id"],
			"is_detach_all_users": desktop["is_detach_all_users"],
			"detach_user_infos":   buildDesktopUserBatchDetachUserInfos(desktop["detach_user_infos"].([]interface{})),
		}
	}

	return desktops
}

func buildDesktopUserBatchDetachUserInfos(rawParams []interface{}) []map[string]interface{} {
	if len(rawParams) == 0 {
		return nil
	}

	userInfos := make([]map[string]interface{}, len(rawParams))
	for i, rawParam := range rawParams {
		userInfo := rawParam.(map[string]interface{})
		userInfos[i] = map[string]interface{}{
			"user_id":    userInfo["user_id"],
			"user_name":  userInfo["user_name"],
			"user_group": userInfo["user_group"],
			"type":       userInfo["type"],
		}
	}

	return userInfos
}

func resourceDesktopUserBatchDetachCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	httpUrl := "v2/{project_id}/desktops/batch-detach"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildDesktopUserBatchDetachBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error detaching users from desktops: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", respBody, "")
	if jobId != nil && jobId.(string) != "" {
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

	return resourceDesktopUserBatchDetachRead(ctx, d, meta)
}

func resourceDesktopUserBatchDetachRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDesktopUserBatchDetachUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDesktopUserBatchDetachDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for detaching users from desktops. 
Deleting this resource will not undo the detachment, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
