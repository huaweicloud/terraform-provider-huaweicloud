package dataarts

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var factoryJobImportNonUpdatableParams = []string{
	"path",
	"workspace_id",
	"target_status",
	"same_name_policy",
}

// @API DataArtsStudio POST /v1/{project_id}/jobs/import
// @API DataArtsStudio GET /v1/{project_id}/system-tasks/{task_id}
func ResourceFactoryJobImport() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFactoryJobImportCreate,
		ReadContext:   resourceFactoryJobImportRead,
		UpdateContext: resourceFactoryJobImportUpdate,
		DeleteContext: resourceFactoryJobImportDelete,

		CustomizeDiff: config.FlexibleForceNew(factoryJobImportNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the job is located.`,
			},

			// Required parameters.
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The OBS path where the job package is stored.`,
			},

			// Optional parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The workspace ID to which the imported jobs belong.`,
			},
			"target_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The target status of imported jobs.`,
			},
			"same_name_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The duplicate name handling policy.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{Internal: true,
						Required: true,
					}),
			},
		},
	}
}

func buildFactoryJobImportBodyParams(d *schema.ResourceData) map[string]interface{} {
	return utils.RemoveNil(map[string]interface{}{
		// Required parameters.
		"path": d.Get("path"),
		// Optional parameters.
		"sameNamePolicy": utils.ValueIgnoreEmpty(d.Get("same_name_policy")),
		"targetStatus":   utils.ValueIgnoreEmpty(d.Get("target_status")),
	})
}

func importFactoryJob(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	var (
		httpUrl     = "v1/{project_id}/jobs/import"
		workspaceId = d.Get("workspace_id").(string)
	)

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildFactoryJobExportRequestMoreHeaders(workspaceId),
		JSONBody:         utils.RemoveNil(buildFactoryJobImportBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func getFactoryJobImportSystemTask(client *golangsdk.ServiceClient, workspaceId, taskId string) (interface{}, error) {
	var httpUrl = "v1/{project_id}/system-tasks/{task_id}"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{task_id}", taskId)
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildFactoryJobExportRequestMoreHeaders(workspaceId),
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func factoryJobImportTaskStateRefreshFunc(client *golangsdk.ServiceClient, workspaceId, taskId string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		taskResp, err := getFactoryJobImportSystemTask(client, workspaceId, taskId)
		if err != nil {
			return nil, "ERROR", err
		}

		status := utils.PathSearch("status", taskResp, "").(string)
		if status == "FAILED" {
			return taskResp, "ERROR", fmt.Errorf("the import task (%s) is failed", taskId)
		}
		if status == "SUCCESSFUL" {
			return taskResp, "COMPLETED", nil
		}
		return taskResp, "PENDING", nil
	}
}

func resourceFactoryJobImportCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts-dlf", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	respBody, err := importFactoryJob(client, d)
	if err != nil {
		return diag.Errorf("error importing DataArts Factory jobs: %s", err)
	}

	taskId := utils.PathSearch("taskId", respBody, "").(string)
	if taskId == "" {
		return diag.Errorf("unable to find the import task ID from the API response")
	}

	workspaceId := d.Get("workspace_id").(string)
	stateConf := &retry.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      factoryJobImportTaskStateRefreshFunc(client, workspaceId, taskId),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the import task (%s) to become successful: %s", taskId, err)
	}

	d.SetId(taskId)

	return resourceFactoryJobImportRead(ctx, d, meta)
}

func resourceFactoryJobImportRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceFactoryJobImportUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceFactoryJobImportDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to import jobs from a specified OBS storage path.
Deleting this resource will not clear the import task, but will only remove the resource information from the tfstate
file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
