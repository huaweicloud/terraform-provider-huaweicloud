package dws

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var workloadQueueUpdateActionNonUpdatableParams = []string{
	"cluster_id",
	"name",
	"configuration",
	"configuration.*.resource_name",
	"configuration.*.resource_value",
	"configuration.*.value_unit",
	"configuration.*.resource_description",
	"logical_cluster_name",
}

// @API DWS PUT /v2/{project_id}/clusters/{cluster_id}/workload/queues/{queue_name}/resources
func ResourceWorkLoadQueueUpdateAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkLoadQueueUpdateActionCreate,
		ReadContext:   resourceWorkLoadQueueUpdateActionRead,
		UpdateContext: resourceWorkLoadQueueUpdateActionUpdate,
		DeleteContext: resourceWorkLoadQueueUpdateActionDelete,

		CustomizeDiff: config.FlexibleForceNew(workloadQueueUpdateActionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the workload queue is located.`,
			},

			// Required parameters.
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the cluster to which the workload queue belongs.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the workload queue to be updated.`,
			},
			"configuration": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        workloadQueueUpdateActionConfigurationSchema(),
				Description: `The list of workload queue resource items to be updated.`,
			},

			// Optional parameters.
			"logical_cluster_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the logical cluster to which the workload queue belongs.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					}),
			},
		},
	}
}

func workloadQueueUpdateActionConfigurationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resource attribute name.`,
			},
			"resource_value": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The resource attribute value.`,
			},
			"value_unit": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The unit of the resource attribute.`,
			},
			"resource_description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the resource attribute.`,
			},
		},
	}
}

func buildWorkloadQueueUpdateActionResourceItemList(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"resource_name":        utils.PathSearch("resource_name", item, "").(string),
			"resource_value":       utils.PathSearch("resource_value", item, 0).(int),
			"value_unit":           utils.ValueIgnoreEmpty(utils.PathSearch("value_unit", item, "").(string)),
			"resource_description": utils.ValueIgnoreEmpty(utils.PathSearch("resource_description", item, "").(string)),
		})
	}

	return result
}

func buildWorkloadQueueUpdateActionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"workload_queue": map[string]interface{}{
			"workload_queue_name":  d.Get("name"),
			"resource_item_list":   buildWorkloadQueueUpdateActionResourceItemList(d.Get("configuration").([]interface{})),
			"logical_cluster_name": utils.ValueIgnoreEmpty(d.Get("logical_cluster_name")),
		},
	}
}

func updateWorkloadQueueResource(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl   = "v2/{project_id}/clusters/{cluster_id}/workload/queues/{queue_name}/resources"
		clusterId = d.Get("cluster_id").(string)
		queueName = d.Get("name").(string)
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{cluster_id}", clusterId)
	updatePath = strings.ReplaceAll(updatePath, "{queue_name}", queueName)

	updateOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildWorkloadQueueUpdateActionBodyParams(d)),
	}

	resp, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	resCode := utils.PathSearch("workload_res_code", respBody, float64(0)).(float64)
	if resCode != 0 {
		resMsg := utils.PathSearch("workload_res_str", respBody, "").(string)
		return errors.New(resMsg)
	}

	return nil
}

func resourceWorkLoadQueueUpdateActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		queueName = d.Get("name").(string)
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	if err := updateWorkloadQueueResource(client, d); err != nil {
		return diag.Errorf("error updating workload queue resources (%s): %s", queueName, err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceWorkLoadQueueUpdateActionRead(ctx, d, meta)
}

func resourceWorkLoadQueueUpdateActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceWorkLoadQueueUpdateActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceWorkLoadQueueUpdateActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource performs a one-time for updating the configuration of workload queue. Deleting
this resource will not revert the configuration on the cluster, but will only remove the resource information
from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
