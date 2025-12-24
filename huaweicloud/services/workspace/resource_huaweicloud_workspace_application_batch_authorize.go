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

var applicationBatchAuthorizeNonUpdatableParams = []string{
	"app_ids",
	"authorization_type",
	"del_users",
	"del_users.*.account",
	"del_users.*.account_type",
	"del_users.*.domain",
	"del_users.*.platform_type",
	"users",
	"users.*.account",
	"users.*.account_type",
	"users.*.domain",
	"users.*.platform_type",
}

// @API Workspace POST /v1/{project_id}/app-center/apps/actions/batch-assign-authorization
func ResourceApplicationBatchAuthorize() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationBatchAuthorizeCreate,
		ReadContext:   resourceApplicationBatchAuthorizeRead,
		UpdateContext: resourceApplicationBatchAuthorizeUpdate,
		DeleteContext: resourceApplicationBatchAuthorizeDelete,

		CustomizeDiff: config.FlexibleForceNew(applicationBatchAuthorizeNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the application batch authorization is located.`,
			},

			// Required parameters.
			"app_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of application IDs to be authorized or unauthorized.`,
			},
			"authorization_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The authorization type.`,
			},

			// Optional parameters.
			"users": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    100,
				Elem:        applicationBatchAuthorizeUserSchema(),
				Description: `The list of users to be authorized.`,
			},
			"del_users": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    100,
				Elem:        applicationBatchAuthorizeUserSchema(),
				Description: `The list of users to be unauthorized.`,
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

func applicationBatchAuthorizeUserSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"account": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The account name.`,
			},
			"account_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The account type. Valid values are "SIMPLE" and "USER_GROUP".`,
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The domain name. Required for user groups.`,
			},
			"platform_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The platform type. Valid values are "AD" and "LOCAL".`,
			},
		},
	}
}

func buildApplicationBatchAuthorizeUsersBodyParams(users []interface{}) []interface{} {
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

func buildApplicationBatchAuthorizeBodyParams(d *schema.ResourceData) map[string]interface{} {
	body := map[string]interface{}{
		"app_ids":            d.Get("app_ids").([]interface{}),
		"authorization_type": d.Get("authorization_type"),
	}

	if v, ok := d.GetOk("users"); ok {
		body["users"] = buildApplicationBatchAuthorizeUsersBodyParams(v.([]interface{}))
	}
	if v, ok := d.GetOk("del_users"); ok {
		body["del_users"] = buildApplicationBatchAuthorizeUsersBodyParams(v.([]interface{}))
	}

	return body
}

func resourceApplicationBatchAuthorizeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/app-center/apps/actions/batch-assign-authorization"
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
		JSONBody: buildApplicationBatchAuthorizeBodyParams(d),
	}

	_, err = client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating application batch authorization: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceApplicationBatchAuthorizeRead(ctx, d, meta)
}

func resourceApplicationBatchAuthorizeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceApplicationBatchAuthorizeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceApplicationBatchAuthorizeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for batch authorizing applications. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
