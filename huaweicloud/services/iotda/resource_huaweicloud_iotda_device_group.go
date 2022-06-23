package iotda

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	v5 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

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
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 64),
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z-_0-9]*$`),
						"Only letters, digits, underscores (_) and hyphens (-) are allowed."),
				),
			},

			"space_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(0, 64),
					validation.StringMatch(regexp.MustCompile(stringRegxp), stringFormatMsg),
				),
			},

			"parent_group_id": {
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

func resourceDeviceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcIoTdaV5Client(region)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	createOpts := model.AddDeviceGroupRequest{
		Body: &model.AddDeviceGroupDto{
			Name:         utils.String(d.Get("name").(string)),
			Description:  utils.StringIgnoreEmpty(d.Get("description").(string)),
			SuperGroupId: utils.StringIgnoreEmpty(d.Get("parent_group_id").(string)),
			AppId:        utils.String(d.Get("space_id").(string)),
		},
	}
	log.Printf("[DEBUG] Create IoTDA device group params: %#v", createOpts)

	resp, err := client.AddDeviceGroup(&createOpts)
	if err != nil {
		return diag.Errorf("error creating IoTDA device group: %s", err)
	}

	if resp.GroupId == nil {
		return diag.Errorf("error creating IoTDA device group: id is not found in API response")
	}

	d.SetId(*resp.GroupId)

	// add device to group
	addIds := d.Get("device_ids").(*schema.Set)
	err = addDevicesToGroup(client, d.Id(), addIds)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceDeviceGroupRead(ctx, d, meta)
}

func resourceDeviceGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcIoTdaV5Client(region)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	detail, err := client.ShowDeviceGroup(&model.ShowDeviceGroupRequest{GroupId: d.Id()})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IoTDA device group")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", detail.Name),
		d.Set("description", detail.Description),
		d.Set("parent_group_id", detail.SuperGroupId),
		setDeviceIdsToState(d, client, d.Id()),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDeviceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcIoTdaV5Client(region)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	if d.HasChanges("name", "description") {
		updateOpts := &model.UpdateDeviceGroupRequest{
			GroupId: d.Id(),
			Body: &model.UpdateDeviceGroupDto{
				Name:        utils.String(d.Get("name").(string)),
				Description: utils.StringIgnoreEmpty(d.Get("description").(string)),
			},
		}
		_, err = client.UpdateDeviceGroup(updateOpts)
		if err != nil {
			return diag.Errorf("error updating IoTDA device group: %s", err)
		}
	}

	if d.HasChange("device_ids") {
		o, n := d.GetChange("device_ids")
		oldIds := o.(*schema.Set)
		newIds := n.(*schema.Set)

		// remove device from group
		deleteIds := oldIds.Difference(newIds)
		err = deleteDevicesFromGroup(client, d.Id(), deleteIds)
		if err != nil {
			return diag.FromErr(err)
		}

		// add device to group
		addIds := newIds.Difference(oldIds)
		err = addDevicesToGroup(client, d.Id(), addIds)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDeviceGroupRead(ctx, d, meta)
}

func resourceDeviceGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcIoTdaV5Client(region)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	// remove devices
	addIds := d.Get("device_ids").(*schema.Set)
	err = deleteDevicesFromGroup(client, d.Id(), addIds)
	if err != nil {
		return diag.FromErr(err)
	}

	deleteOpts := &model.DeleteDeviceGroupRequest{
		GroupId: d.Id(),
	}
	_, err = client.DeleteDeviceGroup(deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting IoTDA device group: %s", err)
	}

	return nil
}

func addDevicesToGroup(client *v5.IoTDAClient, groupId string, addIds *schema.Set) error {
	for _, v := range addIds.List() {
		_, err := client.CreateOrDeleteDeviceInGroup(&model.CreateOrDeleteDeviceInGroupRequest{
			GroupId:  groupId,
			ActionId: "addDevice",
			DeviceId: v.(string),
		})

		if err != nil {
			return fmt.Errorf("error adding device=%s to IoTDA device group=%s: %s", v.(string), groupId, err)
		}
	}
	return nil
}

func deleteDevicesFromGroup(client *v5.IoTDAClient, groupId string, addIds *schema.Set) error {
	for _, v := range addIds.List() {
		_, err := client.CreateOrDeleteDeviceInGroup(&model.CreateOrDeleteDeviceInGroupRequest{
			GroupId:  groupId,
			ActionId: "removeDevice",
			DeviceId: v.(string),
		})

		if err != nil {
			return fmt.Errorf("error deleting device=%s from IoTDA device group=%s: %s", v.(string), groupId, err)
		}
	}
	return nil
}

func setDeviceIdsToState(d *schema.ResourceData, client *v5.IoTDAClient, groupId string) error {
	var rst []string
	var marker *string
	for {
		resp, err := client.ShowDevicesInGroup(&model.ShowDevicesInGroupRequest{GroupId: groupId, Marker: marker})
		if err != nil {
			return fmt.Errorf("error setting the device_ids: %s", err)
		}
		if resp.Devices == nil || len(*resp.Devices) == 0 {
			break
		}
		for _, v := range *resp.Devices {
			rst = append(rst, *v.DeviceId)
		}
		marker = resp.Page.Marker
	}

	return d.Set("device_ids", rst)
}
