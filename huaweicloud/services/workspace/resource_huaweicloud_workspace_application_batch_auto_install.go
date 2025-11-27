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

var applicationBatchAutoInstallNonUpdatableParams = []string{
	"app_ids",
	"assign_scope",
	"users",
}

// @API Workspace POST /v1/{project_id}/app-center/apps/actions/batch-auto-install
func ResourceApplicationBatchAutoInstall() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationBatchAutoInstallCreate,
		ReadContext:   resourceApplicationBatchAutoInstallRead,
		UpdateContext: resourceApplicationBatchAutoInstallUpdate,
		DeleteContext: resourceApplicationBatchAutoInstallDelete,

		CustomizeDiff: config.FlexibleForceNew(applicationBatchAutoInstallNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the applications are located.`,
			},

			// Required parameters.
			"app_ids": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				MaxItems:    50,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of application IDs to be automatically installed.`,
			},
			"assign_scope": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The assignment scope.`,
			},

			// Optional parameters.
			"users": {
				Type:        schema.TypeList,
				Optional:    true,
				MinItems:    1,
				MaxItems:    50,
				Elem:        applicationBatchAutoInstallUserSchema(),
				Description: `The list of users. Required when assign_scope is ASSIGN_USER.`,
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

func applicationBatchAutoInstallUserSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"account": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The account name. The account format must be: account(group).`,
			},
			"account_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The account type.`,
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The domain name. Required for user groups, and defaults to local.com if not specified.`,
			},
			"platform_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The platform type.`,
			},
		},
	}
}

func buildApplicationBatchAutoInstallUsersBodyParams(users []interface{}) []interface{} {
	if len(users) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(users))
	for _, user := range users {
		result = append(result, utils.RemoveNil(map[string]interface{}{
			"account":       utils.PathSearch("account", user, ""),
			"account_type":  utils.PathSearch("account_type", user, ""),
			"domain":        utils.ValueIgnoreEmpty(utils.PathSearch("domain", user, "")),
			"platform_type": utils.ValueIgnoreEmpty(utils.PathSearch("platform_type", user, "")),
		}))
	}

	return result
}

func buildApplicationBatchAutoInstallBodyParams(d *schema.ResourceData) map[string]interface{} {
	body := map[string]interface{}{
		"app_ids":      d.Get("app_ids").([]interface{}),
		"assign_scope": d.Get("assign_scope"),
	}

	if v, ok := d.GetOk("users"); ok {
		body["users"] = buildApplicationBatchAutoInstallUsersBodyParams(v.([]interface{}))
	}

	return body
}

func resourceApplicationBatchAutoInstallCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/app-center/apps/actions/batch-auto-install"
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
		JSONBody: buildApplicationBatchAutoInstallBodyParams(d),
	}

	_, err = client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating application batch auto install: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceApplicationBatchAutoInstallRead(ctx, d, meta)
}

func resourceApplicationBatchAutoInstallRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceApplicationBatchAutoInstallUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceApplicationBatchAutoInstallDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for batch auto installing applications. Deleting this resource will
	not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
