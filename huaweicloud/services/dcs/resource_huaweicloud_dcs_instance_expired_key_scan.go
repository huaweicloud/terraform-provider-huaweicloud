package dcs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var dcsInstanceExpiredKeyScanNonUpdatableParams = []string{
	"instance_id",
}

// @API DCS POST /v2/{project_id}/instances/{instance_id}/auto-expire/scan
// @API DCS GET /v2/{project_id}/instances/{instance_id}/auto-expire/histories
func ResourceDcsInstanceExpiredKeyScan() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcsInstanceExpiredKeyScanCreate,
		UpdateContext: resourceDcsInstanceExpiredKeyScanUpdate,
		ReadContext:   resourceDcsInstanceExpiredKeyScanRead,
		DeleteContext: resourceDcsInstanceExpiredKeyScanDelete,

		CustomizeDiff: config.FlexibleForceNew(dcsInstanceExpiredKeyScanNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceDcsInstanceExpiredKeyScanCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}/auto-expire/scan"
		product = "dcs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error scanning DCS expire key: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := utils.PathSearch("id", createRespBody, "").(string)
	if taskID == "" {
		return diag.Errorf("error creating DCS instance expired key: task ID is not found in API response")
	}

	d.SetId(taskID)

	if err = waitForExpiredKeyScanTaskComplete(ctx, client, instanceID, taskID, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for DCS instance expired key scan task (%s) to complete: %s", taskID, err)
	}

	return resourceDcsInstanceExpiredKeyScanRead(ctx, d, meta)
}

func resourceDcsInstanceExpiredKeyScanUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDcsInstanceExpiredKeyScanRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "dcs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	record, err := getInstanceExpiredKeyScanTask(client, d.Get("instance_id").(string), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DCS expire key scan records")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instance_id", utils.PathSearch("instance_id", record, nil)),
		d.Set("status", utils.PathSearch("status", record, nil)),
		d.Set("scan_type", utils.PathSearch("scan_type", record, nil)),
		d.Set("num", utils.PathSearch("num", record, nil)),
		d.Set("created_at", utils.PathSearch("created_at", record, nil)),
		d.Set("started_at", utils.PathSearch("started_at", record, nil)),
		d.Set("finished_at", utils.PathSearch("finished_at", record, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDcsInstanceExpiredKeyScanDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DCS instance expired key scan resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func waitForExpiredKeyScanTaskComplete(ctx context.Context, client *golangsdk.ServiceClient, instanceID, taskID string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"pending"},
		Target:       []string{"success"},
		Refresh:      instanceExpiredKeyScanTaskStatusRefreshFunc(client, instanceID, taskID),
		Timeout:      timeout,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for DCS instance expired key scan task (%s) to complete: %s", taskID, err)
	}
	return nil
}

func instanceExpiredKeyScanTaskStatusRefreshFunc(client *golangsdk.ServiceClient, instanceID, taskID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		record, err := getInstanceExpiredKeyScanTask(client, instanceID, taskID)
		if err != nil {
			return nil, "failed", err
		}

		status := utils.PathSearch("status", record, "").(string)
		if status == "failed" {
			return nil, "", errors.New("unable to get DCS instance expired key scan task status")
		}
		if status == "success" {
			return record, "success", nil
		}

		return record, "pending", nil
	}
}

func getInstanceExpiredKeyScanTask(client *golangsdk.ServiceClient, instanceID, taskID string) (interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}/auto-expire/histories"
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceID)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{})
	if err != nil {
		return nil, err
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return nil, err
	}

	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return nil, err
	}

	record := utils.PathSearch(fmt.Sprintf("records[?id=='%s']|[0]", taskID), listRespBody, nil)
	if record == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return record, nil
}
