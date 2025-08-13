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

var actionHTTPMethodMap = map[string]string{
	"change-image": "POST",
	"reinstall":    "POST",
}

var appServerActionNonUpdatableParams = []string{"action_type", "server_id", "content"}

// @API Workspace POST /v1/{project_id}/app-servers/{server_id}/actions/change-image
// @API Workspace POST /v1/{project_id}/app-servers/{server_id}/actions/reinstall
func ResourceAppServerAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppServerActionCreate,
		ReadContext:   resourceAppServerActionRead,
		UpdateContext: resourceAppServerActionUpdate,
		DeleteContext: resourceAppServerActionDelete,

		CustomizeDiff: config.FlexibleForceNew(appServerActionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the APP server to be operated is located.`,
			},

			// Required parameter(s).
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The operation (action) type for the APP server.`,
				ValidateFunc: validation.StringInSlice([]string{
					"change-image",
					"reinstall",
				}, false),
			},
			"server_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the server to be operated.`,
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The JSON string content for the operation (action) request.`,
			},

			// Optional parameter(s).
			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The maximum number of retries for the operation (action) when encountering 409 conflict errors.`,
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

func resourceAppServerActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v1/{project_id}/app-servers/{server_id}/actions/{type}"
		actionType = d.Get("type").(string)
		serverId   = d.Get("server_id").(string)
		content    = d.Get("content").(string)
		maxRetries = d.Get("max_retries").(int)
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	httpMethod, exists := actionHTTPMethodMap[actionType]
	if !exists {
		return diag.Errorf("unsupported operation (action) type: %s", actionType)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{server_id}", serverId)
	createPath = strings.ReplaceAll(createPath, "{type}", actionType)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.StringToJson(content),
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate resource ID: %s", err)
	}
	d.SetId(randUUID)

	for i := 0; i < maxRetries+1; i++ {
		_, err = client.Request(httpMethod, createPath, &opt)
		if err == nil {
			break
		}

		if _, ok := err.(golangsdk.ErrDefault409); ok {
			// lintignore:R018
			time.Sleep(30 * time.Second)
			continue
		}
		if i < 1 {
			return diag.Errorf("error executing APP server operation (action: %s): %s", actionType, err)
		}
		return diag.Errorf("after %d retries, the APP server operation (action: %s) still reports an error: %s",
			i, actionType, err)
	}

	return resourceAppServerActionRead(ctx, d, meta)
}

func resourceAppServerActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppServerActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppServerActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to operate APP server. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
