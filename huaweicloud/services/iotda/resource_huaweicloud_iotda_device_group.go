package iotda

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IoTDA POST /v5/iot/{project_id}/device-group/{group_id}/action
// @API IoTDA GET /v5/iot/{project_id}/device-group/{group_id}/devices
// @API IoTDA DELETE /v5/iot/{project_id}/device-group/{group_id}
// @API IoTDA GET /v5/iot/{project_id}/device-group/{group_id}
// @API IoTDA PUT /v5/iot/{project_id}/device-group/{group_id}
// @API IoTDA POST /v5/iot/{project_id}/device-group
func ResourceDeviceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeviceGroupCreate,
		ReadContext:   resourceDeviceGroupRead,
		UpdateContext: resourceDeviceGroupUpdate,
		DeleteContext: resourceDeviceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"space_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parent_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"dynamic_group_rule": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"device_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func buildCreateDeviceGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":               d.Get("name"),
		"description":        utils.ValueIgnoreEmpty(d.Get("description")),
		"super_group_id":     utils.ValueIgnoreEmpty(d.Get("parent_group_id")),
		"app_id":             d.Get("space_id"),
		"group_type":         utils.ValueIgnoreEmpty(d.Get("type")),
		"dynamic_group_rule": utils.ValueIgnoreEmpty(d.Get("dynamic_group_rule")),
	}
}

func resourceDeviceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		product   = "iotda"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	createPath := client.Endpoint + "v5/iot/{project_id}/device-group"
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateDeviceGroupBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating IoTDA device group: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	groupId := utils.PathSearch("group_id", createRespBody, "").(string)
	if groupId == "" {
		return diag.Errorf("error creating IoTDA device group: ID is not found in API response")
	}

	d.SetId(groupId)

	// Add devices to group.
	addDeviceIds := d.Get("device_ids").(*schema.Set)
	err = addOrRemoveDevicesFromGroup(client, "addDevice", groupId, addDeviceIds)
	if err != nil {
		return diag.Errorf("error adding devices to IoTDA device group in creation operation: %s", err)
	}

	return resourceDeviceGroupRead(ctx, d, meta)
}

func buildAddOrRemoveDevicesFromGroupQueryParams(action, deviceId string) string {
	return fmt.Sprintf("?action_id=%v&device_id=%v", action, deviceId)
}

func addOrRemoveDevicesFromGroup(client *golangsdk.ServiceClient, action, groupId string, addDeviceIds *schema.Set) error {
	requestPath := client.Endpoint + "v5/iot/{project_id}/device-group/{group_id}/action"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{group_id}", groupId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for _, deviceId := range addDeviceIds.List() {
		currentPath := requestPath + buildAddOrRemoveDevicesFromGroupQueryParams(action, deviceId.(string))
		_, err := client.Request("POST", currentPath, &requestOpt)
		if err != nil {
			return fmt.Errorf("error operating device (%s) in the IoTDA device group (%s): %s", deviceId, groupId, err)
		}
	}

	return nil
}

func resourceDeviceGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		product   = "iotda"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	getPath := client.Endpoint + "v5/iot/{project_id}/device-group/{group_id}"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{group_id}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		// When the resource does not exist, the query API will return `404` error code.
		return common.CheckDeletedDiag(d, err, "error retrieving IoTDA device group")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
		d.Set("parent_group_id", utils.PathSearch("super_group_id", getRespBody, nil)),
		d.Set("type", utils.PathSearch("group_type", getRespBody, nil)),
		d.Set("dynamic_group_rule", utils.PathSearch("dynamic_group_rule", getRespBody, nil)),
	)

	devicesIds, err := getDeviceIdsFromGroup(client, d.Id())
	if err != nil {
		return diag.Errorf("error setting device_ids field: %s", err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("device_ids", devicesIds),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildQueryDevicesFromGroupQueryParams(marker string) string {
	if marker != "" {
		return fmt.Sprintf("?marker=%v", marker)
	}

	return ""
}

func getDeviceIdsFromGroup(client *golangsdk.ServiceClient, groupId string) ([]string, error) {
	var (
		marker     string
		devicesIds = make([]string, 0)
	)

	getPath := client.Endpoint + "v5/iot/{project_id}/device-group/{group_id}/devices"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{group_id}", groupId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := getPath + buildQueryDevicesFromGroupQueryParams(marker)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving devices in IoTDA device group: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return nil, err
		}

		devicesResp := utils.PathSearch("devices", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(devicesResp) == 0 {
			break
		}

		for _, v := range devicesResp {
			deviceId := utils.PathSearch("device_id", v, "").(string)
			if deviceId != "" {
				devicesIds = append(devicesIds, deviceId)
			}
		}

		marker = utils.PathSearch("page.marker", getRespBody, "").(string)
		if marker == "" {
			break
		}
	}

	return devicesIds, nil
}

func buildUpdateDeviceGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
}

func resourceDeviceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		product   = "iotda"
		groupId   = d.Id()
	)

	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	if d.HasChanges("name", "description") {
		updatePath := client.Endpoint + "v5/iot/{project_id}/device-group/{group_id}"
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{group_id}", groupId)
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateDeviceGroupBodyParams(d)),
		}

		_, err := client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating IoTDA device group: %s", err)
		}
	}

	if d.HasChange("device_ids") {
		o, n := d.GetChange("device_ids")
		oldIds := o.(*schema.Set)
		newIds := n.(*schema.Set)

		// Remove devices from group.
		deleteIds := oldIds.Difference(newIds)
		err = addOrRemoveDevicesFromGroup(client, "removeDevice", groupId, deleteIds)
		if err != nil {
			return diag.Errorf("error deleting devices from IoTDA device group in update operation: %s", err)
		}

		// Add devices to group.
		addIds := newIds.Difference(oldIds)
		err = addOrRemoveDevicesFromGroup(client, "addDevice", groupId, addIds)
		if err != nil {
			return diag.Errorf("error adding devices to IoTDA device group in update operation: %s", err)
		}
	}

	return resourceDeviceGroupRead(ctx, d, meta)
}

func resourceDeviceGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		product   = "iotda"
		groupId   = d.Id()
	)

	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	// Remove devices from static device group, dynamic device group does not require this operation.
	if d.Get("type").(string) == "STATIC" {
		deleteDeviceIds := d.Get("device_ids").(*schema.Set)
		err = addOrRemoveDevicesFromGroup(client, "removeDevice", groupId, deleteDeviceIds)
		if err != nil {
			return diag.Errorf("error deleting devices from IoTDA device group in deletion operation: %s", err)
		}
	}

	deletePath := client.Endpoint + "v5/iot/{project_id}/device-group/{group_id}"
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{group_id}", groupId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		// When the resource does not exist, the delete API will return `404` error code.
		return common.CheckDeletedDiag(d, err, "error deleting IoTDA device group")
	}

	return nil
}
