package cts

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdkerr"
	client "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cts/v3"
	cts "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cts/v3/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CTS POST /v3/{project_id}/tracker
// @API CTS PUT /v3/{project_id}/tracker
// @API CTS GET /v3/{project_id}/trackers
func ResourceCTSTracker() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCTSTrackerCreate,
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
			"exclude_service": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
			"organization_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"delete_tracker": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": common.TagsSchema(),
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
			"agency_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCTSTrackerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	ctsClient, err := cfg.HcCtsV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	resourceID := "system"
	tracker, err := getSystemTracker(ctsClient)
	if err == nil && tracker != nil {
		log.Print("[DEBUG] the system tracker already exists, update the configuration")
		if tracker.Id != nil {
			resourceID = *tracker.Id
		}

		d.SetId(resourceID)
		return resourceCTSTrackerUpdate(ctx, d, meta)
	}

	var statusCode int
	// check if the error is raised by golangsdk.
	if _, ok := err.(golangsdk.ErrDefault404); ok {
		statusCode = http.StatusNotFound
		// check if the error is raised by huaweicloud-sdk-go-v3.
	} else if responseErr, ok := err.(*sdkerr.ServiceResponseError); ok {
		statusCode = responseErr.StatusCode
	}

	if statusCode != http.StatusNotFound {
		return diag.Errorf("error retrieving CTS tracker: %s", err)
	}

	resourceID, err = createSystemTracker(d, ctsClient)
	if err != nil {
		return diag.Errorf("error creating CTS tracker: %s", err)
	}

	d.SetId(resourceID)

	if rawTag := d.Get("tags").(map[string]interface{}); len(rawTag) > 0 {
		tagList := expandResourceTags(rawTag)
		_, err = ctsClient.BatchCreateResourceTags(buildCreateTagOpt(tagList, resourceID))
		if err != nil {
			return diag.Errorf("error creating CTS tracker tags: %s", err)
		}
	}

	// disable status if necessary
	if enabled := d.Get("enabled").(bool); !enabled {
		if err := updateSystemTrackerStatus(ctsClient, "disabled"); err != nil {
			return diag.Errorf("failed to disable CTS system tracker: %s", err)
		}
	}
	return resourceCTSTrackerRead(ctx, d, meta)
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
			BucketName:      utils.String(d.Get("bucket_name").(string)),
			FilePrefixName:  utils.String(d.Get("file_prefix").(string)),
			IsSortByService: utils.Bool(d.Get("is_sort_by_service").(bool)),
		}
		if v, ok := d.GetOk("compress_type"); ok {
			compressType := cts.GetTrackerObsInfoCompressTypeEnum().GZIP
			if v.(string) != "gzip" {
				compressType = cts.GetTrackerObsInfoCompressTypeEnum().JSON
			}
			obsInfo.CompressType = &compressType
		}

		trackerType := cts.GetUpdateTrackerRequestBodyTrackerTypeEnum().SYSTEM
		agencyName := cts.GetUpdateTrackerRequestBodyAgencyNameEnum().CTS_ADMIN_TRUST
		updateBody := cts.UpdateTrackerRequestBody{
			TrackerName:           "system",
			TrackerType:           trackerType,
			IsLtsEnabled:          utils.Bool(d.Get("lts_enabled").(bool)),
			IsOrganizationTracker: utils.Bool(d.Get("organization_enabled").(bool)),
			IsSupportValidate:     utils.Bool(d.Get("validate_file").(bool)),
			ObsInfo:               &obsInfo,
			AgencyName:            &agencyName,
		}

		if v, ok := d.GetOk("exclude_service"); ok {
			updateBody.ManagementEventSelector = buildManagementEventSelector(v.(*schema.Set).List())
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

		if d.HasChange("tags") {
			err = updateResourceTags(ctsClient, d)
			if err != nil {
				return diag.Errorf("error updating CTS tracker tags: %s", err)
			}
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

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("name", ctsTracker.TrackerName),
		d.Set("organization_enabled", ctsTracker.IsOrganizationTracker),
		d.Set("validate_file", ctsTracker.IsSupportValidate),
		d.Set("kms_id", ctsTracker.KmsId),
	)

	if ctsTracker.Lts != nil {
		mErr = multierror.Append(mErr, d.Set("lts_enabled", ctsTracker.Lts.IsLtsEnabled))
	}

	if ctsTracker.AgencyName != nil {
		mErr = multierror.Append(mErr, d.Set("agency_name", ctsTracker.AgencyName.Value()))
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
			mErr = multierror.Append(mErr, d.Set("transfer_enabled", true))
		} else {
			mErr = multierror.Append(mErr, d.Set("transfer_enabled", false))
		}
	}

	if ctsTracker.ManagementEventSelector != nil {
		mErr = multierror.Append(mErr, d.Set("exclude_service", ctsTracker.ManagementEventSelector.ExcludeService))
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

func resourceCTSTrackerDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	ctsClient, err := cfg.HcCtsV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	trackerCanBeDeleted := d.Get("delete_tracker").(bool)
	if trackerCanBeDeleted {
		trackerName := d.Get("name").(string)
		trackerType := cts.GetDeleteTrackerRequestTrackerTypeEnum().SYSTEM
		deleteOpts := cts.DeleteTrackerRequest{
			TrackerName: &trackerName,
			TrackerType: &trackerType,
		}

		_, err = ctsClient.DeleteTracker(&deleteOpts)
		if err != nil {
			return common.CheckDeletedDiag(d, convertExpected403ErrInto404Err(err, "CTS.0013"), "error deleting CTS system tracker")
		}
		return nil
	}

	if err := updateSystemTrackerStatus(ctsClient, "disabled"); err != nil {
		return diag.Errorf("failed to disable CTS system tracker: %s", err)
	}

	compressType := cts.GetTrackerObsInfoCompressTypeEnum().JSON
	obsInfo := cts.TrackerObsInfo{
		BucketName:      utils.String(""),
		FilePrefixName:  utils.String(""),
		IsSortByService: utils.Bool(false),
		CompressType:    &compressType,
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

	oldRaw, _ := d.GetChange("tags")
	if oldTags := oldRaw.(map[string]interface{}); len(oldTags) > 0 {
		oldTagList := expandResourceTags(oldTags)
		_, err = ctsClient.BatchDeleteResourceTags(buildDeleteTagOpt(oldTagList, d.Id()))
		if err != nil {
			return diag.Errorf("falied to delete CTS system tracker tags: %s", err)
		}
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

func createSystemTracker(d *schema.ResourceData, ctsClient *client.CtsClient) (string, error) {
	obsInfo := cts.TrackerObsInfo{
		BucketName:      utils.String(d.Get("bucket_name").(string)),
		FilePrefixName:  utils.String(d.Get("file_prefix").(string)),
		IsSortByService: utils.Bool(d.Get("is_sort_by_service").(bool)),
	}

	if v, ok := d.GetOk("compress_type"); ok {
		compressType := cts.GetTrackerObsInfoCompressTypeEnum().GZIP
		if v.(string) != "gzip" {
			compressType = cts.GetTrackerObsInfoCompressTypeEnum().JSON
		}
		obsInfo.CompressType = &compressType
	}

	trackerType := cts.GetCreateTrackerRequestBodyTrackerTypeEnum().SYSTEM
	agencyName := cts.GetCreateTrackerRequestBodyAgencyNameEnum().CTS_ADMIN_TRUST
	reqBody := cts.CreateTrackerRequestBody{
		TrackerName:           "system",
		TrackerType:           trackerType,
		IsLtsEnabled:          utils.Bool(d.Get("lts_enabled").(bool)),
		IsOrganizationTracker: utils.Bool(d.Get("organization_enabled").(bool)),
		IsSupportValidate:     utils.Bool(d.Get("validate_file").(bool)),
		ObsInfo:               &obsInfo,
		AgencyName:            &agencyName,
	}

	if v, ok := d.GetOk("exclude_service"); ok {
		reqBody.ManagementEventSelector = buildManagementEventSelector(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("kms_id"); ok {
		encryption := true
		reqBody.KmsId = utils.String(v.(string))
		reqBody.IsSupportTraceFilesEncryption = &encryption
	}

	log.Printf("[DEBUG] creating system CTS tracker options: %#v", reqBody)
	createOpts := cts.CreateTrackerRequest{
		Body: &reqBody,
	}

	resp, err := ctsClient.CreateTracker(&createOpts)
	if err != nil {
		return "", err
	}
	if resp.Id == nil {
		return "", fmt.Errorf("ID is not found in API response")
	}

	return *resp.Id, nil
}

func buildManagementEventSelector(rawServices []interface{}) *cts.ManagementEventSelector {
	if len(rawServices) == 0 {
		return nil
	}

	services := utils.ExpandToStringList(rawServices)
	selector := cts.ManagementEventSelector{
		ExcludeService: &services,
	}

	return &selector
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
	agencyName := cts.GetUpdateTrackerRequestBodyAgencyNameEnum().CTS_ADMIN_TRUST
	statusOpts := cts.UpdateTrackerRequestBody{
		TrackerName: "system",
		TrackerType: trackerType,
		Status:      enabledStatus,
		AgencyName:  &agencyName,
	}
	statusReq := cts.UpdateTrackerRequest{
		Body: &statusOpts,
	}

	_, err := c.UpdateTracker(&statusReq)
	return err
}

func convertExpected403ErrInto404Err(err error, errCode string) error {
	if responseErr, ok := err.(*sdkerr.ServiceResponseError); ok {
		if responseErr.StatusCode == http.StatusForbidden && responseErr.ErrorCode == errCode {
			return golangsdk.ErrDefault404{}
		}
	}
	return err
}
