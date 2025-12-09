package workspace

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
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

func desktopVolumeSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"iops": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The IOPS of the volume.`,
			},
			"throughput": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The throughput of the volume.`,
			},
			"kms_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The KMS key ID used to encrypt the volume.`,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"device": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// @API Workspace POST /v2/{project_id}/desktops/rebuild
// @API Workspace POST /v2/{project_id}/desktops/resize
// @API Workspace GET /v2/{project_id}/desktops/{desktop_id}
// @API Workspace DELETE /v2/{project_id}/desktops/{desktop_id}
// @API Workspace POST /v2/{project_id}/desktops/{id}/tags/action
// @API Workspace POST /v2/{project_id}/volumes
// @API Workspace GET /v2/{project_id}/desktops/{desktop_id}/networks
// @API Workspace PUT /v2/{project_id}/desktops/{desktop_id}/networks
// @API Workspace POST /v2/{project_id}/desktops/action
// @API Workspace POST /v2/{project_id}/desktops
// @API Workspace POST /v2/{project_id}/volumes/expand
// @API Workspace GET /v2/{project_id}/workspace-sub-jobs
func ResourceDesktop() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDesktopCreate,
		ReadContext:   resourceDesktopRead,
		UpdateContext: resourceDesktopUpdate,
		DeleteContext: resourceDesktopDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"image_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"market", "gold", "private",
				}, false),
			},
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_email": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_group": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"sudo", "default", "administrators", "users",
				}, false),
			},
			"root_volume": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     desktopVolumeSchemaResource(),
			},
			"data_volume": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     desktopVolumeSchemaResource(),
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nic": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"name": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				DiffSuppressFunc: utils.SuppressCaseDiffs(),
			},
			"email_notification": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"tags": common.TagsSchema(),
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"delete_user": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"power_action": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"power_action_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildDesktopRootVolume(volumes []interface{}) map[string]interface{} {
	if len(volumes) < 1 {
		return nil
	}

	volume := volumes[0].(map[string]interface{})
	return map[string]interface{}{
		"type":       volume["type"].(string),
		"size":       volume["size"].(int),
		"iops":       utils.ValueIgnoreEmpty(volume["iops"].(int)),
		"throughput": utils.ValueIgnoreEmpty(volume["throughput"].(int)),
		"kms_id":     utils.ValueIgnoreEmpty(volume["kms_id"].(string)),
	}
}

func buildDesktopDataVolumes(volumes []interface{}) []map[string]interface{} {
	if len(volumes) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(volumes))
	for i, val := range volumes {
		volume := val.(map[string]interface{})
		result[i] = map[string]interface{}{
			"type":       volume["type"].(string),
			"size":       volume["size"].(int),
			"iops":       utils.ValueIgnoreEmpty(volume["iops"].(int)),
			"throughput": utils.ValueIgnoreEmpty(volume["throughput"].(int)),
			"kms_id":     utils.ValueIgnoreEmpty(volume["kms_id"].(string)),
		}
	}
	return result
}

func buildDesktopNics(nics []interface{}) []map[string]interface{} {
	if len(nics) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(nics))
	for i, val := range nics {
		result[i] = map[string]interface{}{
			"subnet_id": utils.PathSearch("network_id", val, nil),
		}
	}
	return result
}

func buildDesktopSecurityGroups(securityGroups *schema.Set) []map[string]interface{} {
	if securityGroups.Len() < 1 {
		return nil
	}

	result := make([]map[string]interface{}, securityGroups.Len())
	for i, val := range securityGroups.List() {
		result[i] = map[string]interface{}{
			"id": val.(string),
		}
	}
	return result
}

func buildDesktopCreateOpts(d *schema.ResourceData, conf *config.Config) map[string]interface{} {
	return map[string]interface{}{
		"desktops": []map[string]interface{}{
			{
				"user_name":     d.Get("user_name").(string),
				"user_email":    d.Get("user_email").(string),
				"user_group":    d.Get("user_group").(string),
				"computer_name": utils.ValueIgnoreEmpty(d.Get("name").(string)),
			},
		},
		"desktop_type":          "DEDICATED",
		"product_id":            d.Get("flavor_id").(string),
		"root_volume":           buildDesktopRootVolume(d.Get("root_volume").([]interface{})),
		"availability_zone":     utils.ValueIgnoreEmpty(d.Get("availability_zone").(string)),
		"image_type":            d.Get("image_type").(string),
		"image_id":              d.Get("image_id").(string),
		"vpc_id":                d.Get("vpc_id").(string),
		"email_notification":    d.Get("email_notification"),
		"data_volumes":          buildDesktopDataVolumes(d.Get("data_volume").([]interface{})),
		"nics":                  buildDesktopNics(d.Get("nic").([]interface{})),
		"security_groups":       buildDesktopSecurityGroups(d.Get("security_groups").(*schema.Set)),
		"tags":                  utils.ValueIgnoreEmpty(utils.ExpandResourceTags(d.Get("tags").(map[string]interface{}))),
		"enterprise_project_id": utils.ValueIgnoreEmpty(conf.GetEnterpriseProjectID(d)),
	}
}

func waitForWorkspaceJobCompleted(ctx context.Context, client *golangsdk.ServiceClient, jobId string, timeout time.Duration) (string, error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"WAITING", "RUNNING"},
		Target:       []string{"SUCCESS"},
		Refresh:      refreshWorkspaceJobFunc(client, jobId),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 15 * time.Second,
	}

	resp, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return "", err
	}

	return utils.PathSearch("entities.desktop_id", resp, "").(string), nil
}

func refreshWorkspaceJobFunc(client *golangsdk.ServiceClient, jobId string, queryParams ...string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			httpUrl  = "v2/{project_id}/workspace-sub-jobs"
			listOpts = golangsdk.RequestOpts{
				KeepResponseBody: true,
				MoreHeaders: map[string]string{
					"Content-Type": "application/json",
				},
			}
		)

		listPath := client.Endpoint + httpUrl
		listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
		listPath = fmt.Sprintf("%s?job_id=%s", listPath, jobId)
		if len(queryParams) > 0 {
			listPath += queryParams[0]
		}

		resp, err := client.Request("GET", listPath, &listOpts)
		if err != nil {
			return resp, "ERROR", err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return resp, "ERROR", err
		}

		jobs := utils.PathSearch("jobs", respBody, make([]interface{}, 0)).([]interface{})
		if len(jobs) < 1 {
			return resp, "", fmt.Errorf("unable to find any job details")
		}

		for _, job := range jobs {
			status := utils.PathSearch("status", job, "").(string)
			if status == "SUCCESS" {
				continue
			}
			return job, status, nil
		}

		return jobs[0], "SUCCESS", nil
	}
}

func resourceDesktopCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf    = meta.(*config.Config)
		httpUrl = "v2/{project_id}/desktops"
	)

	client, err := conf.NewServiceClient("workspace", conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildDesktopCreateOpts(d, conf)),
	}
	resp, err := client.Request("POST", createPath, &createOpts)
	if err != nil {
		return diag.Errorf("error creating Workspace desktop: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find job ID from API response")
	}

	desktopId, err := waitForWorkspaceJobCompleted(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the job (%s) completed: %s", jobId, err)
	}
	log.Printf("[DEBUG] The job (%s) has been completed", jobId)

	d.SetId(desktopId)

	if action, ok := d.GetOk("power_action"); ok {
		if action == "os-start" {
			log.Printf("[WARN] the power action (os-start) is invalid after desktop created")
		} else if err = updateDesktopPowerAction(ctx, client, d, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceDesktopRead(ctx, d, meta)
}

func flattenDesktopRootVolume(volume interface{}) []map[string]interface{} {
	if volume == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"type":       utils.PathSearch("type", volume, nil),
			"size":       utils.PathSearch("size", volume, nil),
			"iops":       utils.PathSearch("iops", volume, nil),
			"throughput": utils.PathSearch("throughput", volume, nil),
			"id":         utils.PathSearch("volume_id", volume, nil),
			"name":       utils.PathSearch("display_name", volume, nil),
			"device":     utils.PathSearch("device", volume, nil),
			"created_at": utils.PathSearch("create_time", volume, nil),
		},
	}
}

func flattenDesktopDataVolumes(volumes []interface{}) []map[string]interface{} {
	if len(volumes) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(volumes))
	for i, volume := range volumes {
		result[i] = map[string]interface{}{
			"type":       utils.PathSearch("type", volume, nil),
			"size":       utils.PathSearch("size", volume, nil),
			"iops":       utils.PathSearch("iops", volume, nil),
			"throughput": utils.PathSearch("throughput", volume, nil),
			"id":         utils.PathSearch("volume_id", volume, nil),
			"name":       utils.PathSearch("display_name", volume, nil),
			"device":     utils.PathSearch("device", volume, nil),
			"created_at": utils.PathSearch("create_time", volume, nil),
		}
	}

	// Since the volumes in the response body are unordered, they are sorted by device.
	sort.Slice(result, func(i, j int) bool {
		a := result[i]
		b := result[j]

		return a["device"].(string) <= b["device"].(string)
	})

	return result
}

func getDesktopNetwork(client *golangsdk.ServiceClient, desktopId string) ([]map[string]interface{}, error) {
	httpUrl := "v2/{project_id}/desktops/{desktop_id}/networks"
	getNetworkPath := client.Endpoint + httpUrl
	getNetworkPath = strings.ReplaceAll(getNetworkPath, "{project_id}", client.ProjectID)
	getNetworkPath = strings.ReplaceAll(getNetworkPath, "{desktop_id}", desktopId)
	getNetworkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	resp, err := client.Request("GET", getNetworkPath, &getNetworkOpt)
	if err != nil {
		return nil, fmt.Errorf("error getting desktop network info: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	network := utils.PathSearch("network_infos", respBody, make([]interface{}, 0)).([]interface{})
	if len(network) < 1 {
		return nil, fmt.Errorf("unable to find any network information under Workspace desktop (%s)", desktopId)
	}

	nic := []map[string]interface{}{
		{
			"network_id": utils.PathSearch("subnet_info.id", network[0], nil),
		},
	}

	return nic, err
}

func GetDesktopById(client *golangsdk.ServiceClient, desktopId string) (interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/desktops/{desktop_id}"
		getOpts = golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		}
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{desktop_id}", desktopId)
	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}
	return utils.PathSearch("desktop", respBody, nil), nil
}

func resourceDesktopRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	desktopId := d.Id()
	respBody, err := GetDesktopById(client, desktopId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Workspace desktop")
	}
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("flavor_id", utils.PathSearch("product.product_id", respBody, nil)),
		d.Set("user_name", utils.PathSearch("user_name", respBody, nil)),
		d.Set("root_volume", flattenDesktopRootVolume(utils.PathSearch("root_volume", respBody, nil))),
		d.Set("data_volume", flattenDesktopDataVolumes(utils.PathSearch("data_volumes", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("availability_zone", utils.PathSearch("availability_zone", respBody, nil)),
		d.Set("user_group", utils.PathSearch("user_group", respBody, nil)),
		d.Set("name", utils.PathSearch("computer_name", respBody, "").(string)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("tags", respBody, nil))),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
	)

	imageId := utils.PathSearch("metadata.\"metering.image_id\"", respBody, "").(string)
	if imageId != "" {
		mErr = multierror.Append(mErr, d.Set("image_id", imageId))
	} else {
		mErr = multierror.Append(mErr, fmt.Errorf("the image_id field does not found in metadata structure"))
	}
	securityGroups := utils.PathSearch("security_groups[*].id", respBody, make([]interface{}, 0)).([]interface{})
	if len(securityGroups) < 1 {
		mErr = multierror.Append(mErr, fmt.Errorf("the security_groups field does not found in API response"))
	} else {
		mErr = multierror.Append(mErr, d.Set("security_groups", securityGroups))
	}

	nicVal, err := getDesktopNetwork(client, desktopId)
	if err != nil {
		// This feature is not available in some region, so use log.Printf to record the error.
		log.Printf("[ERROR] %s", err)
	} else {
		mErr = multierror.Append(mErr, d.Set("nic", nicVal))
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting desktop fields: %s", err)
	}
	return nil
}

func updateDesktopFlavor(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v2/{project_id}/desktops/resize"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"desktops": []map[string]interface{}{
				{
					"desktop_id": d.Id(),
				},
			},
			"product_id": d.Get("flavor_id"),
			"mode":       "STOP_DESKTOP",
		},
	}
	resp, err := client.Request("POST", updatePath, &opts)
	if err != nil {
		return fmt.Errorf("error updating desktop product: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	_, err = waitForWorkspaceJobCompleted(ctx, client, jobId, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("error waiting for the update product job (%s) completed: %s", jobId, err)
	}
	log.Printf("[DEBUG] The update product job (%s) has been completed", jobId)
	return nil
}

const (
	RootVolume = "root_volume"
	DataVolume = "data_volume"
)

func collectVolumeChangeErrors(d *schema.ResourceData, volumeRole string) error {
	var (
		oldVal, newVal = d.GetChange(volumeRole)
		oldVolumes     = oldVal.([]interface{})
		newVolumes     = newVal.([]interface{})
		errs           *multierror.Error
		oldLen         = len(oldVolumes)
		newLen         = len(newVolumes)
	)

	if newLen < oldLen {
		errs = multierror.Append(errs, errors.New("the number of volumes cannot be reduced"))
		return errs.ErrorOrNil()
	}

	for i, val := range oldVolumes {
		oldVolume := val.(map[string]interface{})
		newVolume := newVolumes[i].(map[string]interface{})
		if newVolume["type"].(string) != oldVolume["type"].(string) {
			errs = multierror.Append(errs, errors.New("volume type does not support updates"))
		}
		if newVolume["size"].(int) < oldVolume["size"].(int) {
			errs = multierror.Append(errs, fmt.Errorf("volume (%v) size (old:%v, new:%v) cannot be smaller than the size before the change",
				oldVolume["name"], oldVolume["size"], newVolume["size"]))
		}
		if newVolume["type"].(string) != "GPSSD2" {
			isQosChange := false
			if oldVolume["iops"] != newVolume["iops"] {
				isQosChange = true
			}
			if oldVolume["throughput"] != newVolume["throughput"] {
				isQosChange = true
			}

			if isQosChange {
				errs = multierror.Append(errs, fmt.Errorf("the type of the volume (%v) is not GPSSD2, cannot set QoS options", oldVolume["name"]))
			}
		}
	}

	if newLen > oldLen {
		for i := oldLen; i < newLen; i++ {
			newVolume := newVolumes[i].(map[string]interface{})
			if newVolume["type"].(string) == "GPSSD2" {
				isQoSNotSet := false
				iops, iopsOk := newVolume["iops"].(int)
				if !iopsOk || iops == 0 {
					isQoSNotSet = true
				}
				throughput, throughputOk := newVolume["throughput"].(int)
				if !throughputOk || throughput == 0 {
					isQoSNotSet = true
				}

				if isQoSNotSet {
					errs = multierror.Append(errs,
						fmt.Errorf("the type of the volume (index number: %d) is GPSSD2, iops and throughput cannot be empty", i))
				}
			}
		}
	}

	return errs.ErrorOrNil()
}

func buildAddDesktopVolumesBody(d *schema.ResourceData, volumeRole string) map[string]interface{} {
	var (
		desktopId      = d.Id()
		oldVal, newVal = d.GetChange(volumeRole)
		oldVolumes     = oldVal.([]interface{})
		newVolumes     = newVal.([]interface{})
		oldLen         = len(oldVolumes)
		newLen         = len(newVolumes)
	)

	if newLen <= oldLen {
		return nil
	}

	volumes := make([]map[string]interface{}, 0, newLen-oldLen)
	for i := oldLen; i < newLen; i++ {
		volume := newVolumes[i].(map[string]interface{})
		volumes = append(volumes, map[string]interface{}{
			"type":       volume["type"],
			"size":       volume["size"],
			"iops":       utils.ValueIgnoreEmpty(volume["iops"]),
			"throughput": utils.ValueIgnoreEmpty(volume["throughput"]),
			"kms_id":     utils.ValueIgnoreEmpty(volume["kms_id"]),
		})
	}

	return map[string]interface{}{
		"addDesktopVolumesReq": []map[string]interface{}{
			{
				"desktop_id": desktopId,
				"volumes":    volumes,
			},
		},
	}
}

func isAddVolumeRetryableError(err error) bool {
	if apiErr, ok := err.(golangsdk.ErrDefault400); ok {
		var respBody interface{}
		if jsonErr := json.Unmarshal(apiErr.Body, &respBody); jsonErr != nil {
			log.Printf("[WARN] failed to unmarshal the response body: %s", jsonErr)
			return false
		}
		// WKS.0407: This operation cannot be performed because disks are being added.
		errCode := utils.PathSearch("error_code", respBody, "")
		if errCode == "WKS.0407" {
			return true
		}
	}
	return false
}

// lintignore:R018
func updateNewVolume(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, volumeRole string) error {
	var httpUrl = "v2/{project_id}/volumes"

	addVolumeBody := buildAddDesktopVolumesBody(d, volumeRole)
	if addVolumeBody == nil {
		return nil
	}

	addVolumePath := client.Endpoint + httpUrl
	addVolumePath = strings.ReplaceAll(addVolumePath, "{project_id}", client.ProjectID)
	addVolumeOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(addVolumeBody),
	}

	var (
		requestResp *http.Response
		retryErr    error
	)
	err := resource.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		if requestResp, retryErr = client.Request("POST", addVolumePath, &addVolumeOpts); retryErr != nil {
			if isAddVolumeRetryableError(retryErr) {
				time.Sleep(5 * time.Minute)
				return resource.RetryableError(
					errors.New("the desktop is adding volumes, the time interval for the each operation must over 5 minutes"))
			}
			return resource.NonRetryableError(retryErr)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("unable to add volumes for desktop: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	_, err = waitForWorkspaceJobCompleted(ctx, client, jobId, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("error waiting for the add volumes job (%s) completed: %s", jobId, err)
	}

	log.Printf("[DEBUG] The add volumes job (%s) has been completed", jobId)
	return nil
}

func buildExpandDesktopVolumesBody(d *schema.ResourceData, volumeRole string) map[string]interface{} {
	var (
		desktopId      = d.Id()
		oldVal, newVal = d.GetChange(volumeRole)
		oldVolumes     = oldVal.([]interface{})
		newVolumes     = newVal.([]interface{})
	)

	expandSlice := make([]map[string]interface{}, 0)
	for i, val := range oldVolumes {
		oldVolume := val.(map[string]interface{})
		newVolume := newVolumes[i].(map[string]interface{})
		if newVolume["size"].(int) > oldVolume["size"].(int) {
			expandSlice = append(expandSlice, map[string]interface{}{
				"desktop_id": desktopId,
				"volume_id":  oldVolume["id"].(string),
				"new_size":   newVolume["size"].(int),
			})
		}
	}

	if len(expandSlice) == 0 {
		return nil
	}

	return map[string]interface{}{
		"desktop_volumes_expansion": expandSlice,
	}
}

func updateVolumeSizeChange(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, volumeRole string) error {
	var httpUrl = "v2/{project_id}/volumes/expand"

	expandBody := buildExpandDesktopVolumesBody(d, volumeRole)
	if expandBody == nil {
		return nil
	}

	expandPath := client.Endpoint + httpUrl
	expandPath = strings.ReplaceAll(expandPath, "{project_id}", client.ProjectID)
	expandOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(expandBody),
	}

	log.Printf("[DEBUG] The new expandOpts is: %#v", expandOpts)
	resp, err := client.Request("POST", expandPath, &expandOpts)
	if err != nil {
		return fmt.Errorf("unable to expand volume size: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	_, err = waitForWorkspaceJobCompleted(ctx, client, jobId, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("error waiting for the job (%s) completed: %s", jobId, err)
	}

	log.Printf("[DEBUG] The expand volume job (%s) has been completed", jobId)
	return nil
}

func buildUpdateDesktopVolumesQosBody(d *schema.ResourceData, volumeRole string) []map[string]interface{} {
	var (
		oldVal, newVal = d.GetChange(volumeRole)
		oldVolumes     = oldVal.([]interface{})
		newVolumes     = newVal.([]interface{})
	)

	qosChangeSlice := make([]map[string]interface{}, 0)
	for i, val := range oldVolumes {
		oldVolume := val.(map[string]interface{})
		newVolume := newVolumes[i].(map[string]interface{})

		if newVolume["type"].(string) != "GPSSD2" {
			continue
		}

		if newVolume["iops"] == oldVolume["iops"] && newVolume["throughput"] == oldVolume["throughput"] {
			continue
		}

		qosChangeSlice = append(qosChangeSlice, map[string]interface{}{
			"volume_ids": []string{oldVolume["id"].(string)},
			"qos": map[string]interface{}{
				"iops":       newVolume["iops"].(int),
				"throughput": newVolume["throughput"].(int),
			},
		})
	}

	return qosChangeSlice
}

func updateVolumeQosChange(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, volumeRole string) error {
	var httpUrl = "v2/{project_id}/volumes/batch-modify-qos"

	qosChangeSlice := buildUpdateDesktopVolumesQosBody(d, volumeRole)
	if len(qosChangeSlice) == 0 {
		return nil
	}

	qosChangePath := client.Endpoint + httpUrl
	qosChangePath = strings.ReplaceAll(qosChangePath, "{project_id}", client.ProjectID)
	for _, qosChange := range qosChangeSlice {
		qosChangeOpts := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
			JSONBody: utils.RemoveNil(qosChange),
		}

		log.Printf("[DEBUG] The new QoSChangeOpts is: %#v", qosChangeOpts)
		resp, err := client.Request("PUT", qosChangePath, &qosChangeOpts)
		if err != nil {
			return fmt.Errorf("unable to change volume QoS: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return err
		}

		jobId := utils.PathSearch("job_id", respBody, "").(string)
		_, err = waitForWorkspaceJobCompleted(ctx, client, jobId, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error waiting for the change volume QoS job (%s) completed: %s", jobId, err)
		}
		log.Printf("[DEBUG] The change volume QoS job (%s) has been completed", jobId)
	}

	return nil
}

func updateDesktopVolumes(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, volumeRole string) error {
	err := collectVolumeChangeErrors(d, volumeRole)
	if err != nil {
		return err
	}

	err = updateNewVolume(ctx, client, d, volumeRole)
	if err != nil {
		return err
	}

	err = updateVolumeSizeChange(ctx, client, d, volumeRole)
	if err != nil {
		return err
	}

	err = updateVolumeQosChange(ctx, client, d, volumeRole)
	if err != nil {
		return err
	}

	return nil
}

func updateDesktopNetwork(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	nicRaw := d.Get("nic").([]interface{})
	if len(nicRaw) < 1 {
		return nil
	}

	var (
		httpUrl   = "v2/{project_id}/desktops/{desktop_id}/networks"
		desktopId = d.Id()
		opts      = golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
			JSONBody: map[string]interface{}{
				"vpc_id":             d.Get("vpc_id").(string),
				"subnet_id":          utils.PathSearch("network_id", nicRaw[0], nil),
				"security_group_ids": utils.ExpandToStringList(d.Get("security_groups").(*schema.Set).List()),
			},
		}
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{desktop_id}", desktopId)
	resp, err := client.Request("PUT", updatePath, &opts)
	if err != nil {
		return fmt.Errorf("error updating the network of the Workspace desktop (%s): %s", desktopId, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	_, err = waitForWorkspaceJobCompleted(ctx, client, jobId, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("error waiting for the job (%s) completed: %s", jobId, err)
	}
	log.Printf("[DEBUG] The job (%s) has been completed", jobId)
	return nil
}

func waitForWorkspaceStatusCompleted(ctx context.Context, client *golangsdk.ServiceClient, desktopId, powerAction string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      desktopStatusRefreshFunc(client, desktopId, powerAction),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
		// The final status of both startup and restart is "ACTIVE". the parameter only applies to "REBOOT" action.
		ContinuousTargetOccurence: 2,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func desktopStatusRefreshFunc(client *golangsdk.ServiceClient, desktopId, powerAction string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetDesktopById(client, desktopId)
		if err != nil {
			return respBody, "", err
		}

		// In statusMap, key represents the power action of the desktop, and value represents the status after the desktop operation is completed.
		// Use statusMap to make a mapping relationship between the power action of the desktop and the final status of the desktop.
		statusMap := map[string]string{
			"os-start":     "ACTIVE",
			"os-stop":      "SHUTOFF",
			"reboot":       "ACTIVE",
			"os-hibernate": "HIBERNATED",
		}

		// TaskStatus variable is always an empty string when the desktop power action is completed.
		// If the desktop power action changes from one state to another, taskStatus is an empty string for a long time,
		// whether a desktop action is completed cannot be determined only by taskStatus.
		taskStatus := utils.PathSearch("task_status", respBody, "").(string)
		status := utils.PathSearch("status", respBody, "").(string)
		if taskStatus == "" && status == statusMap[powerAction] {
			return respBody, "COMPLETED", nil
		}
		return respBody, "PENDING", nil
	}
}

func updateDesktopPowerAction(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, timeout time.Duration) error {
	var (
		httpUrl   = "v2/{project_id}/desktops/action"
		desktopId = d.Id()
		action    = d.Get("power_action").(string)
		params    = map[string]interface{}{
			"desktop_ids": []string{desktopId},
			"op_type":     action,
			"type":        utils.ValueIgnoreEmpty(d.Get("power_action_type")),
		}
		opts = golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
			JSONBody: utils.RemoveNil(params),
		}
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	_, err := client.Request("POST", updatePath, &opts)
	if err != nil {
		return fmt.Errorf("error updating the power action of the Workspace desktop (%s): %s", desktopId, err)
	}

	err = waitForWorkspaceStatusCompleted(ctx, client, desktopId, action, timeout)
	if err != nil {
		return fmt.Errorf("error waiting for power action (%s) for desktop (%s) failed: %s", action, desktopId, err)
	}
	return nil
}

func updateDesktopImage(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl     = "v2/{project_id}/desktops/rebuild"
		desktopId   = d.Id()
		rebuildOpts = golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
			JSONBody: map[string]interface{}{
				"desktop_ids": []string{desktopId},
				"image_type":  d.Get("image_type"),
				"image_id":    d.Get("image_id"),
			},
		}
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	resp, err := client.Request("POST", updatePath, &rebuildOpts)
	if err != nil {
		return fmt.Errorf("error rebuild Workspace desktop: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	_, err = waitForWorkspaceJobCompleted(ctx, client, jobId, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("error waiting for the rebuild desktop job (%s) completed: %s", jobId, err)
	}
	log.Printf("[DEBUG] The rebuild desktop job (%s) has been completed", jobId)
	return nil
}

func resourceDesktopUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	desktopId := d.Id()
	if d.HasChange("flavor_id") {
		if err = updateDesktopFlavor(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("image_type", "image_id") {
		if err = updateDesktopImage(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("root_volume", "data_volume") {
		if err = updateDesktopVolumes(ctx, client, d, RootVolume); err != nil {
			return diag.FromErr(err)
		}
		if err = updateDesktopVolumes(ctx, client, d, DataVolume); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		err = utils.UpdateResourceTags(client, d, "desktops", desktopId)
		if err != nil {
			return diag.Errorf("error updating tags of Workspace desktop (%s): %s", desktopId, err)
		}
	}

	if d.HasChange("nic") {
		err = updateDesktopNetwork(ctx, client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   desktopId,
			ResourceType: "workspace-desktop",
			RegionId:     region,
			ProjectId:    client.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("power_action") {
		err = updateDesktopPowerAction(ctx, client, d, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDesktopRead(ctx, d, meta)
}

func waitForDesktopDeleted(ctx context.Context, client *golangsdk.ServiceClient, desktopId string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"ACTIVE", "DELETING", "SHUTOFF", "HIBERNATED"},
		Target:       []string{"DELETED"},
		Refresh:      refreshDesktopStatusFunc(client, desktopId),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func refreshDesktopStatusFunc(client *golangsdk.ServiceClient, desktopId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := GetDesktopById(client, desktopId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return "deleted", "DELETED", nil
			}
			return resp, "ERROR", err
		}
		// During the removal process of desktop, the workspace service cannot perceive the ECS mechine and the API will
		// return an empty status.
		status := utils.PathSearch("status", resp, "").(string)
		if status == "" {
			return resp, "DELETING", nil
		}
		// The uppercase characters is the default format for attribute 'status' in the API response.
		return resp, strings.ToUpper(status), nil
	}
}

func waitForDesktopUserDeleted(ctx context.Context, client *golangsdk.ServiceClient, userName string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"ACTIVE"},
		Target:       []string{"DELETED"},
		Refresh:      refreshDesktopUserStatusFunc(client, userName),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func refreshDesktopUserStatusFunc(client *golangsdk.ServiceClient, userName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			httpUrl  = "v2/{project_id}/users"
			listOpts = golangsdk.RequestOpts{
				KeepResponseBody: true,
				MoreHeaders: map[string]string{
					"Content-Type": "application/json",
				},
			}
		)
		listPath := client.Endpoint + httpUrl
		listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
		listPath = fmt.Sprintf("%s?user_name=%s", listPath, userName)
		resp, err := client.Request("GET", listPath, &listOpts)
		if err != nil {
			return resp, "ERROR", err
		}
		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return resp, "ERROR", err
		}

		if len(utils.PathSearch("users", respBody, make([]interface{}, 0)).([]interface{})) < 1 {
			return resp, "DELETED", nil
		}
		return resp, "ACTIVE", nil
	}
}

func resourceDesktopDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.NewServiceClient("workspace", conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	var (
		httpUrl      = "v2/{project_id}/desktops/{desktop_id}"
		desktopId    = d.Id()
		isDeleteUser = d.Get("delete_user").(bool)
		opts         = golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		}
	)

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{desktop_id}", desktopId)
	deletePath = fmt.Sprintf("%s?delete_users=%v&email_notification=%v", deletePath, isDeleteUser, d.Get("email_notification").(bool))
	_, err = client.Request("DELETE", deletePath, &opts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting desktop (%s)", desktopId))
	}
	// Make sure the desktop has been deleted.
	err = waitForDesktopDeleted(ctx, client, desktopId, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("unable to delete desktop (%s): %s", desktopId, err)
	}

	if isDeleteUser {
		// Make sure the related user has been deleted.
		userName := d.Get("user_name").(string)
		err = waitForDesktopUserDeleted(ctx, client, userName, d.Timeout(schema.TimeoutDelete))
		if err != nil {
			return diag.Errorf("unable to delete user under desktop (%s): %s", desktopId, err)
		}
	}
	return nil
}
