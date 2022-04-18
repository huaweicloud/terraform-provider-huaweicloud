package cts

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

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
func ResourceCTSDataTracker() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCTSDataTrackerCreate,
		ReadContext:   resourceCTSDataTrackerRead,
		UpdateContext: resourceCTSDataTrackerUpdate,
		DeleteContext: resourceCTSDataTrackerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				ValidateFunc: validation.All(
					validation.StringLenBetween(0, 64),
					validation.StringMatch(regexp.MustCompile(`^[\.\-_A-Za-z0-9]+$`),
						"only letters, numbers, hyphens (-), underscores (_), and periods (.) are allowed"),
				),
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
			"lts_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

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
		},
	}

}

func buildCreateRequestBody(d *schema.ResourceData) *cts.CreateTrackerRequestBody {
	trackerType := cts.GetCreateTrackerRequestBodyTrackerTypeEnum().DATA
	reqBody := cts.CreateTrackerRequestBody{
		TrackerType:       trackerType,
		TrackerName:       d.Get("name").(string),
		IsLtsEnabled:      utils.Bool(d.Get("lts_enabled").(bool)),
		IsSupportValidate: utils.Bool(d.Get("validate_file").(bool)),
		DataBucket:        buildDataBucketOpts(d),
		ObsInfo:           buildTransferBucketOpts(d),
	}

	log.Printf("[DEBUG] creating data CTS tracker options: %#v", reqBody)
	return &reqBody
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
		BucketName:     utils.String(d.Get("bucket_name").(string)),
		FilePrefixName: utils.String(d.Get("file_prefix").(string)),
	}
	if v, ok := d.GetOk("obs_retention_period"); ok {
		lifecycle := int32(v.(int))
		transferCfg.BucketLifecycle = &lifecycle
	}

	return &transferCfg
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

	if _, err := ctsClient.CreateTracker(&createOpts); err != nil {
		return diag.Errorf("error creating data CTS tracker: %s", err)
	}

	trackerName := d.Get("name").(string)
	d.SetId(trackerName)

	// disable status if necessary
	if enabled := d.Get("enabled").(bool); !enabled {
		err = updateDataTrackerStatus(ctsClient, trackerName, "disabled")
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

	trackerName := d.Id()
	// update status firstly
	if d.HasChange("enabled") {
		status := "enabled"
		if enabled := d.Get("enabled").(bool); !enabled {
			status = "disabled"
		}

		err = updateDataTrackerStatus(ctsClient, trackerName, status)
		if err != nil {
			return diag.Errorf("error updating CTS tracker status: %s", err)
		}
	}

	// update other configurations
	if d.HasChangeExcept("enabled") {
		trackerType := cts.GetUpdateTrackerRequestBodyTrackerTypeEnum().DATA
		updateReq := cts.UpdateTrackerRequestBody{
			TrackerName:       trackerName,
			TrackerType:       trackerType,
			IsLtsEnabled:      utils.Bool(d.Get("lts_enabled").(bool)),
			IsSupportValidate: utils.Bool(d.Get("validate_file").(bool)),
			DataBucket:        buildDataBucketOpts(d),
		}

		if d.HasChanges("bucket_name", "file_prefix", "obs_retention_period") {
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

	trackerName := d.Id()
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

	d.Set("region", region)
	d.Set("name", ctsTracker.TrackerName)
	d.Set("lts_enabled", ctsTracker.Lts.IsLtsEnabled)
	d.Set("validate_file", ctsTracker.IsSupportValidate)

	if ctsTracker.DataBucket != nil {
		d.Set("data_bucket", ctsTracker.DataBucket.DataBucketName)

		if ctsTracker.DataBucket.DataEvent != nil {
			operations := make([]string, len(*ctsTracker.DataBucket.DataEvent))
			for i, event := range *ctsTracker.DataBucket.DataEvent {
				operations[i] = formatValue(event)
			}
			d.Set("data_operation", operations)
		}
	}

	if ctsTracker.ObsInfo != nil {
		bucketName := ctsTracker.ObsInfo.BucketName
		d.Set("bucket_name", bucketName)
		d.Set("file_prefix", ctsTracker.ObsInfo.FilePrefixName)

		if *bucketName != "" {
			d.Set("transfer_enabled", true)
			d.Set("obs_retention_period", ctsTracker.ObsInfo.BucketLifecycle)
		} else {
			d.Set("transfer_enabled", false)
		}
	}

	if ctsTracker.TrackerType != nil {
		d.Set("type", formatValue(ctsTracker.TrackerType))
	}
	if ctsTracker.Status != nil {
		status := formatValue(ctsTracker.Status)
		d.Set("status", status)
		d.Set("enabled", status == "enabled")
	}

	return nil
}

func resourceCTSDataTrackerDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	ctsClient, err := cfg.HcCtsV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	trackerName := d.Id()
	trackerType := cts.GetDeleteTrackerRequestTrackerTypeEnum().DATA
	deleteOpts := cts.DeleteTrackerRequest{
		TrackerName: &trackerName,
		TrackerType: &trackerType,
	}

	_, err = ctsClient.DeleteTracker(&deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting CTS data tracker %s: %s", trackerName, err)
	}

	return nil
}

func updateDataTrackerStatus(c *client.CtsClient, name, status string) error {
	enabledStatus := new(cts.UpdateTrackerRequestBodyStatus)
	if err := enabledStatus.UnmarshalJSON([]byte(status)); err != nil {
		return fmt.Errorf("failed to parse status %s: %s", status, err)
	}

	trackerType := cts.GetUpdateTrackerRequestBodyTrackerTypeEnum().DATA
	statusOpts := cts.UpdateTrackerRequestBody{
		TrackerName: name,
		TrackerType: trackerType,
		Status:      enabledStatus,
	}
	statusReq := cts.UpdateTrackerRequest{
		Body: &statusOpts,
	}

	_, err := c.UpdateTracker(&statusReq)
	return err
}
