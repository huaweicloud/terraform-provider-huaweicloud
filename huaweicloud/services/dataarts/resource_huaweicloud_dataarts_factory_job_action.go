package dataarts

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var actionResourceNotFoundCodes = "DLF.0819" // The workspace ID does not exist.

// @API DataArtsStudio POST /v1/{project_id}/jobs/{job_name}/start
// @API DataArtsStudio POST /v1/{project_id}/jobs/{job_name}/stop
// @API DataArtsStudio GET /v1/{project_id}/jobs
func ResourceFactoryJobAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFactoryJobActionCreate,
		ReadContext:   resourceFactoryJobActionRead,
		UpdateContext: resourceFactoryJobActionUpdate,
		DeleteContext: resourceFactoryJobActionDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specified the ID of the workspace.`,
			},
			"job_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specified the name of the job.`,
			},
			"process_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specified the type of the job.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specified the action type of the job.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current status of the job.`,
			},
		},
	}
}

func resourceFactoryJobActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dataarts-dlf", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	jobName := d.Get("job_name").(string)
	workspaceId := d.Get("workspace_id").(string)
	_, err = doActionJob(client, workspaceId, jobName, d.Get("action").(string))
	if err != nil {
		return diag.Errorf("unable to operate status of the job (%s): %s", jobName, err)
	}

	d.SetId(jobName)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      jobStateRefreshFunc(client, workspaceId, jobName, d.Get("process_type").(string)),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the status of job (%s) to become running: %s", jobName, err)
	}
	return resourceFactoryJobActionRead(ctx, d, meta)
}

func doActionJob(client *golangsdk.ServiceClient, workspaceId, jobName, action string) (*http.Response, error) {
	httpUrl := "v1/{project_id}/jobs/{job_name}/{action}"
	actionPath := client.Endpoint + httpUrl
	actionPath = strings.ReplaceAll(actionPath, "{project_id}", client.ProjectID)
	actionPath = strings.ReplaceAll(actionPath, "{job_name}", jobName)
	actionPath = strings.ReplaceAll(actionPath, "{action}", action)

	actionOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"workspace": workspaceId,
		},
		OkCodes: []int{
			204,
		},
	}

	return client.Request("POST", actionPath, &actionOpts)
}

func jobStateRefreshFunc(client *golangsdk.ServiceClient, workspaceId, jobName, jobType string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := getJobByName(client, workspaceId, jobName, jobType)
		if err != nil {
			return respBody, "ERROR", err
		}

		statusResp := utils.PathSearch("status", respBody, "").(string)
		if utils.StrSliceContains([]string{"EXCEPTION"}, statusResp) {
			return respBody, "ERROR", fmt.Errorf("unexpect status (%s)", statusResp)
		}

		if utils.StrSliceContains([]string{"NORMAL", "STOPPED", "SCHEDULING"}, statusResp) {
			return respBody, "COMPLETED", nil
		}
		// Intermediate state: "STARTING" "STOPPING"
		return respBody, "PENDING", nil
	}
}

func getJobByName(client *golangsdk.ServiceClient, workspaceId, jobName, jobType string) (interface{}, error) {
	// The maximum value of limit is 100.
	httpUrl := "v1/{project_id}/jobs?limit=100"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	// The job_name field supports fuzzy matching.
	getPath = fmt.Sprintf("%s&jobName=%v&jobType=%v", getPath, jobName, jobType)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"workspace": workspaceId,
		},
	}

	offset := 0
	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", getPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &getOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		jobs := utils.PathSearch("jobs", respBody, make([]interface{}, 0)).([]interface{})
		if len(jobs) < 1 {
			break
		}

		job := utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0]", jobName), jobs, nil)
		if job != nil {
			return job, nil
		}
		offset += len(jobs)
	}

	return nil, golangsdk.ErrDefault404{}
}

func resourceFactoryJobActionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("dataarts-dlf", region)
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	jobName := d.Get("job_name").(string)
	job, err := getJobByName(client, d.Get("workspace_id").(string), jobName, d.Get("process_type").(string))
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", actionResourceNotFoundCodes),
			fmt.Sprintf("job (%s) not found: %s", jobName, err))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("job_name", utils.PathSearch("name", job, nil)),
		d.Set("status", utils.PathSearch("status", job, nil)),
		d.Set("action", paserAction(utils.PathSearch("status", job, "").(string))),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving the fields of the job action: %s", err)
	}
	return nil
}

func paserAction(action string) string {
	if utils.StrSliceContains([]string{"NORMAL", "SCHEDULING"}, action) {
		return "start"
	}
	return "stop"
}

func resourceFactoryJobActionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dataarts-dlf", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	if d.HasChange("action") {
		jobName := d.Get("job_name").(string)
		workspaceId := d.Get("workspace_id").(string)
		_, err = doActionJob(client, workspaceId, jobName, d.Get("action").(string))
		if err != nil {
			return diag.Errorf("error updating DataArts job (%s) status: %s", jobName, err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:      []string{"PENDING"},
			Target:       []string{"COMPLETED"},
			Refresh:      jobStateRefreshFunc(client, workspaceId, jobName, d.Get("process_type").(string)),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        10 * time.Second,
			PollInterval: 20 * time.Second,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for the status of job (%s) to become running: %s", jobName, err)
		}
	}
	return resourceFactoryJobActionRead(ctx, d, meta)
}

func resourceFactoryJobActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for changing job status. Deleting this resource will
not change the current status, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
