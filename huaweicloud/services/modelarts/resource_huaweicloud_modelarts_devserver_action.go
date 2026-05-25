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
	}
)

// @API ModelArts PUT /v1/{project_id}/dev-servers/{id}/start
// @API ModelArts PUT /v1/{project_id}/dev-servers/{id}/stop
// @API ModelArts PUT /v1/{project_id}/dev-servers/{id}/reboot
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

func operateDevServerAction(client *golangsdk.ServiceClient, devServerId, action string) error {
	httpUrl := "v1/{project_id}/dev-servers/{id}/{action}"
	httpUrl = strings.ReplaceAll(httpUrl, "{project_id}", client.ProjectID)
	httpUrl = strings.ReplaceAll(httpUrl, "{id}", devServerId)
	httpUrl = strings.ReplaceAll(httpUrl, "{action}", action)

	actionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	if action == "start" {
		// For the "start" type, the request body parameter must be specified, so set an empty object for it.
		actionOpt.JSONBody = map[string]interface{}{}
	}

	_, err := client.Request("POST", httpUrl, &actionOpt)
	return err
}

func resourceDevServerActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                = meta.(*config.Config)
		region             = cfg.GetRegion(d)
		devServerId        = d.Get("devserver_id").(string)
		action             = d.Get("action").(string)
		actionCompletedMap = map[string][]string{
			"start":  {"RUNNING"},
			"stop":   {"STOPPED"},
			"reboot": {"RUNNING"},
		}
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	err = operateDevServerAction(client, devServerId, action)
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
