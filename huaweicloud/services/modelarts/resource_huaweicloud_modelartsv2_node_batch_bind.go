package modelarts

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var v2NodeBatchBindNonUpdatableParams = []string{
	"pool_id",
	"nodes",
	"drain",
}

// @API ModelArts POST /v2/{project_id}/pools/{pool_name}/nodes/batch-bind
func ResourceV2NodeBatchBind() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2NodeBatchBindCreate,
		ReadContext:   resourceV2NodeBatchBindRead,
		UpdateContext: resourceV2NodeBatchBindUpdate,
		DeleteContext: resourceV2NodeBatchBindDelete,

		CustomizeDiff: config.FlexibleForceNew(v2NodeBatchBindNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the resource pool is located.`,
			},

			// Required parameters.
			"pool_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the resource pool to which the nodes belong.`,
			},
			"nodes": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `The list of nodes to be bound to the logical sub-pool.`,
				Elem:        nodeBatchBindNodeSchema(),
			},

			// Optional parameters.
			"drain": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to drain the nodes during the bind operation.`,
			},
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

func nodeBatchBindNodeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the node to be bound.`,
			},
			"quota_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the logical sub-pool to which the node is bound.`,
			},
		},
	}
}

func buildNodeBatchBindNodes(nodes []interface{}) []map[string]interface{} {
	if len(nodes) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(nodes))
	for i, node := range nodes {
		result[i] = map[string]interface{}{
			"name":      utils.PathSearch("name", node, nil),
			"quotaName": utils.ValueIgnoreEmpty(utils.PathSearch("quota_name", node, nil)),
		}
	}
	return result
}

func buildNodeBatchBindBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"nodes": buildNodeBatchBindNodes(d.Get("nodes").([]interface{})),
		"drain": utils.ValueIgnoreEmpty(d.Get("drain")),
	}
}

func resourceV2NodeBatchBindCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		poolId = d.Get("pool_id").(string)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	httpUrl := "v2/{project_id}/pools/{pool_name}/nodes/batch-bind"
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{pool_name}", poolId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildNodeBatchBindBodyParams(d)),
	}

	_, err = client.Request("POST", requestPath, &opt)
	if err != nil {
		return diag.Errorf("error binding nodes to the logical sub-pool under the resource pool (%s): %s", poolId, err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	return nil
}

func resourceV2NodeBatchBindRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV2NodeBatchBindUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV2NodeBatchBindDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for batch binding nodes to the logical sub-pool.
Deleting this resource will not clear the corresponding request record, but will only remove the resource information
from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
