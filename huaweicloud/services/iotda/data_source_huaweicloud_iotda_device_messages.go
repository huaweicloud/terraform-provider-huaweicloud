package iotda

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IoTDA GET /v5/iot/{project_id}/devices/{device_id}/messages
func DataSourceDeviceMessages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDeviceMessageRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"device_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"messages": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"encoding": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payload_format": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"topic": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"properties": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"correlation_data": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"response_topic": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"user_properties": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"prop_key": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"prop_value": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"error_code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"error_msg": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"created_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"finished_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDeviceMessageRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
	)

	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	listOpts := model.ListDeviceMessagesRequest{
		DeviceId: d.Get("device_id").(string),
	}

	listResp, listErr := client.ListDeviceMessages(&listOpts)
	if listErr != nil {
		return diag.Errorf("error querying IoTDA device messages: %s", listErr)
	}

	uuID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(uuID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("messages", flattenDeviceMessages(listResp.Messages)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDeviceMessages(messages *[]model.DeviceMessage) []interface{} {
	if messages == nil {
		return nil
	}

	if len(*messages) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(*messages))
	for _, v := range *messages {
		rst = append(rst, map[string]interface{}{
			"id":             v.MessageId,
			"name":           v.Name,
			"message":        v.Message,
			"encoding":       v.Encoding,
			"payload_format": v.PayloadFormat,
			"topic":          v.Topic,
			"properties":     flattenDeviceMessageProperties(v.Properties),
			"status":         v.Status,
			"error_info":     flattenDeviceMessageErrorInfo(v.ErrorInfo),
			"created_time":   v.CreatedTime,
			"finished_time":  v.FinishedTime,
		})
	}

	return rst
}
