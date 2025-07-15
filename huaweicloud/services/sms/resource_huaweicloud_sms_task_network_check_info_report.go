package sms

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

var taskNetworkCheckInfoNonUpdatableParams = []string{"task_id", "network_delay", "network_jitter", "migration_speed",
	"loss_percentage", "cpu_usage", "mem_usage", "evaluation_result", "domain_connectivity", "destination_connectivity"}

// @API SMS POST /v3/{task_id}/update-network-check-info
func ResourceTaskNetworkCheckInfoReport() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaskNetworkCheckInfoReportCreate,
		ReadContext:   resourceTaskNetworkCheckInfoReportRead,
		UpdateContext: resourceTaskNetworkCheckInfoReportUpdate,
		DeleteContext: resourceTaskNetworkCheckInfoReportDelete,

		CustomizeDiff: config.FlexibleForceNew(taskNetworkCheckInfoNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"task_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network_delay": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"network_jitter": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"migration_speed": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"loss_percentage": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"cpu_usage": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"mem_usage": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"evaluation_result": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_connectivity": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"destination_connectivity": {
				Type:     schema.TypeBool,
				Optional: true,
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

func resourceTaskNetworkCheckInfoReportCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.SmsV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SMS client: %s", err)
	}

	err = updateTaskNetworkCheckInfo(d, client)
	if err != nil {
		return diag.Errorf("error updating SMS task network check info: %s", err)
	}

	d.SetId(d.Get("task_id").(string))

	return nil
}

func updateTaskNetworkCheckInfo(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	httpUrl := "v3/{task_id}/update-network-check-info"

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{task_id}", d.Get("task_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateTaskNetworkCheckInfoBodyParams(d)),
	}

	_, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return err
	}
	return nil
}

func buildUpdateTaskNetworkCheckInfoBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"network_delay":            d.Get("network_delay"),
		"network_jitter":           d.Get("network_jitter"),
		"migration_speed":          d.Get("migration_speed"),
		"loss_percentage":          d.Get("loss_percentage"),
		"cpu_usage":                d.Get("cpu_usage"),
		"mem_usage":                d.Get("mem_usage"),
		"evaluation_result":        d.Get("evaluation_result"),
		"domain_connectivity":      utils.ValueIgnoreEmpty(d.Get("domain_connectivity")),
		"destination_connectivity": utils.ValueIgnoreEmpty(d.Get("destination_connectivity")),
	}

	return bodyParams
}

func resourceTaskNetworkCheckInfoReportRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTaskNetworkCheckInfoReportUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTaskNetworkCheckInfoReportDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting the update task network check info resource is not supported. The update task network check" +
		" info resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
