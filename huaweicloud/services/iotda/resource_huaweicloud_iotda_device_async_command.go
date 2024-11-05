package iotda

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IoTDA POST /v5/iot/{project_id}/devices/{device_id}/async-commands
// @API IoTDA GET /v5/iot/{project_id}/devices/{device_id}/async-commands/{command_id}
func ResourceDeviceAsyncCommand() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeviceAsyncCommandCreate,
		ReadContext:   resourceDeviceAsyncCommandRead,
		DeleteContext: resourceDeviceAsyncCommandDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"device_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"send_strategy": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"paras": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"expire_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"result": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sent_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delivered_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"response_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildDeviceAsyncCommandBodyParams(d *schema.ResourceData) *model.CreateAsyncCommandRequest {
	commandParas := utils.ValueIgnoreEmpty(d.Get("paras"))
	bodyParams := model.CreateAsyncCommandRequest{
		DeviceId: d.Get("device_id").(string),
		Body: &model.AsyncDeviceCommandRequest{
			ServiceId:    utils.StringIgnoreEmpty(d.Get("service_id").(string)),
			CommandName:  utils.StringIgnoreEmpty(d.Get("name").(string)),
			SendStrategy: d.Get("send_strategy").(string),
			//nolint:gosec
			ExpireTime: utils.Int32IgnoreEmpty(int32(d.Get("expire_time").(int))),
		},
	}
	if commandParas != nil {
		bodyParams.Body.Paras = &commandParas
	}

	return &bodyParams
}

func resourceDeviceAsyncCommandCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	createOpts := buildDeviceAsyncCommandBodyParams(d)
	resp, err := client.CreateAsyncCommand(createOpts)
	if err != nil {
		return diag.Errorf("error creating device asynchronous command : %s", err)
	}

	if resp == nil || resp.CommandId == nil {
		return diag.Errorf("error creating device asynchronous command: ID is not found in API response")
	}

	d.SetId(*resp.CommandId)

	return resourceDeviceAsyncCommandRead(ctx, d, meta)
}

func resourceDeviceAsyncCommandRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	isDerived := WithDerivedAuth(c, region)
	client, err := c.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	deviceId := d.Get("device_id").(string)

	response, err := client.ShowAsyncDeviceCommand(&model.ShowAsyncDeviceCommandRequest{DeviceId: deviceId, CommandId: d.Id()})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving device asynchronous command")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("device_id", response.DeviceId),
		d.Set("service_id", response.ServiceId),
		d.Set("name", response.CommandName),
		d.Set("status", response.Status),
		d.Set("result", response.Result),
		d.Set("send_strategy", response.SendStrategy),
		d.Set("expire_time", response.ExpireTime),
		d.Set("sent_time", response.SentTime),
		d.Set("delivered_time", response.DeliveredTime),
		d.Set("response_time", response.ResponseTime),
		d.Set("created_at", response.CreatedTime),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDeviceAsyncCommandDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Delete()' method because the resource is a action resource.
	return nil
}
