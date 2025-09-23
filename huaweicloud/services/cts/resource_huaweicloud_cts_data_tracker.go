package cts

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// ResourceCTSDataTracker is the impl of huaweicloud_cts_data_tracker
// @API CTS POST /v3/{project_id}/tracker
// @API CTS PUT /v3/{project_id}/tracker
// @API CTS DELETE /v3/{project_id}/trackers
// @API CTS GET /v3/{project_id}/trackers
// @API CTS POST /v3/{project_id}/{resource_type}/{resource_id}/tags/create
// @API CTS DELETE /v3/{project_id}/{resource_type}/{resource_id}/tags/delete
func ResourceCTSDataTracker() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCTSDataTrackerCreate,
		ReadContext:   resourceCTSDataTrackerRead,
		UpdateContext: resourceCTSDataTrackerUpdate,
		DeleteContext: resourceCTSDataTrackerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCTSDataTrackerImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringNotInSlice([]string{"system", "system-trace"}, false),
			},
			"data_bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_operation": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 2,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"READ", "WRITE"}, false),
				},
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"bucket_name"},
			},
			"obs_retention_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"bucket_name"},
				ValidateFunc: validation.IntInSlice([]int{0, 30, 60, 90, 180, 1095}),
			},
			"validate_file": {
				Type:         schema.TypeBool,
				Optional:     true,
				RequiredWith: []string{"bucket_name"},
			},
			"compress_type": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"bucket_name"},
			},
			"is_sort_by_service": {
				Type:         schema.TypeBool,
				Optional:     true,
				RequiredWith: []string{"bucket_name"},
				Default:      true,
			},
			"lts_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"tags": common.TagsSchema(),
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"transfer_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"agency_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"detail": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stream_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_topic_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_authorized_bucket": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func buildCreateRequestBody(d *schema.ResourceData) map[string]interface{} {
	trackerType := "data"
	agencyName := "cts_admin_trust"
	reqBody := map[string]interface{}{
		"tracker_type":        trackerType,
		"tracker_name":        d.Get("name").(string),
		"is_lts_enabled":      d.Get("lts_enabled").(bool),
		"is_support_validate": d.Get("validate_file").(bool),
		"data_bucket":         buildDataBucketOpts(d),
		"obs_info":            buildTransferBucketOpts(d),
		"agency_name":         agencyName,
	}

	log.Printf("[DEBUG] creating data CTS tracker options: %#v", reqBody)
	return reqBody
}

func buildDataBucketOpts(d *schema.ResourceData) map[string]interface{} {
	dataBucketCfg := map[string]interface{}{
		"data_bucket_name": d.Get("data_bucket").(string),
	}

	rawOperations, ok := d.GetOk("data_operation")
	if !ok {
		rawOperations = []interface{}{"READ", "WRITE"}
	}

	dataBucketCfg["data_event"] = utils.ExpandToStringList(rawOperations.([]interface{}))
	return dataBucketCfg
}

func buildTransferBucketOpts(d *schema.ResourceData) map[string]interface{} {
	transferCfg := map[string]interface{}{
		"bucket_name":        utils.ValueIgnoreEmpty(d.Get("bucket_name").(string)),
		"file_prefix_name":   utils.ValueIgnoreEmpty(d.Get("file_prefix").(string)),
		"is_sort_by_service": d.Get("is_sort_by_service").(bool),
	}
	if compressType, ok := d.GetOk("compress_type"); ok {
		if compressType.(string) != "gzip" {
			transferCfg["compress_type"] = "json"
		} else {
			transferCfg["compress_type"] = "gzip"
		}
	}
	if v, ok := d.GetOk("obs_retention_period"); ok {
		transferCfg["bucket_lifecycle"] = v.(int)
	}

	return transferCfg
}

func updateTags(ctsClient *golangsdk.ServiceClient, d *schema.ResourceData) error {
	oldRaw, newRaw := d.GetChange("tags")
	id := d.Id()

	if oldTags := oldRaw.(map[string]interface{}); len(oldTags) > 0 {
		err := deleteTags(ctsClient, oldTags, id)
		if err != nil {
			return err
		}
	}

	if newTags := newRaw.(map[string]interface{}); len(newTags) > 0 {
		err := createTags(ctsClient, newTags, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func createTags(ctsClient *golangsdk.ServiceClient, tags map[string]interface{}, id string) error {
	if len(tags) == 0 {
		return nil
	}

	createTagsHttpUrl := "v3/{project_id}/{resource_type}/{resource_id}/tags/create"
	createTagsPath := ctsClient.Endpoint + createTagsHttpUrl
	createTagsPath = strings.ReplaceAll(createTagsPath, "{project_id}", ctsClient.ProjectID)
	createTagsPath = strings.ReplaceAll(createTagsPath, "{resource_type}", "cts-tracker")
	createTagsPath = strings.ReplaceAll(createTagsPath, "{resource_id}", id)
	createTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"tags": utils.ExpandResourceTags(tags),
		},
	}

	_, err := ctsClient.Request("POST", createTagsPath, &createTagsOpt)
	return err
}

func deleteTags(ctsClient *golangsdk.ServiceClient, tags map[string]interface{}, id string) error {
	if len(tags) == 0 {
		return nil
	}

	deleteTagsHttpUrl := "v3/{project_id}/{resource_type}/{resource_id}/tags/delete"
	deleteTagsPath := ctsClient.Endpoint + deleteTagsHttpUrl
	deleteTagsPath = strings.ReplaceAll(deleteTagsPath, "{project_id}", ctsClient.ProjectID)
	deleteTagsPath = strings.ReplaceAll(deleteTagsPath, "{resource_type}", "cts-tracker")
	deleteTagsPath = strings.ReplaceAll(deleteTagsPath, "{resource_id}", id)
	deleteTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	deleteTagsOpt.JSONBody = map[string]interface{}{
		"tags": utils.ExpandResourceTags(tags),
	}

	_, err := ctsClient.Request("DELETE", deleteTagsPath, &deleteTagsOpt)
	return err
}

func resourceCTSDataTrackerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	ctsClient, err := cfg.NewServiceClient("cts", region)
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	createHttpUrl := "v3/{project_id}/tracker"
	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}
	createOpts.JSONBody = utils.RemoveNil(buildCreateRequestBody(d))

	createPath := ctsClient.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", ctsClient.ProjectID)
	resp, err := ctsClient.Request("POST", createPath, &createOpts)
	if err != nil {
		return diag.Errorf("error creating CTS data tracker: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating CTS data tracker: ID is not found in API response")
	}
	d.SetId(id)

	if rawTag := d.Get("tags").(map[string]interface{}); len(rawTag) > 0 {
		err = createTags(ctsClient, rawTag, id)
		if err != nil {
			return diag.Errorf("error creating CTS tracker tags: %s", err)
		}
	}

	// disable status if necessary
	if enabled := d.Get("enabled").(bool); !enabled {
		err = updateDataTrackerStatus(ctsClient, d)
		if err != nil {
			return diag.Errorf("failed to disable CTS data tracker: %s", err)
		}
	}

	return resourceCTSDataTrackerRead(ctx, d, meta)
}

func resourceCTSDataTrackerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	ctsClient, err := cfg.NewServiceClient("cts", region)
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	// update status firstly
	if d.HasChange("enabled") {
		err = updateDataTrackerStatus(ctsClient, d)
		if err != nil {
			return diag.Errorf("error updating CTS tracker status: %s", err)
		}
	}

	// update other configurations
	if d.HasChangeExcept("enabled") {
		trackerType := "data"
		agencyName := "cts_admin_trust"
		trackerName := d.Get("name").(string)
		updateOpts := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json;charset=UTF-8",
			},
		}
		updateBody := map[string]interface{}{
			"tracker_name":        trackerName,
			"tracker_type":        trackerType,
			"is_lts_enabled":      d.Get("lts_enabled").(bool),
			"is_support_validate": d.Get("validate_file").(bool),
			"data_bucket":         buildDataBucketOpts(d),
			"aegncy_name":         agencyName,
		}

		if d.HasChanges("bucket_name", "file_prefix", "obs_retention_period", "compress_type", "is_sort_by_service") {
			updateBody["obs_info"] = buildTransferBucketOpts(d)
		}

		log.Printf("[DEBUG] updating CTS tracker options: %#v", updateBody)
		updateOpts.JSONBody = utils.RemoveNil(updateBody)
		err = updateTracker(ctsClient, &updateOpts)
		if err != nil {
			return diag.Errorf("error updating CTS tracker: %s", err)
		}

		if d.HasChange("tags") {
			err = updateTags(ctsClient, d)
			if err != nil {
				return diag.Errorf("error updating CTS tracker tags: %s", err)
			}
		}
	}

	return resourceCTSDataTrackerRead(ctx, d, meta)
}

func resourceCTSDataTrackerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	ctsClient, err := cfg.NewServiceClient("cts", region)
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	getTrackerHttpUrl := "v3/{project_id}/trackers?tracker_name={tracker_name}&tracker_type={tracker_type}"
	trackerName := d.Get("name").(string)
	trackerType := "data"
	getTrackerOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getTrackerPath := ctsClient.Endpoint + getTrackerHttpUrl
	getTrackerPath = strings.ReplaceAll(getTrackerPath, "{project_id}", ctsClient.ProjectID)
	getTrackerPath = strings.ReplaceAll(getTrackerPath, "{tracker_name}", trackerName)
	getTrackerPath = strings.ReplaceAll(getTrackerPath, "{tracker_type}", trackerType)
	response, err := ctsClient.Request("GET", getTrackerPath, &getTrackerOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CTS data tracker")
	}

	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}

	tracker := utils.PathSearch("trackers|[0]", respBody, nil)
	if tracker == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving CTS data tracker")
	}

	id := utils.PathSearch("id", tracker, "").(string)
	if id == "" {
		return diag.Errorf("error retrieve CTS data tracker: ID is not found in API response")
	}

	d.SetId(id)

	var mErr *multierror.Error
	bucketName := utils.PathSearch("obs_info.bucket_name", tracker, "").(string)
	status := utils.PathSearch("status", tracker, "").(string)
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("tracker_name", tracker, nil)),
		d.Set("validate_file", utils.PathSearch("is_support_validate", tracker, nil)),
		d.Set("lts_enabled", utils.PathSearch("lts.is_lts_enabled", tracker, nil)),
		d.Set("agency_name", utils.PathSearch("agency_name", tracker, nil)),
		d.Set("data_bucket", utils.PathSearch("data_bucket.data_bucket_name", tracker, nil)),
		d.Set("data_operation", utils.PathSearch("data_bucket.data_event", tracker, nil)),
		d.Set("bucket_name", bucketName),
		d.Set("compress_type", utils.PathSearch("obs_info.compress_type", tracker, nil)),
		d.Set("file_prefix", utils.PathSearch("obs_info.file_prefix_name", tracker, nil)),
		d.Set("is_sort_by_service", utils.PathSearch("obs_info.is_sort_by_service", tracker, nil)),
		d.Set("type", utils.PathSearch("tracker_type", tracker, nil)),
		d.Set("status", status),
		d.Set("enabled", status == "enabled"),
		d.Set("create_time", utils.PathSearch("create_time", tracker, nil)),
		d.Set("detail", utils.PathSearch("detail", tracker, nil)),
		d.Set("domain_id", utils.PathSearch("domain_id", tracker, nil)),
		d.Set("group_id", utils.PathSearch("group_id", tracker, nil)),
		d.Set("stream_id", utils.PathSearch("stream_id", tracker, nil)),
		d.Set("log_group_name", utils.PathSearch("lts.log_group_name", tracker, nil)),
		d.Set("log_topic_name", utils.PathSearch("lts.log_topic_name", tracker, nil)),
		d.Set("is_authorized_bucket", utils.PathSearch("obs_info.is_authorized_bucket", tracker, false)),
		d.Set("tags", d.Get("tags")),
	)

	if bucketName != "" {
		mErr = multierror.Append(
			mErr,
			d.Set("transfer_enabled", true),
			d.Set("obs_retention_period", utils.PathSearch("obs_info.bucket_lifecycle", tracker, nil)),
		)
	} else {
		mErr = multierror.Append(mErr, d.Set("transfer_enabled", false))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCTSDataTrackerDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	ctsClient, err := cfg.NewServiceClient("cts", region)
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	deleteHttpUrl := "v3/{project_id}/trackers?tracker_name={tracker_name}&tracker_type={tracker_type}"
	trackerName := d.Get("name").(string)
	trackerType := "data"
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	deletePath := ctsClient.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", ctsClient.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{tracker_name}", trackerName)
	deletePath = strings.ReplaceAll(deletePath, "{tracker_type}", trackerType)
	_, err = ctsClient.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected403ErrInto404Err(err, "error_code", "CTS.0013"),
			fmt.Sprintf("error deleting CTS data tracker %s", trackerName),
		)
	}

	return nil
}

func updateDataTrackerStatus(ctsClient *golangsdk.ServiceClient, d *schema.ResourceData) error {
	status := "enabled"
	if enabled := d.Get("enabled").(bool); !enabled {
		status = "disabled"
	}

	trackerType := "data"
	agencyName := "cts_admin_trust"
	name := d.Get("name").(string)
	statusOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}
	statusOpts.JSONBody = map[string]interface{}{
		"tracker_name": name,
		"tracker_type": trackerType,
		"status":       status,
		"data_bucket":  buildDataBucketOpts(d),
		"agency_name":  agencyName,
	}

	err := updateTracker(ctsClient, &statusOpts)
	return err
}

func updateTracker(ctsClient *golangsdk.ServiceClient, opt *golangsdk.RequestOpts) error {
	updateHttpUrl := "v3/{project_id}/tracker"
	updatePath := ctsClient.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", ctsClient.ProjectID)

	_, err := ctsClient.Request("PUT", updatePath, opt)
	return err
}

func resourceCTSDataTrackerImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	trackerName := d.Id()
	d.Set("name", trackerName)

	return []*schema.ResourceData{d}, nil
}
