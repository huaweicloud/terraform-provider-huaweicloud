// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product ModelArts
// ---------------------------------------------------------------

package modelarts

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbc"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	billingModePostPaid = "0"
	billingModePrePaid  = "1"
)

// @API ModelArts POST /v2/{project_id}/pools
// @API ModelArts GET /v1/{project_id}/orders
// @API ModelArts GET /v1/{project_id}/orders/{order_name}
// @API ModelArts GET /v2/{project_id}/pools/{id}
// @API ModelArts PATCH /v2/{project_id}/pools/{id}
// @API ModelArts DELETE /v2/{project_id}/pools/{id}
// @API ModelArts GET /v2/{project_id}/pools/{pool_name}/nodes
func ResourceModelartsResourcePool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceModelartsResourcePoolCreate,
		ReadContext:   resourceModelartsResourcePoolRead,
		UpdateContext: resourceModelartsResourcePoolUpdate,
		DeleteContext: resourceModelartsResourcePoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceResourcePoolImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the resource pool.`,
			},
			"scope": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc(
					`List of job types supported by the resource pool.`,
					utils.SchemaDescInput{
						Required: true,
					},
				),
			},
			"metadata": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"annotations": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
							Description:  `The annotations of the resource pool, in JSON format.`,
						},
					}},
				Description: `The metadata of the resource pool.`,
			},
			"resources": {
				Type:        schema.TypeList,
				Elem:        modelartsResourcePoolResourceFlavorSchema(),
				Required:    true,
				Description: `List of resource specifications in the resource pool.`,
			},
			"network_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"vpc_id"},
				Description:  `The ModelArts network ID of the resource pool.`,
			},
			"prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The prefix of the user-defined node name of the resource pool.`,
			},
			"vpc_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				RequiredWith: []string{"subnet_id", "clusters", "user_login"},
				Description:  `The VPC ID.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The network ID of a subnet.`,
			},
			"clusters": {
				Type:        schema.TypeList,
				Elem:        modelartsResourcePoolClustersSchema(),
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The list of the CCE clusters.`,
			},
			"user_login": {
				Type:        schema.TypeList,
				Elem:        modelartsResourcePoolUserLoginSchema(),
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: `The user login info of the resource pool.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Workspace ID, which defaults to 0.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the resource pool.`,
			},
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the resource pool.`,
			},
			"resource_pool_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource ID of the resource pool.`,
			},
			"server_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of service IDs corresponding to the latest expanded nodes.`,
			},
			// Internal attributes(s).
			"resources_order_origin": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_pool": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(`The name of the node pool under resource pool.`,
								utils.SchemaDescInput{Internal: true},
							),
						},
						"flavor_id": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(
								`The flavor of the node pool.`,
								utils.SchemaDescInput{Internal: true},
							),
						},
						"creating_step": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(
								`The creating step of the node pool, in JSON format.`,
								utils.SchemaDescInput{Internal: true},
							),
						},
					},
				},
				Description: utils.SchemaDesc(
					`The script configuration value after the last change, according to which the resources are sorted.`,
					utils.SchemaDescInput{Internal: true},
				),
			},
		},
	}
}

func modelartsResourcePoolResourceFlavorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"flavor_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resource flavor ID.`,
			},
			"count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Number of resources of the corresponding flavors.`,
			},
			"node_pool": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The name of resource pool nodes.`,
			},
			"max_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The max number of resources of the corresponding flavors.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The VPC ID.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The network ID of a subnet.`,
			},
			"security_group_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `The security group IDs.`,
			},
			"azs": {
				Type:        schema.TypeSet,
				Elem:        modelartsResourcePoolResourcesAzsSchema(),
				Optional:    true,
				Computed:    true,
				Description: `AZs for resource pool nodes.`,
			},
			"taints": {
				Type:        schema.TypeSet,
				Elem:        modelartsResourcePoolResourcesTaintSchema(),
				Optional:    true,
				Computed:    true,
				Description: `The taints added to nodes.`,
			},
			"labels": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `The labels of resource pool.`,
			},
			"tags": common.TagsSchema(),
			"extend_params": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(_, o, n string, _ *schema.ResourceData) bool {
					// The current SuppressMapDiffs method just only supports object type sub-parameters, and does not
					// support list type sub-parameters.
					return utils.ContainsAllKeyValues(utils.TryMapValueAnalysis(o), utils.TryMapValueAnalysis(n))
				},
				Description: `The extend parameters of the resource pool.`,
			},
			"root_volume": {
				Type:        schema.TypeList,
				Elem:        modelartsResourcePoolResourcesRootVolumeSchema(),
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `The root volume of the resource pool nodes.`,
			},
			"data_volumes": {
				Type:        schema.TypeSet,
				Elem:        modelartsResourcePoolResourcesDataVolumeSchema(),
				Optional:    true,
				Computed:    true,
				Description: `The data volumes of the resource pool nodes.`,
			},
			"volume_group_configs": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"volume_group": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the volume group.`,
						},
						"docker_thin_pool": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: `The percentage of container volumes to data volumes on resource pool nodes.`,
						},
						"lvm_config": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lv_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The LVM write mode.`,
									},
									"path": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `The volume mount path.`,
									},
								},
							},
							Description: `The configuration of the LVM management.`,
						},
						"types": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The storage types of the volume group.`,
						},
					},
				},
				Description: `The extend configurations of the volume groups.`,
			},
			"os": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        modelartsResourcePoolResourcesOsSchema(),
				Description: `The image information for the specified OS.`,
			},
			"driver": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        modelartsResourcePoolResourcesDriverSchema(),
				Description: `The driver information.`,
			},
			"creating_step": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"step": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The creation step of the resource pool nodes.`,
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The type of the resource pool nodes.`,
						},
					},
				},
				Description: `The creation step configuration of the resource pool nodes.`,
			},
			// Deprecated parameter(s).
			"post_install": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc(
					`The script to be executed after security.`,
					utils.SchemaDescInput{
						Deprecated: true,
					},
				),
			},
		},
	}
	return &sc
}

func modelartsResourcePoolResourcesTaintSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The key of the taint.`,
			},
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The value of the taint.`,
			},
			"effect": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The effect of the taint.`,
			},
		},
	}
	return &sc
}

func modelartsResourcePoolResourcesAzsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"az": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The AZ name.`,
			},
			"count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Number of nodes.`,
			},
		},
	}
	return &sc
}

func modelartsResourcePoolResourcesRootVolumeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"volume_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the root volume.`,
			},
			"size": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The size of the root volume.`,
			},
			"extend_params": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsJSON,
				Description: utils.SchemaDesc(
					`The extend parameters of the root volume, in JSON format.`,
					utils.SchemaDescInput{
						Computed: true,
					},
				),
			},
		},
	}
	return &sc
}

func modelartsResourcePoolResourcesDataVolumeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"volume_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the data volume.`,
			},
			"size": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The size of the data volume.`,
			},
			"extend_params": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(_, o, n string, _ *schema.ResourceData) bool {
					// The current SuppressMapDiffs method just only supports object type sub-parameters, and does not
					// support list type sub-parameters.
					return utils.ContainsAllKeyValues(utils.TryMapValueAnalysis(o), utils.TryMapValueAnalysis(n))
				},
				Description: `The extend parameters of the data volume.`,
			},
			"count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The count of the current data volume configuration.`,
			},
		},
	}
	return &sc
}

func modelartsResourcePoolResourcesOsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The OS name of the image.`,
			},
			"image_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The image ID.`,
			},
			"image_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The image type.`,
			},
		},
	}
}

func modelartsResourcePoolResourcesDriverSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The driver version.`,
			},
		},
	}
}

func modelartsResourcePoolClustersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"provider_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the CCE cluster.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the CCE cluster.`,
			},
		},
	}
	return &sc
}

func modelartsResourcePoolUserLoginSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: `The password of the login user.`,
			},
			"key_pair_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The key pair name of the login user.`,
			},
		},
	}
	return &sc
}

func waitingForResourcePoolOrderStatusCompleted(ctx context.Context, client *golangsdk.ServiceClient, resourcePoolId, orderName string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshResourcePoolOrderStatus(client, resourcePoolId, orderName, []string{"successed", "partialSuccessed"}),
		Timeout:      timeout,
		Delay:        20 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

// This interface only supports querying resource pool orders, and it does not return the node name of the current operation.
func getResourcePoolOrdersByResourcePoolId(client *golangsdk.ServiceClient, resourcePoolId string) ([]interface{}, error) {
	// Currently, the maximum total number of orders is 500.
	// The query results are sorted in descending order according to the order creation time.
	var (
		httpUrl = "v1/{project_id}/orders?limit=500"
		getOpt  = golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
	)

	httpUrl = client.Endpoint + httpUrl
	httpUrl = strings.ReplaceAll(httpUrl, "{project_id}", client.ProjectID)
	httpUrl = fmt.Sprintf("%s&involvedName=%s", httpUrl, resourcePoolId)
	// Sometimes, the database may prioritize historical orders and push the latest orders out of the list, resulting
	// in inaccurate query results.
	// TODO: Using `since` parameter and `until` parameter to limit the number of results returned.

	resp, err := client.Request("GET", httpUrl, &getOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func refreshResourcePoolOrderStatus(client *golangsdk.ServiceClient, resourcePoolId, orderName string, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		orders, err := getResourcePoolOrdersByResourcePoolId(client, resourcePoolId)
		if err != nil {
			return nil, "ERROR", err
		}

		resourceOrder := utils.PathSearch(fmt.Sprintf("[?orderName == '%s']|[0]", orderName), orders, nil)
		if resourceOrder == nil {
			return nil, "ERROR", fmt.Errorf("unable to find any orders under the resource pool (%s) from API response", resourcePoolId)
		}

		status := utils.PathSearch("phase", resourceOrder, "").(string)
		if status == "failed" {
			return nil, "ERROR",
				fmt.Errorf("the order (%s) failed: %s", orderName, utils.PathSearch("failReason", resourceOrder, "").(string))
		}

		if !utils.StrSliceContains(targets, status) {
			return resourceOrder, "PENDING", nil
		}

		return resourceOrder, "COMPLETED", nil
	}
}

func getResourcePoolNodeNamesByOrderName(client *golangsdk.ServiceClient, orderName string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/orders/{order_name}"
		getOpt  = golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
	)
	httpUrl = client.Endpoint + httpUrl
	httpUrl = strings.ReplaceAll(httpUrl, "{project_id}", client.ProjectID)
	httpUrl = strings.ReplaceAll(httpUrl, "{order_name}", orderName)

	resp, err := client.Request("GET", httpUrl, &getOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	failedMessages := make([]string, 0)
	failedNodes := utils.PathSearch("items[?status == 'failed']", respBody, make([]interface{}, 0)).([]interface{})
	for _, v := range failedNodes {
		failedMessages = append(failedMessages,
			fmt.Sprintf("resource name: %s, fail message: %s",
				utils.PathSearch("resourceName", v, "").(string),
				utils.PathSearch("message", v, "").(string),
			))
	}

	if len(failedMessages) > 0 {
		log.Printf("[WARN] Action failed nodes: %s", strings.Join(failedMessages, "|"))
	}

	// For partial success nodes, only return the resource name of the successful node.
	// statuses: processing, succeed, failed
	return utils.PathSearch("items[?status == 'succeed'].resourceName", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func waitForNodesDriverStatusCompleted(ctx context.Context, client *golangsdk.ServiceClient, resourcePoolId string, nodeNames []interface{},
	timeout time.Duration) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshNodesDriverStatus(client, resourcePoolId, nodeNames),
		Timeout:      timeout,
		PollInterval: 30 * time.Second,
	}
	return stateConf.WaitForStateContext(ctx)
}

func refreshNodesDriverStatus(client *golangsdk.ServiceClient, resourcePoolId string, nodeNames []interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		nodes, err := listV2ResourcePoolNodes(client, resourcePoolId)
		if err != nil {
			return nil, "ERROR", err
		}

		// Driver statuses: Running, Abnormal, Pending, Updating
		// In the resource pool creation phase, `nodeNames` contains the name of the resource pool in addition to the node being operated.
		for _, nodeName := range nodeNames {
			nodeDriverStatus := utils.PathSearch(fmt.Sprintf("[?metadata.name == '%s'].status.driver.phase|[0]", nodeName), nodes, "").(string)
			if nodeDriverStatus != "" && !utils.StrSliceContains([]string{"Running", "Abnormal"}, nodeDriverStatus) {
				return nodes, "PENDING", nil
			}
		}

		return nodes, "COMPLETED", nil
	}
}

func getServerIdsByNodeNames(actionNodeNames []interface{}, actionNodes []interface{}) []interface{} {
	serverIds := make([]interface{}, 0)
	for _, nodeName := range actionNodeNames {
		serverId := utils.PathSearch(
			fmt.Sprintf(`[?metadata.name == '%s']|[0].metadata.labels."os.modelarts.node/underlying.instance.id"`, nodeName),
			actionNodes,
			nil,
		)
		// During the creation phase, `actionNodes` contains not only nodes information but also resource pool information.
		// So, we need to filter out the resource pool information.
		if serverId == nil {
			log.Printf("[WARN] Unable to find server ID for node (%s)", nodeName)
			continue
		}
		serverIds = append(serverIds, serverId)
	}

	return serverIds
}

func resourceModelartsResourcePoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		createResourcePoolHttpUrl = "v2/{project_id}/pools"
		createResourcePoolProduct = "modelarts"
	)
	client, err := cfg.NewServiceClient(createResourcePoolProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	createResourcePoolPath := client.Endpoint + createResourcePoolHttpUrl
	createResourcePoolPath = strings.ReplaceAll(createResourcePoolPath, "{project_id}", client.ProjectID)

	createResourcePoolOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	createResourcePoolOpt.JSONBody = utils.RemoveNil(buildCreateResourcePoolBodyParams(d))
	createResourcePoolResp, err := client.Request("POST", createResourcePoolPath, &createResourcePoolOpt)
	if err != nil {
		return diag.Errorf("error creating Modelarts resource pool: %s", err)
	}

	createResourcePoolRespBody, err := utils.FlattenResponse(createResourcePoolResp)
	if err != nil {
		return diag.FromErr(err)
	}

	resourcePoolId := utils.PathSearch("metadata.name", createResourcePoolRespBody, "").(string)
	if resourcePoolId == "" {
		return diag.Errorf("unable to find resource pool ID in API response")
	}
	d.SetId(resourcePoolId)

	if err = d.Set("resources_order_origin", refreshResourcesOrderOrigin(d.GetRawConfig())); err != nil {
		log.Printf("[ERROR] error saving resources_order_origin field for creating resource pool: %s", err)
	}

	orderName := utils.PathSearch(`metadata.labels."os.modelarts/order.name"`, createResourcePoolRespBody, "").(string)
	if orderName == "" {
		return diag.Errorf("unable to find order name of resource pool (%s) in API response", resourcePoolId)
	}

	err = waitingForResourcePoolOrderStatusCompleted(ctx, client, resourcePoolId, orderName, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the order status of resource pool (%s) creation to complete: %s", resourcePoolId, err)
	}

	actionNodeNames, err := getResourcePoolNodeNamesByOrderName(client, orderName)
	if err != nil {
		return diag.Errorf("error getting the node names by order name (%s) for creating resource pool (%s): %s", orderName, resourcePoolId, err)
	}

	actionNodes, err := waitForNodesDriverStatusCompleted(ctx, client, resourcePoolId, actionNodeNames, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the nodes driver status under the resource pool (%s) to complete: %s", resourcePoolId, err)
	}

	if err := d.Set("server_ids", getServerIdsByNodeNames(actionNodeNames, actionNodes.([]interface{}))); err != nil {
		log.Printf("[ERROR] error setting the server IDs for creating resource pool (%s): %s", resourcePoolId, err)
	}

	return resourceModelartsResourcePoolRead(ctx, d, meta)
}

func buildCreateResourcePoolBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"apiVersion": "v2",
		"kind":       "Pool",
		"metadata":   buildCreateResourcePoolMetaDataBodyParams(d),
		"spec":       buildCreateResourcePoolSpecBodyParams(d),
	}
	return bodyParams
}

func buildCreateResourcePoolMetaDataBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"labels":      buildCreateResourcePoolMetaDataLabelsBodyParams(d),
		"annotations": buildCreateResourcePoolMetaDataAnnotationsBodyParams(d),
	}
	return params
}

func buildCreateResourcePoolMetaDataLabelsBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"os.modelarts/name":         d.Get("name"),
		"os.modelarts/node.prefix":  utils.ValueIgnoreEmpty(d.Get("prefix")),
		"os.modelarts/workspace.id": utils.ValueIgnoreEmpty(d.Get("workspace_id")),
	}
	return params
}

func buildCreateResourcePoolMetaDataAnnotationsBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := make(map[string]interface{})
	if annotations, ok := d.GetOk("metadata.0.annotations"); ok {
		params = utils.StringToJson(annotations.(string)).(map[string]interface{})
	}

	if description, ok := d.GetOk("description"); ok {
		params["os.modelarts/description"] = description
	}

	if d.Get("charging_mode") == "prePaid" {
		params["os.modelarts/billing.mode"] = "1"
		if d.Get("period_unit") == "month" {
			params["os.modelarts/period.type"] = "2"
		} else {
			params["os.modelarts/period.type"] = "3"
		}
		params["os.modelarts/period.num"] = strconv.Itoa(d.Get("period").(int))
		if d.Get("auto_renew").(string) == "true" {
			params["os.modelarts/auto.renew"] = "1"
		}
		params["os.modelarts/auto.pay"] = "1"
	}
	return params
}

func buildCreateResourcePoolSpecBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"type":      "Dedicate",
		"scope":     utils.ValueIgnoreEmpty(d.Get("scope").(*schema.Set).List()),
		"resources": buildCreateResourcePoolSpecResources(d.Get("resources").([]interface{})),
		"userLogin": buildCreateResourcePoolSpecUserLoginBodyParams(d),
		"network":   buildCreateResourcePoolSpecNetworkBodyParams(d),
		"clusters":  buildCreateResourcePoolSpecClustersBodyParams(d),
	}
	return params
}

func buildCreateResourcePoolSpecUserLoginBodyParams(d *schema.ResourceData) map[string]interface{} {
	userLoginRaw := d.Get("user_login").([]interface{})
	if len(userLoginRaw) < 1 {
		return nil
	}
	return map[string]interface{}{
		"password":    utils.ValueIgnoreEmpty(d.Get("user_login.0.password")),
		"keyPairName": utils.ValueIgnoreEmpty(d.Get("user_login.0.key_pair_name")),
	}
}

func buildCreateResourcePoolSpecNetworkBodyParams(d *schema.ResourceData) map[string]interface{} {
	network := make(map[string]interface{})
	if networkId, ok := d.GetOk("network_id"); ok {
		network["name"] = networkId
	} else {
		network["vpcId"] = d.Get("vpc_id")
		network["subnetId"] = d.Get("subnet_id")
	}
	return network
}

func buildCreateResourcePoolSpecClustersBodyParams(d *schema.ResourceData) []interface{} {
	if clusters, ok := d.GetOk("clusters"); ok {
		providerIDs := make([]interface{}, 0, len(clusters.([]interface{})))
		for _, clusterRaw := range clusters.([]interface{}) {
			cluster := clusterRaw.(map[string]interface{})
			providerIDs = append(providerIDs, map[string]interface{}{
				"providerId": cluster["provider_id"],
			})
		}
		return providerIDs
	}
	return nil
}

func buildCreateResourcePoolSpecResources(resources []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(resources))
	for _, v := range resources {
		result = append(result, map[string]interface{}{
			"flavor":   utils.ValueIgnoreEmpty(utils.PathSearch("flavor_id", v, nil)),
			"count":    utils.ValueIgnoreEmpty(utils.PathSearch("count", v, nil)),
			"nodePool": utils.ValueIgnoreEmpty(utils.PathSearch("node_pool", v, nil)),
			"maxCount": utils.ValueIgnoreEmpty(utils.PathSearch("max_count", v, nil)),
			"azs": buildResourcePoolResourcesAzs(utils.PathSearch("azs", v,
				schema.NewSet(schema.HashString, nil)).(*schema.Set)),
			"network": buildResourcePoolSpecResourcesNetworkBodyParams(v),
			"taints": buildResourcePoolResourcesTaints(utils.PathSearch("taints", v,
				schema.NewSet(schema.HashString, nil)).(*schema.Set)),
			"tags": utils.ValueIgnoreEmpty(utils.ExpandResourceTags(utils.PathSearch("tags", v,
				make(map[string]interface{})).(map[string]interface{}))),
			"labels": utils.ValueIgnoreEmpty(utils.PathSearch("labels", v, nil)),
			"extendParams": buildCreateResourcePoolResourcesExtendParamsBodyParams(
				utils.PathSearch("extend_params", v, "{}").(string),
				utils.PathSearch("post_install", v, "").(string),
			),
			"rootVolume": buildResourcePoolResourcesRootVolume(utils.PathSearch("root_volume", v,
				make([]interface{}, 0)).([]interface{})),
			"dataVolumes": buildCreateResourcePoolResourcesDataVolumes(utils.PathSearch("data_volumes", v,
				schema.NewSet(schema.HashString, nil)).(*schema.Set)),
			"volumeGroupConfigs": buildResourcePoolResourcesVolumeGroupConfigs(
				make([]interface{}, 0),
				utils.PathSearch("volume_group_configs", v, schema.NewSet(schema.HashString, nil)).(*schema.Set).List(),
			),
			"os":     buildResourcePoolResourcesOsInfo(utils.PathSearch("os", v, make([]interface{}, 0)).([]interface{})),
			"driver": buildResourcePoolResourcesDriver(utils.PathSearch("driver", v, make([]interface{}, 0)).([]interface{})),
			"creatingStep": buildResourcePoolResourcesCreatingStep(
				utils.PathSearch("creating_step", v, make([]interface{}, 0)).([]interface{})),
		})
	}
	return result
}

func buildCreateResourcePoolResourcesExtendParamsBodyParams(extendParams, postInstall string) map[string]interface{} {
	result := make(map[string]interface{})

	if postInstall != "" {
		result["post_install"] = postInstall
	}

	if objExtendParams := utils.TryMapValueAnalysis(utils.StringToJson(extendParams)); len(objExtendParams) > 0 {
		for k, v := range objExtendParams {
			result[k] = v
		}
	}
	return result
}

func buildCreateResourcePoolResourcesDataVolumes(dataVolumes *schema.Set) []map[string]interface{} {
	if dataVolumes.Len() < 1 {
		return nil
	}

	result := make([]map[string]interface{}, dataVolumes.Len())
	for i, dataVolume := range dataVolumes.List() {
		result[i] = map[string]interface{}{
			"volumeType":   utils.PathSearch("volume_type", dataVolume, nil),
			"size":         utils.PathSearch("size", dataVolume, nil),
			"extendParams": utils.StringToJson(utils.PathSearch("extend_params", dataVolume, "").(string)),
			"count":        utils.ValueIgnoreEmpty(utils.PathSearch("count", dataVolume, nil)),
		}
	}

	return result
}

func buildResourcePoolSpecResourcesNetworkBodyParams(resourceRaw interface{}) map[string]interface{} {
	if resourceRaw != nil && utils.PathSearch("vpc_id", resourceRaw, nil) != nil {
		return utils.RemoveNil(map[string]interface{}{
			"vpc":    utils.ValueIgnoreEmpty(utils.PathSearch("vpc_id", resourceRaw, nil)),
			"subnet": utils.ValueIgnoreEmpty(utils.PathSearch("subnet_id", resourceRaw, nil)),
			"securityGroups": utils.ValueIgnoreEmpty(utils.ExpandToStringListBySet(utils.PathSearch("security_group_ids", resourceRaw,
				schema.NewSet(schema.HashString, nil)).(*schema.Set))),
		})
	}
	return nil
}

func buildResourcePoolResourcesExtendParamsBodyParams(oldExtendParams, newExtendParams, postInstall string) map[string]interface{} {
	extendParams := utils.TryMapValueAnalysis(utils.StringToJson(oldExtendParams))
	if postInstall != "" {
		extendParams["post_install"] = postInstall
	}

	if objExtendParams := utils.TryMapValueAnalysis(utils.StringToJson(newExtendParams)); len(objExtendParams) > 0 {
		for k, v := range objExtendParams {
			extendParams[k] = v
		}
	}
	return extendParams
}

func buildResourcePoolResourcesRootVolume(rootVolumes []interface{}) map[string]interface{} {
	if len(rootVolumes) < 1 {
		return nil
	}

	rootVolume := rootVolumes[0]
	return map[string]interface{}{
		"volumeType":   utils.PathSearch("volume_type", rootVolume, nil),
		"size":         utils.PathSearch("size", rootVolume, nil),
		"extendParams": utils.StringToJson(utils.PathSearch("extend_params", rootVolume, "").(string)),
	}
}

func buildResourcePoolResourcesVolumeGroupConfigs(oldVolumeGroupConfigs, newVolumeGroupConfigs []interface{}) []map[string]interface{} {
	if len(oldVolumeGroupConfigs) < 1 {
		result := make([]map[string]interface{}, 0, len(newVolumeGroupConfigs))
		for _, volumeGroupConfig := range newVolumeGroupConfigs {
			result = append(result, map[string]interface{}{
				"volumeGroup":    utils.PathSearch("volume_group", volumeGroupConfig, nil),
				"dockerThinPool": utils.ValueIgnoreEmpty(utils.PathSearch("docker_thin_pool", volumeGroupConfig, nil)),
				"lvmConfig": buildResourceVolumeGroupConfigsLvmConfig(utils.PathSearch("lvm_config", volumeGroupConfig,
					make([]interface{}, 0)).([]interface{})),
				"types": utils.ValueIgnoreEmpty(utils.PathSearch("types", volumeGroupConfig,
					make([]interface{}, 0)).([]interface{})),
			})
		}
		return result
	}

	result := make([]map[string]interface{}, 0, len(oldVolumeGroupConfigs))
	for _, volumeGroupConfig := range oldVolumeGroupConfigs {
		newVolumeGroupConfig := utils.PathSearch(fmt.Sprintf("[?volume_group=='%s']|[0]",
			utils.PathSearch("volume_group", volumeGroupConfig, "").(string)), newVolumeGroupConfigs, make(map[string]interface{}))

		elem := map[string]interface{}{
			// Required parameter.
			"volumeGroup": utils.PathSearch("volume_group", volumeGroupConfig, nil),
		}

		if dockerThinPool := utils.PathSearch("docker_thin_pool", newVolumeGroupConfig, 0).(int); dockerThinPool != 0 {
			elem["dockerThinPool"] = dockerThinPool
		} else {
			elem["dockerThinPool"] = utils.ValueIgnoreEmpty(utils.PathSearch("docker_thin_pool", volumeGroupConfig, nil))
		}

		if lvmConfigs := utils.PathSearch("lvm_config", newVolumeGroupConfig, make([]interface{}, 0)).([]interface{}); len(lvmConfigs) > 0 {
			elem["lvmConfig"] = buildResourceVolumeGroupConfigsLvmConfig(lvmConfigs)
		} else {
			elem["lvmConfig"] = utils.ValueIgnoreEmpty(buildResourceVolumeGroupConfigsLvmConfig(
				utils.PathSearch("lvm_config", volumeGroupConfig, make([]interface{}, 0)).([]interface{})))
		}

		if types := utils.PathSearch("types", newVolumeGroupConfig, make([]interface{}, 0)).([]interface{}); len(types) > 0 {
			elem["types"] = utils.ValueIgnoreEmpty(types)
		} else {
			elem["types"] = utils.ValueIgnoreEmpty(utils.PathSearch("types", volumeGroupConfig, make([]interface{}, 0)).([]interface{}))
		}

		result = append(result, elem)
	}

	return result
}

func buildResourceVolumeGroupConfigsLvmConfig(lvmConfigs []interface{}) map[string]interface{} {
	if len(lvmConfigs) < 1 {
		return nil
	}

	lvmConfig := lvmConfigs[0]
	return map[string]interface{}{
		"lvType": utils.PathSearch("lv_type", lvmConfig, nil),
		"path":   utils.ValueIgnoreEmpty(utils.PathSearch("path", lvmConfig, nil)),
	}
}

func buildResourcePoolResourcesOsInfo(osInfos []interface{}) map[string]interface{} {
	// All parameters are as the optional behavior.
	if len(osInfos) < 1 || osInfos[0] == nil {
		return nil
	}

	osInfo := osInfos[0]
	return map[string]interface{}{
		"name":      utils.ValueIgnoreEmpty(utils.PathSearch("name", osInfo, nil)),
		"imageId":   utils.ValueIgnoreEmpty(utils.PathSearch("image_id", osInfo, nil)),
		"imageType": utils.ValueIgnoreEmpty(utils.PathSearch("image_type", osInfo, nil)),
	}
}

func buildResourcePoolResourcesDriver(drivers []interface{}) map[string]interface{} {
	// All parameters are as the optional behavior.
	if len(drivers) < 1 || drivers[0] == nil {
		return nil
	}

	driver := drivers[0]
	return map[string]interface{}{
		"version": utils.ValueIgnoreEmpty(utils.PathSearch("version", driver, nil)),
	}
}

func buildResourcePoolResourcesCreatingStep(creatingSteps []interface{}) map[string]interface{} {
	if len(creatingSteps) < 1 {
		return nil
	}

	return map[string]interface{}{
		"type": utils.ValueIgnoreEmpty(utils.PathSearch("type", creatingSteps[0], nil)),
		"step": utils.ValueIgnoreEmpty(utils.PathSearch("step", creatingSteps[0], nil)),
	}
}

func resourceModelartsResourcePoolRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	getModelartsResourcePoolRespBody, err := queryResourcePool(cfg, region, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Modelarts resource pool")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch(`metadata.labels."os.modelarts/name"`,
			getModelartsResourcePoolRespBody, nil)),
		d.Set("prefix", utils.PathSearch(`metadata.labels."os.modelarts/node.prefix"`,
			getModelartsResourcePoolRespBody, nil)),
		d.Set("workspace_id", utils.PathSearch(`metadata.labels."os.modelarts/workspace.id"`,
			getModelartsResourcePoolRespBody, nil)),
		d.Set("description", utils.PathSearch(`metadata.annotations."os.modelarts/description"`,
			getModelartsResourcePoolRespBody, nil)),
		d.Set("scope", utils.PathSearch("spec.scope", getModelartsResourcePoolRespBody, nil)),
		d.Set("resources", flattenGetResourcePoolResponseBodyResources(
			utils.PathSearch("spec.resources", getModelartsResourcePoolRespBody, make([]interface{}, 0)).([]interface{}),
			d.Get("resources_order_origin").([]interface{})),
		),

		d.Set("network_id", utils.PathSearch("spec.network.name", getModelartsResourcePoolRespBody, nil)),
		d.Set("vpc_id", utils.PathSearch("spec.network.vpcId", getModelartsResourcePoolRespBody, nil)),
		d.Set("subnet_id", utils.PathSearch("spec.network.subnetId", getModelartsResourcePoolRespBody, nil)),
		d.Set("clusters", flattenGetResourcePoolClusterBodyResources(getModelartsResourcePoolRespBody)),
		d.Set("status", utils.PathSearch("status.phase", getModelartsResourcePoolRespBody, nil)),
		d.Set("resource_pool_id", utils.PathSearch(`metadata.labels."os.modelarts/resource.id"`,
			getModelartsResourcePoolRespBody, nil)),
	)
	chargingMode := utils.PathSearch(`metadata.annotations."os.modelarts/billing.mode"`,
		getModelartsResourcePoolRespBody, nil)
	if chargingMode == "1" {
		mErr = multierror.Append(mErr, d.Set("charging_mode", "prePaid"))
	}

	keyPairName := utils.PathSearch("spec.userLogin.keyPairName", getModelartsResourcePoolRespBody, nil)
	if keyPairName != nil {
		rst := make([]interface{}, 0, 1)
		rst = append(rst, map[string]interface{}{
			"key_pair_name": keyPairName,
		})
		mErr = multierror.Append(mErr, d.Set("user_login", rst))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func queryResourcePool(cfg *config.Config, region string, d *schema.ResourceData) (interface{}, error) {
	var (
		getModelartsResourcePoolHttpUrl = "v2/{project_id}/pools/{id}"
		getModelartsResourcePoolProduct = "modelarts"
	)
	getModelartsResourcePoolClient, err := cfg.NewServiceClient(getModelartsResourcePoolProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating ModelArts client: %s", err)
	}

	getModelartsResourcePoolPath := getModelartsResourcePoolClient.Endpoint + getModelartsResourcePoolHttpUrl
	getModelartsResourcePoolPath = strings.ReplaceAll(getModelartsResourcePoolPath, "{project_id}",
		getModelartsResourcePoolClient.ProjectID)
	getModelartsResourcePoolPath = strings.ReplaceAll(getModelartsResourcePoolPath, "{id}", d.Id())

	getModelartsResourcePoolOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getModelartsResourcePoolResp, err := getModelartsResourcePoolClient.Request("GET", getModelartsResourcePoolPath,
		&getModelartsResourcePoolOpt)

	if err != nil {
		return nil, err
	}

	getModelartsResourcePoolRespBody, err := utils.FlattenResponse(getModelartsResourcePoolResp)
	if err != nil {
		return nil, err
	}
	return getModelartsResourcePoolRespBody, nil
}

func flattenResourcePoolResourcesRootVolume(rootVolume interface{}) []map[string]interface{} {
	if rootVolume == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"volume_type":   utils.PathSearch("volumeType", rootVolume, nil),
			"size":          utils.PathSearch("size", rootVolume, nil),
			"extend_params": utils.JsonToString(utils.PathSearch("extendParams", rootVolume, nil)),
		},
	}
}

func flattenResourcePoolResourcesDataVolumes(dataVolumes []interface{}) []map[string]interface{} {
	if len(dataVolumes) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(dataVolumes))
	for _, dataVolume := range dataVolumes {
		result = append(result, map[string]interface{}{
			"volume_type":   utils.PathSearch("volumeType", dataVolume, nil),
			"size":          utils.PathSearch("size", dataVolume, nil),
			"extend_params": utils.JsonToString(utils.PathSearch("extendParams", dataVolume, nil)),
			"count":         utils.PathSearch("count", dataVolume, nil),
		})
	}

	return result
}

func flattenResourcePoolVolumeGroupConfigsLvmConfig(lvmConfig interface{}) []map[string]interface{} {
	if lvmConfig == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"lv_type": utils.PathSearch("lvType", lvmConfig, nil),
			"path":    utils.PathSearch("path", lvmConfig, nil),
		},
	}
}

func flattenResourcePoolResourcesVolumeGroupConfigs(volumeGroupConfigs []interface{}) []map[string]interface{} {
	if len(volumeGroupConfigs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(volumeGroupConfigs))
	for _, volumeGroupConfig := range volumeGroupConfigs {
		result = append(result, map[string]interface{}{
			"volume_group":     utils.PathSearch("volumeGroup", volumeGroupConfig, nil),
			"docker_thin_pool": utils.PathSearch("dockerThinPool", volumeGroupConfig, nil),
			"lvm_config": flattenResourcePoolVolumeGroupConfigsLvmConfig(utils.PathSearch("lvmConfig",
				volumeGroupConfig, nil)),
			"types": utils.PathSearch("types", volumeGroupConfig, make([]interface{}, 0)),
		})
	}
	return result
}

func flattenResourcePoolResourcesOsInfo(osInfo interface{}) []map[string]interface{} {
	if osInfo == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"name":       utils.PathSearch("name", osInfo, nil),
			"image_id":   utils.PathSearch("imageId", osInfo, nil),
			"image_type": utils.PathSearch("iamgeType", osInfo, nil),
		},
	}
}

func flattenResourcePoolResourcesDriver(driver interface{}) []map[string]interface{} {
	if driver == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"version": utils.PathSearch("version", driver, nil),
		},
	}
}

func flattenResourcePoolResourcesCreatingStep(creatingStep interface{}) []map[string]interface{} {
	if creatingStep == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"step": utils.PathSearch("step", creatingStep, nil),
			"type": utils.PathSearch("type", creatingStep, nil),
		},
	}
}

func orderResourcesByResourcesOrderOrigin(resources, resourcesOrderOrigin []interface{}) []interface{} {
	if len(resourcesOrderOrigin) == 0 {
		return resources
	}

	sortedResources := make([]interface{}, 0)
	// According to the `resources_order_origin` to sort the `resources.
	for _, v := range resourcesOrderOrigin {
		// Find matching resource in resources array based on flavor, node_pool and creating_step.
		_, index := findResourceByFlavorAndNodePoolAndCreatingStep(
			resources,
			utils.PathSearch("flavor", v, "").(string),
			utils.PathSearch("node_pool", v, "").(string),
			utils.PathSearch("creating_step", v, "").(string),
		)

		// -1 means no matching resource is found in resources.
		if index == -1 {
			continue
		}

		// Add the found resource to the sorted resources list.
		sortedResources = append(sortedResources, resources[index])
		// Remove the processed resource from the original array.
		resources = append(resources[:index], resources[index+1:]...)
	}

	// Add any remaining unsorted resources to the end of the sorted list.
	sortedResources = append(sortedResources, resources...)
	return sortedResources
}

func flattenGetResourcePoolResponseBodyResources(respResources, resourcesOrderOrigin []interface{}) []interface{} {
	if len(respResources) == 0 {
		return nil
	}

	sortedResources := orderResourcesByResourcesOrderOrigin(respResources, resourcesOrderOrigin)
	rst := make([]interface{}, 0, len(sortedResources))
	for _, v := range sortedResources {
		rst = append(rst, map[string]interface{}{
			"flavor_id":          utils.PathSearch("flavor", v, nil),
			"count":              utils.PathSearch("count", v, nil),
			"max_count":          utils.PathSearch("maxCount", v, nil),
			"node_pool":          utils.PathSearch("nodePool", v, nil),
			"vpc_id":             utils.PathSearch("network.vpc", v, nil),
			"subnet_id":          utils.PathSearch("network.subnet", v, nil),
			"security_group_ids": utils.PathSearch("network.securityGroups", v, nil),
			"azs":                flattenResourcePoolResourcesAzs(v),
			"taints":             flattenResourcePoolResourcesTaints(v),
			"labels":             utils.PathSearch("labels", v, nil),
			"tags":               flattenResourcePoolResourcesTags(v),
			"extend_params":      utils.JsonToString(utils.PathSearch("extendParams", v, nil)),
			"root_volume":        flattenResourcePoolResourcesRootVolume(utils.PathSearch("rootVolume", v, nil)),
			"data_volumes": flattenResourcePoolResourcesDataVolumes(utils.PathSearch("dataVolumes",
				v, make([]interface{}, 0)).([]interface{})),
			"volume_group_configs": flattenResourcePoolResourcesVolumeGroupConfigs(utils.PathSearch("volumeGroupConfigs",
				v, make([]interface{}, 0)).([]interface{})),
			"os":            flattenResourcePoolResourcesOsInfo(utils.PathSearch("os", v, nil)),
			"driver":        flattenResourcePoolResourcesDriver(utils.PathSearch("driver", v, nil)),
			"creating_step": flattenResourcePoolResourcesCreatingStep(utils.PathSearch("creatingStep", v, nil)),
			// Deprecated parameter(s).
			"post_install": utils.PathSearch("extendParams.post_install", v, nil),
		})
	}
	return rst
}

func flattenGetResourcePoolClusterBodyResources(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("spec.clusters", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"provider_id": utils.PathSearch("providerId", v, nil),
			"name":        utils.PathSearch("name", v, nil),
		})
	}
	return rst
}

func flattenResourcePoolResourcesAzs(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("azs", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"az":    utils.PathSearch("az", v, nil),
			"count": utils.PathSearch("count", v, nil),
		})
	}
	return rst
}

func flattenResourcePoolResourcesTaints(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("taints", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"key":    utils.PathSearch("key", v, nil),
			"value":  utils.PathSearch("value", v, nil),
			"effect": utils.PathSearch("effect", v, nil),
		})
	}
	return rst
}

func flattenResourcePoolResourcesTags(resp interface{}) map[string]interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("tags", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make(map[string]interface{})
	for _, v := range curArray {
		key := utils.PathSearch("key", v, "").(string)
		value := utils.PathSearch("value", v, "").(string)
		rst[key] = value
	}
	return rst
}

func refreshResourcesOrderOrigin(rawConfig cty.Value) []interface{} {
	raw := getConfigFileResources(rawConfig)
	if raw == nil {
		return make([]interface{}, 0)
	}

	resources := raw.(cty.Value)
	result := make([]interface{}, resources.LengthInt())
	for i, resourceElem := range resources.AsValueSlice() {
		result[i] = map[string]interface{}{
			"flavor_id":     getConfigFileStringValueByKey(resourceElem, "flavor_id"),
			"node_pool":     getConfigFileStringValueByKey(resourceElem, "node_pool"),
			"creating_step": getResourcesCreatingStepFromConfigFile(resourceElem),
		}
	}

	return result
}

func resourceModelartsResourcePoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		oldRaw, newRaw = d.GetChange("resources")
		resourcePoolId = d.Id()
	)

	updateResourcePoolChanges := []string{
		"description",
		"scope",
		"resources",
		"metadata",
	}

	if d.HasChanges(updateResourcePoolChanges...) {
		var (
			updateResourcePoolHttpUrl = "v2/{project_id}/pools/{id}"
			updateResourcePoolProduct = "modelarts"
		)
		client, err := cfg.NewServiceClient(updateResourcePoolProduct, region)
		if err != nil {
			return diag.Errorf("error creating ModelArts client: %s", err)
		}

		updateResourcePoolPath := client.Endpoint + updateResourcePoolHttpUrl
		updateResourcePoolPath = strings.ReplaceAll(updateResourcePoolPath, "{project_id}", client.ProjectID)
		updateResourcePoolPath = strings.ReplaceAll(updateResourcePoolPath, "{id}", resourcePoolId)

		updateResourcePoolOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
			MoreHeaders: map[string]string{"Content-Type": "application/merge-patch+json"},
		}

		updateResourcePoolOpt.JSONBody = utils.RemoveNil(buildUpdateResourcePoolBodyParams(d))
		updateModelartsResourcePoolResp, err := client.Request("PATCH", updateResourcePoolPath, &updateResourcePoolOpt)
		if err != nil {
			return diag.Errorf("error updating Modelarts resource pool: %s", err)
		}

		if err = d.Set("resources_order_origin", refreshResourcesOrderOrigin(d.GetRawConfig())); err != nil {
			log.Printf("[ERROR] error saving resources_order_origin field for updating resource pool: %s", err)
		}

		updateRespBody, err := utils.FlattenResponse(updateModelartsResourcePoolResp)
		if err != nil {
			return diag.FromErr(err)
		}

		// Except count and workspace_id, all other fields are updated synchronously.
		if !isAnyNodeScalling(oldRaw.([]interface{}), newRaw.([]interface{})) && !d.HasChange("workspace_id") {
			return resourceModelartsResourcePoolRead(ctx, d, meta)
		}

		// Only `count` and `workspace_id` are updated, we need to wait for the order to complete.
		orderName := utils.PathSearch(`metadata.labels."os.modelarts/order.name"`, updateRespBody, "").(string)
		if orderName == "" {
			return diag.Errorf("unable to find order name of resource pool (%s) in API response", resourcePoolId)
		}

		if err = waitingForResourcePoolOrderStatusCompleted(ctx, client, resourcePoolId, orderName, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.Errorf("error waiting for the order status of resource pool (%s) update to complete: %s", resourcePoolId, err)
		}

		// After node number has been increased for any node pool, make sure all node drivers were upgrade completed.
		if isAnyNodeScalingUp(oldRaw, newRaw) {
			actionNodeNames, err := getResourcePoolNodeNamesByOrderName(client, orderName)
			if err != nil {
				return diag.Errorf("error getting the node names by order name (%s) for updating resource pool (%s): %s",
					orderName,
					resourcePoolId,
					err,
				)
			}

			actionNodes, err := waitForNodesDriverStatusCompleted(ctx, client, resourcePoolId, actionNodeNames, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.Errorf("error waiting for the nodes driver status under the resource pool (%s) to complete: %s", resourcePoolId, err)
			}

			if err := d.Set("server_ids", getServerIdsByNodeNames(actionNodeNames, actionNodes.([]interface{}))); err != nil {
				log.Printf("[ERROR] error setting the server IDs for updating resource pool (%s): %s", resourcePoolId, err)
			}
		}
	}

	return resourceModelartsResourcePoolRead(ctx, d, meta)
}

func isAnyNodeScalling(oldResource, newResource []interface{}) bool {
	if len(oldResource) != len(newResource) {
		return true
	}

	for _, v := range newResource {
		matchedOldResource, _ := findResourceByFlavorAndNodePoolAndCreatingStep(oldResource,
			utils.PathSearch("flavor_id", v, "").(string),
			utils.PathSearch("node_pool", v, "").(string),
			utils.JsonToString(utils.PathSearch("creating_step[0]", v, nil)),
		)

		oldCount := utils.PathSearch("count", matchedOldResource, 0).(int)
		newCount := utils.PathSearch("count", v, 0).(int)
		if oldCount != newCount {
			return true
		}
	}
	return false
}

func buildUpdateResourcePoolBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"metadata": buildUpdateResourcePoolMetaDataBodyParams(d),
		"spec":     buildUpdateResourcePoolSpecBodyParams(d),
	}
	return bodyParams
}

func buildUpdateResourcePoolMetaDataBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"annotations": buildUpdateResourcePoolMetaDataAnnotationsBodyParams(d),
	}
	return params
}

func buildUpdateResourcePoolMetaDataAnnotationsBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := make(map[string]interface{})
	if annotations, ok := d.GetOk("metadata.0.annotations"); ok {
		params = utils.StringToJson(annotations.(string)).(map[string]interface{})
	}

	if d.Get("charging_mode") == "prePaid" {
		// Only apply the billing mode related parameters when the billing mode isn't set to post-paid in the annotations manually.
		if utils.PathSearch(`"os.modelarts/billing.mode"`, params, billingModePostPaid).(string) != billingModePostPaid {
			params["os.modelarts/order.id"] = ""
			params["os.modelarts/auto.pay"] = "1"
		}
	}

	if description, ok := d.GetOk("description"); ok {
		params["os.modelarts/description"] = description
	}

	// If the node pools are not increased, delete the billing mode related parameters.
	if !isAnyNodeScalingUp(d.GetChange("resources")) {
		delete(params, "os.modelarts/billing.mode")
		delete(params, "os.modelarts/period.num")
		delete(params, "os.modelarts/period.type")
		delete(params, "os.modelarts/auto.renew")
		delete(params, "os.modelarts/promotion.info")
		delete(params, "os.modelarts/service.console.url")
		delete(params, "os.modelarts/flavor.resource.ids")
		delete(params, "os.modelarts/order.id")
		delete(params, "os.modelarts/auto.pay")
	}

	return params
}

func isAnyNodeScalingUp(oldRaw, newRaw interface{}) bool {
	for i, v := range newRaw.([]interface{}) {
		oldCount := utils.PathSearch(fmt.Sprintf("[%d].count", i), oldRaw, 0).(int)
		newCount := utils.PathSearch("count", v, 0).(int)
		if newCount > oldCount {
			return true
		}
	}
	return false
}

func buildUpdateResourcePoolSpecBodyParams(d *schema.ResourceData) map[string]interface{} {
	oldResourcesVal, _ := d.GetChange("resources")
	rawConfig := d.GetRawConfig()
	params := map[string]interface{}{
		"scope":     utils.ValueIgnoreEmpty(d.Get("scope").(*schema.Set).List()),
		"resources": buildUpdateResourcePoolSpecResources(rawConfig, oldResourcesVal.([]interface{})),
	}
	return params
}

func getConfigFileResources(rawConfig cty.Value) interface{} {
	if rawConfig.IsNull() || !rawConfig.IsKnown() || !rawConfig.Type().IsObjectType() {
		return nil
	}

	return rawConfig.GetAttr("resources")
}

func isRawConfigListExist(elem cty.Value, key string) bool {
	if !elem.Type().HasAttribute(key) {
		return false
	}

	attr := elem.GetAttr(key)
	if attr.IsNull() || !attr.IsKnown() || !attr.Type().IsListType() || attr.LengthInt() < 1 {
		return false
	}

	return true
}

func getRawConfigSetValueByKey(elem cty.Value, key string) interface{} {
	if !elem.Type().HasAttribute(key) {
		return nil
	}

	raw := elem.GetAttr(key)
	if raw.IsNull() || !raw.IsKnown() || !raw.Type().IsSetType() || raw.LengthInt() < 1 {
		return nil
	}

	return raw
}

// Get the value of string type from the raw config file.
func getConfigFileStringValueByKey(elem cty.Value, key string) string {
	if !elem.Type().HasAttribute(key) {
		return ""
	}

	raw := elem.GetAttr(key)
	if raw.IsNull() || !raw.IsKnown() || raw.Type() != cty.String {
		return ""
	}

	return raw.AsString()
}

// Get the value of int type from the raw config file.
func getConfigFileIntValueByKey(elem cty.Value, key string) int {
	if !elem.Type().HasAttribute(key) {
		return 0
	}

	raw := elem.GetAttr(key)
	if raw.IsNull() || !raw.IsKnown() || raw.Type() != cty.Number {
		return 0
	}

	rawValue, _ := raw.AsBigFloat().Int64()
	return int(rawValue)
}

func getResourcesCreatingStepFromConfigFile(resourceElem cty.Value) string {
	var configFileCreatringStep string
	if resourceElem.Type().HasAttribute("creating_step") {
		creatingStepElem := resourceElem.GetAttr("creating_step")
		if !creatingStepElem.IsNull() && creatingStepElem.IsKnown() && creatingStepElem.Type().IsListType() && creatingStepElem.LengthInt() > 0 {
			configFileCreatringStep = utils.JsonToString(map[string]interface{}{
				"step": getConfigFileIntValueByKey(creatingStepElem.Index(cty.NumberIntVal(0)), "step"),
				"type": getConfigFileStringValueByKey(creatingStepElem.Index(cty.NumberIntVal(0)), "type"),
			})
		}
	}
	return configFileCreatringStep
}

func getMatchedResourceFromConfigfile(newResource cty.Value, oldResources []interface{}) interface{} {
	var (
		newFlavor        = getConfigFileStringValueByKey(newResource, "flavor_id")
		newNodePool      = getConfigFileStringValueByKey(newResource, "node_pool")
		newCreatringStep = getResourcesCreatingStepFromConfigFile(newResource)
	)

	resources, _ := findResourceByFlavorAndNodePoolAndCreatingStep(oldResources, newFlavor, newNodePool, newCreatringStep)
	return resources
}

func findResourceByFlavorAndNodePoolAndCreatingStep(oldResources []interface{}, flavor string, nodePool string,
	creatingStep string) (interface{}, int) {
	for index, oldResource := range oldResources {
		var (
			oldNodePool     = utils.PathSearch("node_pool", oldResource, "").(string)
			oldFlavor       = utils.PathSearch("flavor_id", oldResource, "").(string)
			oldCreatingStep = utils.JsonToString(utils.PathSearch("creating_step[0]", oldResource, nil))
		)

		if nodePool != "" &&
			oldNodePool == nodePool &&
			oldFlavor == flavor &&
			oldCreatingStep == creatingStep {
			return oldResource, index
		}

		pattern := regexp.MustCompile(fmt.Sprintf(`^%s-default|$`, oldFlavor))
		if nodePool == "" &&
			pattern.MatchString(oldNodePool) &&
			oldFlavor == flavor && oldCreatingStep == creatingStep {
			return oldResource, index
		}
	}

	return nil, -1
}

// Get the value of string type from the old resource by config file key.
func getStringValueByConfigFileKey(elem cty.Value, key string, oldResource interface{}) interface{} {
	if !elem.Type().HasAttribute(key) {
		return utils.ValueIgnoreEmpty(utils.PathSearch(key, oldResource, nil))
	}

	raw := elem.GetAttr(key)
	if raw.IsNull() || !raw.IsKnown() || raw.Type() != cty.String {
		return utils.ValueIgnoreEmpty(utils.PathSearch(key, oldResource, nil))
	}

	return utils.ValueIgnoreEmpty(raw.AsString())
}

// Get the value of int type from the old resource by config file key.
func getIntValueByConfigFileKey(elem cty.Value, key string, oldResource interface{}) interface{} {
	num := getConfigFileIntValueByKey(elem, key)
	if num != 0 {
		return num
	}

	return utils.PathSearch(key, oldResource, nil)
}

func buildSecurityGroupIdsByConfigFileKey(elem cty.Value, key string, oldResource interface{}) interface{} {
	raw := getRawConfigSetValueByKey(elem, key)
	if raw == nil {
		return utils.ValueIgnoreEmpty(utils.ExpandToStringListBySet(utils.PathSearch(key, oldResource,
			schema.NewSet(schema.HashString, nil)).(*schema.Set)))
	}

	securityGroupIdsRaw := raw.(cty.Value)
	securityGroupIds := make([]string, securityGroupIdsRaw.LengthInt())
	for i, securityGroupId := range securityGroupIdsRaw.AsValueSlice() {
		securityGroupIds[i] = securityGroupId.AsString()
	}
	return utils.ValueIgnoreEmpty(securityGroupIds)
}

func buildUpdateResourcePoolResourceLabels(elem cty.Value, key string, oldResource interface{}) interface{} {
	if !elem.Type().HasAttribute(key) {
		return utils.ValueIgnoreEmpty(utils.PathSearch(key, oldResource, nil))
	}

	raw := elem.GetAttr(key)
	if raw.IsNull() || !raw.IsKnown() || !raw.Type().IsMapType() {
		return utils.ValueIgnoreEmpty(utils.PathSearch(key, oldResource, nil))
	}

	labels := make(map[string]interface{})
	for k, v := range raw.AsValueMap() {
		labels[k] = v.AsString()
	}
	return utils.ValueIgnoreEmpty(labels)
}

func buildUpdateResourcePoolResourceAzs(resourceElem cty.Value, oldResource interface{}) interface{} {
	raw := getRawConfigSetValueByKey(resourceElem, "azs")
	if raw == nil {
		return buildResourcePoolResourcesAzs(utils.PathSearch("azs", oldResource, schema.NewSet(schema.HashString, nil)).(*schema.Set))
	}

	azs := raw.(cty.Value)
	result := make([]map[string]interface{}, azs.LengthInt())
	for i, az := range azs.AsValueSlice() {
		result[i] = map[string]interface{}{
			"az":    utils.ValueIgnoreEmpty(getConfigFileStringValueByKey(az, "az")),
			"count": utils.ValueIgnoreEmpty(getConfigFileIntValueByKey(az, "count")),
		}
	}
	return result
}

func buildUpdateResourcePoolResourceNetwork(resourceElem cty.Value, oldResource interface{}) interface{} {
	return utils.RemoveNil(map[string]interface{}{
		"vpc":            getStringValueByConfigFileKey(resourceElem, "vpc_id", oldResource),
		"subnet":         getStringValueByConfigFileKey(resourceElem, "subnet_id", oldResource),
		"securityGroups": buildSecurityGroupIdsByConfigFileKey(resourceElem, "security_group_ids", oldResource),
	})
}

func buildUpdateResourcePoolResourceTaints(resourceElem cty.Value, oldResource interface{}) interface{} {
	raw := getRawConfigSetValueByKey(resourceElem, "taints")
	if raw == nil {
		return buildResourcePoolResourcesTaints(utils.PathSearch("taints", oldResource, schema.NewSet(schema.HashString, nil)).(*schema.Set))
	}

	taints := raw.(cty.Value)
	result := make([]map[string]interface{}, taints.LengthInt())
	for i, taint := range taints.AsValueSlice() {
		result[i] = map[string]interface{}{
			"key":    taint.GetAttr("key").AsString(),
			"effect": taint.GetAttr("effect").AsString(),
			"value":  utils.ValueIgnoreEmpty(getConfigFileStringValueByKey(taint, "value")),
		}
	}
	return result
}

func buildUpdateResourcePoolResourceTags(resourceElem cty.Value) interface{} {
	if !resourceElem.Type().HasAttribute("tags") {
		return nil
	}

	tags := resourceElem.GetAttr("tags")
	if tags.IsNull() || !tags.IsKnown() || !tags.Type().IsMapType() {
		return nil
	}

	result := make([]map[string]interface{}, 0)
	for k, v := range tags.AsValueMap() {
		tagMap := map[string]interface{}{
			"key": k,
		}
		if v.Type() == cty.String && !v.IsNull() && v.IsKnown() {
			tagMap["value"] = v.AsString()
		}

		result = append(result, tagMap)
	}
	return result
}

func buildUpdateResourcePoolResourceRootVolume(resourceElem cty.Value, oldResource interface{}) interface{} {
	if !isRawConfigListExist(resourceElem, "root_volume") {
		return buildResourcePoolResourcesRootVolume(utils.PathSearch("root_volume", oldResource, make([]interface{}, 0)).([]interface{}))
	}

	rootVolume := resourceElem.GetAttr("root_volume")
	return map[string]interface{}{
		"volumeType": rootVolume.Index(cty.NumberIntVal(0)).GetAttr("volume_type").AsString(),
		"size":       rootVolume.Index(cty.NumberIntVal(0)).GetAttr("size").AsString(),
		"extendParams": buildResourcePoolResourcesExtendParamsBodyParams(
			utils.PathSearch("root_volume[0].extend_params", oldResource, "").(string),
			getConfigFileStringValueByKey(rootVolume.Index(cty.NumberIntVal(0)), "extend_params"),
			"",
		),
	}
}

func buildUpdateResourcePoolResourceDataVolumes(resourceElem cty.Value, oldResource interface{}) interface{} {
	raw := getRawConfigSetValueByKey(resourceElem, "data_volumes")
	if raw == nil {
		return buildCreateResourcePoolResourcesDataVolumes(utils.PathSearch("data_volumes", oldResource,
			schema.NewSet(schema.HashString, nil)).(*schema.Set))
	}

	dataVolumes := resourceElem.GetAttr("data_volumes")
	result := make([]map[string]interface{}, dataVolumes.LengthInt())
	for i, dataVolume := range dataVolumes.AsValueSlice() {
		result[i] = map[string]interface{}{
			"volumeType": dataVolume.GetAttr("volume_type").AsString(),
			"size":       dataVolume.GetAttr("size").AsString(),
			"extendParams": buildResourcePoolResourcesExtendParamsBodyParams(
				utils.PathSearch(fmt.Sprintf("data_volumes[%d].extend_params", i), oldResource, "").(string),
				getConfigFileStringValueByKey(dataVolume, "extend_params"),
				"",
			),
			"count": utils.ValueIgnoreEmpty(getConfigFileIntValueByKey(dataVolume, "count")),
		}
	}
	return result
}

func buildUpdateResourcePoolResourceOs(resourceElem cty.Value, oldResource interface{}) interface{} {
	if !isRawConfigListExist(resourceElem, "os") {
		return buildResourcePoolResourcesOsInfo(utils.PathSearch("os", oldResource, make([]interface{}, 0)).([]interface{}))
	}

	os := resourceElem.GetAttr("os")
	return map[string]interface{}{
		"name":      utils.ValueIgnoreEmpty(getConfigFileStringValueByKey(os.Index(cty.NumberIntVal(0)), "name")),
		"imageId":   utils.ValueIgnoreEmpty(getConfigFileStringValueByKey(os.Index(cty.NumberIntVal(0)), "image_id")),
		"imageType": utils.ValueIgnoreEmpty(getConfigFileStringValueByKey(os.Index(cty.NumberIntVal(0)), "image_type")),
	}
}

func buildUpdateResourcePoolResourceDriver(resourceElem cty.Value, oldResource interface{}) interface{} {
	if !isRawConfigListExist(resourceElem, "driver") {
		return buildResourcePoolResourcesDriver(utils.PathSearch("driver", oldResource, make([]interface{}, 0)).([]interface{}))
	}

	driver := resourceElem.GetAttr("driver")
	return utils.RemoveNil(map[string]interface{}{
		"version": utils.ValueIgnoreEmpty(getConfigFileStringValueByKey(driver.Index(cty.NumberIntVal(0)), "version")),
	})
}

func buildUpdateResourcePoolResourceCreatingStep(resourceElem cty.Value, oldResource interface{}) interface{} {
	if !isRawConfigListExist(resourceElem, "creating_step") {
		return buildResourcePoolResourcesCreatingStep(utils.PathSearch("creating_step", oldResource, make([]interface{}, 0)).([]interface{}))
	}

	creatingStep := resourceElem.GetAttr("creating_step")
	return map[string]interface{}{
		"type": creatingStep.Index(cty.NumberIntVal(0)).GetAttr("type").AsString(),
		"step": getConfigFileIntValueByKey(creatingStep.Index(cty.NumberIntVal(0)), "step"),
	}
}

func buildUpdateResourcePoolSpecResources(rawConfig cty.Value, oldResources []interface{}) []map[string]interface{} {
	resources := getConfigFileResources(rawConfig)
	if resources == nil {
		return nil
	}

	newResources := resources.(cty.Value)
	result := make([]map[string]interface{}, 0, newResources.LengthInt())
	for _, newResource := range newResources.AsValueSlice() {
		matchedOldResource := getMatchedResourceFromConfigfile(newResource, oldResources)
		resourceMap := map[string]interface{}{
			// Required parameters.
			"flavor": getStringValueByConfigFileKey(newResource, "flavor_id", matchedOldResource),
			"count":  getIntValueByConfigFileKey(newResource, "count", matchedOldResource),
			// Only Optional parameter(s).
			"tags": buildUpdateResourcePoolResourceTags(newResource),
			// The parameters of the Computed behavior.
			"nodePool": getStringValueByConfigFileKey(newResource, "node_pool", matchedOldResource),
			"maxCount": getIntValueByConfigFileKey(newResource, "max_count", matchedOldResource),
			"azs":      buildUpdateResourcePoolResourceAzs(newResource, matchedOldResource),
			"network":  buildUpdateResourcePoolResourceNetwork(newResource, matchedOldResource),
			"taints":   buildUpdateResourcePoolResourceTaints(newResource, matchedOldResource),
			"labels":   buildUpdateResourcePoolResourceLabels(newResource, "labels", matchedOldResource),
			"extendParams": buildResourcePoolResourcesExtendParamsBodyParams(
				utils.PathSearch("extend_params", matchedOldResource, "{}").(string),
				getConfigFileStringValueByKey(newResource, "extend_params"),
				getConfigFileStringValueByKey(newResource, "post_install"),
			),
			"rootVolume":  buildUpdateResourcePoolResourceRootVolume(newResource, matchedOldResource),
			"dataVolumes": buildUpdateResourcePoolResourceDataVolumes(newResource, matchedOldResource),
			"volumeGroupConfigs": buildResourcePoolResourcesVolumeGroupConfigs(
				utils.PathSearch("volume_group_configs", matchedOldResource, schema.NewSet(schema.HashString, nil)).(*schema.Set).List(),
				buildUpdateResourcePoolResourceVolumeGroupConfigs(newResource),
			),
			"os":           buildUpdateResourcePoolResourceOs(newResource, matchedOldResource),
			"driver":       buildUpdateResourcePoolResourceDriver(newResource, matchedOldResource),
			"creatingStep": buildUpdateResourcePoolResourceCreatingStep(newResource, matchedOldResource),
		}
		result = append(result, resourceMap)
	}
	return result
}

func buildUpdateResourcePoolResourceVolumeGroupConfigs(resourceElem cty.Value) []interface{} {
	raw := getRawConfigSetValueByKey(resourceElem, "volume_group_configs")
	if raw == nil {
		return make([]interface{}, 0)
	}

	volumeGroupConfigs := raw.(cty.Value)
	result := make([]interface{}, volumeGroupConfigs.LengthInt())
	for i, volumeGroupConfigElem := range volumeGroupConfigs.AsValueSlice() {
		result[i] = map[string]interface{}{
			"volume_group":     volumeGroupConfigElem.GetAttr("volume_group").AsString(),
			"docker_thin_pool": getConfigFileIntValueByKey(volumeGroupConfigElem, "docker_thin_pool"),
			"lvm_config":       buildUpdateResourceVolumeGroupConfigsLvmConfig(volumeGroupConfigElem),
			"types":            buildUpdateResourcePoolResourceVolumeGroupConfigsTypes(volumeGroupConfigElem),
		}
	}
	return result
}

func buildUpdateResourceVolumeGroupConfigsLvmConfig(resourceElem cty.Value) []interface{} {
	if !isRawConfigListExist(resourceElem, "lvm_config") {
		return make([]interface{}, 0)
	}

	lvmConfig := resourceElem.GetAttr("lvm_config")
	return []interface{}{
		map[string]interface{}{
			"lv_type": lvmConfig.Index(cty.NumberIntVal(0)).GetAttr("lv_type").AsString(),
			"path":    getConfigFileStringValueByKey(lvmConfig.Index(cty.NumberIntVal(0)), "path"),
		},
	}
}

func buildUpdateResourcePoolResourceVolumeGroupConfigsTypes(elem cty.Value) []interface{} {
	if !elem.Type().HasAttribute("types") {
		return nil
	}

	raw := elem.GetAttr("types")
	if raw.IsNull() || !raw.IsKnown() || !raw.Type().IsListType() || raw.LengthInt() < 1 {
		return nil
	}

	types := make([]interface{}, raw.LengthInt())
	for i, typeElem := range raw.AsValueSlice() {
		types[i] = typeElem.AsString()
	}
	return types
}

func buildResourcePoolResourcesAzs(azSet *schema.Set) []map[string]interface{} {
	if azSet.Len() < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, azSet.Len())
	for _, az := range azSet.List() {
		result = append(result, map[string]interface{}{
			"az":    utils.ValueIgnoreEmpty(utils.PathSearch("az", az, nil)),
			"count": utils.ValueIgnoreEmpty(utils.PathSearch("count", az, nil)),
		})
	}

	return result
}

func buildResourcePoolResourcesTaints(taintSet *schema.Set) []map[string]interface{} {
	if taintSet.Len() < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, taintSet.Len())
	for _, taint := range taintSet.List() {
		result = append(result, map[string]interface{}{
			"key":    utils.ValueIgnoreEmpty(utils.PathSearch("key", taint, nil)),
			"value":  utils.ValueIgnoreEmpty(utils.PathSearch("value", taint, nil)),
			"effect": utils.ValueIgnoreEmpty(utils.PathSearch("effect", taint, nil)),
		})
	}

	return result
}

func resourceModelartsResourcePoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		resourcePoolName = d.Id()
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	bssClient, err := cfg.NewServiceClient("bssv2", region)
	if err != nil {
		return diag.Errorf("error creating BSS client: %s", err)
	}

	// If there are nodes in the prepaid billing mode under the resource pool (pre-paid or post-paid), we must unsubscribe the nodes first.
	if err := unsubscribePrePaidBillingNodes(ctx, client, bssClient, resourcePoolName, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error unsubscribing nodes under specified resource pool (%s): %s", resourcePoolName, err)
	}

	// When there is no node in the resource pool, the resource pool will be automatically deleted.
	_, err = queryResourcePool(cfg, region, d)
	if _, ok := err.(golangsdk.ErrDefault404); ok {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting resource pool (%s)", resourcePoolName))
	}

	if d.Get("charging_mode").(string) == "prePaid" {
		resourcePoolId := d.Get("resource_pool_id")
		if resourcePoolId == nil {
			return diag.Errorf("error getting resource ID from the resource pool (%s)", d.Id())
		}
		if err := common.UnsubscribePrePaidResource(d, cfg, []string{resourcePoolId.(string)}); err != nil {
			return diag.Errorf("error unsubscribing Modelarts resource pool: %s", err)
		}
	} else {
		err := deleteResourcePool(client, resourcePoolName)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = deleteResourcePoolWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the Modelarts resource pool (%s) deletion to complete: %s", d.Id(), err)
	}
	return nil
}

func deleteResourcePool(client *golangsdk.ServiceClient, resourcePoolName string) error {
	deleteHttpUrl := "v2/{project_id}/pools/{pool_name}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{pool_name}", resourcePoolName)

	deleteResourcePoolOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err := client.Request("DELETE", deletePath, &deleteResourcePoolOpt)
	if err != nil {
		return fmt.Errorf("error deleting Modelarts resource pool: %s", err)
	}
	return nil
}

func unsubscribePrePaidBillingNodes(ctx context.Context, client, bssClient *golangsdk.ServiceClient, resourcePoolName string,
	timeout time.Duration) error {
	nodes, err := listV2ResourcePoolNodes(client, resourcePoolName)
	if err != nil {
		return fmt.Errorf("error querying nodes under specified resource pool (%s): %s", resourcePoolName, err)
	}

	// Obtain the node IDs list that are in the pre-paid billing mode.
	deleteNodeIds := utils.PathSearch(
		fmt.Sprintf(`[?metadata.annotations."os.modelarts/billing.mode"=='%s'].metadata.labels."os.modelarts/resource.id"`, billingModePrePaid),
		nodes, make([]interface{}, 0)).([]interface{})

	if len(deleteNodeIds) == 0 {
		return nil
	}

	// Unsubscribe the pre-paid billing nodes.
	err = cbc.UnsubscribePrePaidResources(bssClient, deleteNodeIds)
	if err != nil {
		return err
	}
	err = cbc.WaitForResourcesUnsubscribed(ctx, bssClient, deleteNodeIds, timeout)
	if err != nil {
		return fmt.Errorf("error waiting for all nodes to be unsubscribed: %s ", err)
	}

	err = waitForV2NodeBatchUnsubscribeCompleted(ctx, client, resourcePoolName, deleteNodeIds, timeout)
	if err != nil {
		return err
	}

	return nil
}

func deleteResourcePoolWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			res, err := queryResourcePool(cfg, region, d)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					res = map[string]string{"code": "COMPLETED"}
					return res, "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			return res, "PENDING", nil
		},
		Timeout:      t,
		Delay:        30 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceResourcePoolImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	mErr := multierror.Append(
		d.Set("resources_order_origin", nil),
		d.Set("server_ids", nil),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
