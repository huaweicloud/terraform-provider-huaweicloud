package css

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/logs/open
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
			StateContext: schema.ImportStatePassthroughContext,
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
			"agency": {
				Type:     schema.TypeString,
				Required: true,
			},
			"base_path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"period": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"log_switch": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceLogSettingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	clusterID := d.Get("cluster_id").(string)
	cssV1Client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	createLogSettingHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/logs/open"
	createLogSettingPath := cssV1Client.Endpoint + createLogSettingHttpUrl
	createLogSettingPath = strings.ReplaceAll(createLogSettingPath, "{project_id}", cssV1Client.ProjectID)
	createLogSettingPath = strings.ReplaceAll(createLogSettingPath, "{cluster_id}", clusterID)

	createLogSettingOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createLogSettingOpt.JSONBody = map[string]interface{}{
		"agency":      d.Get("agency").(string),
		"logBasePath": d.Get("base_path").(string),
		"logBucket":   d.Get("bucket").(string),
	}
	_, err = cssV1Client.Request("POST", createLogSettingPath, &createLogSettingOpt)
	if err != nil {
		return diag.Errorf("error opening CSS cluster log function: %s", err)
	}

	d.SetId(clusterID)

	if _, ok := d.GetOk("period"); ok {
		err := openLogAutoBackup(d, cssV1Client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceLogSettingRead(ctx, d, meta)
}

func resourceLogSettingRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	cssV1Client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	getLogSettingHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/logs/settings"
	getLogSettingPath := cssV1Client.Endpoint + getLogSettingHttpUrl
	getLogSettingPath = strings.ReplaceAll(getLogSettingPath, "{project_id}", cssV1Client.ProjectID)
	getLogSettingPath = strings.ReplaceAll(getLogSettingPath, "{cluster_id}", d.Id())

	getLogSettingPathOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getLogSettingResp, err := cssV1Client.Request("GET", getLogSettingPath, &getLogSettingPathOpt)
	if err != nil {
		// The cluster does not exist, http code is 403, key/value of error code is errCode/CSS.0015.
		err = common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015")
		return common.CheckDeletedDiag(d, err, "erorr retrieving CSS cluster log setting")
	}

	getLogSettingRespBody, err := utils.FlattenResponse(getLogSettingResp)
	if err != nil {
		return diag.Errorf("erorr retrieving CSS cluster log setting: %s", err)
	}

	logSetting := utils.PathSearch("logConfiguration", getLogSettingRespBody, nil)
	updateAt := utils.PathSearch("updateAt", logSetting, float64(0)).(float64)
	autoEnabled := utils.PathSearch("autoEnable", logSetting, false).(bool)
	var period string
	if autoEnabled {
		period = utils.PathSearch("period", logSetting, nil).(string)
	}
	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("cluster_id", d.Id()),
		d.Set("agency", utils.PathSearch("agency", logSetting, nil)),
		d.Set("base_path", utils.PathSearch("basePath", logSetting, nil)),
		d.Set("bucket", utils.PathSearch("obsBucket", logSetting, nil)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(int64(updateAt)/1000, false)),
		d.Set("auto_enabled", autoEnabled),
		d.Set("period", period),
		d.Set("log_switch", utils.PathSearch("logSwitch", logSetting, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceLogSettingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	cssV1Client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	logBaseSettingChanges := []string{
		"agency",
		"base_path",
		"bucket",
	}

	if d.HasChanges(logBaseSettingChanges...) {
		err = updateLogBaseSetting(d, cssV1Client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("period") {
		if _, ok := d.GetOk("period"); ok {
			err = openLogAutoBackup(d, cssV1Client)
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			err = closeLogAutoBackup(d.Id(), cssV1Client)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceLogSettingRead(ctx, d, meta)
}

func resourceLogSettingDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	cssV1Client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	deleteLogSettingUrl := "v1.0/{project_id}/clusters/{cluster_id}/logs/close"
	deleteLogSettingPath := cssV1Client.Endpoint + deleteLogSettingUrl
	deleteLogSettingPath = strings.ReplaceAll(deleteLogSettingPath, "{project_id}", cssV1Client.ProjectID)
	deleteLogSettingPath = strings.ReplaceAll(deleteLogSettingPath, "{cluster_id}", d.Id())

	deleteLogSettingOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = cssV1Client.Request("PUT", deleteLogSettingPath, &deleteLogSettingOpt)
	if err != nil {
		// "CSS.0015": The cluster does not exist. Status code is 403.
		err = common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015")
		// "CSS.0004": Invalid operation. Status code is 415.
		// {"errCode":"CSS.0004","externalMessage":"CSS.0004 : Invalid operation. (Illegal operation)"}
		err = common.ConvertUndefinedErrInto404Err(err, 415, "errCode", "CSS.0004")
		return common.CheckDeletedDiag(d, err, "error closing CSS cluster log function")
	}

	return nil
}

func updateLogBaseSetting(d *schema.ResourceData, cssV1Client *golangsdk.ServiceClient) error {
	updateLogSettingHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/logs/settings"
	updateLogSettingPath := cssV1Client.Endpoint + updateLogSettingHttpUrl
	updateLogSettingPath = strings.ReplaceAll(updateLogSettingPath, "{project_id}", cssV1Client.ProjectID)
	updateLogSettingPath = strings.ReplaceAll(updateLogSettingPath, "{cluster_id}", d.Id())

	updateLogSettingOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	updateLogSettingOpt.JSONBody = map[string]interface{}{
		"agency":      d.Get("agency").(string),
		"logBasePath": d.Get("base_path").(string),
		"logBucket":   d.Get("bucket").(string),
	}
	_, err := cssV1Client.Request("POST", updateLogSettingPath, &updateLogSettingOpt)
	if err != nil {
		return fmt.Errorf("error updating CSS cluster log setting: %s", err)
	}
	return nil
}

func openLogAutoBackup(d *schema.ResourceData, cssV1Client *golangsdk.ServiceClient) error {
	openLogAutoBackupHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/logs/policy/update"
	openLogAutoBackupPath := cssV1Client.Endpoint + openLogAutoBackupHttpUrl
	openLogAutoBackupPath = strings.ReplaceAll(openLogAutoBackupPath, "{project_id}", cssV1Client.ProjectID)
	openLogAutoBackupPath = strings.ReplaceAll(openLogAutoBackupPath, "{cluster_id}", d.Id())

	openLogAutoBackupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	openLogAutoBackupOpt.JSONBody = map[string]interface{}{
		"period": d.Get("period"),
	}

	_, err := cssV1Client.Request("POST", openLogAutoBackupPath, &openLogAutoBackupOpt)
	if err != nil {
		return fmt.Errorf("error opening CSS cluster log auto backup policy: %s", err)
	}
	return nil
}

func closeLogAutoBackup(clusterID string, cssV1Client *golangsdk.ServiceClient) error {
	closeLogAutoBackupHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/logs/policy/close"
	closeLogAutoBackupPath := cssV1Client.Endpoint + closeLogAutoBackupHttpUrl
	closeLogAutoBackupPath = strings.ReplaceAll(closeLogAutoBackupPath, "{project_id}", cssV1Client.ProjectID)
	closeLogAutoBackupPath = strings.ReplaceAll(closeLogAutoBackupPath, "{cluster_id}", clusterID)

	closeLogAutoBackupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err := cssV1Client.Request("PUT", closeLogAutoBackupPath, &closeLogAutoBackupOpt)
	if err != nil {
		return fmt.Errorf("error closing CSS cluster log auto backup policy: %s", err)
	}
	return nil
}
