package evs

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/block_devices"
	ecsjobs "github.com/chnsz/golangsdk/openstack/ecs/v1/jobs"
	"github.com/chnsz/golangsdk/openstack/evs/v2/cloudvolumes"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EVS GET /v1/{project_id}/jobs/{job_id}
// @API EVS POST /v2.1/{project_id}/cloudvolumes/{volume_id}/action
// @API EVS POST /v2/{project_id}/cloudvolumes/{volume_id}/tags/action
// @API EVS GET /v2/{project_id}/cloudvolumes/{volume_id}
// @API EVS PUT /v2/{project_id}/cloudvolumes/{volume_id}
// @API EVS DELETE /v2/{project_id}/cloudvolumes/{id}
// @API EVS POST /v2.1/{project_id}/cloudvolumes
// @API EVS PUT /v5/{project_id}/cloudvolumes/{volume_id}/qos
// @API ECS DELETE /v1/{project_id}/cloudservers/{serverId}/detachvolume/{volumeId}
// @API ECS GET /v1/{project_id}/jobs/{jobId}
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/suscriptions/resources/query
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{resource_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{resource_id}
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
func ResourceEvsVolume() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEvsVolumeCreate,
		ReadContext:   resourceEvsVolumeRead,
		UpdateContext: resourceEvsVolumeUpdate,
		DeleteContext: resourceEvsVolumeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(180 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"iops": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"throughput": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"device_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "VBD",
				ValidateFunc: validation.StringInSlice([]string{"VBD", "SCSI"}, true),
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"size", "backup_id"},
			},
			"backup_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"kms_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"multiattach": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"dedicated_storage_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),
			"auto_pay":      common.SchemaAutoPay(nil),
			"tags":          common.TagsSchema(),
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"cascade": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"attachment": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"device": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: resourceVolumeAttachmentHash,
			},
			"wwn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dedicated_storage_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVolumeAttachmentHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if m["instance_id"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["instance_id"].(string)))
	}
	return hashcode.String(buf.String())
}

func buildVolumeMetadataBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"create_for_volume_id": "true",
	}

	if v, ok := d.GetOk("kms_id"); ok {
		bodyParams["__system__cmkid"] = v
		bodyParams["__system__encrypted"] = "1"
	}

	if d.Get("device_type").(string) == "SCSI" {
		bodyParams["hw:passthrough"] = "true"
	}

	return bodyParams
}

func buildVolumeBodyParams(cfg *config.Config, d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"availability_zone":     d.Get("availability_zone"),
		"volume_type":           d.Get("volume_type"),
		"name":                  utils.ValueIgnoreEmpty(d.Get("name")),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"size":                  utils.ValueIgnoreEmpty(d.Get("size")),
		"backup_id":             utils.ValueIgnoreEmpty(d.Get("backup_id")),
		"snapshot_id":           utils.ValueIgnoreEmpty(d.Get("snapshot_id")),
		"imageRef":              utils.ValueIgnoreEmpty(d.Get("image_id")),
		"iops":                  utils.ValueIgnoreEmpty(d.Get("iops")),
		"throughput":            utils.ValueIgnoreEmpty(d.Get("throughput")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		"tags":                  utils.ValueIgnoreEmpty(utils.ExpandToStringMap(d.Get("tags").(map[string]interface{}))),
		"metadata":              utils.ValueIgnoreEmpty(buildVolumeMetadataBodyParams(d)),
	}

	if d.Get("multiattach").(bool) {
		bodyParams["multiattach"] = true
	}

	return bodyParams
}

func buildBssParamBodyParams(d *schema.ResourceData) map[string]interface{} {
	if v, ok := d.GetOk("charging_mode"); ok && v == "prePaid" {
		return map[string]interface{}{
			"chargingMode": d.Get("charging_mode"),
			"periodType":   utils.ValueIgnoreEmpty(d.Get("period_unit")),
			"periodNum":    utils.ValueIgnoreEmpty(d.Get("period")),
			"isAutoRenew":  utils.ValueIgnoreEmpty(d.Get("auto_renew")),
			"isAutoPay":    utils.ValueIgnoreEmpty(common.GetAutoPay(d)),
		}
	}

	return nil
}

func buildSchedulerBodyParams(d *schema.ResourceData) map[string]interface{} {
	if v, ok := d.GetOk("dedicated_storage_id"); ok {
		return map[string]interface{}{
			"dedicated_storage_id": v,
		}
	}

	return nil
}

func buildCreateVolumeBodyParams(cfg *config.Config, d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"volume":                     buildVolumeBodyParams(cfg, d),
		"server_id":                  utils.ValueIgnoreEmpty(d.Get("server_id")),
		"bssParam":                   buildBssParamBodyParams(d),
		"OS-SCH-HNT:scheduler_hints": buildSchedulerBodyParams(d),
	}

	return bodyParams
}

func getJobDetail(client *golangsdk.ServiceClient, jobID string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/jobs/{job_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", jobID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error querying EVS job detail: %s", err)
	}

	return utils.FlattenResponse(resp)
}

func waitingForVolumeJobSuccess(ctx context.Context, client *golangsdk.ServiceClient, jobID string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := getJobDetail(client, jobID)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("status", respBody, "").(string)
			if status == "" {
				return respBody, "ERROR", fmt.Errorf("status is not found in EVS job (%s) detail API response", jobID)
			}

			if status == "SUCCESS" {
				return respBody, "COMPLETED", nil
			}

			if status == "FAIL" {
				return respBody, status, fmt.Errorf("the EVS job (%s) status is FAIL, the fail reason is: %s",
					jobID, utils.PathSearch("fail_reason", respBody, "").(string))
			}

			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func GetVolumeDetail(client *golangsdk.ServiceClient, volumeID string) (interface{}, error) {
	requestPath := client.Endpoint + "v2/{project_id}/cloudvolumes/{volume_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{volume_id}", volumeID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func waitingForEvsVolumeComplete(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	errorStatuses := []string{"error", "error_restoring", "error_extending", "error_deleting", "error_rollbacking"}
	successStatuses := []string{"available", "in-use"}
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := GetVolumeDetail(client, d.Id())
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("volume.status", respBody, "").(string)
			if status == "" {
				return respBody, "ERROR", fmt.Errorf("status is not found in EVS volume (%s) detail API response", d.Id())
			}

			if utils.StrSliceContains(errorStatuses, status) {
				return respBody, status, fmt.Errorf("unexpect status (%s)", status)
			}

			if utils.StrSliceContains(successStatuses, status) {
				return respBody, "COMPLETED", nil
			}

			return respBody, "PENDING", nil
		},
		Timeout:                   timeout,
		Delay:                     10 * time.Second,
		PollInterval:              10 * time.Second,
		ContinuousTargetOccurence: 2,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceEvsVolumeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2.1/{project_id}/cloudvolumes"
		product = "evs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateVolumeBodyParams(cfg, d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating EVS volume: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	volumeID := utils.PathSearch("volume_ids|[0]", respBody, "").(string)
	if volumeID == "" {
		return diag.Errorf("error creating EVS volume: ID is not found in API response")
	}
	d.SetId(volumeID)

	if orderID := utils.PathSearch("order_id", respBody, "").(string); orderID != "" {
		bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		if err = common.WaitOrderComplete(ctx, bssClient, orderID, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.Errorf("the order is not completed while creating EVS volume (%s): %v", d.Id(), err)
		}
		if _, err = common.WaitOrderAllResourceComplete(ctx, bssClient, orderID, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.FromErr(err)
		}
	}

	if jobID := utils.PathSearch("job_id", respBody, "").(string); jobID != "" {
		if err := waitingForVolumeJobSuccess(ctx, client, jobID, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.Errorf("error waiting for EVS volume (%s) job success: %s", d.Id(), err)
		}
	}

	if err := waitingForEvsVolumeComplete(ctx, client, d, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for the creation of EVS volume (%s) to complete: %s", d.Id(), err)
	}

	return resourceEvsVolumeRead(ctx, d, meta)
}

func setVolumeChargingMode(d *schema.ResourceData, meta interface{}) error {
	if utils.PathSearch("volume.metadata.orderID", meta, "").(string) != "" {
		return d.Set("charging_mode", "prePaid")
	}
	return nil
}

func flattenVolumeDeviceType(respBody interface{}) interface{} {
	if utils.PathSearch("volume.metadata.\"hw:passthrough\"", respBody, "").(string) == "true" {
		return "SCSI"
	}

	return "VBD"
}

func flattenVolumeAttachment(respBody interface{}) interface{} {
	attachments := utils.PathSearch("volume.attachments", respBody, make([]interface{}, 0)).([]interface{})
	result := make([]map[string]interface{}, len(attachments))
	for i, attachment := range attachments {
		result[i] = map[string]interface{}{
			"id":          utils.PathSearch("attachment_id", attachment, nil),
			"instance_id": utils.PathSearch("server_id", attachment, nil),
			"device":      utils.PathSearch("device", attachment, nil),
		}
	}

	return result
}

func resourceEvsVolumeRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "evs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	respBody, err := GetVolumeDetail(client, d.Id())
	if err != nil {
		// When the resource does not exist, calling the query API will return a `404` status code.
		return common.CheckDeletedDiag(d, err, "error retrieving EVS volume")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("volume.name", respBody, nil)),
		d.Set("size", utils.PathSearch("volume.size", respBody, nil)),
		d.Set("description", utils.PathSearch("volume.description", respBody, nil)),
		d.Set("availability_zone", utils.PathSearch("volume.availability_zone", respBody, nil)),
		d.Set("snapshot_id", utils.PathSearch("volume.snapshot_id", respBody, nil)),
		d.Set("volume_type", utils.PathSearch("volume.volume_type", respBody, nil)),
		d.Set("iops", utils.PathSearch("volume.iops.total_val", respBody, nil)),
		d.Set("throughput", utils.PathSearch("volume.throughput.total_val", respBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("volume.enterprise_project_id", respBody, nil)),
		d.Set("wwn", utils.PathSearch("volume.wwn", respBody, nil)),
		d.Set("multiattach", utils.PathSearch("volume.multiattach", respBody, nil)),
		d.Set("tags", utils.PathSearch("volume.tags", respBody, nil)),
		d.Set("dedicated_storage_id", utils.PathSearch("volume.dedicated_storage_id", respBody, nil)),
		d.Set("dedicated_storage_name", utils.PathSearch("volume.dedicated_storage_name", respBody, nil)),
		d.Set("status", utils.PathSearch("volume.status", respBody, nil)),
		setVolumeChargingMode(d, respBody),
		d.Set("device_type", flattenVolumeDeviceType(respBody)),
		d.Set("image_id", utils.PathSearch("volume.volume_image_metadata.image_id", respBody, nil)),
		d.Set("attachment", flattenVolumeAttachment(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateVolumeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"description": d.Get("description"),
	}

	return map[string]interface{}{
		"volume": bodyParams,
	}
}

func updateVolume(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v2/{project_id}/cloudvolumes/{volume_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{volume_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateVolumeBodyParams(d)),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	return err
}

func buildUpdateVolumeSizeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"os-extend": map[string]interface{}{
			"new_size": d.Get("size"),
		},
	}

	if d.Get("charging_mode").(string) == "prePaid" {
		bodyParams["bssParam"] = map[string]interface{}{
			"isAutoPay": common.GetAutoPay(d),
		}
	}

	return bodyParams
}

func extendVolumeSize(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, cfg *config.Config) error {
	requestPath := client.Endpoint + "v2.1/{project_id}/cloudvolumes/{volume_id}/action"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{volume_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUpdateVolumeSizeBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	if orderID := utils.PathSearch("order_id", respBody, "").(string); orderID != "" {
		bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
		if err != nil {
			return fmt.Errorf("error creating BSS v2 client: %s", err)
		}

		if err = common.WaitOrderComplete(ctx, bssClient, orderID, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return fmt.Errorf("error waiting for the order (%s) to complete: %s", orderID, err)
		}
	}

	if jobID := utils.PathSearch("job_id", respBody, "").(string); jobID != "" {
		if err := waitingForVolumeJobSuccess(ctx, client, jobID, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return fmt.Errorf("error waiting for the job (%s) to succeed: %s", jobID, err)
		}
	}

	if err := waitingForEvsVolumeComplete(ctx, client, d, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return fmt.Errorf("error waiting for the EVS volume to complete: %s", err)
	}

	return nil
}

func buildUpdateVolumeTypeBodyParams(d *schema.ResourceData) map[string]interface{} {
	osRetypeParam := map[string]interface{}{
		"new_type": d.Get("volume_type"),
	}

	if d.Get("volume_type") == "GPSSD2" {
		osRetypeParam["iops"] = d.Get("iops")
		osRetypeParam["throughput"] = d.Get("throughput")
	}

	// Currently, EVS does not support changing the disk type to ESSD2.
	// However, the documentation does not clearly state this limitation, so keep this code.
	if d.Get("volume_type") == "ESSD2" {
		osRetypeParam["iops"] = d.Get("iops")
	}

	bodyParams := map[string]interface{}{
		"os-retype": osRetypeParam,
	}

	if d.Get("charging_mode").(string) == "prePaid" {
		bodyParams["bssParam"] = map[string]interface{}{
			"isAutoPay": common.GetAutoPay(d),
		}
	}

	return bodyParams
}

func updateVolumeType(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, cfg *config.Config) error {
	// Interface constraints: QoS can be updated only when the volume status is available or in-use
	if err := waitingForEvsVolumeComplete(ctx, client, d, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return fmt.Errorf("error waiting for the EVS volume to complete before updating volume type: %s", err)
	}

	requestPath := client.Endpoint + "v2/{project_id}/volumes/{volume_id}/retype"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{volume_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUpdateVolumeTypeBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	if orderID := utils.PathSearch("order_id", respBody, "").(string); orderID != "" {
		bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
		if err != nil {
			return fmt.Errorf("error creating BSS v2 client: %s", err)
		}

		if err = common.WaitOrderComplete(ctx, bssClient, orderID, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return fmt.Errorf("error waiting for the order (%s) to complete: %s", orderID, err)
		}
	}

	if jobID := utils.PathSearch("job_id", respBody, "").(string); jobID != "" {
		if err := waitingForVolumeJobSuccess(ctx, client, jobID, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return fmt.Errorf("error waiting for the job (%s) to succeed: %s", jobID, err)
		}
	}

	if err := waitingForEvsVolumeComplete(ctx, client, d, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return fmt.Errorf("error waiting for the EVS volume to complete after updating volume type: %s", err)
	}

	return nil
}

func buildUpdateVolumeQosBodyParams(d *schema.ResourceData) map[string]interface{} {
	qosParam := map[string]interface{}{
		"iops":       d.Get("iops"),
		"throughput": utils.ValueIgnoreEmpty(d.Get("throughput")),
	}

	return map[string]interface{}{
		"qos_modify": qosParam,
	}
}

func updateVolumeQoS(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	// Interface constraints: QoS can be updated only when the volume status is available or in-use
	if err := waitingForEvsVolumeComplete(ctx, client, d, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return fmt.Errorf("error waiting for the EVS volume to complete before updating volume qos: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/cloudvolumes/{volume_id}/qos"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{volume_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateVolumeQosBodyParams(d)),
	}

	resp, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	if jobID := utils.PathSearch("job_id", respBody, "").(string); jobID != "" {
		if err := waitingForVolumeJobSuccess(ctx, client, jobID, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return fmt.Errorf("error waiting for the job (%s) to succeed: %s", jobID, err)
		}
	}

	if err := waitingForEvsVolumeComplete(ctx, client, d, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return fmt.Errorf("error waiting for the EVS volume to complete after updating volume qos: %s", err)
	}

	return nil
}

func resourceEvsVolumeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "evs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	if d.HasChanges("name", "description") {
		if err := updateVolume(client, d); err != nil {
			return diag.Errorf("error updating EVS volume (%s): %s", d.Id(), err)
		}
	}

	if d.HasChange("tags") {
		if err := utils.UpdateResourceTags(client, d, "cloudvolumes", d.Id()); err != nil {
			return diag.Errorf("error updating EVS volume (%s) tags: %s", d.Id(), err)
		}
	}

	if d.HasChange("size") {
		if err := extendVolumeSize(ctx, client, d, cfg); err != nil {
			return diag.Errorf("error extending EVS volume (%s) size: %s", d.Id(), err)
		}
	}

	// Changing this field requires prerequisites, see the documentation for details.
	// Changing this field may use the fields "iops" and "throughput", so put this change before "iops" and "throughput".
	if d.HasChange("volume_type") {
		if err := updateVolumeType(ctx, client, d, cfg); err != nil {
			return diag.Errorf("error updating EVS volume (%s) type: %s", d.Id(), err)
		}
	}

	if d.HasChanges("iops", "throughput") {
		if err := updateVolumeQoS(ctx, client, d); err != nil {
			return diag.Errorf("error updating EVS volume (%s) qos: %s", d.Id(), err)
		}
	}

	if d.HasChange("auto_renew") {
		bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), d.Id()); err != nil {
			return diag.Errorf("error updating the auto-renew of the volume (%s): %s", d.Id(), err)
		}
	}

	return resourceEvsVolumeRead(ctx, d, meta)
}

func resourceEvsVolumeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	evsV2Client, err := cfg.BlockStorageV2Client(region)
	if err != nil {
		return diag.Errorf("eError creating block storage v2 client: %s", err)
	}

	v, err := cloudvolumes.Get(evsV2Client, d.Id()).Extract()
	if err != nil {
		// Before deleting a resource, check if the resource exists first,
		// if resource does not exist, perform checkDeleted processing.
		// When the resource does not exist, calling the query API will return a `404` status code.
		return common.CheckDeletedDiag(d, err, "EVS volume")
	}

	// Make sure this volume is detached from all instances before deleting.
	if len(v.Attachments) > 0 {
		log.Printf("[DEBUG] Start to detaching volumes.")
		computeClient, err := cfg.ComputeV1Client(region)
		if err != nil {
			return diag.Errorf("error creating ECS v2 client: %s", err)
		}
		for _, attachment := range v.Attachments {
			log.Printf("[DEBUG] The attachment is: %v", attachment)
			opts := block_devices.DetachOpts{
				ServerId: attachment.ServerID,
			}
			job, err := block_devices.Detach(computeClient, attachment.VolumeID, opts)
			if err != nil {
				return diag.FromErr(err)
			}
			stateConf := &resource.StateChangeConf{
				Pending:    []string{"RUNNING"},
				Target:     []string{"SUCCESS", "NOTFOUND"},
				Refresh:    AttachmentJobRefreshFunc(computeClient, job.ID),
				Timeout:    10 * time.Minute,
				Delay:      10 * time.Second,
				MinTimeout: 3 * time.Second,
			}
			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.Get("charging_mode").(string) == "prePaid" {
		err = common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()})
		if err != nil {
			return diag.Errorf("error unsubscribing EVS volume : %s", err)
		}
	} else {
		// The snapshots associated with the disk are deleted together with the EVS disk if cascade value is true
		deleteOpts := cloudvolumes.DeleteOpts{
			Cascade: d.Get("cascade").(bool),
		}
		// It's possible that this volume was used as a boot device and is currently
		// in a "deleting" state from when the instance was terminated.
		// If this is true, just move on. It'll eventually delete.
		if v.Status != "deleting" {
			if err := cloudvolumes.Delete(evsV2Client, d.Id(), deleteOpts).ExtractErr(); err != nil {
				return common.CheckDeletedDiag(d, err, "EVS volume")
			}
		}
	}

	// Wait for the volume to delete before moving on.
	log.Printf("[DEBUG] Waiting for the EVS volume (%s) to delete", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshVolumeDeleteFunc(evsV2Client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the EVS volume (%s) to delete: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func AttachmentJobRefreshFunc(c *golangsdk.ServiceClient, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := ecsjobs.Get(c, jobId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return resp, "NOTFOUND", nil
			}
			return resp, "ERROR", err
		}

		return resp, resp.Status, nil
	}
}

func CloudVolumeRefreshFunc(c *golangsdk.ServiceClient, volumeId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		response, err := cloudvolumes.Get(c, volumeId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return response, "deleted", nil
			}
			return response, "ERROR", err
		}
		if response != nil {
			return response, response.Status, nil
		}
		return response, "ERROR", nil
	}
}

func refreshVolumeDeleteFunc(c *golangsdk.ServiceClient, volumeId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		response, err := cloudvolumes.Get(c, volumeId).Extract()
		if err != nil {
			var errDefault404 golangsdk.ErrDefault404
			if errors.As(err, &errDefault404) {
				return "deleted", "COMPLETED", nil
			}
			return response, "ERROR", err
		}
		return response, "PENDING", nil
	}
}

func refreshVolumeStatusFunc(c *golangsdk.ServiceClient, volumeId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		response, err := cloudvolumes.Get(c, volumeId).Extract()
		if err != nil {
			var errDefault404 golangsdk.ErrDefault404
			if errors.As(err, &errDefault404) {
				return response, "deleted", nil
			}
			return response, "ERROR", err
		}
		if response == nil {
			return response, "ERROR", nil
		}

		errorStatus := []string{"error", "error_restoring", "error_extending", "error_deleting", "error_rollbacking"}
		status := response.Status
		if utils.StrSliceContains(errorStatus, status) {
			return response, status, fmt.Errorf("unexpect status (%s)", status)
		}

		if utils.StrSliceContains([]string{"available", "in-use"}, status) {
			return response, "COMPLETED", nil
		}
		return response, "PENDING", nil
	}
}
