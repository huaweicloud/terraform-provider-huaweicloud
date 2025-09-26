package modelarts

import (
	"context"
	"fmt"
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

var v2NodeBatchRebootNonUpdatableParams = []string{
	"resource_pool_name",
	"node_names",
}

// @API ModelArts POST /v2/{project_id}/pools/{pool_name}/nodes/batch-reboot
// @API ModelArts GET /v2/{project_id}/jobs/{job_id}
func ResourceV2NodeBatchReboot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2NodeBatchRebootCreate,
		ReadContext:   resourceV2NodeBatchRebootRead,
		UpdateContext: resourceV2NodeBatchRebootUpdate,
		DeleteContext: resourceV2NodeBatchRebootDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(v2NodeBatchRebootNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the resource nodes are located.`,
			},

			// Required parameters.
			"resource_pool_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resource pool name to which the resource nodes belong.`,
			},
			"node_names": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The name list of resource nodes to be rebooted.`,
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

func getV2JobById(client *golangsdk.ServiceClient, jobId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/jobs/{job_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{job_id}", jobId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func refreshV2JobStatus(client *golangsdk.ServiceClient, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		nodes, err := getV2JobById(client, jobId)
		if err != nil {
			return nil, "ERROR", err
		}
		status := utils.PathSearch("phase", nodes, "").(string)
		if utils.IsStrContainsSliceElement(status, []string{"Failed"}, false, true) {
			return nodes, "ERROR", fmt.Errorf("unexpected status: %s", status)
		}
		if status == "Success" {
			return nodes, "COMPLETED", nil
		}
		return nodes, "PENDING", nil
	}
}

func waitForV2JobCompleted(ctx context.Context, client *golangsdk.ServiceClient, jobId string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshV2JobStatus(client, jobId),
		Timeout:      timeout,
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func rebootV2ResourcePoolNodes(ctx context.Context, client *golangsdk.ServiceClient, resourcePoolName string,
	nodeNames []interface{}, timeout time.Duration) error {
	httpUrl := "v2/{project_id}/pools/{pool_name}/nodes/batch-reboot"
	actionPath := client.Endpoint + httpUrl
	actionPath = strings.ReplaceAll(actionPath, "{project_id}", client.ProjectID)
	actionPath = strings.ReplaceAll(actionPath, "{pool_name}", resourcePoolName)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"nodeNames": nodeNames,
		},
	}

	requestResp, err := client.Request("POST", actionPath, &opt)
	if err != nil {
		return fmt.Errorf("error executing batch reboot operation for the specified nodes (%v): %s", nodeNames, err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return err
	}
	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return fmt.Errorf("unable to find job ID under the resource pool (%s) in API response", resourcePoolName)
	}
	err = waitForV2JobCompleted(ctx, client, jobId, timeout)
	if err != nil {
		return fmt.Errorf("error waiting for the job status of resource pool (%s) creation to complete: %s", resourcePoolName, err)
	}
	return nil
}

func resourceV2NodeBatchRebootCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		resourcePoolName = d.Get("resource_pool_name").(string)
		nodeNames        = d.Get("node_names").([]interface{})
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	err = rebootV2ResourcePoolNodes(ctx, client, resourcePoolName, nodeNames, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	return resourceV2NodeBatchRebootRead(ctx, d, meta)
}

func resourceV2NodeBatchRebootRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV2NodeBatchRebootUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV2NodeBatchRebootDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for batch reboot the ModelArts nodes. Deleting this
resource will not clear the corresponding request record, but will only remove the resource information from the tfstate
file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
