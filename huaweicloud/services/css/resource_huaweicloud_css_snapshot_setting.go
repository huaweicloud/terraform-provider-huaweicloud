package css

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var snapshotSettingNonUpdatableParams = []string{"cluster_id"}

// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/index_snapshot/setting
// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}
// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}/index_snapshot/policy
// @API CSS DELETE /v1.0/{project_id}/clusters/{cluster_id}/index_snapshots
func ResourceSnapshotSetting() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSnapshotSettingCreate,
		ReadContext:   resourceSnapshotSettingRead,
		UpdateContext: resourceSnapshotSettingUpdate,
		DeleteContext: resourceSnapshotSettingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(snapshotSettingNonUpdatableParams),

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
			},
			"agency": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"base_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"max_snapshot_bytes_per_seconds": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"max_restore_bytes_per_seconds": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"indices": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"prefix": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"period": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"keepday": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"frequency": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"delete_auto": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"snapshot_cmk_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildSnapshotSettingRequestBody(d *schema.ResourceData) map[string]interface{} {
	requestBody := map[string]interface{}{
		"agency":                         d.Get("agency").(string),
		"bucket":                         d.Get("bucket").(string),
		"base_path":                      utils.ValueIgnoreEmpty(d.Get("base_path").(string)),
		"max_snapshot_bytes_per_seconds": utils.ValueIgnoreEmpty(d.Get("max_snapshot_bytes_per_seconds").(string)),
		"max_restore_bytes_per_seconds":  utils.ValueIgnoreEmpty(d.Get("max_restore_bytes_per_seconds").(string)),
		"enable":                         utils.ValueIgnoreEmpty(d.Get("enable").(string)),
		"delete_auto":                    utils.ValueIgnoreEmpty(d.Get("delete_auto").(string)),
	}

	// Set the snapshot policy when enable is true
	enable, err := strconv.ParseBool(d.Get("enable").(string))
	if err != nil {
		log.Printf("[ERROR] error parsing 'enable' field to Boolean: %s", err)
	}
	if enable {
		requestBody["indices"] = utils.ValueIgnoreEmpty(d.Get("indices").(string))
		requestBody["prefix"] = utils.ValueIgnoreEmpty(d.Get("prefix").(string))
		requestBody["period"] = utils.ValueIgnoreEmpty(d.Get("period").(string))
		requestBody["keepday"] = utils.ValueIgnoreEmpty(d.Get("keepday").(int))
		requestBody["frequency"] = utils.ValueIgnoreEmpty(d.Get("frequency").(string))
	}

	return requestBody
}

func resourceSnapshotSettingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	clusterID := d.Get("cluster_id").(string)
	cssClient, err := conf.NewServiceClient("css", region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	createSnapshotSettingHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/index_snapshot/setting"
	createSnapshotSettingPath := cssClient.Endpoint + createSnapshotSettingHttpUrl
	createSnapshotSettingPath = strings.ReplaceAll(createSnapshotSettingPath, "{project_id}", cssClient.ProjectID)
	createSnapshotSettingPath = strings.ReplaceAll(createSnapshotSettingPath, "{cluster_id}", clusterID)

	createSnapshotSettingOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildSnapshotSettingRequestBody(d)),
	}

	_, err = cssClient.Request("POST", createSnapshotSettingPath, &createSnapshotSettingOpt)
	if err != nil {
		return diag.Errorf("error opening CSS cluster snapshot function: %s", err)
	}

	d.SetId(clusterID)

	return resourceSnapshotSettingRead(ctx, d, meta)
}

func getClusterSnapshotAutoPolicy(client *golangsdk.ServiceClient, clusterID string) (interface{}, error) {
	getSnapshotSettingHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/index_snapshot/policy"

	getSnapshotSettingPath := client.Endpoint + getSnapshotSettingHttpUrl
	getSnapshotSettingPath = strings.ReplaceAll(getSnapshotSettingPath, "{project_id}", client.ProjectID)
	getSnapshotSettingPath = strings.ReplaceAll(getSnapshotSettingPath, "{cluster_id}", clusterID)

	getSnapshotSettingPathOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getSnapshotSettingResp, err := client.Request("GET", getSnapshotSettingPath, &getSnapshotSettingPathOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getSnapshotSettingResp)
}

func resourceSnapshotSettingRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf   = meta.(*config.Config)
		region = conf.GetRegion(d)
	)

	cssClient, err := conf.NewServiceClient("css", region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	clusterDetail, err := getClusterDetails(cssClient, d.Id())
	if err != nil {
		// "CSS.0015": The cluster does not exist. Status code is 403.
		return common.CheckDeletedDiag(d,
			common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015"), "error getting CSS cluster")
	}

	backupAvailable := utils.PathSearch("backupAvailable", clusterDetail, false).(bool)
	// If backupAvailable is false, the snapshot function is closed.
	if !backupAvailable {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "CSS cluster snapshot is closed")
	}

	getSnapshotSettingRespBody, err := getClusterSnapshotAutoPolicy(cssClient, d.Id())
	if err != nil {
		// "CSS.0015": The cluster does not exist. Status code is 403.
		return common.CheckDeletedDiag(d,
			common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015"), "error getting CSS cluster snapshot policy")
	}

	deleteAuto := d.Get("delete_auto")

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("cluster_id", d.Id()),
		d.Set("agency", utils.PathSearch("agency", getSnapshotSettingRespBody, nil)),
		d.Set("base_path", utils.PathSearch("basePath", getSnapshotSettingRespBody, nil)),
		d.Set("bucket", utils.PathSearch("bucket", getSnapshotSettingRespBody, nil)),
		d.Set("max_snapshot_bytes_per_seconds", utils.PathSearch("maxSnapshotBytesPerSeconds", getSnapshotSettingRespBody, nil)),
		d.Set("max_restore_bytes_per_seconds", utils.PathSearch("maxRestoreBytesPerSeconds", getSnapshotSettingRespBody, nil)),
		d.Set("enable", utils.PathSearch("enable", getSnapshotSettingRespBody, nil)),
		d.Set("indices", utils.PathSearch("indices", getSnapshotSettingRespBody, nil)),
		d.Set("prefix", utils.PathSearch("prefix", getSnapshotSettingRespBody, nil)),
		d.Set("period", utils.PathSearch("period", getSnapshotSettingRespBody, nil)),
		d.Set("keepday", utils.PathSearch("keepday", getSnapshotSettingRespBody, nil)),
		d.Set("frequency", utils.PathSearch("frequency", getSnapshotSettingRespBody, nil)),
		d.Set("snapshot_cmk_id", utils.PathSearch("snapshotCmkId", getSnapshotSettingRespBody, nil)),
		d.Set("delete_auto", deleteAuto),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSnapshotSettingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	cssClient, err := conf.NewServiceClient("css", region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	baseSettingChanges := []string{
		"agency",
		"bucket",
		"base_path",
		"max_snapshot_bytes_per_seconds",
		"max_restore_bytes_per_seconds",
		"enable",
		"indices",
		"prefix",
		"period",
		"keepday",
		"frequency",
		"delete_auto",
	}

	if d.HasChanges(baseSettingChanges...) {
		err = updateSnapshotBaseSetting(d, cssClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSnapshotSettingRead(ctx, d, meta)
}

func resourceSnapshotSettingDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	cssClient, err := conf.NewServiceClient("css", region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	deleteSnapshotSettingUrl := "v1.0/{project_id}/clusters/{cluster_id}/index_snapshots"
	deleteSnapshotSettingPath := cssClient.Endpoint + deleteSnapshotSettingUrl
	deleteSnapshotSettingPath = strings.ReplaceAll(deleteSnapshotSettingPath, "{project_id}", cssClient.ProjectID)
	deleteSnapshotSettingPath = strings.ReplaceAll(deleteSnapshotSettingPath, "{cluster_id}", d.Id())

	deleteSnapshotSettingOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = cssClient.Request("DELETE", deleteSnapshotSettingPath, &deleteSnapshotSettingOpt)
	if err != nil {
		// "CSS.0015": The cluster does not exist. Status code is 403.
		err = common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015")
		// "CSS.0004": Invalid operation. Status code is 415.
		// {"errCode":"CSS.0004","externalMessage":"CSS.0004 : Invalid operation. (Illegal operation)"}
		err = common.ConvertUndefinedErrInto404Err(err, 415, "errCode", "CSS.0004")
		return common.CheckDeletedDiag(d, err, "error closing CSS cluster snapshot function")
	}
	return nil
}

func updateSnapshotBaseSetting(d *schema.ResourceData, cssClient *golangsdk.ServiceClient) error {
	updateSnapshotSettingHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/index_snapshot/setting"
	updateSnapshotSettingPath := cssClient.Endpoint + updateSnapshotSettingHttpUrl
	updateSnapshotSettingPath = strings.ReplaceAll(updateSnapshotSettingPath, "{project_id}", cssClient.ProjectID)
	updateSnapshotSettingPath = strings.ReplaceAll(updateSnapshotSettingPath, "{cluster_id}", d.Id())

	updateSnapshotSettingOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildSnapshotSettingRequestBody(d)),
	}

	_, err := cssClient.Request("POST", updateSnapshotSettingPath, &updateSnapshotSettingOpt)
	if err != nil {
		return fmt.Errorf("error updating CSS cluster snapshot setting: %s", err)
	}
	return nil
}
