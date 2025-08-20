package evs

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var v3VolumeNonUpdatableParams = []string{
	"volume_type",
	"availability_zone",
	"disaster_recovery_azs",
	"image_id",
	"multiattach",
	"snapshot_id",
	"iops",
	"throughput",
	"dedicated_storage_id",
}

// @API EVS POST /v3/{project_id}/volumes
// @API EVS GET /v3/{project_id}/volumes/{volume_id}
// @API EVS PUT /v3/{project_id}/volumes/{volume_id}
// @API EVS POST /v3/{project_id}/volumes/{volume_id}/action
// @API EVS DELETE /v3/{project_id}/volumes/{volume_id}
// @API EVS POST /v2/{project_id}/cloudvolumes/{volume_id}/tags/action
// @API EVS GET /v2/{project_id}/cloudvolumes/{volume_id}/tags
// @API ECS DELETE /v1/{project_id}/cloudservers/{server_id}/detachvolume/{volume_id}
// @API ECS GET /v1/{project_id}/jobs/{job_id}
func ResourceV3Volume() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV3VolumeCreate,
		ReadContext:   resourceV3VolumeRead,
		UpdateContext: resourceV3VolumeUpdate,
		DeleteContext: resourceV3VolumeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(v3VolumeNonUpdatableParams),
			config.MergeDefaultTags(),
		),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"disaster_recovery_azs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// `description` can be left blank.
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Named `imageRef` in the API documentation, which does not comply with the schema specification.
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// Field `metadata` is defined as a map here, although it is an object type in the API documentation.
			// Some fields in the object structure have names that do not conform to the schema specification,
			// defining it as a map is simpler and more convenient.
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"multiattach": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"dedicated_storage_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// `cascade` parameter used for deletion.
			"cascade": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": common.TagsSchema(),
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			// Attributes.
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attachments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attached_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attachment_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"device": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"bootable": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_image_metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"iops_attribute": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"frozened": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"total_val": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"volume_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"throughput_attribute": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"frozened": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"total_val": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"volume_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildV3VolumeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"volume_type":           d.Get("volume_type"),
		"availability_zone":     utils.ValueIgnoreEmpty(d.Get("availability_zone")),
		"disaster_recovery_azs": utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("disaster_recovery_azs").([]interface{}))),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"image_id":              utils.ValueIgnoreEmpty(d.Get("image_id")),
		"metadata":              utils.ValueIgnoreEmpty(utils.ExpandToStringMap(d.Get("metadata").(map[string]interface{}))),
		"multiattach":           d.Get("multiattach"),
		"name":                  utils.ValueIgnoreEmpty(d.Get("name")),
		"size":                  utils.ValueIgnoreEmpty(d.Get("size")),
		"snapshot_id":           utils.ValueIgnoreEmpty(d.Get("snapshot_id")),
		"iops":                  utils.ValueIgnoreEmpty(d.Get("iops")),
		"throughput":            utils.ValueIgnoreEmpty(d.Get("throughput")),
	}

	return bodyParams
}

func buildV3VolumeSchedulerBodyParams(d *schema.ResourceData) map[string]interface{} {
	if v, ok := d.GetOk("dedicated_storage_id"); ok {
		return map[string]interface{}{
			"dedicated_storage_id": v,
		}
	}

	return nil
}

func buildCreateV3VolumeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"volume":                     buildV3VolumeBodyParams(d),
		"OS-SCH-HNT:scheduler_hints": buildV3VolumeSchedulerBodyParams(d),
	}

	return bodyParams
}

func GetV3VolumeDetail(client *golangsdk.ServiceClient, volumeID string) (interface{}, error) {
	requestPath := client.Endpoint + "v3/{project_id}/volumes/{volume_id}"
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

func waitingForV3VolumeComplete(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	errorStatuses := []string{"error", "error_restoring", "error_extending", "error_deleting", "error_rollbacking"}
	successStatuses := []string{"available", "in-use"}
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := GetV3VolumeDetail(client, d.Id())
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("volume.status", respBody, "").(string)
			if status == "" {
				return respBody, "ERROR", fmt.Errorf(
					"status is not found in EVS v3 volume (%s) detail API response", d.Id())
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
		Delay:                     5 * time.Second,
		PollInterval:              5 * time.Second,
		ContinuousTargetOccurence: 2,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func resourceV3VolumeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/volumes"
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
		JSONBody:         utils.RemoveNil(buildCreateV3VolumeBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating EVS v3 volume: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	volumeID := utils.PathSearch("volume.id", respBody, "").(string)
	if volumeID == "" {
		return diag.Errorf("error creating EVS v3 volume: ID is not found in API response")
	}

	d.SetId(volumeID)

	if err := waitingForV3VolumeComplete(ctx, client, d, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for EVS v3 volume (%s) creation to complete: %s", volumeID, err)
	}

	if err := utils.CreateResourceTags(client, d, "cloudvolumes", volumeID); err != nil {
		return diag.Errorf("error setting tags for EVS v3 volume (%s): %s", volumeID, err)
	}

	return resourceV3VolumeRead(ctx, d, meta)
}

func resourceV3VolumeRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "evs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	respBody, err := GetV3VolumeDetail(client, d.Id())
	if err != nil {
		// When the resource does not exist, calling the query API will return a `404` status code.
		return common.CheckDeletedDiag(d, err, "error retrieving EVS v3 volume")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("volume_type", utils.PathSearch("volume.volume_type", respBody, nil)),
		d.Set("availability_zone", utils.PathSearch("volume.availability_zone", respBody, nil)),
		d.Set("description", utils.PathSearch("volume.description", respBody, nil)),
		d.Set("image_id", utils.PathSearch("volume.volume_image_metadata.image_id", respBody, nil)),
		d.Set("metadata", utils.ExpandToStringMap(utils.PathSearch("volume.metadata", respBody,
			make(map[string]interface{})).(map[string]interface{}))),
		d.Set("multiattach", utils.PathSearch("volume.multiattach", respBody, nil)),
		d.Set("name", utils.PathSearch("volume.name", respBody, nil)),
		d.Set("size", utils.PathSearch("volume.size", respBody, nil)),
		d.Set("snapshot_id", utils.PathSearch("volume.snapshot_id", respBody, nil)),
		d.Set("iops", utils.PathSearch("volume.iops.total_val", respBody, nil)),
		d.Set("throughput", utils.PathSearch("volume.throughput.total_val", respBody, nil)),
		utils.SetResourceTagsToState(d, client, "cloudvolumes", d.Id()),
		d.Set("tags", d.Get("tags")),
		d.Set("links", flattenV3VolumeLinksAttribute(respBody)),
		d.Set("status", utils.PathSearch("volume.status", respBody, nil)),
		d.Set("attachments", flattenV3VolumeAttachmentsAttribute(respBody)),
		d.Set("bootable", utils.PathSearch("volume.bootable", respBody, nil)),
		d.Set("created_at", utils.PathSearch("volume.created_at", respBody, nil)),
		d.Set("volume_image_metadata", utils.ExpandToStringMap(utils.PathSearch(
			"volume.volume_image_metadata", respBody, make(map[string]interface{})).(map[string]interface{}))),
		d.Set("iops_attribute", flattenV3VolumeIopsAttribute(respBody)),
		d.Set("throughput_attribute", flattenV3VolumeThroughputAttribute(respBody)),
		d.Set("updated_at", utils.PathSearch("volume.updated_at", respBody, nil)),
		d.Set("snapshot_policy_id", utils.PathSearch("volume.snapshot_policy_id", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenV3VolumeLinksAttribute(respBody interface{}) interface{} {
	linksAttribute := utils.PathSearch("volume.links", respBody, make([]interface{}, 0)).([]interface{})
	if len(linksAttribute) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(linksAttribute))
	for _, v := range linksAttribute {
		rst = append(rst, map[string]interface{}{
			"href": utils.PathSearch("href", v, nil),
			"rel":  utils.PathSearch("rel", v, nil),
		})
	}

	return rst
}

func flattenV3VolumeAttachmentsAttribute(respBody interface{}) interface{} {
	attachments := utils.PathSearch("volume.attachments", respBody, make([]interface{}, 0)).([]interface{})
	result := make([]map[string]interface{}, len(attachments))
	for i, attachment := range attachments {
		result[i] = map[string]interface{}{
			"attached_at":   utils.PathSearch("attached_at", attachment, nil),
			"attachment_id": utils.PathSearch("attachment_id", attachment, nil),
			"device":        utils.PathSearch("device", attachment, nil),
			"host_name":     utils.PathSearch("host_name", attachment, nil),
			"id":            utils.PathSearch("id", attachment, nil),
			"server_id":     utils.PathSearch("server_id", attachment, nil),
			"volume_id":     utils.PathSearch("volume_id", attachment, nil),
		}
	}

	return result
}

func flattenV3VolumeIopsAttribute(respBody interface{}) interface{} {
	iopsAttribute := utils.PathSearch("volume.iops", respBody, nil)
	if iopsAttribute == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"frozened":  utils.PathSearch("frozened", iopsAttribute, nil),
			"id":        utils.PathSearch("id", iopsAttribute, nil),
			"total_val": utils.PathSearch("total_val", iopsAttribute, nil),
			"volume_id": utils.PathSearch("volume_id", iopsAttribute, nil),
		},
	}
}

func flattenV3VolumeThroughputAttribute(respBody interface{}) interface{} {
	throughputAttribute := utils.PathSearch("volume.throughput", respBody, nil)
	if throughputAttribute == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"frozened":  utils.PathSearch("frozened", throughputAttribute, nil),
			"id":        utils.PathSearch("id", throughputAttribute, nil),
			"total_val": utils.PathSearch("total_val", throughputAttribute, nil),
			"volume_id": utils.PathSearch("volume_id", throughputAttribute, nil),
		},
	}
}

func buildUpdateV3VolumeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"description": d.Get("description"),
		"metadata":    utils.ValueIgnoreEmpty(utils.ExpandToStringMap(d.Get("metadata").(map[string]interface{}))),
	}

	return map[string]interface{}{
		"volume": bodyParams,
	}
}

func updateV3Volume(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v3/{project_id}/volumes/{volume_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{volume_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateV3VolumeBodyParams(d)),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	return err
}

func buildExtendV3VolumeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"os-extend": map[string]interface{}{
			"new_size": d.Get("size"),
		},
	}

	return bodyParams
}

func extendV3Volume(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v3/{project_id}/volumes/{volume_id}/action"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{volume_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildExtendV3VolumeBodyParams(d),
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error extending EVS v3 volume: %s", err)
	}

	if err := waitingForV3VolumeComplete(ctx, client, d, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return fmt.Errorf("error waiting for the EVS v3 volume extension to complete: %s", err)
	}

	return nil
}

func resourceV3VolumeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "evs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	if d.HasChanges("name", "description", "metadata") {
		if err := updateV3Volume(client, d); err != nil {
			return diag.Errorf("error updating EVS v3 volume (%s): %s", d.Id(), err)
		}
	}

	if d.HasChange("tags") {
		if err := utils.UpdateResourceTags(client, d, "cloudvolumes", d.Id()); err != nil {
			return diag.Errorf("error updating EVS v3 volume (%s) tags: %s", d.Id(), err)
		}
	}

	if d.HasChange("size") {
		if err := extendV3Volume(ctx, client, d); err != nil {
			return diag.Errorf("error extending EVS v3 volume (%s): %s", d.Id(), err)
		}
	}

	return resourceV3VolumeRead(ctx, d, meta)
}

func buildDeleteV3VolumeQueryParams(d *schema.ResourceData) string {
	if d.Get("cascade").(bool) {
		return "?cascade=true"
	}

	return ""
}

func deleteV3Volume(client *golangsdk.ServiceClient, getRespBody interface{}, d *schema.ResourceData) error {
	status := utils.PathSearch("status", getRespBody, "").(string)
	if status == "deleting" {
		return nil
	}

	requestPath := client.Endpoint + "v3/{project_id}/volumes/{volume_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{volume_id}", d.Id())
	requestPath += buildDeleteV3VolumeQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err := client.Request("DELETE", requestPath, &requestOpt)

	return err
}

func waitingForV3VolumeDelete(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := GetV3VolumeDetail(client, d.Id())
			if err != nil {
				var errDefault404 golangsdk.ErrDefault404
				if errors.As(err, &errDefault404) {
					return "deleted", "COMPLETED", nil
				}

				return respBody, "ERROR", err
			}

			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func detachV3Volume(ctx context.Context, getRespBody interface{}, d *schema.ResourceData, cfg *config.Config) error {
	attachments := utils.PathSearch("volume.attachments", getRespBody, make([]interface{}, 0)).([]interface{})
	computeClient, err := cfg.NewServiceClient("ecs", cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating ECS client: %s", err)
	}

	for _, attachment := range attachments {
		serverID := utils.PathSearch("server_id", attachment, "").(string)
		volumeID := utils.PathSearch("volume_id", attachment, "").(string)
		if serverID == "" || volumeID == "" {
			log.Printf("[WARN] field `server_id` (%s) or `volume_id` (%s) is empty in API response", serverID, volumeID)
			continue
		}

		requestPath := computeClient.Endpoint + "v1/{project_id}/cloudservers/{server_id}/detachvolume/{volume_id}"
		requestPath = strings.ReplaceAll(requestPath, "{project_id}", computeClient.ProjectID)
		requestPath = strings.ReplaceAll(requestPath, "{server_id}", serverID)
		requestPath = strings.ReplaceAll(requestPath, "{volume_id}", volumeID)
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		resp, err := computeClient.Request("DELETE", requestPath, &requestOpt)
		if err != nil {
			return err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return err
		}

		jobID := utils.PathSearch("job_id", respBody, "").(string)
		if jobID == "" {
			return errors.New("field `job_id` is empty in ECS detach API response")
		}

		if err := waitingForDetachJobSuccess(ctx, computeClient, jobID, d.Timeout(schema.TimeoutDelete)); err != nil {
			return fmt.Errorf("error waiting for the detach job (%s) to succeed: %s", jobID, err)
		}
	}

	return nil
}

func waitingForDetachJobSuccess(ctx context.Context, client *golangsdk.ServiceClient, jobID string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := getV3VolumeDetachJobDetail(client, jobID)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("status", respBody, "").(string)
			if status == "" {
				return respBody, "ERROR", fmt.Errorf("status is not found in ECS job (%s) detail API response", jobID)
			}

			if status == "SUCCESS" {
				return respBody, "COMPLETED", nil
			}

			if status == "FAIL" {
				return respBody, status, fmt.Errorf("the ECS job (%s) status is FAIL, the fail reason is: %s",
					jobID, utils.PathSearch("fail_reason", respBody, "").(string))
			}

			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func getV3VolumeDetachJobDetail(client *golangsdk.ServiceClient, jobID string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/jobs/{job_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", jobID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error querying ECS job detail: %s", err)
	}

	return utils.FlattenResponse(resp)
}

func resourceV3VolumeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "evs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	getRespBody, err := GetV3VolumeDetail(client, d.Id())
	if err != nil {
		// When the resource does not exist, calling the query API will return a `404` status code.
		return common.CheckDeletedDiag(d, err, "error retrieving EVS v3 volume before deleting it")
	}

	// If a volume is attached to an instance, it must be detached from the instance before deletion.
	if err := detachV3Volume(ctx, getRespBody, d, cfg); err != nil {
		return diag.Errorf("error detaching EVS v3 volume: %s", err)
	}

	if err := deleteV3Volume(client, getRespBody, d); err != nil {
		// When the resource does not exist, calling the delete API will return a `404` status code.
		return common.CheckDeletedDiag(d, err, "error deleting EVS v3 volume")
	}

	if err := waitingForV3VolumeDelete(ctx, client, d, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for the EVS v3 volume (%s) to be deleted: %s", d.Id(), err)
	}

	return nil
}
