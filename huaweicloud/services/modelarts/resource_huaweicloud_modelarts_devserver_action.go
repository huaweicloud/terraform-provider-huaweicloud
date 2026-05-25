package modelarts

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	devserverActionNonUpdatableParams = []string{
		"devserver_id",
		"action",
		"admin_pass",
		"key_pair_name",
		"image_id",
		"user_data",
	}
	actionCompletedMap = map[string][]string{
		"start":       {"RUNNING"},
		"stop":        {"STOPPED"},
		"reboot":      {"RUNNING"},
		"changeos":    {"RUNNING"},
		"reinstallos": {"RUNNING"},
	}
)

// @API ModelArts PUT /v1/{project_id}/dev-servers/{id}/start
// @API ModelArts PUT /v1/{project_id}/dev-servers/{id}/stop
// @API ModelArts PUT /v1/{project_id}/dev-servers/{id}/reboot
// @API ModelArts POST /v1/{project_id}/dev-servers/{id}/changeos
// @API ModelArts POST /v1/{project_id}/dev-servers/{id}/reinstallos
// @API ModelArts GET /v1/{project_id}/dev-servers/{id}
func ResourceDevServerAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDevServerActionCreate,
		ReadContext:   resourceDevServerActionRead,
		UpdateContext: resourceDevServerActionUpdate,
		DeleteContext: resourceDevServerActionDelete,

		CustomizeDiff: config.FlexibleForceNew(devserverActionNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			// Required parameters.
			"devserver_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the DevServer.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The action type of the DevServer.`,
			},

			// Optional parameters.
			"admin_pass": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `The login password of the DevServer.`,
			},
			"key_pair_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The key pair name used to log in to the DevServer.`,
			},
			"image_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The image ID used to change the operating system.`,
			},
			"user_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The user data to be injected into the DevServer during the OS operation.`,
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

func parseDevServerActionRequestMethod(action string) string {
	switch action {
	case "changeos", "reinstallos":
		return "POST"
	default:
		return "PUT"
	}
}

func buildDevServerActionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"admin_pass":    utils.ValueIgnoreEmpty(d.Get("admin_pass")),
		"key_pair_name": utils.ValueIgnoreEmpty(d.Get("key_pair_name")),
		"image_id":      utils.ValueIgnoreEmpty(d.Get("image_id")),
		"user_data":     utils.ValueIgnoreEmpty(d.Get("user_data")),
	}
}

func operateDevServerAction(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		devServerId = d.Get("devserver_id").(string)
		action      = d.Get("action").(string)
		httpUrl     = "v1/{project_id}/dev-servers/{id}/{action}"
	)

	actionPath := client.Endpoint + httpUrl
	actionPath = strings.ReplaceAll(actionPath, "{project_id}", client.ProjectID)
	actionPath = strings.ReplaceAll(actionPath, "{id}", devServerId)
	actionPath = strings.ReplaceAll(actionPath, "{action}", action)

	actionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	switch action {
	case "start":
		// For the "start" type, the request body parameter must be specified, so set an empty object for it.
		actionOpt.JSONBody = map[string]interface{}{}
	case "changeos", "reinstallos":
		actionOpt.JSONBody = utils.RemoveNil(buildDevServerActionBodyParams(d))
	}

	_, err := client.Request(parseDevServerActionRequestMethod(action), actionPath, &actionOpt)
	return err
}

func resourceDevServerActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		devServerId = d.Get("devserver_id").(string)
		action      = d.Get("action").(string)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	err = operateDevServerAction(client, d)
	if err != nil {
		return diag.Errorf("unable to %s DevServer (%s): %s", action, devServerId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshDevServerStatusFunc(client, devServerId, actionCompletedMap[action]),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for DevServer (%s) to %s completed: %s", devServerId, action, err)
	}

	d.SetId(devServerId)

	return resourceDevServerActionRead(ctx, d, meta)
}

func resourceDevServerActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDevServerActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDevServerActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for operating the DevServer. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
