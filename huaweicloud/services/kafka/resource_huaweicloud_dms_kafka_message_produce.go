package kafka

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Kafka POST /v2/{project_id}/instances/{instance_id}/messages/action
func ResourceDmsKafkaMessageProduce() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsKafkaMessageProduceCreate,
		ReadContext:   resourceDmsKafkaMessageProduceRead,
		DeleteContext: resourceDmsKafkaMessageProduceDelete,

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
			"body": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"property_list": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceDmsKafkaMessageProduceCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	createHttpUrl := "v2/{project_id}/instances/{instance_id}/messages/action?action_id={action_id}"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))
	createPath = strings.ReplaceAll(createPath, "{action_id}", "send")

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateKafkaMessageProduceBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error producing kafka topic message: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	return nil
}

func buildCreateKafkaMessageProduceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"topic":         d.Get("topic"),
		"body":          d.Get("body"),
		"property_list": buildCreateMessageBodyParamsPropertyList(d.Get("property_list").([]interface{})),
	}
	return bodyParams
}

func buildCreateMessageBodyParamsPropertyList(rawParams []interface{}) []map[string]interface{} {
	if len(rawParams) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0, len(rawParams))
	for _, val := range rawParams {
		raw := val.(map[string]interface{})
		params := map[string]interface{}{
			"name":  raw["name"],
			"value": raw["value"],
		}
		rst = append(rst, params)
	}

	return rst
}

func resourceDmsKafkaMessageProduceRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDmsKafkaMessageProduceDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting resource is not supported. The resource is only removed from the state, the message remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
