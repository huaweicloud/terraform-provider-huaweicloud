package iotda

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	iotdav5 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"

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
func ResourceDevice() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceDeviceCreate,
		UpdateContext: ResourceDeviceUpdate,
		DeleteContext: ResourceDeviceDelete,
		ReadContext:   ResourceDeviceRead,

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

func ResourceDeviceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	isDerived := WithDerivedAuth(c, region)
	client, err := c.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	createOpts := buildDeviceCreateParams(d)
	resp, err := client.AddDevice(createOpts)
	if err != nil {
		return diag.Errorf("error creating IoTDA device: %s", err)
	}

	if resp.DeviceId == nil {
		return diag.Errorf("error creating IoTDA device: id is not found in API response")
	}

	d.SetId(*resp.DeviceId)

	// Set Secondary Secret.
	if v, ok := d.Get("secondary_secret").(string); ok && v != "" {
		err = resetDeviceSecondarySecret(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Set Secondary Fingerprint.
	if v, ok := d.Get("secondary_fingerprint").(string); ok && v != "" {
		err = resetDeviceSecondaryFingerprint(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// bind tags
	err = bindDeviceTags(client, d.Id(), nil, d.Get("tags").(map[string]interface{}))
	if err != nil {
		return diag.Errorf("error binding tags when creating IoTDA device: %s", err)
	}

	// freeze device
	isFrozen := d.Get("frozen").(bool)
	err = updateDeviceStatus(client, d.Id(), utils.StringValue(resp.Status), isFrozen)
	if err != nil {
		return diag.FromErr(err)
	}

	return ResourceDeviceRead(ctx, d, meta)
}

func ResourceDeviceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	isDerived := WithDerivedAuth(c, region)
	client, err := c.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	response, err := client.ShowDevice(&model.ShowDeviceRequest{DeviceId: d.Id()})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IoTDA device")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", response.DeviceName),
		d.Set("device_id", response.DeviceId),
		d.Set("node_id", response.NodeId),
		d.Set("product_id", response.ProductId),
		d.Set("gateway_id", response.GatewayId),
		d.Set("description", response.Description),
		d.Set("space_id", response.AppId),
		d.Set("status", response.Status),
		d.Set("node_type", response.NodeType),
		d.Set("tags", flattenTags(response.Tags)),
		d.Set("frozen", utils.StringValue(response.Status) == deviceStatusFrozen),
	)

	if response.AuthInfo != nil {
		mErr = multierror.Append(mErr,
			d.Set("secret", response.AuthInfo.Secret),
			d.Set("secondary_secret", response.AuthInfo.SecondarySecret),
			d.Set("fingerprint", response.AuthInfo.Fingerprint),
			d.Set("secondary_fingerprint", response.AuthInfo.SecondaryFingerprint),
			d.Set("secure_access", response.AuthInfo.SecureAccess),
			d.Set("auth_type", response.AuthInfo.AuthType),
		)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func ResourceDeviceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	isDerived := WithDerivedAuth(c, region)
	client, err := c.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	// Update name,desc,secure_access
	if d.HasChanges("name", "description", "secure_access") {
		updateOpts := buildDeviceUpdateParams(d)
		_, err = client.UpdateDevice(updateOpts)
		if err != nil {
			return diag.Errorf("error updating IoTDA device: %s", err)
		}
	}

	// Reset Device Primary Secret.
	if d.HasChange("secret") {
		_, err = client.ResetDeviceSecret(&model.ResetDeviceSecretRequest{
			DeviceId: d.Id(),
			ActionId: "resetSecret",
			Body: &model.ResetDeviceSecret{
				Secret:          utils.StringIgnoreEmpty(d.Get("secret").(string)),
				ForceDisconnect: utils.Bool(d.Get("force_disconnect").(bool)),
				SecretType:      utils.String("PRIMARY"),
			},
		})
		if err != nil {
			return diag.Errorf("error updating the primary secret of IoTDA device: %s", err)
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
		_, err = client.ResetFingerprint(&model.ResetFingerprintRequest{
			DeviceId: d.Id(),
			Body: &model.ResetFingerprint{
				Fingerprint:     utils.StringIgnoreEmpty(d.Get("fingerprint").(string)),
				ForceDisconnect: utils.Bool(d.Get("force_disconnect").(bool)),
				FingerprintType: utils.String("PRIMARY"),
			},
		})
		if err != nil {
			return diag.Errorf("error updating the primary fingerprint of IoTDA device: %s", err)
		}
	}

	// Reset Secondary Fingerprint.
	if d.HasChange("secondary_fingerprint") {
		err = resetDeviceSecondaryFingerprint(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// frozen
	if d.HasChange("frozen") {
		isFrozen := d.Get("frozen").(bool)
		status := d.Get("status").(string)
		err = updateDeviceStatus(client, d.Id(), status, isFrozen)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// tags
	if d.HasChange("tags") {
		o, n := d.GetChange("tags")
		err = bindDeviceTags(client, d.Id(), o.(map[string]interface{}), n.(map[string]interface{}))
		if err != nil {
			return diag.Errorf("error updating the tags of IoTDA device: %s", err)
		}
	}

	return ResourceDeviceRead(ctx, d, meta)
}

func resetDeviceSecondarySecret(client *iotdav5.IoTDAClient, d *schema.ResourceData) error {
	_, err := client.ResetDeviceSecret(&model.ResetDeviceSecretRequest{
		DeviceId: d.Id(),
		ActionId: "resetSecret",
		Body: &model.ResetDeviceSecret{
			Secret:          utils.StringIgnoreEmpty(d.Get("secondary_secret").(string)),
			ForceDisconnect: utils.Bool(d.Get("force_disconnect").(bool)),
			SecretType:      utils.String("SECONDARY"),
		},
	})
	if err != nil {
		return fmt.Errorf("error updating the secondary secret of IoTDA device: %s", err)
	}

	return nil
}

func resetDeviceSecondaryFingerprint(client *iotdav5.IoTDAClient, d *schema.ResourceData) error {
	_, err := client.ResetFingerprint(&model.ResetFingerprintRequest{
		DeviceId: d.Id(),
		Body: &model.ResetFingerprint{
			Fingerprint:     utils.StringIgnoreEmpty(d.Get("secondary_fingerprint").(string)),
			ForceDisconnect: utils.Bool(d.Get("force_disconnect").(bool)),
			FingerprintType: utils.String("SECONDARY"),
		},
	})
	if err != nil {
		return fmt.Errorf("error updating the secondary fingerprint of IoTDA device: %s", err)
	}

	return nil
}

func ResourceDeviceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	isDerived := WithDerivedAuth(c, region)
	client, err := c.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	deleteOpts := &model.DeleteDeviceRequest{DeviceId: d.Id()}
	_, err = client.DeleteDevice(deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting IoTDA device: %s", err)
	}

	return nil
}

func bindDeviceTags(client *iotdav5.IoTDAClient, deviceId string, oMap, nMap map[string]interface{}) error {
	// remove old tags
	if len(oMap) > 0 {
		keys := ExpandKeyOfTags(oMap)
		_, err := client.UntagDevice(&model.UntagDeviceRequest{
			Body: &model.UnbindTagsDto{
				ResourceType: "device",
				ResourceId:   deviceId,
				TagKeys:      keys,
			},
		})
		if err != nil {
			return err
		}
	}

	// set new tags
	if len(nMap) > 0 {
		taglist := ExpandResourceTags(nMap)
		_, err := client.TagDevice(&model.TagDeviceRequest{
			Body: &model.BindTagsDto{
				ResourceType: "device",
				ResourceId:   deviceId,
				Tags:         taglist,
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func ExpandResourceTags(tagmap map[string]interface{}) []model.TagV5Dto {
	var taglist []model.TagV5Dto

	for k, v := range tagmap {
		tag := model.TagV5Dto{
			TagKey:   k,
			TagValue: utils.StringIgnoreEmpty(v.(string)),
		}
		taglist = append(taglist, tag)
	}
	return taglist
}

func ExpandKeyOfTags(tagmap map[string]interface{}) []string {
	var taglist []string
	for k := range tagmap {
		taglist = append(taglist, k)
	}
	return taglist
}

func buildDeviceCreateParams(d *schema.ResourceData) *model.AddDeviceRequest {
	req := model.AddDeviceRequest{
		Body: &model.AddDevice{
			DeviceId:    utils.StringIgnoreEmpty(d.Get("device_id").(string)),
			NodeId:      d.Get("node_id").(string),
			DeviceName:  utils.StringIgnoreEmpty(d.Get("name").(string)),
			ProductId:   d.Get("product_id").(string),
			GatewayId:   utils.StringIgnoreEmpty(d.Get("gateway_id").(string)),
			Description: utils.StringIgnoreEmpty(d.Get("description").(string)),
			AppId:       utils.StringIgnoreEmpty(d.Get("space_id").(string)),
			AuthInfo:    buildAuthInfo(d.Get("secret").(string), d.Get("fingerprint").(string), d.Get("secure_access").(bool)),
		},
	}
	return &req
}

func buildDeviceUpdateParams(d *schema.ResourceData) *model.UpdateDeviceRequest {
	req := model.UpdateDeviceRequest{
		DeviceId: d.Id(),
		Body: &model.UpdateDevice{
			DeviceName:  utils.StringIgnoreEmpty(d.Get("name").(string)),
			Description: utils.StringIgnoreEmpty(d.Get("description").(string)),
			AuthInfo: &model.AuthInfoWithoutSecret{
				SecureAccess: utils.Bool(d.Get("secure_access").(bool)),
			},
		},
	}
	return &req
}

func buildAuthInfo(secret, fingerprint string, secureAccess bool) *model.AuthInfo {
	if len(secret) > 0 {
		return &model.AuthInfo{
			AuthType:     utils.String("SECRET"),
			Secret:       &secret,
			SecureAccess: &secureAccess,
		}
	}

	if len(fingerprint) > 0 {
		return &model.AuthInfo{
			AuthType:     utils.String("CERTIFICATES"),
			Fingerprint:  &fingerprint,
			SecureAccess: &secureAccess,
		}
	}

	return nil
}

func flattenTags(s *[]model.TagV5Dto) map[string]interface{} {
	if s == nil {
		return nil
	}

	rst := make(map[string]interface{})
	for _, v := range *s {
		rst[v.TagKey] = v.TagValue
	}

	return rst
}

func updateDeviceStatus(client *iotdav5.IoTDAClient, deviceId, status string, isFrozen bool) error {
	if isFrozen && status != deviceStatusFrozen {
		_, err := client.FreezeDevice(&model.FreezeDeviceRequest{
			DeviceId: deviceId,
		})
		if err != nil {
			return fmt.Errorf("error freezing IoTDA device: %s", err)
		}
	}

	if !isFrozen && status == deviceStatusFrozen {
		_, err := client.UnfreezeDevice(&model.UnfreezeDeviceRequest{
			DeviceId: deviceId,
		})
		if err != nil {
			return fmt.Errorf("error unfreezing IoTDA device: %s", err)
		}
	}
	return nil
}
