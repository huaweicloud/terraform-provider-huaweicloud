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

const (
	deviceStatusFrozen = "FROZEN"
)

// @API IoTDA POST /v5/iot/{project_id}/devices/{device_id}/unfreeze
// @API IoTDA DELETE /v5/iot/{project_id}/devices/{device_id}
// @API IoTDA GET /v5/iot/{project_id}/devices/{device_id}
// @API IoTDA PUT /v5/iot/{project_id}/devices/{device_id}
// @API IoTDA POST /v5/iot/{project_id}/devices
// @API IoTDA POST /v5/iot/{project_id}/tags/bind-resource
// @API IoTDA POST /v5/iot/{project_id}/tags/unbind-resource
// @API IoTDA POST /v5/iot/{project_id}/devices/{device_id}/action
// @API IoTDA POST /v5/iot/{project_id}/devices/{device_id}/freeze
// @API IoTDA POST /v5/iot/{project_id}/devices/{device_id}/reset-fingerprint
// @API IoTDA PUT /v5/iot/{project_id}/devices/{device_id}/shadow
func ResourceDevice() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceDeviceCreate,
		UpdateContext: ResourceDeviceUpdate,
		DeleteContext: ResourceDeviceDelete,
		ReadContext:   ResourceDeviceRead,

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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"space_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"device_id": { // keep same with console
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"secret": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Sensitive:     true,
				ConflictsWith: []string{"fingerprint", "secondary_fingerprint"},
			},
			"secondary_secret": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Sensitive:     true,
				ConflictsWith: []string{"fingerprint", "secondary_fingerprint"},
			},
			"fingerprint": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"secret", "secondary_secret"},
			},
			"secondary_fingerprint": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"secret", "secondary_secret"},
			},
			"secure_access": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"force_disconnect": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"extension_info": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"shadow": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"desired": {
							Type:     schema.TypeMap,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"tags": common.TagsSchema(),
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"frozen": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auth_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateDeviceAuthInfo(d *schema.ResourceData) map[string]interface{} {
	if secret, ok := d.GetOk("secret"); ok {
		return map[string]interface{}{
			"auth_type":     "SECRET",
			"secret":        secret,
			"secure_access": d.Get("secure_access"),
		}
	}

	if fingerprint, ok := d.GetOk("fingerprint"); ok {
		return map[string]interface{}{
			"auth_type":     "CERTIFICATES",
			"fingerprint":   fingerprint,
			"secure_access": d.Get("secure_access"),
		}
	}

	return nil
}

func buildCreateDeviceBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"device_id":      utils.ValueIgnoreEmpty(d.Get("device_id")),
		"node_id":        d.Get("node_id"),
		"device_name":    utils.ValueIgnoreEmpty(d.Get("name")),
		"product_id":     d.Get("product_id"),
		"gateway_id":     utils.ValueIgnoreEmpty(d.Get("gateway_id")),
		"description":    utils.ValueIgnoreEmpty(d.Get("description")),
		"app_id":         utils.ValueIgnoreEmpty(d.Get("space_id")),
		"auth_info":      buildCreateDeviceAuthInfo(d),
		"extension_info": utils.ValueIgnoreEmpty(d.Get("extension_info")),
	}
}

func createDevice(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	requestPath := client.Endpoint + "v5/iot/{project_id}/devices"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateDeviceBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA device: %s", err)
	}

	return utils.FlattenResponse(resp)
}

func buildResetDeviceSecondarySecretBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"secret":           utils.ValueIgnoreEmpty(d.Get("secondary_secret")),
		"force_disconnect": d.Get("force_disconnect"),
		"secret_type":      "SECONDARY",
	}
}

func resetDeviceSecondarySecret(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v5/iot/{project_id}/devices/{device_id}/action?action_id=resetSecret"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{device_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildResetDeviceSecondarySecretBodyParams(d)),
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error updating the device secondary secret of IoTDA device: %s", err)
	}

	return nil
}

func buildResetSecondaryFingerprintBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"fingerprint":      utils.ValueIgnoreEmpty(d.Get("secondary_fingerprint")),
		"force_disconnect": d.Get("force_disconnect"),
		"fingerprint_type": "SECONDARY",
	}
}

func resetSecondaryFingerprint(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v5/iot/{project_id}/devices/{device_id}/reset-fingerprint"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{device_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildResetSecondaryFingerprintBodyParams(d)),
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error updating the secondary fingerprint of IoTDA device: %s", err)
	}

	return nil
}

func buildCreateTagsBodyParams(tagsRaw map[string]interface{}) []map[string]interface{} {
	var tags []map[string]interface{}
	for k, v := range tagsRaw {
		tags = append(tags, map[string]interface{}{
			"tag_key":   k,
			"tag_value": utils.ValueIgnoreEmpty(v),
		})
	}

	return tags
}

func buildBindingDeviceTagsBodyParams(d *schema.ResourceData, tags map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"resource_type": "device",
		"resource_id":   d.Id(),
		"tags":          buildCreateTagsBodyParams(tags),
	}
}

func bindingDeviceTags(client *golangsdk.ServiceClient, d *schema.ResourceData, tags map[string]interface{}) error {
	requestPath := client.Endpoint + "v5/iot/{project_id}/tags/bind-resource"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildBindingDeviceTagsBodyParams(d, tags)),
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error binding tags of IoTDA device: %s", err)
	}

	return nil
}

func freezeDevice(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v5/iot/{project_id}/devices/{device_id}/freeze"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{device_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{200, 201, 204},
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error freezing IoTDA device: %s", err)
	}

	return nil
}

func unfreezeDevice(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v5/iot/{project_id}/devices/{device_id}/unfreeze"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{device_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{200, 201, 204},
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error unfreezing IoTDA device: %s", err)
	}

	return nil
}

func updateFrozenStatus(client *golangsdk.ServiceClient, d *schema.ResourceData, status string) error {
	isFrozen := d.Get("frozen").(bool)
	if isFrozen && status != deviceStatusFrozen {
		return freezeDevice(client, d)
	}

	if !isFrozen && status == deviceStatusFrozen {
		return unfreezeDevice(client, d)
	}

	return nil
}

func buildDeviceShadowDesiredDataBodyParams(d *schema.ResourceData) map[string]interface{} {
	rawArray := d.Get("shadow").([]interface{})
	shadows := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		shadows = append(shadows, map[string]interface{}{
			"service_id": utils.PathSearch("service_id", v, nil),
			"desired":    utils.PathSearch("desired", v, nil),
		})
	}

	return map[string]interface{}{
		"shadow": shadows,
	}
}

func updateDeviceShadowDesiredData(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v5/iot/{project_id}/devices/{device_id}/shadow"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{device_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildDeviceShadowDesiredDataBodyParams(d),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error setting device shadow data for the device: %s", err)
	}

	return nil
}

func ResourceDeviceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "iotda"
	)

	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	// Create new device.
	respBody, err := createDevice(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	deviceID := utils.PathSearch("device_id", respBody, "").(string)
	if deviceID == "" {
		return diag.Errorf("error creating IoTDA device: ID is not found in API response")
	}
	d.SetId(deviceID)

	// Set Secondary Secret.
	if _, ok := d.GetOk("secondary_secret"); ok {
		if err := resetDeviceSecondarySecret(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	// Set Secondary Fingerprint.
	if _, ok := d.GetOk("secondary_fingerprint"); ok {
		if err := resetSecondaryFingerprint(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	// Bind tags
	if tagRaw := d.Get("tags").(map[string]interface{}); len(tagRaw) > 0 {
		if err := bindingDeviceTags(client, d, tagRaw); err != nil {
			return diag.FromErr(err)
		}
	}

	// Update device frozen status
	status := utils.PathSearch("status", respBody, "").(string)
	if err := updateFrozenStatus(client, d, status); err != nil {
		return diag.FromErr(err)
	}

	// Set device shadow data
	if shadowInfo := d.Get("shadow").([]interface{}); len(shadowInfo) > 0 {
		if err := updateDeviceShadowDesiredData(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return ResourceDeviceRead(ctx, d, meta)
}

func flattenDeviceTags(tags []interface{}) map[string]interface{} {
	if len(tags) == 0 {
		return nil
	}

	rst := make(map[string]interface{})
	for _, v := range tags {
		tagKey := utils.PathSearch("tag_key", v, "").(string)
		tagValue := utils.PathSearch("tag_value", v, nil)
		rst[tagKey] = tagValue
	}

	return rst
}

func flattenFrozenAttribute(status string) bool {
	return status == deviceStatusFrozen
}

func ResourceDeviceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/iot/{project_id}/devices/{device_id}"
		product = "iotda"
	)

	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{device_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IoTDA device")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("device_name", respBody, nil)),
		d.Set("device_id", utils.PathSearch("device_id", respBody, nil)),
		d.Set("node_id", utils.PathSearch("node_id", respBody, nil)),
		d.Set("product_id", utils.PathSearch("product_id", respBody, nil)),
		d.Set("gateway_id", utils.PathSearch("gateway_id", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("space_id", utils.PathSearch("app_id", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("node_type", utils.PathSearch("node_type", respBody, nil)),
		d.Set("tags", flattenDeviceTags(utils.PathSearch("tags", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("frozen", flattenFrozenAttribute(utils.PathSearch("status", respBody, "").(string))),
		d.Set("secret", utils.PathSearch("auth_info.secret", respBody, nil)),
		d.Set("secondary_secret", utils.PathSearch("auth_info.secondary_secret", respBody, nil)),
		d.Set("fingerprint", utils.PathSearch("auth_info.fingerprint", respBody, nil)),
		d.Set("secondary_fingerprint", utils.PathSearch("auth_info.secondary_fingerprint", respBody, nil)),
		d.Set("secure_access", utils.PathSearch("auth_info.secure_access", respBody, nil)),
		d.Set("auth_type", utils.PathSearch("auth_info.auth_type", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateDeviceAuthInfoBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"secure_access": d.Get("secure_access"),
	}
}

func buildUpdateDeviceBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"device_name": utils.ValueIgnoreEmpty(d.Get("name")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"auth_info":   buildUpdateDeviceAuthInfoBodyParams(d),
	}
}

func updateDevice(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v5/iot/{project_id}/devices/{device_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{device_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateDeviceBodyParams(d)),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error updating IoTDA device: %s", err)
	}

	return nil
}

func buildResetDeviceSecretBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"secret":           utils.ValueIgnoreEmpty(d.Get("secret")),
		"force_disconnect": d.Get("force_disconnect"),
		"secret_type":      "PRIMARY",
	}
}

func resetDeviceSecret(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v5/iot/{project_id}/devices/{device_id}/action?action_id=resetSecret"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{device_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildResetDeviceSecretBodyParams(d)),
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error updating the primary secret of IoTDA device: %s", err)
	}

	return nil
}

func buildResetFingerprintBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"fingerprint":      utils.ValueIgnoreEmpty(d.Get("fingerprint")),
		"force_disconnect": d.Get("force_disconnect"),
		"fingerprint_type": "PRIMARY",
	}
}

func resetFingerprint(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v5/iot/{project_id}/devices/{device_id}/reset-fingerprint"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{device_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildResetFingerprintBodyParams(d)),
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error updating the primary fingerprint of IoTDA device: %s", err)
	}

	return nil
}

func buildDeleteTagsBodyParams(tags map[string]interface{}) []string {
	var tagKeys []string
	for k := range tags {
		tagKeys = append(tagKeys, k)
	}

	return tagKeys
}

func buildUnbindingDeviceTagsBodyParams(d *schema.ResourceData, tags map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"resource_type": "device",
		"resource_id":   d.Id(),
		"tag_keys":      buildDeleteTagsBodyParams(tags),
	}
}

func unbindingDeviceTags(client *golangsdk.ServiceClient, d *schema.ResourceData, tags map[string]interface{}) error {
	requestPath := client.Endpoint + "v5/iot/{project_id}/tags/unbind-resource"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUnbindingDeviceTagsBodyParams(d, tags),
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error unbinding tags of IoTDA device: %s", err)
	}

	return nil
}

func updateDeviceTags(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	o, n := d.GetChange("tags")
	oMap := o.(map[string]interface{})
	nMap := n.(map[string]interface{})

	// remove old tags
	if len(oMap) > 0 {
		if err := unbindingDeviceTags(client, d, oMap); err != nil {
			return err
		}
	}

	// set new tags
	if len(nMap) > 0 {
		if err := bindingDeviceTags(client, d, nMap); err != nil {
			return err
		}
	}

	return nil
}

func ResourceDeviceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "iotda"
	)

	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	// Update name,desc,secure_access
	if d.HasChanges("name", "description", "secure_access", "extension_info") {
		if err := updateDevice(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	// Reset Device Primary Secret.
	if d.HasChange("secret") {
		if err := resetDeviceSecret(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	// Reset Device Secondary Secret.
	if d.HasChange("secondary_secret") {
		err = resetDeviceSecondarySecret(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Reset Primary Fingerprint.
	if d.HasChange("fingerprint") {
		if err := resetFingerprint(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	// Reset Secondary Fingerprint.
	if d.HasChange("secondary_fingerprint") {
		err = resetSecondaryFingerprint(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// frozen
	if d.HasChange("frozen") {
		if err := updateFrozenStatus(client, d, d.Get("status").(string)); err != nil {
			return diag.FromErr(err)
		}
	}

	// tags
	if d.HasChange("tags") {
		if err := updateDeviceTags(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	// Update device shadow data
	if d.HasChange("shadow") {
		if err := updateDeviceShadowDesiredData(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return ResourceDeviceRead(ctx, d, meta)
}

func ResourceDeviceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/iot/{project_id}/devices/{device_id}"
		product = "iotda"
	)

	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{device_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting IoTDA device")
	}

	return nil
}
