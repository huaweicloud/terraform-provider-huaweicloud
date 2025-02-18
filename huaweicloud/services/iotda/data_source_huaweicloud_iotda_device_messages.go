package iotda

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IoTDA GET /v5/iot/{project_id}/devices/{device_id}/messages
func DataSourceDeviceMessages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDeviceMessagesRead,

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

func dataSourceDeviceMessagesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		product   = "iotda"
		httpUrl   = "v5/iot/{project_id}/devices/{device_id}/messages"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{device_id}", d.Get("device_id").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving IoTDA device messages: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	messagesResp := utils.PathSearch("messages", getRespBody, make([]interface{}, 0)).([]interface{})
	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("messages", flattenDeviceMessages(messagesResp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDeviceMessages(messagesResp []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(messagesResp))
	for _, v := range messagesResp {
		rst = append(rst, map[string]interface{}{
			"id":             utils.PathSearch("message_id", v, nil),
			"name":           utils.PathSearch("name", v, nil),
			"message":        utils.PathSearch("message", v, nil),
			"encoding":       utils.PathSearch("encoding", v, nil),
			"payload_format": utils.PathSearch("payload_format", v, nil),
			"topic":          utils.PathSearch("topic", v, nil),
			"properties":     flattenDeviceMessageProperties(utils.PathSearch("properties", v, nil)),
			"status":         utils.PathSearch("status", v, nil),
			"error_info":     flattenDeviceMessageErrorInfo(utils.PathSearch("error_info", v, nil)),
			"created_time":   utils.PathSearch("created_time", v, nil),
			"finished_time":  utils.PathSearch("finished_time", v, nil),
		})
	}

	return rst
}
