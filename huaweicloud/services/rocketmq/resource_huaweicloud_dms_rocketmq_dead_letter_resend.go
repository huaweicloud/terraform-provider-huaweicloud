package rocketmq

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

// @API RocketMQ POST /v2/{engine}/{project_id}/instances/{instance_id}/messages/deadletter-resend
func ResourceDmsRocketMQDeadLetterResend() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsRocketMQDeadLetterResendCreate,
		ReadContext:   resourceDmsRocketMQDeadLetterResendRead,
		DeleteContext: resourceDmsRocketMQDeadLetterResendDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"topic": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"message_id_list": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resend_results": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the resend results.`,
				Elem:        rocketMQDeadLetterResendResultsSchema(),
			},
		},
	}
}

func rocketMQDeadLetterResendResultsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"message_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the message ID.`,
			},
			"error_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the error code.`,
			},
			"error_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the error message.`,
			},
		},
	}
	return &sc
}

func resourceDmsRocketMQDeadLetterResendCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	createHttpUrl := "v2/{engine}/{project_id}/instances/{instance_id}/messages/deadletter-resend"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{engine}", "reliability")
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateRocketMQDeadLetterResendBodyParams(d),
	}

	// 200 even failed, have to return response body
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error resending RocketMQ dead letter messages: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("resend_results", flattenRocketMQMessageResendResults(
			utils.PathSearch("resend_results", createRespBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRocketMQMessageResendResults(paramsList []interface{}) interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		rst = append(rst, map[string]interface{}{
			"message_id":    utils.PathSearch("msg_id", params, nil),
			"error_code":    utils.PathSearch("error_code", params, nil),
			"error_message": utils.PathSearch("error_message", params, nil),
		})
	}
	return rst
}

func buildCreateRocketMQDeadLetterResendBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"topic":       d.Get("topic"),
		"msg_id_list": d.Get("message_id_list"),
	}
	return bodyParams
}

func resourceDmsRocketMQDeadLetterResendRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDmsRocketMQDeadLetterResendDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
