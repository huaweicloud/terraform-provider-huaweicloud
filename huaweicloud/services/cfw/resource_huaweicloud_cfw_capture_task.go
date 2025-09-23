package cfw

import (
	"context"
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

var captureTaskNonUpdatableParams = []string{
	"fw_instance_id", "name", "duration", "max_packets",
	"destination", "destination.*.address", "destination.*.address_type",
	"source", "source.*.address", "source.*.address_type",
	"service", "service.*.dest_port", "service.*.protocol", "service.*.source_port",
}

// @API CFW POST /v1/{project_id}/capture-task
// @API CFW GET /v1/{project_id}/capture-task
// @API CFW POST /v1/{project_id}/capture-task/batch-delete
// @API CFW POST /v1/{project_id}/capture-task/stop
func ResourceCaptureTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCaptureTaskCreate,
		UpdateContext: resourceCaptureTaskUpdate,
		ReadContext:   resourceCaptureTaskRead,
		DeleteContext: resourceCaptureTaskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCaptureTaskImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(captureTaskNonUpdatableParams),
			func(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
				// A capture task can only be stopped once.
				oldVal, newVal := d.GetChange("stop_capture")
				if oldVal.(bool) && !newVal.(bool) {
					return fmt.Errorf("stop_capture can only change from false to true")
				}
				return nil
			},
		),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"fw_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the firewall instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The capture task name.`,
			},
			"duration": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The capture task duration.`,
			},
			"max_packets": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The maximum number of packets captured.`,
			},
			"destination": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `The destination configuration.`,
				Elem:        captureAddressDtoSchema(),
			},
			"source": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `The source configuration.`,
				Elem:        captureAddressDtoSchema(),
			},
			"service": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem:        captureServiceDtoSchema(),
				Description: `The service configuration.`,
			},
			"stop_capture": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to stop the capture.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The status of the capture task.`,
			},
			"task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the capture task.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the capture task.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time of the capture task.`,
			},
		},
	}
}

func captureAddressDtoSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The address.`,
			},
			"address_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The address type.`,
			},
		},
	}
	return &sc
}

func captureServiceDtoSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"dest_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The destination port.`,
			},
			"protocol": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The protocol type.`,
			},
			"source_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The source port.`,
			},
		},
	}
	return &sc
}

func resourceCaptureTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var (
		httpUrl = "v1/{project_id}/capture-task"
		product = "cfw"
	)
	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path += fmt.Sprintf("?fw_instance_id=%v", d.Get("fw_instance_id").(string))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	if _, ok := d.GetOk("stop_capture"); ok {
		return diag.Errorf("error creating capture task: stop_capture can't be set to true during creation")
	}
	opt.JSONBody = utils.RemoveNil(buildCreateCaptureTaskBodyParams(d))
	_, err = client.Request("POST", path, &opt)
	if err != nil {
		return diag.Errorf("error creating capture task: %s", err)
	}

	name := d.Get("name").(string)
	d.SetId(name)

	return resourceCaptureTaskRead(ctx, d, meta)
}

func buildCreateCaptureTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"destination": buildCreateCaptureTaskRequestBodyAddressDto(d.Get("destination")),
		"source":      buildCreateCaptureTaskRequestBodyAddressDto(d.Get("source")),
		"service":     buildCreateCaptureTaskRequestBodyServiceDto(d.Get("service")),
		"duration":    d.Get("duration"),
		"max_packets": d.Get("max_packets"),
	}
}

func buildCreateCaptureTaskRequestBodyAddressDto(rawParams interface{}) map[string]interface{} {
	rawArray := rawParams.([]interface{})
	raw := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"address":      raw["address"],
		"address_type": raw["address_type"],
		"type":         0,
	}
}

func buildCreateCaptureTaskRequestBodyServiceDto(rawParams interface{}) map[string]interface{} {
	rawArray := rawParams.([]interface{})
	raw := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"dest_port":   utils.ValueIgnoreEmpty(raw["dest_port"]),
		"protocol":    raw["protocol"],
		"source_port": utils.ValueIgnoreEmpty(raw["source_port"]),
	}
}

func resourceCaptureTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var mErr *multierror.Error

	client, err := conf.NewServiceClient("cfw", region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	name := d.Id()
	instanceID := d.Get("fw_instance_id").(string)
	capture, err := GetCaptureTask(client, name, instanceID)
	if err != nil {
		return common.CheckDeletedDiag(d, parseError(err), "error retrieving capture task")
	}

	createdDateStr := utils.PathSearch("created_date", capture, "").(string)
	createdTimestamp := utils.ConvertTimeStrToNanoTimestamp(createdDateStr, "2006/01/02 15:04:05") / 1000
	modifiedDateStr := utils.PathSearch("modified_date", capture, "").(string)
	modifiedTimestamp := utils.ConvertTimeStrToNanoTimestamp(modifiedDateStr, "2006/01/02 15:04:05") / 1000

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", capture, nil)),
		d.Set("destination", flattenGetDestinationAddressDto(capture)),
		d.Set("source", flattenGetSourceAddressDto(capture)),
		d.Set("service", flattenGetServiceDto(capture)),
		d.Set("duration", utils.PathSearch("duration", capture, nil)),
		d.Set("max_packets", utils.PathSearch("max_packets", capture, nil)),
		d.Set("status", utils.PathSearch("status", capture, nil)),
		d.Set("task_id", utils.PathSearch("task_id", capture, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(createdTimestamp, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(modifiedTimestamp, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetCaptureTask(client *golangsdk.ServiceClient, name, instanceID string) (interface{}, error) {
	httpUrl := "v1/{project_id}/capture-task"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	offset := 0
	for {
		path := fmt.Sprintf("%s?fw_instance_id=%s&limit=100&offset=%d", path, instanceID, offset)
		resp, err := client.Request("GET", path, &opt)

		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		findCaptureStr := fmt.Sprintf("data.records[?name=='%s']|[0]", name)
		capture := utils.PathSearch(findCaptureStr, respBody, nil)
		if capture != nil {
			return capture, nil
		}

		offset += 100
		total := utils.PathSearch("data.total", respBody, float64(0))
		if int(total.(float64)) <= offset {
			return nil, golangsdk.ErrDefault404{}
		}
	}
}

func flattenGetServiceDto(resp interface{}) []interface{} {
	var rst []interface{}
	if resp == nil {
		return rst
	}

	return []interface{}{
		map[string]interface{}{
			"dest_port":   utils.PathSearch("dest_port", resp, nil),
			"protocol":    utils.PathSearch("protocol", resp, nil),
			"source_port": utils.PathSearch("source_port", resp, nil),
		},
	}
}

func flattenGetSourceAddressDto(resp interface{}) []interface{} {
	var rst []interface{}
	if resp == nil {
		return rst
	}

	return []interface{}{
		map[string]interface{}{
			"address":      utils.PathSearch("source_address", resp, nil),
			"address_type": utils.PathSearch("source_address_type", resp, nil),
		},
	}
}

func flattenGetDestinationAddressDto(resp interface{}) []interface{} {
	var rst []interface{}
	if resp == nil {
		return rst
	}

	return []interface{}{
		map[string]interface{}{
			"address":      utils.PathSearch("dest_address", resp, nil),
			"address_type": utils.PathSearch("dest_address_type", resp, nil),
		},
	}
}

func resourceCaptureTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	updateCaptureTaskhasChanges := []string{
		"stop_capture",
	}

	client, err := conf.NewServiceClient("cfw", region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	if d.HasChanges(updateCaptureTaskhasChanges...) {
		err = stopCaptureTask(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
		instanceID := d.Get("fw_instance_id").(string)
		name := d.Get("name").(string)
		taskID := d.Get("task_id").(string)
		if err := waitCaptureTaskStopped(ctx, client, name, instanceID, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.Errorf("error waiting for CFW capture task %s to stop: %s", taskID, err)
		}
	}

	return resourceCaptureTaskRead(ctx, d, meta)
}

func stopCaptureTask(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v1/{project_id}/capture-task/stop"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path += fmt.Sprintf("?fw_instance_id=%v", d.Get("fw_instance_id").(string))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	opt.JSONBody = map[string]interface{}{
		"task_id": d.Get("task_id").(string),
	}
	_, err := client.Request("POST", path, &opt)
	if err != nil {
		return fmt.Errorf("error stopping CFW capture task: %s", err)
	}

	return nil
}

func waitCaptureTaskStopped(ctx context.Context, client *golangsdk.ServiceClient, name, instanceID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETE"},
		Refresh: func() (interface{}, string, error) {
			capture, err := GetCaptureTask(client, name, instanceID)
			if err != nil {
				return nil, "ERROR", err
			}

			status := int(utils.PathSearch("status", capture, float64(0)).(float64))
			// Status is terminating or in progress.
			if status == 5 || status == 2 {
				return capture, "PENDING", nil
			}

			return capture, "COMPLETE", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceCaptureTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.NewServiceClient("cfw", region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	status := d.Get("status").(int)
	name := d.Get("name").(string)
	instanceID := d.Get("fw_instance_id").(string)
	// Status is not abnormal, completed, or terminated.
	if status != 0 && status != 1 && status != 4 {
		// Status is in progress.
		if status == 2 {
			err = stopCaptureTask(client, d)
			if err != nil {
				log.Printf("[DEBUG] failed to stop capture task: %#v", err)
			}
		}

		err = waitCaptureTaskStopped(ctx, client, name, instanceID, d.Timeout(schema.TimeoutDelete))
		if err != nil {
			return common.CheckDeletedDiag(d, parseError(err), "error waiting for CFW capture task to stop")
		}
	}

	httpUrl := "v1/{project_id}/capture-task/batch-delete"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path += fmt.Sprintf("?fw_instance_id=%v", d.Get("fw_instance_id").(string))

	taskID := d.Get("task_id").(string)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"task_ids": []string{taskID},
		},
	}
	_, err = client.Request("POST", path, &opt)
	if err != nil {
		return common.CheckDeletedDiag(d, parseError(err), "error deleting capture task")
	}

	// Deletion API call does not guranatee the deletion of the resource.
	// Call the detail API to check if the resource is actually deleted.
	_, err = GetCaptureTask(client, name, instanceID)

	return common.CheckDeletedDiag(d, parseError(err), "error deleting capture task")
}

func resourceCaptureTaskImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <fw_instance_id>/<name>")
	}

	d.Set("fw_instance_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
