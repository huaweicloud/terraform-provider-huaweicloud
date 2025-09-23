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

// @API RocketMQ POST /v2/{engine}/{project_id}/instances/{instance_id}/messages/resend
func ResourceDmsRocketMQConsumptionVerify() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsRocketMQConsumptionVerifyCreate,
		ReadContext:   resourceDmsRocketMQConsumptionVerifyRead,
		DeleteContext: resourceDmsRocketMQConsumptionVerifyDelete,

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
			"group": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"topic": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"message_id_list": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				ForceNew: true,
			},
			"resend_results": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the verify results.`,
				Elem:        rocketMQDeadLetterResendResultsSchema(),
			},
		},
	}
}

func resourceDmsRocketMQConsumptionVerifyCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	createHttpUrl := "v2/{engine}/{project_id}/instances/{instance_id}/messages/resend"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{engine}", "reliability")
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateRocketMQConsumptionVerifyBodyParams(d)),
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	// return 200 even failed, the error message is in response body
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating RocketMQ consumption verify: %s", err)
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

func buildCreateRocketMQConsumptionVerifyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"group":       utils.ValueIgnoreEmpty(d.Get("group")),
		"topic":       utils.ValueIgnoreEmpty(d.Get("topic")),
		"client_id":   utils.ValueIgnoreEmpty(d.Get("client_id")),
		"msg_id_list": utils.ValueIgnoreEmpty(d.Get("message_id_list")),
	}
	return bodyParams
}

func resourceDmsRocketMQConsumptionVerifyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDmsRocketMQConsumptionVerifyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
