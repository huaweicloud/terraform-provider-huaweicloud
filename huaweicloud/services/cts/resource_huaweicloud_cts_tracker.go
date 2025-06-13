package cts

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CTS POST /v3/{project_id}/tracker
// @API CTS PUT /v3/{project_id}/tracker
// @API CTS GET /v3/{project_id}/trackers
// @API CTS DELETE /v3/{project_id}/trackers
// @API CTS POST /v3/{project_id}/{resource_type}/{resource_id}/tags/create
// @API CTS DELETE /v3/{project_id}/{resource_type}/{resource_id}/tags/delete
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

func resourceCTSTrackerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	ctsClient, err := cfg.NewServiceClient("cts", region)
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	resourceID := "system"
	tracker, err := GetSystemTracker(ctsClient)
	if err == nil && tracker != nil {
		log.Print("[DEBUG] the system tracker already exists, update the configuration")
		if id := utils.PathSearch("id", tracker, "").(string); id != "" {
			resourceID = id
		}
		d.SetId(resourceID)
		return resourceCTSTrackerUpdate(ctx, d, meta)
	}

	if _, ok := err.(golangsdk.ErrDefault404); !ok {
		return diag.Errorf("error retrieving CTS tracker: %s", err)
	}

	resourceID, err = createSystemTracker(d, ctsClient)
	if err != nil {
		return diag.Errorf("error creating CTS tracker: %s", err)
	}

	d.SetId(resourceID)

	if rawTag := d.Get("tags").(map[string]interface{}); len(rawTag) > 0 {
		err = createTags(ctsClient, rawTag, resourceID)
		if err != nil {
			return diag.Errorf("error creating CTS tracker tags: %s", err)
		}
	}

	// disable status if necessary
	if enabled := d.Get("enabled").(bool); !enabled {
		if err := updateSystemTrackerStatus(ctsClient, "disabled"); err != nil {
			return diag.Errorf("failed to disable CTS tracker: %s", err)
		}
	}
	return resourceCTSTrackerRead(ctx, d, meta)
}

func resourceCTSTrackerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	ctsClient, err := cfg.NewServiceClient("cts", region)
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
		obsInfo := map[string]interface{}{
			"bucket_name":        utils.ValueIgnoreEmpty(d.Get("bucket_name").(string)),
			"file_prefix_name":   utils.ValueIgnoreEmpty(d.Get("file_prefix").(string)),
			"is_sort_by_service": d.Get("is_sort_by_service").(bool),
		}
		if compressType, ok := d.GetOk("compress_type"); ok {
			if compressType.(string) != "gzip" {
				obsInfo["compress_type"] = "json"
			} else {
				obsInfo["compress_type"] = "gzip"
			}
		}

		trackerType := "system"
		agencyName := "cts_admin_trust"
		updateBody := map[string]interface{}{
			"tracker_name":            "system",
			"tracker_type":            trackerType,
			"is_lts_enabled":          d.Get("lts_enabled").(bool),
			"is_organization_tracker": d.Get("organization_enabled").(bool),
			"is_support_validate":     d.Get("validate_file").(bool),
			"obs_info":                obsInfo,
			"agency_name":             agencyName,
		}

		if v, ok := d.GetOk("exclude_service"); ok {
			updateBody["management_event_selector"] = buildManagementEventSelector(v.(*schema.Set).List())
		}

		encryption := false
		if v, ok := d.GetOk("kms_id"); ok {
			encryption = true
			updateBody["kms_id"] = v.(string)
		}
		updateBody["is_support_trace_files_encryption"] = encryption

		log.Printf("[DEBUG] updating CTS tracker options: %#v", updateBody)
		statusOpts := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json;charset=UTF-8",
			},
		}

		statusOpts.JSONBody = utils.RemoveNil(updateBody)
		err := updateTracker(ctsClient, &statusOpts)
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

	return resourceCTSTrackerRead(ctx, d, meta)
}

func buildManagementEventSelector(rawServices []interface{}) map[string]interface{} {
	if len(rawServices) == 0 {
		return nil
	}

	services := utils.ExpandToStringList(rawServices)
	return map[string]interface{}{
		"exclude_service": services,
	}
}

func resourceCTSTrackerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	ctsClient, err := cfg.NewServiceClient("cts", region)
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	ctsTracker, err := GetSystemTracker(ctsClient)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CTS tracker")
	}

	id := utils.PathSearch("id", ctsTracker, "system").(string)
	d.SetId(id)

	var mErr *multierror.Error
	bucketName := utils.PathSearch("obs_info.bucket_name", ctsTracker, "").(string)
	status := utils.PathSearch("status", ctsTracker, "").(string)
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("tracker_name", ctsTracker, nil)),
		d.Set("organization_enabled", utils.PathSearch("is_organization_tracker", ctsTracker, false)),
		d.Set("validate_file", utils.PathSearch("is_support_validate", ctsTracker, false)),
		d.Set("lts_enabled", utils.PathSearch("lts.is_lts_enabled", ctsTracker, nil)),
		d.Set("agency_name", utils.PathSearch("agency_name", ctsTracker, nil)),
		d.Set("kms_id", utils.PathSearch("kms_id", ctsTracker, nil)),
		d.Set("bucket_name", bucketName),
		d.Set("file_prefix", utils.PathSearch("obs_info.file_prefix_name", ctsTracker, nil)),
		d.Set("is_sort_by_service", utils.PathSearch("obs_info.is_sort_by_service", ctsTracker, nil)),
		d.Set("compress_type", utils.PathSearch("obs_info.compress_type", ctsTracker, nil)),
		d.Set("exclude_service", utils.PathSearch("management_event_selector.exclude_service", ctsTracker, nil)),
		d.Set("type", utils.PathSearch("tracker_type", ctsTracker, nil)),
		d.Set("status", status),
		d.Set("enabled", status == "enabled"),
		d.Set("create_time", utils.PathSearch("create_time", ctsTracker, nil)),
		d.Set("detail", utils.PathSearch("detail", ctsTracker, nil)),
		d.Set("domain_id", utils.PathSearch("domain_id", ctsTracker, nil)),
		d.Set("group_id", utils.PathSearch("group_id", ctsTracker, nil)),
		d.Set("stream_id", utils.PathSearch("stream_id", ctsTracker, nil)),
		d.Set("log_group_name", utils.PathSearch("lts.log_group_name", ctsTracker, nil)),
		d.Set("log_topic_name", utils.PathSearch("lts.log_topic_name", ctsTracker, nil)),
		d.Set("is_authorized_bucket", utils.PathSearch("obs_info.is_authorized_bucket", ctsTracker, false)),
	)

	if bucketName != "" {
		mErr = multierror.Append(mErr, d.Set("transfer_enabled", true))
	} else {
		mErr = multierror.Append(mErr, d.Set("transfer_enabled", false))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCTSTrackerDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	ctsClient, err := cfg.NewServiceClient("cts", region)
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	trackerCanBeDeleted := d.Get("delete_tracker").(bool)
	if trackerCanBeDeleted {
		deleteHttpUrl := "v3/{project_id}/trackers?tracker_name={tracker_name}&tracker_type={tracker_type}"
		trackerName := d.Get("name").(string)
		trackerType := "system"

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
				fmt.Sprintf("error deleting CTS tracker %s", trackerName),
			)
		}
		return nil
	}

	if err := updateSystemTrackerStatus(ctsClient, "disabled"); err != nil {
		return diag.Errorf("failed to disable CTS tracker: %s", err)
	}

	compressType := "json"
	obsInfo := map[string]interface{}{
		"bucket_name":        "",
		"file_prefix_name":   "",
		"is_sort_by_service": false,
		"compress_type":      compressType,
	}

	updateBody := map[string]interface{}{
		"tracker_name":                      "system",
		"tracker_type":                      "system",
		"is_lts_enabled":                    false,
		"is_support_validate":               false,
		"is_support_trace_files_encryption": false,
		"kms_id":                            "",
		"obs_info":                          obsInfo,
	}

	log.Printf("[DEBUG] updating CTS tracker to default configuration: %#v", updateBody)
	updateOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}

	updateOpts.JSONBody = updateBody

	err = updateTracker(ctsClient, &updateOpts)
	if err != nil {
		return diag.Errorf("falied to reset CTS tracker: %s", err)
	}

	oldRaw, _ := d.GetChange("tags")
	if oldTags := oldRaw.(map[string]interface{}); len(oldTags) > 0 {
		err = deleteTags(ctsClient, oldTags, d.Id())
		if err != nil {
			return diag.Errorf("falied to delete CTS tracker tags: %s", err)
		}
	}

	return nil
}

func resourceCTSTrackerImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	name := d.Id()
	d.Set("name", name)
	return []*schema.ResourceData{d}, nil
}

func createSystemTracker(d *schema.ResourceData, ctsClient *golangsdk.ServiceClient) (string, error) {
	obsInfo := map[string]interface{}{
		"bucket_name":        utils.ValueIgnoreEmpty(d.Get("bucket_name").(string)),
		"file_prefix_name":   utils.ValueIgnoreEmpty(d.Get("file_prefix").(string)),
		"is_sort_by_service": d.Get("is_sort_by_service").(bool),
	}

	if compressType, ok := d.GetOk("compress_type"); ok {
		if compressType.(string) != "gzip" {
			obsInfo["compress_type"] = "json"
		} else {
			obsInfo["compress_type"] = "gzip"
		}
	}

	trackerType := "system"
	agencyName := "cts_admin_trust"
	reqBody := map[string]interface{}{
		"tracker_name":            "system",
		"tracker_type":            trackerType,
		"is_lts_enabled":          d.Get("lts_enabled").(bool),
		"is_organization_tracker": d.Get("organization_enabled").(bool),
		"is_support_validate":     d.Get("validate_file").(bool),
		"obs_info":                obsInfo,
		"agency_name":             agencyName,
	}

	if v, ok := d.GetOk("exclude_service"); ok {
		reqBody["management_event_selector"] = buildManagementEventSelector(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("kms_id"); ok {
		encryption := true
		reqBody["kms_id"] = v.(string)
		reqBody["is_support_trace_files_encryption"] = encryption
	}

	log.Printf("[DEBUG] creating CTS tracker options: %#v", reqBody)
	createHttpUrl := "v3/{project_id}/tracker"
	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}
	createOpts.JSONBody = utils.RemoveNil(reqBody)

	createPath := ctsClient.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", ctsClient.ProjectID)
	resp, err := ctsClient.Request("POST", createPath, &createOpts)
	if err != nil {
		return "", err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return "", err
	}

	id := utils.PathSearch("id", respBody, "").(string)
	if id == "" {
		return "", errors.New("error creating CTS tracker: ID is not found in API response")
	}

	return id, nil
}

func GetSystemTracker(ctsClient *golangsdk.ServiceClient) (interface{}, error) {
	getTrackerHttpUrl := "v3/{project_id}/trackers?tracker_name={tracker_name}"
	trackerName := "system"
	getTrackerOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getTrackerPath := ctsClient.Endpoint + getTrackerHttpUrl
	getTrackerPath = strings.ReplaceAll(getTrackerPath, "{project_id}", ctsClient.ProjectID)
	getTrackerPath = strings.ReplaceAll(getTrackerPath, "{tracker_name}", trackerName)
	response, err := ctsClient.Request("GET", getTrackerPath, &getTrackerOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return nil, err
	}

	searchPath := fmt.Sprintf("trackers[?project_id=='%s']|[0]", ctsClient.ProjectID)
	tracker := utils.PathSearch(searchPath, respBody, nil)
	if tracker == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return tracker, nil
}

func updateSystemTrackerStatus(ctsClient *golangsdk.ServiceClient, status string) error {
	trackerType := "system"
	agencyName := "cts_admin_trust"
	statusOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}
	statusOpts.JSONBody = map[string]interface{}{
		"tracker_name": "system",
		"tracker_type": trackerType,
		"status":       status,
		"agency_name":  &agencyName,
	}
	err := updateTracker(ctsClient, &statusOpts)
	return err
}
