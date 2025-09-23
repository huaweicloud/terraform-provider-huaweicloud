package rocketmq

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nodeBatchRestartNonUpdatableParams = []string{"instance_id", "nodes"}

// @API RocketMQ POST /v2/{project_id}/{engine}/instances/{instance_id}/restart
func ResourceNodeBatchRestart() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNodeBatchRestartCreate,
		ReadContext:   resourceNodeBatchRestartRead,
		UpdateContext: resourceNodeBatchRestartUpdate,
		DeleteContext: resourceNodeBatchRestartDelete,

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(nodeBatchRestartNonUpdatableParams),
		),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The region in which to restart the nodes of a RocketMQ instance.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the RocketMQ instance to be restarted.`,
			},
			"nodes": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of node names to be restarted.`,
			},

			// Internal
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceNodeBatchRestartCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	httpUrl := "v2/{project_id}/rocketmq/instances/{instance_id}/restart"
	httpUrl = strings.ReplaceAll(httpUrl, "{project_id}", client.ProjectID)
	httpUrl = strings.ReplaceAll(httpUrl, "{instance_id}", d.Get("instance_id").(string))
	restartPath := client.Endpoint + httpUrl

	nodesRaw := d.Get("nodes").([]interface{})
	nodes := make([]string, len(nodesRaw))
	for i, node := range nodesRaw {
		nodes[i] = node.(string)
	}

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"nodes": nodes,
		},
	}

	_, err = client.Request("POST", restartPath, &opt)
	if err != nil {
		return diag.Errorf("error restarting nodes of the RocketMQ instance: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate resource ID: %s", err)
	}
	d.SetId(randUUID)

	return resourceNodeBatchRestartRead(ctx, d, meta)
}

func resourceNodeBatchRestartRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceNodeBatchRestartUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
func resourceNodeBatchRestartDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to restart nodes of a RocketMQ instance. Deleting
this resource will not clear the restart operation record, but will only remove the resource information from the
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
