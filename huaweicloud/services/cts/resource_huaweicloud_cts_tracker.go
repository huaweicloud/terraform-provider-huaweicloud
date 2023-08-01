package cts

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	client "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cts/v3"
	cts "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cts/v3/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceCTSTracker() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCTSTrackerUpdate,
		ReadContext:   resourceCTSTrackerRead,
		UpdateContext: resourceCTSTrackerUpdate,
		DeleteContext: resourceCTSTrackerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCTSTrackerImportState,
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
			"validate_file": {
				Type:         schema.TypeBool,
				Optional:     true,
				RequiredWith: []string{"bucket_name"},
			},
			"kms_id": {
				Type:         schema.TypeString,
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

			"name": {
				Type:     schema.TypeString,
				Computed: true,
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

func resourceCTSTrackerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	ctsClient, err := cfg.HcCtsV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	// update status firstly
	if d.IsNewResource() || d.HasChange("enabled") {
		status := "enabled"
		if enabled := d.Get("enabled").(bool); !enabled {
			status = "disabled"
		}

		if err := updateSystemTrackerStatus(ctsClient, status); err != nil {
			return diag.Errorf("error updating CTS tracker status: %s", err)
		}
	}

	// update other configurations
	if d.IsNewResource() || d.HasChangeExcept("enabled") {
		obsInfo := cts.TrackerObsInfo{
			BucketName:     utils.String(d.Get("bucket_name").(string)),
			FilePrefixName: utils.String(d.Get("file_prefix").(string)),
		}

		trackerType := cts.GetUpdateTrackerRequestBodyTrackerTypeEnum().SYSTEM
		updateBody := cts.UpdateTrackerRequestBody{
			TrackerName:       "system",
			TrackerType:       trackerType,
			IsLtsEnabled:      utils.Bool(d.Get("lts_enabled").(bool)),
			IsSupportValidate: utils.Bool(d.Get("validate_file").(bool)),
			ObsInfo:           &obsInfo,
		}

		var encryption bool
		if v, ok := d.GetOk("kms_id"); ok {
			encryption = true
			updateBody.KmsId = utils.String(v.(string))
		}
		updateBody.IsSupportTraceFilesEncryption = &encryption

		log.Printf("[DEBUG] updating CTS tracker options: %#v", updateBody)
		updateOpts := cts.UpdateTrackerRequest{
			Body: &updateBody,
		}

		_, err = ctsClient.UpdateTracker(&updateOpts)
		if err != nil {
			return diag.Errorf("error updating CTS tracker: %s", err)
		}
	}

	return resourceCTSTrackerRead(ctx, d, meta)
}

func resourceCTSTrackerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	ctsClient, err := cfg.HcCtsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	ctsTracker, err := getSystemTracker(ctsClient)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CTS tracker")
	}

	if ctsTracker.Id != nil {
		d.SetId(*ctsTracker.Id)
	} else {
		d.SetId("system")
	}

	d.Set("region", region)
	d.Set("name", ctsTracker.TrackerName)
	d.Set("lts_enabled", ctsTracker.Lts.IsLtsEnabled)
	d.Set("validate_file", ctsTracker.IsSupportValidate)
	d.Set("kms_id", ctsTracker.KmsId)

	if ctsTracker.ObsInfo != nil {
		bucketName := ctsTracker.ObsInfo.BucketName
		d.Set("bucket_name", bucketName)
		d.Set("file_prefix", ctsTracker.ObsInfo.FilePrefixName)
		if *bucketName != "" {
			d.Set("transfer_enabled", true)
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

func resourceCTSTrackerDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	ctsClient, err := cfg.HcCtsV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	if err := updateSystemTrackerStatus(ctsClient, "disabled"); err != nil {
		return diag.Errorf("failed to disable CTS system tracker: %s", err)
	}

	obsInfo := cts.TrackerObsInfo{
		BucketName:     utils.String(""),
		FilePrefixName: utils.String(""),
	}

	updateBody := cts.UpdateTrackerRequestBody{
		TrackerName:                   "system",
		TrackerType:                   cts.GetUpdateTrackerRequestBodyTrackerTypeEnum().SYSTEM,
		IsLtsEnabled:                  utils.Bool(false),
		IsSupportValidate:             utils.Bool(false),
		IsSupportTraceFilesEncryption: utils.Bool(false),
		KmsId:                         utils.String(""),
		ObsInfo:                       &obsInfo,
	}

	log.Printf("[DEBUG] updating CTS tracker to default configuration: %#v", updateBody)
	updateOpts := cts.UpdateTrackerRequest{
		Body: &updateBody,
	}

	_, err = ctsClient.UpdateTracker(&updateOpts)
	if err != nil {
		return diag.Errorf("falied to reset CTS system tracker: %s", err)
	}

	return nil
}

func resourceCTSTrackerImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	name := d.Id()
	d.Set("name", name)
	return []*schema.ResourceData{d}, nil
}

func formatValue(i interface{}) string {
	jsonRaw, err := json.Marshal(i)
	if err != nil {
		log.Printf("[WARN] failed to marshal %#v: %s", i, err)
		return ""
	}

	return strings.Trim(string(jsonRaw), `"`)
}

func getSystemTracker(ctsClient *client.CtsClient) (*cts.TrackerResponseBody, error) {
	name := "system"
	listOpts := &cts.ListTrackersRequest{
		TrackerName: &name,
	}

	response, err := ctsClient.ListTrackers(listOpts)
	if err != nil {
		return nil, err
	}

	if response.Trackers == nil || len(*response.Trackers) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	allTrackers := *response.Trackers
	return &allTrackers[0], nil
}

func updateSystemTrackerStatus(c *client.CtsClient, status string) error {
	enabledStatus := new(cts.UpdateTrackerRequestBodyStatus)
	if err := enabledStatus.UnmarshalJSON([]byte(status)); err != nil {
		return fmt.Errorf("failed to parse status %s: %s", status, err)
	}

	trackerType := cts.GetUpdateTrackerRequestBodyTrackerTypeEnum().SYSTEM
	statusOpts := cts.UpdateTrackerRequestBody{
		TrackerName: "system",
		TrackerType: trackerType,
		Status:      enabledStatus,
	}
	statusReq := cts.UpdateTrackerRequest{
		Body: &statusOpts,
	}

	_, err := c.UpdateTracker(&statusReq)
	return err
}
