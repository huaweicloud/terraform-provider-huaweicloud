package dcs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var instanceExpiredKeyScanTaskNonUpdatableParams = []string{"instance_id"}

// @API DCS POST /v2/{project_id}/instances/{instance_id}/scan-expire-keys-task
// @API DCS GET /v2/{project_id}/instances/{instance_id}/auto-expire/histories
func ResourceDcsInstanceExpiredKeyScanTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcsInstanceExpiredKeyScanTaskCreate,
		ReadContext:   resourceDcsInstanceExpiredKeyScanTaskRead,
		UpdateContext: resourceDcsInstanceExpiredKeyScanTaskUpdate,
		DeleteContext: resourceDcsInstanceExpiredKeyScanTaskDelete,

		CustomizeDiff: config.FlexibleForceNew(instanceExpiredKeyScanTaskNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scan_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"started_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"finished_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDcsInstanceExpiredKeyScanTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}/scan-expire-keys-task"
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DCS instance(%s) expired key scan task: %s", instanceId, err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}
	taskId := utils.PathSearch("id", createRespBody, "").(string)
	if taskId == "" {
		return diag.Errorf("error creating DCS instance(%s) expired key scan task: id is not found in API response", instanceId)
	}

	d.SetId(taskId)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"pending"},
		Target:       []string{"success"},
		Refresh:      instanceExpiredKeyScanTaskRefreshFunc(client, instanceId, taskId),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for expired key scan task(%s) to be completed: %s ", taskId, err)
	}
	return resourceDcsInstanceExpiredKeyScanTaskRead(ctx, d, meta)
}

func instanceExpiredKeyScanTaskRefreshFunc(client *golangsdk.ServiceClient, instanceId, taskId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getRespBody, err := getInstanceExpiredKeyScanTaskHistories(client, instanceId)
		if err != nil {
			return nil, "error", err
		}
		task := utils.PathSearch(fmt.Sprintf("records[?id=='%s']|[0]", taskId), getRespBody, nil)
		if task == nil {
			return nil, "error", golangsdk.ErrDefault404{}
		}

		status := utils.PathSearch("status", task, "").(string)
		if status == "success" || status == "failed" {
			return getRespBody, status, nil
		}

		return getRespBody, "pending", nil
	}
}

func getInstanceExpiredKeyScanTaskHistories(client *golangsdk.ServiceClient, instanceId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/instances/{instance_id}/auto-expire/histories"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	getResp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, err
	}

	getRespJson, err := json.Marshal(getResp)
	if err != nil {
		return nil, err
	}
	var getRespBody interface{}
	err = json.Unmarshal(getRespJson, &getRespBody)
	if err != nil {
		return nil, err
	}
	return getRespBody, nil
}

func resourceDcsInstanceExpiredKeyScanTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	getRespBody, err := getInstanceExpiredKeyScanTaskHistories(client, d.Get("instance_id").(string))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting DCS instance expired key scan task")
	}
	task := utils.PathSearch(fmt.Sprintf("records[?id=='%s']|[0]", d.Id()), getRespBody, nil)
	if task == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error getting DCS instance expired key scan task")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instance_id", utils.PathSearch("instance_id", task, nil)),
		d.Set("status", utils.PathSearch("status", task, nil)),
		d.Set("scan_type", utils.PathSearch("scan_type", task, nil)),
		d.Set("num", utils.PathSearch("num", task, nil)),
		d.Set("created_at", utils.PathSearch("created_at", task, nil)),
		d.Set("started_at", utils.PathSearch("started_at", task, nil)),
		d.Set("finished_at", utils.PathSearch("finished_at", task, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDcsInstanceExpiredKeyScanTaskUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDcsInstanceExpiredKeyScanTaskDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DCS instance expired key scan task resource is not supported. The resource is only removed " +
		"from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
