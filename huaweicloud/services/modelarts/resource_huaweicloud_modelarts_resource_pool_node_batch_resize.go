package modelarts

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var resourcePoolNodeBatchResizeNonUpdatableParams = []string{
	"resource_pool_name",
	"nodes",
	"nodes.*.batch_uid",
	"nodes.*.delete_node_names",
	"source",
	"source.*.node_pool",
	"source.*.flavor",
	"source.*.creating_step",
	"source.*.creating_step.*.type",
	"source.*.creating_step.*.step",
	"target",
	"target.*.node_pool",
	"target.*.flavor",
	"target.*.creating_step",
	"target.*.creating_step.*.type",
	"target.*.creating_step.*.step",
	"billing",
}

// @API ModelArts POST /v2/{project_id}/pools/{pool_name}/nodes/batch-resize
// @API ModelArts GET /v1/{project_id}/orders
// @API ModelArts GET /v1/{project_id}/orders/{order_name}
func ResourceResourcePoolNodeBatchResize() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceResourcePoolNodeBatchResizeCreate,
		ReadContext:   resourceResourcePoolNodeBatchResizeRead,
		UpdateContext: resourceResourcePoolNodeBatchResizeUpdate,
		DeleteContext: resourceResourcePoolNodeBatchResizeDelete,

		CustomizeDiff: config.FlexibleForceNew(resourcePoolNodeBatchResizeNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the resource pool is located.`,
			},
			"resource_pool_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resource pool name to which the resource nodes belong.`,
			},
			"nodes": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `The list of nodes to be scaled.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"batch_uid": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The batch UID of the node.`,
						},
						"delete_node_names": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `The list of nodes to be deleted.`,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"source": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem:        resourceResourcePoolNodeBatchResizeNodePoolSchema(),
				Description: `The configuration of the source node pool to which the node to be scaled belongs.`,
			},
			"target": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem:        resourceResourcePoolNodeBatchResizeNodePoolSchema(),
				Description: `The configuration of the target node pool to which the node to be scaled belongs.`,
			},
			"billing": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
				Description:  `Whether to automatically pay, in JSON format.`,
			},
			// Internal attribute(s).
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"server_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of service IDs corresponding to the currently upgraded specification nodes.`,
			},
		},
	}
}

func resourceResourcePoolNodeBatchResizeNodePoolSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_pool": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the node pool.`,
			},
			"flavor": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The flavor of the node pool.`,
			},
			"creating_step": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The type of the nodes.`,
						},
						"step": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The step number of the nodes.`,
						},
					},
				},
				Description: `The creating step of the node pool.`,
			},
		},
	}
}

func buildResourcePoolNodeBatchResizeNodes(nodes []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, len(nodes))
	for i, node := range nodes {
		result[i] = map[string]interface{}{
			"batchUID": utils.PathSearch("batch_uid", node, nil),
			"deleteNodeNames": utils.ValueIgnoreEmpty(utils.ExpandToStringList(utils.PathSearch("delete_node_names", node,
				make([]interface{}, 0)).([]interface{}))),
		}
	}
	return result
}

func buildResourcePoolNodeBatchResizeNodePool(nodePoolConfigs []interface{}) map[string]interface{} {
	nodePoolConfig := nodePoolConfigs[0]
	return map[string]interface{}{
		"nodePool": utils.PathSearch("node_pool", nodePoolConfig, nil),
		"flavor":   utils.PathSearch("flavor", nodePoolConfig, nil),
		"creatingStep": map[string]interface{}{
			"type": utils.PathSearch("creating_step[0].type", nodePoolConfig, nil),
			"step": utils.PathSearch("creating_step[0].step", nodePoolConfig, nil),
		},
	}
}

func buildResourcePoolNodeBatchResizeBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"nodes":   buildResourcePoolNodeBatchResizeNodes(d.Get("nodes").([]interface{})),
		"source":  buildResourcePoolNodeBatchResizeNodePool(d.Get("source").([]interface{})),
		"target":  buildResourcePoolNodeBatchResizeNodePool(d.Get("target").([]interface{})),
		"billing": utils.StringToJson(d.Get("billing").(string)),
	}
}

func resourceResourcePoolNodeBatchResizeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg              = meta.(*config.Config)
		httpUrl          = "v2/{project_id}/pools/{pool_name}/nodes/batch-resize"
		resourcePoolName = d.Get("resource_pool_name").(string)
	)

	client, err := cfg.NewServiceClient("modelarts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	actionPath := client.Endpoint + httpUrl
	actionPath = strings.ReplaceAll(actionPath, "{project_id}", client.ProjectID)
	actionPath = strings.ReplaceAll(actionPath, "{pool_name}", resourcePoolName)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildResourcePoolNodeBatchResizeBodyParams(d)),
	}

	resp, err := client.Request("POST", actionPath, &opt)
	if err != nil {
		return diag.Errorf("error adjusting nodes specifications under the resource pool (%s): %s", resourcePoolName, err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	repBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	orderName := utils.PathSearch("orderName", repBody, "").(string)
	if orderName == "" {
		return diag.Errorf("unable to find order name under the resource pool (%s) in API response", resourcePoolName)
	}

	err = waitingForResourcePoolOrderStatusCompleted(ctx, client, resourcePoolName, orderName, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the order status of scaling nodes under the resource pool (%s) to complete: %s", resourcePoolName, err)
	}

	actionNodeNames, err := getResourcePoolNodeNamesByOrderName(client, orderName)
	if err != nil {
		return diag.Errorf("error getting the scaling node names by order name (%s): %s", orderName, err)
	}

	actionNodes, err := waitForNodesDriverStatusCompleted(ctx, client, resourcePoolName, actionNodeNames, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the scaling nodes driver status under the resource pool (%s) to complete: %s", resourcePoolName, err)
	}

	if err := d.Set("server_ids", getServerIdsByNodeNames(actionNodeNames, actionNodes.([]interface{}))); err != nil {
		log.Printf("[ERROR] error setting the server IDs for updating resource pool (%s): %s", resourcePoolName, err)
	}

	return nil
}

func resourceResourcePoolNodeBatchResizeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceResourcePoolNodeBatchResizeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceResourcePoolNodeBatchResizeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time operation resource for batch adjust the step size of resource pool nodes. Deleting this resource
	will not clear the corresponding request records, only removing the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
