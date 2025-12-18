package workspace

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var batchActionHTTPMethodMap = map[string]string{
	"batch-change-image":  "POST",
	"batch-reinstall":     "POST",
	"batch-rejoin-domain": "PATCH",
	"batch-update-tsvi":   "PATCH",
	"batch-maint":         "PATCH",
	"batch-reboot":        "PATCH",
	"batch-start":         "PATCH",
	"batch-stop":          "PATCH",
}

var appServerBatchActionNonUpdatableParams = []string{"type", "content"}

// @API Workspace POST /v1/{project_id}/app-servers/actions/batch-change-image
// @API Workspace POST /v1/{project_id}/app-servers/actions/batch-reinstall
// @API Workspace PATCH /v1/{project_id}/app-servers/actions/batch-rejoin-domain
// @API Workspace PATCH /v1/{project_id}/app-servers/actions/batch-update-tsvi
// @API Workspace PATCH /v1/{project_id}/app-servers/actions/batch-maint
// @API Workspace PATCH /v1/{project_id}/app-servers/actions/batch-reboot
// @API Workspace PATCH /v1/{project_id}/app-servers/actions/batch-start
// @API Workspace PATCH /v1/{project_id}/app-servers/actions/batch-stop
// @API Workspace GET /v2/{project_id}/job/{job_id}
// @API Workspace GET /v1/{project_id}/app-servers
func ResourceAppServerBatchAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppServerBatchActionCreate,
		ReadContext:   resourceAppServerBatchActionRead,
		UpdateContext: resourceAppServerBatchActionUpdate,
		DeleteContext: resourceAppServerBatchActionDelete,

		CustomizeDiff: config.FlexibleForceNew(appServerBatchActionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the APP servers to be batch operated are located.`,
			},

			// Required parameter(s).
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The batch operation (action) type for the APP servers.`,
				ValidateFunc: validation.StringInSlice([]string{
					"batch-change-image",
					"batch-reinstall",
					"batch-rejoin-domain",
					"batch-update-tsvi",
					"batch-maint",
					"batch-reboot",
					"batch-start",
					"batch-stop",
				}, false),
			},
			"content": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
				Description:  `The JSON string content for the batch operation (action) request.`,
			},

			// Optional parameter(s).
			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The maximum number of retries for the batch operation (action) when encountering 409 conflict errors.`,
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

func resourceAppServerBatchActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		httpUrl       = "v1/{project_id}/app-servers/actions/{type}"
		actionType    = d.Get("type").(string)
		actionServers = utils.StringToJson(d.Get("content").(string))
		maxRetries    = d.Get("max_retries").(int)
		timeout       = d.Timeout(schema.TimeoutCreate)
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	httpMethod, exists := batchActionHTTPMethodMap[actionType]
	if !exists {
		return diag.Errorf("unsupported operation (action) type: %s", actionType)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{type}", actionType)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: actionServers,
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate resource ID: %s", err)
	}
	d.SetId(randUUID)

	var resp *http.Response
	for i := 0; i < maxRetries+1; i++ {
		resp, err = client.Request(httpMethod, createPath, &opt)
		if err == nil {
			break
		}

		if _, ok := err.(golangsdk.ErrDefault409); ok {
			// lintignore:R018
			time.Sleep(30 * time.Second)
			continue
		}
		if i < 1 {
			return diag.Errorf("error executing APP server batch operation (action: %s): %s", actionType, err)
		}
		return diag.Errorf("after %d retries, the APP server batch operation (action: %s) still reports an error: %s",
			i, actionType, err)
	}

	// When the operation type is `batch-rejoin-domain`, the server status is REGISTERED, but the task of joining the domain
	// is not actually completed. Therefore, we cannot use a query of the entire server list to check if the task is complete,
	// we use the task ID to check the status.
	if utils.StrSliceContains([]string{"batch-rejoin-domain", "batch-update-tsvi"}, actionType) {
		repsBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		jobId := utils.PathSearch("job_id", repsBody, "").(string)
		if jobId == "" {
			return diag.Errorf("unable to find job ID from API response")
		}

		_, err = waitForAppServerJobCompleted(ctx, client, timeout, jobId)
		if err != nil {
			return diag.Errorf("error waiting for the APP server batch operation (action: %s) job (%s) completed: %s",
				actionType, jobId, err)
		}

		return nil
	}

	// a. `batch-reinstall` and `batch-change-image` response item format: [{"job_id": "xxx", "server_id": "xxx"}].
	// b. `batch-maint`, `batch-reboot`, `batch-start`, and `batch-stop` not return job IDs.
	// For both of these cases, we uniformly use the API to query the list of all servers to check the status.
	if utils.StrSliceContains([]string{"batch-reinstall", "batch-change-image", "batch-maint", "batch-reboot",
		"batch-start", "batch-stop"}, actionType) {
		serverIds := utils.PathSearch("items", actionServers, make([]interface{}, 0)).([]interface{})
		if actionType == "batch-change-image" || actionType == "batch-reinstall" {
			serverIds = utils.PathSearch("server_ids", actionServers, make([]interface{}, 0)).([]interface{})
		}

		err = waitForAppServerBatchActionCompleted(ctx, client, utils.ExpandToStringList(serverIds), timeout)
		if err != nil {
			return diag.Errorf("error waiting for the APP server batch operation (action: %s) completed: %s", actionType, err)
		}

		return nil
	}

	return nil
}

func waitForAppServerBatchActionCompleted(ctx context.Context, client *golangsdk.ServiceClient, serverIds []string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshAppServerBatchActionStatusFunc(client, serverIds),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 15 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func refreshAppServerBatchActionStatusFunc(client *golangsdk.ServiceClient, serverIds []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		// Some action not return job ID, so we use query all server list API to check the status.
		servers, err := listAppServers(client)
		if err != nil {
			return servers, "ERROR", err
		}

		for _, server := range servers {
			serverId := utils.PathSearch("id", server, "").(string)
			if utils.StrSliceContains(serverIds, serverId) {
				status := utils.PathSearch("status", server, "").(string)
				if !utils.StrSliceContains([]string{"REGISTERED", "STOPPED", "FREEZE", "NONE"}, status) {
					return servers, "PENDING", nil
				}
			}
		}

		return servers, "COMPLETED", nil
	}
}

func resourceAppServerBatchActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppServerBatchActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppServerBatchActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to batch operate APP servers. Deleting this resource
will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
