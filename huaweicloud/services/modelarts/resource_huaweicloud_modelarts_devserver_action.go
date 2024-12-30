package modelarts

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts PUT /v1/{project_id}/dev-servers/{id}/start
// @API ModelArts PUT /v1/{project_id}/dev-servers/{id}/stop
// @API ModelArts GET /v1/{project_id}/dev-servers/{id}
func ResourceDevServerAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDevServerActionCreate,
		ReadContext:   resourceDevServerActionRead,
		DeleteContext: resourceDevServerActionDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"devserver_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the DevServer.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The action type of the DevServer.`,
			},
		},
	}
}

func resourceDevServerActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                = meta.(*config.Config)
		region             = cfg.GetRegion(d)
		httpUrl            = "v1/{project_id}/dev-servers/{id}/{action}"
		devServerId        = d.Get("devserver_id").(string)
		action             = d.Get("action").(string)
		actionCompletedMap = map[string]string{
			"start": "RUNNING",
			"stop":  "STOPPED",
		}
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	actionPath := client.Endpoint + httpUrl
	actionPath = strings.ReplaceAll(actionPath, "{project_id}", client.ProjectID)
	actionPath = strings.ReplaceAll(actionPath, "{id}", devServerId)
	actionPath = strings.ReplaceAll(actionPath, "{action}", action)
	actionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	if action == "start" {
		// For the "start" type, the request body parameter must be specified, so set an empty object for it.
		actionOpt.JSONBody = map[string]interface{}{}
	}

	_, err = client.Request("PUT", actionPath, &actionOpt)
	if err != nil {
		return diag.Errorf("unable to %s DevServer (%s): %s", action, devServerId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshDevServerActionStatusFunc(client, devServerId, actionCompletedMap[action]),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for DevServer (%s) to %s completed: %s", devServerId, action, err)
	}

	d.SetId(devServerId)

	return nil
}

func refreshDevServerActionStatusFunc(client *golangsdk.ServiceClient, devServerId string, target string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetDevServerById(client, devServerId)
		if err != nil {
			return nil, "ERROR", err
		}

		status := utils.PathSearch("status", respBody, "").(string)
		if utils.StrSliceContains([]string{"START_FAILED", "ERROR", "STOP_FAILED"}, status) {
			return respBody, "ERROR", fmt.Errorf("unexpected status (%s)", status)
		}

		if status == target {
			return respBody, "COMPLETED", nil
		}
		return "continue", "PENDING", nil
	}
}

func resourceDevServerActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
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
