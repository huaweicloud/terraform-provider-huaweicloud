package css

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}/logs/settings
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/logs/collect
// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}/logs/records
func ResourceManualLogBackup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceManualLogBackupCreate,
		ReadContext:   resourceManualLogBackupRead,
		DeleteContext: resourceManualLogBackupDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the region in which to create the resource.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the CSS cluster.`,
			},
			"job_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the log backup job.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the log backup job.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the log backup job.`,
			},
			"log_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The storage path of backed up logs in the OBS bucket.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"finished_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The end time.`,
			},
			"failed_msg": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The error information.`,
			},
		},
	}
}

func resourceManualLogBackupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	clusterID := d.Get("cluster_id").(string)
	client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	// Check whether opening log function.
	err = checkOpeningLogBackupFunc(client, clusterID)
	if err != nil {
		return diag.FromErr(err)
	}

	// manual backup logs
	err = createLogBackup(client, clusterID)
	if err != nil {
		return diag.FromErr(err)
	}
	// There can be at most one log backup task for the same period.
	// So we get the created job ID through the "RUNNING" status.
	expression := "clusterLogRecord[?status=='RUNNING'] | [0]"
	resp, err := getLogBackupJob(client, clusterID, expression)
	if err != nil {
		return diag.FromErr(err)
	}
	id := utils.PathSearch("jobId", resp, "").(string)
	if id == "" {
		return diag.Errorf("not found job ID of the manual log backup, resp: %v", resp)
	}
	d.SetId(id)

	// Check whether manual backup log has completed.
	err = checkLogBackupCompleted(ctx, client, clusterID, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceManualLogBackupRead(ctx, d, meta)
}

func resourceManualLogBackupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	clusterID := d.Get("cluster_id").(string)
	client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	expression := fmt.Sprintf("clusterLogRecord[?jobId=='%s'] | [0]", d.Id())
	resp, err := getLogBackupJob(client, clusterID, expression)
	if err != nil {
		return diag.FromErr(err)
	}

	createdAt := utils.PathSearch("createAt", resp, float64(0)).(float64) / 1000
	finishedAt := utils.PathSearch("finishedAt", resp, float64(0)).(float64) / 1000
	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("job_id", utils.PathSearch("jobId", resp, nil)),
		d.Set("type", utils.PathSearch("jobTypes", resp, nil)),
		d.Set("status", utils.PathSearch("status", resp, nil)),
		d.Set("log_path", utils.PathSearch("logPath", resp, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(createdAt), false)),
		d.Set("finished_at", utils.FormatTimeStampRFC3339(int64(finishedAt), false)),
		d.Set("failed_msg", utils.PathSearch("failedMsg", resp, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceManualLogBackupDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting manual log backup resource is not supported. The manual log backup resource" +
		"is only removed from the state, the manual log backup content and record remain in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func checkOpeningLogBackupFunc(client *golangsdk.ServiceClient, clusterID string) error {
	getLogSettingHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/logs/settings"
	getLogSettingPath := client.Endpoint + getLogSettingHttpUrl
	getLogSettingPath = strings.ReplaceAll(getLogSettingPath, "{project_id}", client.ProjectID)
	getLogSettingPath = strings.ReplaceAll(getLogSettingPath, "{cluster_id}", clusterID)

	getLogSettingPathOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getLogSettingResp, err := client.Request("GET", getLogSettingPath, &getLogSettingPathOpt)
	if err != nil {
		return err
	}

	getLogSettingRespBody, err := utils.FlattenResponse(getLogSettingResp)
	if err != nil {
		return err
	}

	logSetting := utils.PathSearch("logConfiguration", getLogSettingRespBody, nil)
	logSwitch := utils.PathSearch("logSwitch", logSetting, false).(bool)

	if logSetting == nil || !logSwitch {
		return fmt.Errorf("before backing up logs, please enable the log backup function")
	}

	return nil
}

func createLogBackup(client *golangsdk.ServiceClient, clusterID string) error {
	createLogBackupHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/logs/collect"
	createLogBackupPath := client.Endpoint + createLogBackupHttpUrl
	createLogBackupPath = strings.ReplaceAll(createLogBackupPath, "{project_id}", client.ProjectID)
	createLogBackupPath = strings.ReplaceAll(createLogBackupPath, "{cluster_id}", clusterID)

	createLogBackupPathOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err := client.Request("POST", createLogBackupPath, &createLogBackupPathOpt)
	if err != nil {
		return fmt.Errorf("error creating manual backup log, err: %s", err)
	}

	return nil
}

func checkLogBackupCompleted(ctx context.Context, client *golangsdk.ServiceClient, clusterId, id string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"RUNNING"},
		Target:       []string{"COMPLETED"},
		Refresh:      logBackupStateRefreshFunc(client, clusterId, id),
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for CSS (%s) to be extend: %s", clusterId, err)
	}
	return nil
}

func logBackupStateRefreshFunc(client *golangsdk.ServiceClient, clusterID, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		expression := fmt.Sprintf("clusterLogRecord[?jobId=='%s'] | [0]", id)
		resp, err := getLogBackupJob(client, clusterID, expression)
		if err != nil {
			return resp, "ERROR", err
		}

		status := utils.PathSearch("status", resp, "").(string)
		if status == "RUNNING" {
			return resp, "RUNNING", nil
		}

		if status == "FAIL" {
			failedMsg := utils.PathSearch("failedMsg", resp, "").(string)
			return resp, "FAIL", fmt.Errorf("manual backup log failed: %s", failedMsg)
		}

		return resp, "COMPLETED", nil
	}
}

func getLogBackupJob(client *golangsdk.ServiceClient, clusterID, expression string) (interface{}, error) {
	getLogBackupJobHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/logs/records"
	getLogBackupJobPath := client.Endpoint + getLogBackupJobHttpUrl
	getLogBackupJobPath = strings.ReplaceAll(getLogBackupJobPath, "{project_id}", client.ProjectID)
	getLogBackupJobPath = strings.ReplaceAll(getLogBackupJobPath, "{cluster_id}", clusterID)

	getLogBackupJobPathOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	currentTotal := 1
	for {
		currentPath := fmt.Sprintf("%s?limit=10&start=%d", getLogBackupJobPath, currentTotal)
		getLogBackupJobResp, err := client.Request("GET", currentPath, &getLogBackupJobPathOpt)
		if err != nil {
			return getLogBackupJobResp, err
		}
		getLogBackupJobRespBody, err := utils.FlattenResponse(getLogBackupJobResp)
		if err != nil {
			return getLogBackupJobRespBody, err
		}
		logBackupJob := utils.PathSearch(expression, getLogBackupJobRespBody, nil)
		if logBackupJob != nil {
			return logBackupJob, nil
		}
		logBackupJobs := utils.PathSearch("clusterLogRecord",
			getLogBackupJobRespBody, make([]interface{}, 0)).([]interface{})
		if len(logBackupJobs) < 10 {
			break
		}
		currentTotal += len(logBackupJobs)
	}

	return nil, golangsdk.ErrDefault404{}
}
