package dataarts

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	factoryScriptExecuteNonUpdatableParams = []string{
		"workspace_id",
		"script_name",
		"params",
	}

	factoryScriptExecuteNotFoundCodes = []string{
		"DLF.20819", // The workspace does not exist.
		"DLF.6201",  // The script does not exist.
	}
)

// @API DataArtsStudio POST /v1/{project_id}/scripts/{script_name}/execute
// @API DataArtsStudio GET /v1/{project_id}/scripts/{script_name}/instances/{instance_id}
func ResourceFactoryScriptExecute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFactoryScriptExecuteCreate,
		ReadContext:   resourceFactoryScriptExecuteRead,
		UpdateContext: resourceFactoryScriptExecuteUpdate,
		DeleteContext: resourceFactoryScriptExecuteDelete,

		CustomizeDiff: config.FlexibleForceNew(factoryScriptExecuteNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the workspace is located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the script belongs.`,
			},
			"script_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the script to be executed.`,
			},

			// Optional parameters.
			"params": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
				Description:  `The execution parameters passed to the script content, in JSON format.`,
			},

			// Attributes.
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The execution status of the script instance.`,
			},
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The message when the script instance execution fails.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func buildFactoryScriptExecuteBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"params": utils.StringToJson(d.Get("params").(string)),
	}
}

func executeFactoryScript(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v1/{project_id}/scripts/{script_name}/execute"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{script_name}", d.Get("script_name").(string))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
		JSONBody:         utils.RemoveNil(buildFactoryScriptExecuteBodyParams(d)),
		OkCodes:          []int{200},
	}

	resp, err := client.Request("POST", path, &opt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(resp)
}

func GetScriptInstanceById(client *golangsdk.ServiceClient, workspaceId, scriptName, instanceId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/scripts/{script_name}/instances/{instance_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{script_name}", scriptName)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": workspaceId},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(getResp)
}

func scriptInstanceStateRefreshFunc(client *golangsdk.ServiceClient, workspaceId, scriptName, instanceId string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetScriptInstanceById(client, workspaceId, scriptName, instanceId)
		if err != nil {
			return nil, "ERROR", err
		}

		status := utils.PathSearch("status", respBody, "").(string)
		targetStatuses := []string{"FINISHED", "FAILED"}
		if utils.StrSliceContains(targetStatuses, status) {
			return respBody, "COMPLETED", nil
		}
		return respBody, "PENDING", nil
	}
}

func waitForScriptExecuteInstanceCompletion(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, instanceId string) error {
	var (
		workspaceId = d.Get("workspace_id").(string)
		scriptName  = d.Get("script_name").(string)
	)

	stateConf := &retry.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      scriptInstanceStateRefreshFunc(client, workspaceId, scriptName, instanceId),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceFactoryScriptExecuteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts-dlf", region)
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	scriptName := d.Get("script_name").(string)
	respBody, err := executeFactoryScript(client, d)
	if err != nil {
		return diag.Errorf("error executing DataArts script (%s): %s", scriptName, err)
	}

	instanceId := utils.PathSearch("instanceId", respBody, "").(string)
	if instanceId == "" {
		return diag.Errorf("unable to find the script execution instance ID from the API response")
	}
	d.SetId(instanceId)

	err = waitForScriptExecuteInstanceCompletion(ctx, d, client, instanceId)
	if err != nil {
		return diag.Errorf("error waiting for the script (%s) instance (%s) to finish: %s", scriptName, instanceId, err)
	}

	return resourceFactoryScriptExecuteRead(ctx, d, meta)
}

func resourceFactoryScriptExecuteRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		scriptName  = d.Get("script_name").(string)
		instanceId  = d.Id()
	)

	client, err := cfg.NewServiceClient("dataarts-dlf", region)
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	respBody, err := GetScriptInstanceById(client, workspaceId, scriptName, instanceId)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", factoryScriptExecuteNotFoundCodes...),
			fmt.Sprintf("error getting DataArts script execute instance (%s)", instanceId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("message", utils.PathSearch("message", respBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceFactoryScriptExecuteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceFactoryScriptExecuteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	msg := `This resource is only a one-time action resource for executing a script. Deleting this resource will not
remove the execution record (only when the script is deleted can the record be cleared), but will only remove the
resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  msg,
		},
	}
}
