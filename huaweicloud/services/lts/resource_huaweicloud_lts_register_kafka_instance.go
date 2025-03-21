package lts

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var registerKafkaNonUpdatableParams = []string{"instance_id", "kafka_name", "connect_info", "connect_info.*.user_name", "connect_info.*.pwd"}

// @API LTS POST /v2/{project_id}/lts/dms/kafka-instance
func ResourceRegisterKafkaInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRegisterKafkaInstanceCreate,
		ReadContext:   resourceRegisterKafkaInstanceRead,
		UpdateContext: resourceRegisterKafkaInstanceUpdate,
		DeleteContext: resourceRegisterKafkaInstanceDelete,

		CustomizeDiff: config.FlexibleForceNew(registerKafkaNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the Kafka instance to be registered to the LTS.`,
			},
			"kafka_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the Kafka instance to be registered to the LTS.`,
			},
			"connect_info": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The name of the SASL_SSL user of the Kafka instance.`,
						},
						"pwd": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The password of the SASL_SSL user of the Kafka instance.`,
						},
					},
				},
				Description: `The connection information of the Kafka instance to be registered to the LTS.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceRegisterKafkaInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/lts/dms/kafka-instance"
	)

	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		JSONBody:         buildCreateRegisterKafkaInstancemBodyParams(d),
	}

	requestResp, err := client.Request("POST", createPath, &createOpts)
	if err != nil {
		return diag.Errorf("unable to register Kafka instance (%s) to LTS: %s", d.Get("instance_id").(string), err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	instanceId := utils.PathSearch("instance_id", respBody, "").(string)
	if instanceId == "" {
		return diag.Errorf("unable to find the ID from the API response")
	}

	d.SetId(instanceId)

	return resourceRegisterKafkaInstanceRead(ctx, d, meta)
}

func buildCreateRegisterKafkaInstancemBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"instance_id": d.Get("instance_id"),
		"kafka_name":  d.Get("kafka_name"),
		// For non-authenticated Kafka, this parameter must be specified as an empty object, otherwise the interface will report an error.
		"connect_info": buildConnectInfo(d.Get("connect_info").([]interface{})),
	}
}

func buildConnectInfo(connectInfo []interface{}) map[string]interface{} {
	if len(connectInfo) == 0 || connectInfo[0] == nil {
		return map[string]interface{}{}
	}

	return map[string]interface{}{
		"user_name": utils.PathSearch("user_name", connectInfo[0], nil),
		"pwd":       utils.PathSearch("pwd", connectInfo[0], nil),
	}
}

func resourceRegisterKafkaInstanceRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRegisterKafkaInstanceUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRegisterKafkaInstanceDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for registering the Kafka instance to LTS. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
