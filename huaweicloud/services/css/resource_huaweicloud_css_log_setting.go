package css

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/logs/open
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/logs/connectivity
// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}/logs/settings
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/logs/settings
// @API CSS PUT /v1.0/{project_id}/clusters/{cluster_id}/logs/close
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/logs/policy/update
// @API CSS PUT /v1.0/{project_id}/clusters/{cluster_id}/logs/policy/close
func ResourceLogSetting() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLogSettingCreate,
		ReadContext:   resourceLogSettingRead,
		UpdateContext: resourceLogSettingUpdate,
		DeleteContext: resourceLogSettingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceLogSettingImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"action": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "base_log_collect",
				ValidateFunc: validation.StringInSlice([]string{"base_log_collect", "real_time_log_collect"}, false),
			},
			"agency": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"base_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bucket": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"period": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"index_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"keep_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 3650),
			},
			"target_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_switch": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"auto_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_ingestion_create_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_ingestion_update_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLogSettingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	clusterID := d.Get("cluster_id").(string)
	cssClient, err := conf.NewServiceClient("css", region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	action := d.Get("action").(string)
	switch action {
	case "real_time_log_collect":
		err := openLogIngestion(ctx, d, cssClient)
		if err != nil {
			return diag.Errorf("error opening CSS cluster (%s) log ingestion: %s", clusterID, err)
		}
	default:
		err = openLogBackup(ctx, d, cssClient)
		if err != nil {
			return diag.Errorf("error opening CSS cluster (%s) log backup: %s", clusterID, err)
		}

		if _, ok := d.GetOk("period"); ok {
			err = openLogAutoBackup(d, cssClient)
			if err != nil {
				return diag.Errorf("error opening CSS cluster (%s) log auto backup: %s", clusterID, err)
			}
		}
	}

	d.SetId(fmt.Sprintf("%s/%s", clusterID, action))
	return resourceLogSettingRead(ctx, d, meta)
}

func resourceLogSettingRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	cssClient, err := conf.NewServiceClient("css", region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	action := d.Get("action").(string)
	clusterID := d.Get("cluster_id").(string)

	switch action {
	case "real_time_log_collect":
		logIngestionSetting, err := getLogIngestionSetting(clusterID, cssClient)
		if err != nil {
			// The cluster does not exist, http code is 403, key/value of error code is errCode/CSS.0015.
			err = common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015")
			return common.CheckDeletedDiag(d, err, "erorr retrieving CSS cluster log ingestion setting")
		}
		if logIngestionSetting == nil {
			errorMsg := fmt.Sprintf("CSS cluster (%s) log ingestion is closed", clusterID)
			return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, errorMsg)
		}

		realTimeLogCollectCreateAt := utils.PathSearch("createAt", logIngestionSetting, float64(0)).(float64)
		realTimeLogCollectUpdateAt := utils.PathSearch("updateAt", logIngestionSetting, float64(0)).(float64)

		mErr := multierror.Append(nil,
			d.Set("region", region),
			d.Set("cluster_id", clusterID),
			d.Set("action", "real_time_log_collect"),
			d.Set("index_prefix", utils.PathSearch("indexPrefix", logIngestionSetting, nil)),
			d.Set("keep_days", utils.PathSearch("keepDays", logIngestionSetting, nil)),
			d.Set("target_cluster_id", utils.PathSearch("targetClusterId", logIngestionSetting, nil)),
			d.Set("status", utils.PathSearch("status", logIngestionSetting, nil)),
			d.Set("log_ingestion_create_at", utils.FormatTimeStampRFC3339(int64(realTimeLogCollectCreateAt)/1000, false)),
			d.Set("log_ingestion_update_at", utils.FormatTimeStampRFC3339(int64(realTimeLogCollectUpdateAt)/1000, false)),
		)

		return diag.FromErr(mErr.ErrorOrNil())
	default:
		logBackupSetting, err := getLogBackupSetting(clusterID, cssClient)
		if err != nil {
			// The cluster does not exist, http code is 403, key/value of error code is errCode/CSS.0015.
			err = common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015")
			return common.CheckDeletedDiag(d, err, "erorr retrieving CSS cluster log ingestion setting")
		}
		logBackupSwitch := utils.PathSearch("logSwitch", logBackupSetting, false).(bool)
		if !logBackupSwitch {
			errorMsg := fmt.Sprintf("CSS cluster (%s) log backup is closed", clusterID)
			return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, errorMsg)
		}

		updateAt := utils.PathSearch("updateAt", logBackupSetting, float64(0)).(float64)
		autoEnabled := utils.PathSearch("autoEnable", logBackupSetting, false).(bool)
		var period string
		if autoEnabled {
			period = utils.PathSearch("period", logBackupSetting, false).(string)
		}

		mErr := multierror.Append(
			nil,
			d.Set("region", region),
			d.Set("cluster_id", clusterID),
			d.Set("action", "base_log_collect"),
			d.Set("agency", utils.PathSearch("agency", logBackupSetting, nil)),
			d.Set("base_path", utils.PathSearch("basePath", logBackupSetting, nil)),
			d.Set("bucket", utils.PathSearch("obsBucket", logBackupSetting, nil)),
			d.Set("updated_at", utils.FormatTimeStampRFC3339(int64(updateAt)/1000, false)),
			d.Set("auto_enabled", autoEnabled),
			d.Set("period", period),
			d.Set("log_switch", logBackupSwitch),
		)

		return diag.FromErr(mErr.ErrorOrNil())
	}
}

func resourceLogSettingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	cssClient, err := conf.NewServiceClient("css", region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	clusterID := d.Get("cluster_id").(string)
	action := d.Get("action").(string)

	switch action {
	case "real_time_log_collect":
		logIngestionSettingChanges := []string{
			"index_prefix",
			"keep_days",
			"target_cluster_id",
		}
		if d.HasChanges(logIngestionSettingChanges...) {
			err := updateLogIngestionSetting(ctx, d, cssClient)
			if err != nil {
				return diag.Errorf("error updating CSS cluster (%s) log ingestion setting: %s", clusterID, err)
			}
		}
		return resourceLogSettingRead(ctx, d, meta)
	default:
		logBackupSettingChanges := []string{
			"agency",
			"base_path",
			"bucket",
		}

		if d.HasChanges(logBackupSettingChanges...) {
			err := updateLogBackupSetting(d, cssClient)
			if err != nil {
				return diag.Errorf("error updating CSS cluster (%s) log backup setting: %s", clusterID, err)
			}
		}

		if d.HasChange("period") {
			period := d.Get("period").(string)
			var err error
			if period != "" {
				err = openLogAutoBackup(d, cssClient)
				if err != nil {
					return diag.Errorf("error updating CSS cluster (%s) log auto backup setting: %s", clusterID, err)
				}
			} else {
				err = closeLogAutoBackup(clusterID, cssClient)
				if err != nil {
					return diag.Errorf("error closing CSS cluster (%s) log auto backup: %s", clusterID, err)
				}
			}
		}
		return resourceLogSettingRead(ctx, d, meta)
	}
}

func resourceLogSettingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	cssClient, err := conf.NewServiceClient("css", region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}
	clusterID := d.Get("cluster_id").(string)
	action := d.Get("action").(string)

	switch action {
	case "real_time_log_collect":
		err := closeLogIngestion(ctx, d, cssClient)
		if err != nil {
			// "CSS.0015": The cluster does not exist. Status code is 403.
			err = common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015")
			// "CSS.0004": Invalid operation. Status code is 415.
			// {"errCode":"CSS.0004","externalMessage":"CSS.0004 : Invalid operation. (Illegal operation)"}
			err = common.ConvertUndefinedErrInto404Err(err, 415, "errCode", "CSS.0004")
			return common.CheckDeletedDiag(d, err, fmt.Sprintf("error close the CSS cluser (%s) log ingestion", clusterID))
		}
		return nil
	default:
		err := closeLogBackup(clusterID, cssClient)
		if err != nil {
			// "CSS.0015": The cluster does not exist. Status code is 403.
			err = common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015")
			// "CSS.0004": Invalid operation. Status code is 415.
			// {"errCode":"CSS.0004","externalMessage":"CSS.0004 : Invalid operation. (Illegal operation)"}
			err = common.ConvertUndefinedErrInto404Err(err, 415, "errCode", "CSS.0004")
			return common.CheckDeletedDiag(d, err, fmt.Sprintf("error close the CSS cluser (%s) log backup", clusterID))
		}
		return nil
	}
}

func openLogBackup(_ context.Context, d *schema.ResourceData, cssClient *golangsdk.ServiceClient) error {
	clusterID := d.Get("cluster_id").(string)
	openLogBackupHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/logs/open"
	openLogBackupPath := cssClient.Endpoint + openLogBackupHttpUrl
	openLogBackupPath = strings.ReplaceAll(openLogBackupPath, "{project_id}", cssClient.ProjectID)
	openLogBackupPath = strings.ReplaceAll(openLogBackupPath, "{cluster_id}", clusterID)

	openLogBackupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	openLogBackupOpt.JSONBody = map[string]interface{}{
		"agency":      d.Get("agency").(string),
		"logBasePath": d.Get("base_path").(string),
		"logBucket":   d.Get("bucket").(string),
	}
	_, err := cssClient.Request("POST", openLogBackupPath, &openLogBackupOpt)
	if err != nil {
		return err
	}
	return nil
}

func getLogBackupSetting(clusterID string, cssClient *golangsdk.ServiceClient) (interface{}, error) {
	getLogBackupSettingHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/logs/settings"
	getLogBackupSettingPath := cssClient.Endpoint + getLogBackupSettingHttpUrl
	getLogBackupSettingPath = strings.ReplaceAll(getLogBackupSettingPath, "{project_id}", cssClient.ProjectID)
	getLogBackupSettingPath = strings.ReplaceAll(getLogBackupSettingPath, "{cluster_id}", clusterID)

	getLogBackupSettingPathOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getLogBackupSettingResp, err := cssClient.Request("GET", getLogBackupSettingPath, &getLogBackupSettingPathOpt)
	if err != nil {
		return nil, err
	}

	getLogBackupSettingRespBody, err := utils.FlattenResponse(getLogBackupSettingResp)
	if err != nil {
		return nil, err
	}
	return utils.PathSearch("logConfiguration", getLogBackupSettingRespBody, nil), nil
}

func updateLogBackupSetting(d *schema.ResourceData, cssClient *golangsdk.ServiceClient) error {
	clusterID := d.Get("cluster_id").(string)
	updateLogSettingHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/logs/settings"
	updateLogSettingPath := cssClient.Endpoint + updateLogSettingHttpUrl
	updateLogSettingPath = strings.ReplaceAll(updateLogSettingPath, "{project_id}", cssClient.ProjectID)
	updateLogSettingPath = strings.ReplaceAll(updateLogSettingPath, "{cluster_id}", clusterID)

	updateLogSettingOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	updateLogSettingOpt.JSONBody = map[string]interface{}{
		"agency":      d.Get("agency").(string),
		"logBasePath": d.Get("base_path").(string),
		"logBucket":   d.Get("bucket").(string),
	}
	_, err := cssClient.Request("POST", updateLogSettingPath, &updateLogSettingOpt)
	if err != nil {
		return err
	}
	return nil
}

func closeLogBackup(clusterID string, cssClient *golangsdk.ServiceClient) error {
	closeLogBackupUrl := "v1.0/{project_id}/clusters/{cluster_id}/logs/close"
	closeLogBackupPath := cssClient.Endpoint + closeLogBackupUrl
	closeLogBackupPath = strings.ReplaceAll(closeLogBackupPath, "{project_id}", cssClient.ProjectID)
	closeLogBackupPath = strings.ReplaceAll(closeLogBackupPath, "{cluster_id}", clusterID)

	closeLogBackupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	_, err := cssClient.Request("PUT", closeLogBackupPath, &closeLogBackupOpt)
	if err != nil {
		return err
	}
	return nil
}

func openLogAutoBackup(d *schema.ResourceData, cssClient *golangsdk.ServiceClient) error {
	clusterID := d.Get("cluster_id").(string)
	openLogAutoBackupHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/logs/policy/update"
	openLogAutoBackupPath := cssClient.Endpoint + openLogAutoBackupHttpUrl
	openLogAutoBackupPath = strings.ReplaceAll(openLogAutoBackupPath, "{project_id}", cssClient.ProjectID)
	openLogAutoBackupPath = strings.ReplaceAll(openLogAutoBackupPath, "{cluster_id}", clusterID)

	openLogAutoBackupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	openLogAutoBackupOpt.JSONBody = map[string]interface{}{
		"period": d.Get("period"),
	}

	_, err := cssClient.Request("POST", openLogAutoBackupPath, &openLogAutoBackupOpt)
	if err != nil {
		return err
	}
	return nil
}

func closeLogAutoBackup(clusterID string, cssClient *golangsdk.ServiceClient) error {
	closeLogAutoBackupHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/logs/policy/close"
	closeLogAutoBackupPath := cssClient.Endpoint + closeLogAutoBackupHttpUrl
	closeLogAutoBackupPath = strings.ReplaceAll(closeLogAutoBackupPath, "{project_id}", cssClient.ProjectID)
	closeLogAutoBackupPath = strings.ReplaceAll(closeLogAutoBackupPath, "{cluster_id}", clusterID)

	closeLogAutoBackupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err := cssClient.Request("PUT", closeLogAutoBackupPath, &closeLogAutoBackupOpt)
	if err != nil {
		return err
	}
	return nil
}

func openLogIngestion(ctx context.Context, d *schema.ResourceData, cssClient *golangsdk.ServiceClient) error {
	clusterID := d.Get("cluster_id").(string)
	targetClusterID := d.Get("target_cluster_id").(string)
	openLogIngestionHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/logs/open"
	openLogIngestionPath := cssClient.Endpoint + openLogIngestionHttpUrl

	openLogIngestionPath = strings.ReplaceAll(openLogIngestionPath, "{project_id}", cssClient.ProjectID)
	openLogIngestionPath = strings.ReplaceAll(openLogIngestionPath, "{cluster_id}", clusterID)
	openLogIngestionPath = fmt.Sprintf("%s?action=real_time_log_collect", openLogIngestionPath)

	err := checkTargetClusterConnectivity(clusterID, targetClusterID, cssClient)
	if err != nil {
		return err
	}

	openLogIngestionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	openLogIngestionOpt.JSONBody = map[string]interface{}{
		"index_prefix":      d.Get("index_prefix").(string),
		"keep_days":         d.Get("keep_days").(int),
		"target_cluster_id": d.Get("target_cluster_id").(string),
	}
	_, err = cssClient.Request("POST", openLogIngestionPath, &openLogIngestionOpt)
	if err != nil {
		return err
	}

	err = configLogIngestionWaitingForCompleted(ctx, clusterID, cssClient, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}
	return nil
}

func checkTargetClusterConnectivity(clusterID string, targetClusterID string, client *golangsdk.ServiceClient) error {
	// Check the connection between two clusters when they are different
	if targetClusterID == clusterID {
		return nil
	}
	connectivityTestHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/logs/connectivity"
	connectivityTestPath := client.Endpoint + connectivityTestHttpUrl
	connectivityTestPath = strings.ReplaceAll(connectivityTestPath, "{project_id}", client.ProjectID)
	connectivityTestPath = strings.ReplaceAll(connectivityTestPath, "{cluster_id}", clusterID)
	connectivityTestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"target_cluster_id": targetClusterID,
		},
	}
	_, err := client.Request("POST", connectivityTestPath, &connectivityTestOpt)
	if err != nil {
		return err
	}
	return nil
}

func configLogIngestionWaitingForCompleted(ctx context.Context, clusterID string, client *golangsdk.ServiceClient, timeout time.Duration) error {
	targetStatus := []string{"200"}
	failedStatus := []string{"300", "302", "303"}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			logIngestionSetting, err := getLogIngestionSetting(clusterID, client)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch(`status`, logIngestionSetting, "").(string)
			if utils.StrSliceContains(targetStatus, status) {
				return logIngestionSetting, "COMPLETED", nil
			}

			// Check for permanent failure states
			if utils.StrSliceContains(failedStatus, status) {
				return logIngestionSetting, status, fmt.Errorf("log ingestion failed with status: %s", status)
			}

			// For other statuses (including 100, 150, 304, 400, 920), continue waiting
			return logIngestionSetting, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func getLogIngestionSetting(clusterID string, cssClient *golangsdk.ServiceClient) (interface{}, error) {
	getLogSettingHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/logs/settings"
	getLogSettingPath := cssClient.Endpoint + getLogSettingHttpUrl
	getLogSettingPath = strings.ReplaceAll(getLogSettingPath, "{project_id}", cssClient.ProjectID)
	getLogSettingPath = strings.ReplaceAll(getLogSettingPath, "{cluster_id}", clusterID)

	getLogIngestionSettingPath := fmt.Sprintf("%s?action=real_time_log_collect", getLogSettingPath)

	getLogIngestionSettingOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getLogIngestionSettingResp, err := cssClient.Request("GET", getLogIngestionSettingPath, &getLogIngestionSettingOpt)
	if err != nil {
		return nil, err
	}
	getLogIngestionSettingRespBody, err := utils.FlattenResponse(getLogIngestionSettingResp)
	if err != nil {
		return nil, err
	}
	return utils.PathSearch("realTimeLogCollectRecord", getLogIngestionSettingRespBody, nil), nil
}

func updateLogIngestionSetting(ctx context.Context, d *schema.ResourceData, cssClient *golangsdk.ServiceClient) error {
	clusterID := d.Get("cluster_id").(string)
	updateLogIngestionHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/logs/settings"
	updateLogIngestionPath := cssClient.Endpoint + updateLogIngestionHttpUrl
	updateLogIngestionPath = strings.ReplaceAll(updateLogIngestionPath, "{project_id}", cssClient.ProjectID)
	updateLogIngestionPath = strings.ReplaceAll(updateLogIngestionPath, "{cluster_id}", clusterID)
	updateLogIngestionPath = fmt.Sprintf("%s?action=real_time_log_collect", updateLogIngestionPath)

	if d.HasChange("target_cluster_id") {
		targetClusterID := d.Get("target_cluster_id").(string)
		err := checkTargetClusterConnectivity(clusterID, targetClusterID, cssClient)
		if err != nil {
			return err
		}
	}

	updateLogIngestionSettingOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	updateLogIngestionSettingOpt.JSONBody = map[string]interface{}{
		"index_prefix":      d.Get("index_prefix").(string),
		"keep_days":         d.Get("keep_days").(int),
		"target_cluster_id": d.Get("target_cluster_id").(string),
	}
	_, err := cssClient.Request("POST", updateLogIngestionPath, &updateLogIngestionSettingOpt)
	if err != nil {
		return err
	}
	err = configLogIngestionWaitingForCompleted(ctx, clusterID, cssClient, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}
	return nil
}

func closeLogIngestion(ctx context.Context, d *schema.ResourceData, cssClient *golangsdk.ServiceClient) error {
	clusterID := d.Get("cluster_id").(string)
	closeLogIngestionHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/logs/close"
	closeLogIngestionPath := cssClient.Endpoint + closeLogIngestionHttpUrl
	closeLogIngestionPath = strings.ReplaceAll(closeLogIngestionPath, "{project_id}", cssClient.ProjectID)
	closeLogIngestionPath = strings.ReplaceAll(closeLogIngestionPath, "{cluster_id}", clusterID)

	closeLogIngestionPath = fmt.Sprintf("%s?action=real_time_log_collect", closeLogIngestionPath)

	closeLogIngestionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err := cssClient.Request("PUT", closeLogIngestionPath, &closeLogIngestionOpt)
	if err != nil {
		return err
	}
	err = closeLogIngestionWaitingForCompleted(ctx, clusterID, cssClient, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	return nil
}

func closeLogIngestionWaitingForCompleted(ctx context.Context, clusterID string, client *golangsdk.ServiceClient, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			setting, err := getLogIngestionSetting(clusterID, client)
			if err != nil {
				return setting, "ERROR", err
			}
			// When log ingestion is closed, the setting will be nil
			if setting == nil {
				return "closed", "COMPLETED", nil
			}
			return setting, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceLogSettingImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	switch len(parts) {
	case 1:
		d.SetId(parts[0] + "/base_log_collect")
		d.Set("cluster_id", parts[0])
		d.Set("action", "base_log_collect")
	case 2:
		d.SetId(d.Id())
		d.Set("cluster_id", parts[0])
		d.Set("action", parts[1])
	default:
		return nil, errors.New("invalid format of import ID, must be <cluster_id> or <cluster_id>/<action>")
	}

	return []*schema.ResourceData{d}, nil
}
