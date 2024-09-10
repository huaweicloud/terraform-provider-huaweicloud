package cts

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	client "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cts/v3"
	cts "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cts/v3/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// ResourceCTSDataTracker is the impl of huaweicloud_cts_data_tracker
// @API CTS POST /v3/{project_id}/tracker
// @API CTS PUT /v3/{project_id}/tracker
// @API CTS DELETE /v3/{project_id}/trackers
// @API CTS GET /v3/{project_id}/trackers
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
		},
	}
}

func buildCreateRequestBody(d *schema.ResourceData) *cts.CreateTrackerRequestBody {
	trackerType := cts.GetCreateTrackerRequestBodyTrackerTypeEnum().DATA
	agencyName := cts.GetCreateTrackerRequestBodyAgencyNameEnum().CTS_ADMIN_TRUST
	reqBody := cts.CreateTrackerRequestBody{
		TrackerType:       trackerType,
		TrackerName:       d.Get("name").(string),
		IsLtsEnabled:      utils.Bool(d.Get("lts_enabled").(bool)),
		IsSupportValidate: utils.Bool(d.Get("validate_file").(bool)),
		DataBucket:        buildDataBucketOpts(d),
		ObsInfo:           buildTransferBucketOpts(d),
		AgencyName:        &agencyName,
	}

	log.Printf("[DEBUG] creating data CTS tracker options: %#v", reqBody)
	return &reqBody
}

func expandResourceTags(tagmap map[string]interface{}) *[]cts.Tags {
	taglist := make([]cts.Tags, 0, len(tagmap))

	for k, v := range tagmap {
		taglist = append(taglist, cts.Tags{
			Key:   utils.String(k),
			Value: utils.String(v.(string)),
		})
	}

	return &taglist
}

func buildCreateTagOpt(taglist *[]cts.Tags, id string) *cts.BatchCreateResourceTagsRequest {
	reqBody := cts.BatchCreateResourceTagsRequestBody{
		Tags: taglist,
	}
	tagOpt := cts.BatchCreateResourceTagsRequest{
		ResourceId:   id,
		ResourceType: cts.GetBatchCreateResourceTagsRequestResourceTypeEnum().CTS_TRACKER,
		Body:         &reqBody,
	}

	return &tagOpt
}

func buildDeleteTagOpt(taglist *[]cts.Tags, id string) *cts.BatchDeleteResourceTagsRequest {
	reqBody := cts.BatchDeleteResourceTagsRequestBody{
		Tags: taglist,
	}
	tagOpt := cts.BatchDeleteResourceTagsRequest{
		ResourceId:   id,
		ResourceType: cts.GetBatchDeleteResourceTagsRequestResourceTypeEnum().CTS_TRACKER,
		Body:         &reqBody,
	}

	return &tagOpt
}

func buildDataBucketOpts(d *schema.ResourceData) *cts.DataBucket {
	dataBucketCfg := cts.DataBucket{
		DataBucketName: utils.String(d.Get("data_bucket").(string)),
	}

	var rawOperations []interface{}
	if v, ok := d.GetOk("data_operation"); ok {
		rawOperations = v.([]interface{})
	} else {
		rawOperations = []interface{}{"READ", "WRITE"}
	}

	operations := make([]cts.DataBucketDataEvent, len(rawOperations))
	dataEvents := cts.GetDataBucketDataEventEnum()
	for i, raw := range rawOperations {
		if op, ok := raw.(string); ok {
			if op == "READ" {
				operations[i] = dataEvents.READ
			} else if op == "WRITE" {
				operations[i] = dataEvents.WRITE
			}
		}
	}
	dataBucketCfg.DataEvent = &operations

	return &dataBucketCfg
}

func buildTransferBucketOpts(d *schema.ResourceData) *cts.TrackerObsInfo {
	transferCfg := cts.TrackerObsInfo{
		BucketName:      utils.String(d.Get("bucket_name").(string)),
		FilePrefixName:  utils.String(d.Get("file_prefix").(string)),
		IsSortByService: utils.Bool(d.Get("is_sort_by_service").(bool)),
	}
	if v, ok := d.GetOk("compress_type"); ok {
		compressType := cts.GetTrackerObsInfoCompressTypeEnum().GZIP
		if v.(string) != "gzip" {
			compressType = cts.GetTrackerObsInfoCompressTypeEnum().JSON
		}
		transferCfg.CompressType = &compressType
	}
	if v, ok := d.GetOk("obs_retention_period"); ok {
		lifecycle := int32(v.(int))
		transferCfg.BucketLifecycle = &lifecycle
	}

	return &transferCfg
}

func updateResourceTags(ctsClient *client.CtsClient, d *schema.ResourceData) error {
	oldRaw, newRaw := d.GetChange("tags")
	id := d.Id()

	if oldTags := oldRaw.(map[string]interface{}); len(oldTags) > 0 {
		oldTagList := expandResourceTags(oldTags)
		_, err := ctsClient.BatchDeleteResourceTags(buildDeleteTagOpt(oldTagList, id))
		if err != nil {
			return err
		}
	}

	if newTags := newRaw.(map[string]interface{}); len(newTags) > 0 {
		newTagsList := expandResourceTags(newTags)
		_, err := ctsClient.BatchCreateResourceTags(buildCreateTagOpt(newTagsList, id))
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceCTSDataTrackerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	ctsClient, err := cfg.HcCtsV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	createOpts := cts.CreateTrackerRequest{
		Body: buildCreateRequestBody(d),
	}

	resp, err := ctsClient.CreateTracker(&createOpts)
	if err != nil {
		return diag.Errorf("error creating CTS data tracker: %s", err)
	}

	if resp.Id == nil {
		return diag.Errorf("error creating CTS data tracker: ID is not found in API response")
	}
	d.SetId(*resp.Id)

	if rawTag := d.Get("tags").(map[string]interface{}); len(rawTag) > 0 {
		tagList := expandResourceTags(rawTag)
		_, err = ctsClient.BatchCreateResourceTags(buildCreateTagOpt(tagList, *resp.Id))
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
	ctsClient, err := cfg.HcCtsV3Client(cfg.GetRegion(d))
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
		trackerType := cts.GetUpdateTrackerRequestBodyTrackerTypeEnum().DATA
		agencyName := cts.GetUpdateTrackerRequestBodyAgencyNameEnum().CTS_ADMIN_TRUST
		trackerName := d.Get("name").(string)
		updateReq := cts.UpdateTrackerRequestBody{
			TrackerName:       trackerName,
			TrackerType:       trackerType,
			IsLtsEnabled:      utils.Bool(d.Get("lts_enabled").(bool)),
			IsSupportValidate: utils.Bool(d.Get("validate_file").(bool)),
			DataBucket:        buildDataBucketOpts(d),
			AgencyName:        &agencyName,
		}

		if d.HasChanges("bucket_name", "file_prefix", "obs_retention_period", "compress_type", "is_sort_by_service") {
			updateReq.ObsInfo = buildTransferBucketOpts(d)
		}

		log.Printf("[DEBUG] updating CTS tracker options: %#v", updateReq)
		updateOpts := cts.UpdateTrackerRequest{
			Body: &updateReq,
		}

		_, err = ctsClient.UpdateTracker(&updateOpts)
		if err != nil {
			return diag.Errorf("error updating CTS tracker: %s", err)
		}

		if d.HasChange("tags") {
			err = updateResourceTags(ctsClient, d)
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
	ctsClient, err := cfg.HcCtsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	trackerName := d.Get("name").(string)
	trackerType := cts.GetListTrackersRequestTrackerTypeEnum().DATA
	listOpts := &cts.ListTrackersRequest{
		TrackerName: &trackerName,
		TrackerType: &trackerType,
	}

	response, err := ctsClient.ListTrackers(listOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CTS data tracker")
	}

	if response.Trackers == nil || len(*response.Trackers) == 0 {
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Resource not found",
				Detail:   fmt.Sprintf("cannot retrieve CTS data tracker %s", trackerName),
			},
		}
	}

	allTrackers := *response.Trackers
	ctsTracker := allTrackers[0]

	if ctsTracker.Id == nil {
		return diag.Errorf("error retrieve CTS data tracker: ID is not found in API response")
	}

	d.SetId(*ctsTracker.Id)

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("name", ctsTracker.TrackerName),
		d.Set("validate_file", ctsTracker.IsSupportValidate),
	)

	if ctsTracker.Lts != nil {
		mErr = multierror.Append(mErr, d.Set("lts_enabled", ctsTracker.Lts.IsLtsEnabled))
	}

	if ctsTracker.AgencyName != nil {
		mErr = multierror.Append(mErr, d.Set("agency_name", ctsTracker.AgencyName.Value()))
	}

	if ctsTracker.DataBucket != nil {
		mErr = multierror.Append(mErr, d.Set("data_bucket", ctsTracker.DataBucket.DataBucketName))

		if ctsTracker.DataBucket.DataEvent != nil {
			operations := make([]string, len(*ctsTracker.DataBucket.DataEvent))
			for i, event := range *ctsTracker.DataBucket.DataEvent {
				operations[i] = formatValue(event)
			}
			mErr = multierror.Append(mErr, d.Set("data_operation", operations))
		}
	}

	if ctsTracker.ObsInfo != nil {
		bucketName := ctsTracker.ObsInfo.BucketName
		mErr = multierror.Append(
			mErr,
			d.Set("bucket_name", bucketName),
			d.Set("file_prefix", ctsTracker.ObsInfo.FilePrefixName),
			d.Set("is_sort_by_service", ctsTracker.ObsInfo.IsSortByService),
		)

		if ctsTracker.ObsInfo.CompressType != nil {
			mErr = multierror.Append(mErr, d.Set("compress_type", formatValue(ctsTracker.ObsInfo.CompressType)))
		}

		if *bucketName != "" {
			mErr = multierror.Append(
				mErr,
				d.Set("transfer_enabled", true),
				d.Set("obs_retention_period", ctsTracker.ObsInfo.BucketLifecycle),
			)
		} else {
			mErr = multierror.Append(mErr, d.Set("transfer_enabled", false))
		}
	}

	if ctsTracker.TrackerType != nil {
		mErr = multierror.Append(mErr, d.Set("type", formatValue(ctsTracker.TrackerType)))
	}
	if ctsTracker.Status != nil {
		status := formatValue(ctsTracker.Status)
		mErr = multierror.Append(
			mErr,
			d.Set("status", status),
			d.Set("enabled", status == "enabled"),
		)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCTSDataTrackerDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	ctsClient, err := cfg.HcCtsV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	trackerName := d.Get("name").(string)
	trackerType := cts.GetDeleteTrackerRequestTrackerTypeEnum().DATA
	deleteOpts := cts.DeleteTrackerRequest{
		TrackerName: &trackerName,
		TrackerType: &trackerType,
	}

	_, err = ctsClient.DeleteTracker(&deleteOpts)
	if err != nil {
		return common.CheckDeletedDiag(d,
			convertExpected403ErrInto404Err(err, "CTS.0013"),
			fmt.Sprintf("error deleting CTS data tracker %s", trackerName),
		)
	}

	return nil
}

func updateDataTrackerStatus(c *client.CtsClient, d *schema.ResourceData) error {
	status := "enabled"
	if enabled := d.Get("enabled").(bool); !enabled {
		status = "disabled"
	}
	enabledStatus := new(cts.UpdateTrackerRequestBodyStatus)
	if err := enabledStatus.UnmarshalJSON([]byte(status)); err != nil {
		return fmt.Errorf("failed to parse status %s: %s", status, err)
	}

	trackerType := cts.GetUpdateTrackerRequestBodyTrackerTypeEnum().DATA
	agencyName := cts.GetUpdateTrackerRequestBodyAgencyNameEnum().CTS_ADMIN_TRUST
	name := d.Get("name").(string)
	statusOpts := cts.UpdateTrackerRequestBody{
		TrackerName: name,
		TrackerType: trackerType,
		Status:      enabledStatus,
		DataBucket:  buildDataBucketOpts(d),
		AgencyName:  &agencyName,
	}
	statusReq := cts.UpdateTrackerRequest{
		Body: &statusOpts,
	}

	_, err := c.UpdateTracker(&statusReq)
	return err
}

func resourceCTSDataTrackerImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	trackerName := d.Id()
	d.Set("name", trackerName)

	return []*schema.ResourceData{d}, nil
}
