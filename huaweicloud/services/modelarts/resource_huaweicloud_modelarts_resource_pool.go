// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product ModelArts
// ---------------------------------------------------------------

package modelarts

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

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

const billingModePostPaid = "0"
const billingModePrePaid = "1"

// @API ModelArts POST /v2/{project_id}/pools
// @API ModelArts GET /v2/{project_id}/pools/{id}
// @API ModelArts PATCH /v2/{project_id}/pools/{id}
// @API ModelArts GET /v2/{project_id}/pools/{pool_name}/nodepools/{nodepool_name}/nodes
// @API ModelArts DELETE /v2/{project_id}/pools/{id}
func ResourceModelartsResourcePool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceModelartsResourcePoolCreate,
		ReadContext:   resourceModelartsResourcePoolRead,
		UpdateContext: resourceModelartsResourcePoolUpdate,
		DeleteContext: resourceModelartsResourcePoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
				Type:        schema.TypeList,
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
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"step": {
							Type:        schema.TypeInt,
							Required:    true,
							ForceNew:    true,
							Description: `The creation step of the resource pool nodes.`,
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: `The type of the resource pool nodes.`,
						},
					},
				},
				Description: `The creation step configuration of the resource pool nodes.`,
			},
			// Internal attribute(s).
			"volume_group_configs_origin": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"volume_group": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(
								`The script configuration value of this change is also the original value used for the new value
next time in the list build method and reorder. The corresponding parameter name is 'volume_group_configs.volume_group'.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
						"docker_thin_pool": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: utils.SchemaDesc(
								`The script configuration value of this change is also the original value used for the new value
next time in the list build method and reorder. The corresponding parameter name is 'volume_group_configs.docker_thin_pool'.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
						"lvm_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lv_type": {
										Type:     schema.TypeString,
										Computed: true,
										Description: utils.SchemaDesc(
											`The script configuration value of this change is also the original value used for the new value
next time in the list build method and reorder. The corresponding parameter name is 'volume_group_configs.lvm_config.lv_type'.`,
											utils.SchemaDescInput{
												Internal: true,
											},
										),
									},
									"path": {
										Type:     schema.TypeString,
										Computed: true,
										Description: utils.SchemaDesc(
											`The script configuration value of this change is also the original value used for the new value
next time in the list build method and reorder. The corresponding parameter name is 'volume_group_configs.lvm_config.path'.`,
											utils.SchemaDescInput{
												Internal: true,
											},
										),
									},
								},
							},
							Description: utils.SchemaDesc(
								`The script configuration value of this change is also the original value used for the new value
next time in the list build method and reorder. The corresponding parameter name is 'volume_group_configs.lvm_config'.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
						"types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Description: utils.SchemaDesc(
								`The script configuration value of this change is also the original value used for the new value
next time in the list build method and reorder. The corresponding parameter name is 'volume_group_configs.types'.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
					},
				},
				Description: utils.SchemaDesc(
					`The script configuration value of this change is also the original value used for the new value
next time in the list build method and reorder. The corresponding parameter name is 'volume_group_configs'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

func scopeStatusRefreshFunc(cfg *config.Config, region string, d *schema.ResourceData, scopes []interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getResourcePoolRespBody, err := queryResourcePool(cfg, region, d)
		if err != nil {
			return getResourcePoolRespBody, "ERROR", err
		}

		for _, scope := range scopes {
			scopeStatus := fmt.Sprintf("status.scope[?scopeType=='%s']|[0].state", scope)
			if utils.PathSearch(scopeStatus, getResourcePoolRespBody, "").(string) != "Enabled" {
				return "No matches found", "PENDING", nil
			}
		}
		return "Matched", "COMPLETED", nil
	}
}

func createResourcePoolWaitingForScopesCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, timeout time.Duration) error {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: scopeStatusRefreshFunc(cfg, region, d, d.Get("scope").(*schema.Set).List()),
		Timeout: timeout,
		// In most cases, the bind operation will be completed immediately, but in a few cases, it needs to wait
		// for a short period of time, and the polling is performed by incrementing the time here.
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the scope statuses are both completed: %s", err)
	}
	return nil
}

func waitForDriverStatusCompleted(ctx context.Context, cfg *config.Config, region string, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      driverStatusRefreshFunc(cfg, region, d),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		PollInterval: 10 * time.Second,
		// In some cases, the following status changes may occur: Upgrading -> Running -> Creating -> Running
		ContinuousTargetOccurence: 2,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func driverStatusRefreshFunc(cfg *config.Config, region string, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resourcePool, err := queryResourcePool(cfg, region, d)
		if err != nil {
			return resourcePool, "ERROR", err
		}

		// Cueerntly, only for GPU driver or NPU driver.
		driverStatuses := utils.PathSearch("status.driver.*.state", resourcePool, make([]interface{}, 0)).([]interface{})
		if len(driverStatuses) == 0 {
			return "No matches found", "COMPLETED", nil
		}

		for _, status := range driverStatuses {
			if status == "Abnormal" {
				return resourcePool, "ERROR", fmt.Errorf("unexpect status (%s)", status)
			}

			if status != "Running" {
				return resourcePool, "PENDING", nil
			}
		}

		return resourcePool, "COMPLETED", nil
	}
}

func resourceModelartsResourcePoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createResourcePoolHttpUrl = "v2/{project_id}/pools"
		createResourcePoolProduct = "modelarts"
	)
	createResourcePoolClient, err := cfg.NewServiceClient(createResourcePoolProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	createResourcePoolPath := createResourcePoolClient.Endpoint + createResourcePoolHttpUrl
	createResourcePoolPath = strings.ReplaceAll(createResourcePoolPath, "{project_id}", createResourcePoolClient.ProjectID)

	createResourcePoolOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	createResourcePoolOpt.JSONBody = utils.RemoveNil(buildCreateResourcePoolBodyParams(d))
	createResourcePoolResp, err := createResourcePoolClient.Request("POST", createResourcePoolPath, &createResourcePoolOpt)
	if err != nil {
		return diag.Errorf("error creating Modelarts resource pool: %s", err)
	}

	createResourcePoolRespBody, err := utils.FlattenResponse(createResourcePoolResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("metadata.name", createResourcePoolRespBody, nil)
	if id == nil {
		return diag.Errorf("error creating Modelarts resource pool: ID is not found in API response")
	}
	d.SetId(id.(string))

	if d.Get("charging_mode") == "prePaid" {
		// wait 30 seconds so that the resource pool can be queried
		// lintignore:R018
		time.Sleep(30 * time.Second)
		resourcePool, err := queryResourcePool(cfg, region, d)
		if err != nil {
			return diag.Errorf("error retrieving Modelarts resource pool: %s", err)
		}
		orderId := utils.PathSearch(`metadata.annotations."os.modelarts/order.id"`, resourcePool, nil)
		if orderId == nil {
			return diag.Errorf("error creating Modelarts resource pool: order ID is not found in API response")
		}
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		err = common.WaitOrderComplete(ctx, bssClient, orderId.(string), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = common.WaitOrderResourceComplete(ctx, bssClient, orderId.(string), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = createResourcePoolWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the Modelarts resource pool (%s) creation to complete: %s", d.Id(), err)
	}

	err = createResourcePoolWaitingForScopesCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the Modelarts resource pool (%s) creation to complete: %s", d.Id(), err)
	}

	err = waitForDriverStatusCompleted(ctx, cfg, region, d)
	if err != nil {
		return diag.Errorf("error waiting for the Modelarts resource pool (%s) driver status to become running: %s", d.Id(), err)
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
	oldResourcesVal, newResourcesVal := d.GetChange("resources")

	params := map[string]interface{}{
		"type":      "Dedicate",
		"scope":     utils.ValueIgnoreEmpty(d.Get("scope").(*schema.Set).List()),
		"resources": buildResourcePoolSpecResources(oldResourcesVal.([]interface{}), newResourcesVal.([]interface{})),
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

func buildResourcePoolSpecResources(oldResources, newResources []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(newResources))

	for i, newResource := range newResources {
		result = append(result, map[string]interface{}{
			"flavor":   utils.ValueIgnoreEmpty(utils.PathSearch("flavor_id", newResource, nil)),
			"count":    utils.ValueIgnoreEmpty(utils.PathSearch("count", newResource, nil)),
			"nodePool": utils.ValueIgnoreEmpty(utils.PathSearch("node_pool", newResource, nil)),
			"maxCount": utils.ValueIgnoreEmpty(utils.PathSearch("max_count", newResource, nil)),
			"azs": buildResourcePoolResourcesAzs(utils.PathSearch("azs", newResource,
				schema.NewSet(schema.HashString, nil)).(*schema.Set)),
			"network": buildResourcePoolSpecResourcesNetworkBodyParams(newResource),
			"taints": buildResourcePoolResourcesTaints(utils.PathSearch("taints", newResource,
				schema.NewSet(schema.HashString, nil)).(*schema.Set)),
			"tags":   utils.ExpandResourceTags(utils.PathSearch("tags", newResource, make(map[string]interface{})).(map[string]interface{})),
			"labels": utils.ValueIgnoreEmpty(utils.PathSearch("labels", newResource, nil)),
			"extendParams": buildResourcePoolResourcesExtendParamsBodyParams(
				utils.PathSearch(fmt.Sprintf("[%d].extend_params", i), oldResources, "{}").(string),
				utils.PathSearch("extend_params", newResource, "{}").(string),
				utils.PathSearch("post_install", newResource, "").(string),
			),
			"rootVolume": buildResourcePoolResourcesRootVolume(utils.PathSearch("root_volume", newResource,
				make([]interface{}, 0)).([]interface{})),
			"dataVolumes": buildResourcePoolResourcesDataVolumes(
				utils.PathSearch(fmt.Sprintf("[%d].data_volumes", i), oldResources, make([]interface{}, 0)).([]interface{}),
				utils.PathSearch("data_volumes", newResource, make([]interface{}, 0)).([]interface{}),
			),
			"volumeGroupConfigs": buildResourcePoolResourcesVolumeGroupConfigs(
				utils.PathSearch(fmt.Sprintf("[%d].volume_group_configs_origin", i), oldResources, make([]interface{}, 0)).([]interface{}),
				utils.PathSearch("volume_group_configs", newResource, schema.NewSet(schema.HashString, nil)).(*schema.Set).List(),
			),
			"os":     buildResourcePoolResourcesOsInfo(utils.PathSearch("os", newResource, make([]interface{}, 0)).([]interface{})),
			"driver": buildResourcePoolResourcesDriver(utils.PathSearch("driver", newResource, make([]interface{}, 0)).([]interface{})),
			"creatingStep": buildResourcePoolResourcesCreatingStep(
				utils.PathSearch("creating_step", newResource, make([]interface{}, 0)).([]interface{})),
		})
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
		"volumeType": utils.PathSearch("volume_type", rootVolume, nil),
		"size":       utils.PathSearch("size", rootVolume, nil),
	}
}

func buildResourcePoolResourcesDataVolumes(oldDataVolumes, newDataVolumes []interface{}) []map[string]interface{} {
	if len(newDataVolumes) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(newDataVolumes))
	for i, dataVolume := range newDataVolumes {
		result = append(result, map[string]interface{}{
			"volumeType": utils.PathSearch("volume_type", dataVolume, nil),
			"size":       utils.PathSearch("size", dataVolume, nil),
			"extendParams": buildResourcePoolResourcesExtendParamsBodyParams(
				utils.PathSearch(fmt.Sprintf("[%d].extend_params", i), oldDataVolumes, "").(string),
				utils.PathSearch("extend_params", dataVolume, "").(string),
				"",
			),
			"count": utils.ValueIgnoreEmpty(utils.PathSearch("count", dataVolume, nil)),
		})
	}

	return result
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

func createResourcePoolWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)

			createResourcePoolWaitingRespBody, err := queryResourcePool(cfg, region, d)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw := utils.PathSearch(`status.phase`, createResourcePoolWaitingRespBody, nil)
			if statusRaw == nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `status.phase`)
			}

			if utils.PathSearch("status.resources.abnormal", createResourcePoolWaitingRespBody, nil) != nil {
				return nil, "ERROR", fmt.Errorf("error creating resource pool: the resource pool is abnormal")
			}

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"Running",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return createResourcePoolWaitingRespBody, "COMPLETED", nil
			}

			pendingStatus := []string{
				"Creating",
				"Waiting",
			}
			if utils.StrSliceContains(pendingStatus, status) {
				return createResourcePoolWaitingRespBody, "PENDING", nil
			}

			return createResourcePoolWaitingRespBody, status, nil
		},
		Timeout:      t,
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
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
		d.Set("resources", flattenGetResourcePoolResponseBodyResources(getModelartsResourcePoolRespBody)),
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
			"volume_type": utils.PathSearch("volumeType", rootVolume, nil),
			"size":        utils.PathSearch("size", rootVolume, nil),
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

func flattenGetResourcePoolResponseBodyResources(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("spec.resources", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
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
			"volume_group_configs_origin": flattenResourcePoolResourcesVolumeGroupConfigs(utils.PathSearch("volumeGroupConfigs",
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

func resourceModelartsResourcePoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

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
		updateResourcePoolClient, err := cfg.NewServiceClient(updateResourcePoolProduct, region)
		if err != nil {
			return diag.Errorf("error creating ModelArts client: %s", err)
		}

		updateResourcePoolPath := updateResourcePoolClient.Endpoint + updateResourcePoolHttpUrl
		updateResourcePoolPath = strings.ReplaceAll(updateResourcePoolPath, "{project_id}", updateResourcePoolClient.ProjectID)
		updateResourcePoolPath = strings.ReplaceAll(updateResourcePoolPath, "{id}", d.Id())

		updateResourcePoolOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
			MoreHeaders: map[string]string{"Content-Type": "application/merge-patch+json"},
		}

		updateResourcePoolOpt.JSONBody = utils.RemoveNil(buildUpdateResourcePoolBodyParams(d))
		updateModelartsResourcePoolResp, err := updateResourcePoolClient.Request("PATCH", updateResourcePoolPath, &updateResourcePoolOpt)
		if err != nil {
			return diag.Errorf("error updating Modelarts resource pool: %s", err)
		}

		// The `nodeBillingMode` indicates the billing mode of the node when scaling the node capacity.
		nodeBillingMode := ""
		if annotations, ok := d.GetOk("metadata.0.annotations"); ok {
			nodeBillingMode = utils.PathSearch(`"os.modelarts/billing.mode"`, utils.StringToJson(annotations.(string)), billingModePostPaid).(string)
		}
		// Only when scaling prepaid type nodes, we need to determine the order status.
		if nodeBillingMode == "1" && d.HasChange("resources") {
			// Whenever any count in resouces changes, the order status needs to be determined.
			oldRaw, newRaw := d.GetChange("resources")
			if isAnyNodeScalling(oldRaw.([]interface{}), newRaw.([]interface{})) {
				updateRespBody, err := utils.FlattenResponse(updateModelartsResourcePoolResp)
				if err != nil {
					return diag.FromErr(err)
				}

				orderId := utils.PathSearch(`metadata.annotations."os.modelarts/order.id"`, updateRespBody, nil)
				if orderId == nil {
					return diag.Errorf("error updating Modelarts resource pool: order ID is not found in API response")
				}

				bssClient, err := cfg.BssV2Client(region)
				if err != nil {
					return diag.Errorf("error creating BSS v2 client: %s", err)
				}
				err = common.WaitOrderComplete(ctx, bssClient, orderId.(string), d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return diag.FromErr(err)
				}
				_, err = common.WaitOrderAllResourceComplete(ctx, bssClient, orderId.(string), d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}

		err = updateResourcePoolWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the Modelarts resource pool (%s) update to complete: %s", d.Id(), err)
		}
	}
	return resourceModelartsResourcePoolRead(ctx, d, meta)
}

func isAnyNodeScalling(oldResource, newResource []interface{}) bool {
	for i, v := range newResource {
		oldCount := utils.PathSearch(fmt.Sprintf("[%d].count", i), oldResource, 0).(int)
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
	if !isAnyNodePoolIncrease(d) {
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

func isAnyNodePoolIncrease(d *schema.ResourceData) bool {
	oldRaw, newRaw := d.GetChange("resources")
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
	oldResourcesVal, newResourcesVal := d.GetChange("resources")

	params := map[string]interface{}{
		"scope":     utils.ValueIgnoreEmpty(d.Get("scope").(*schema.Set).List()),
		"resources": buildResourcePoolSpecResources(oldResourcesVal.([]interface{}), newResourcesVal.([]interface{})),
	}
	return params
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

func updateResourcePoolWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)

			getResourcePoolRespBody, err := queryResourcePool(cfg, region, d)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw := utils.PathSearch(`status.phase`, getResourcePoolRespBody, nil)
			if statusRaw == nil {
				return nil, "ERROR", fmt.Errorf("error parse 'status.phase' from response body")
			}

			status := fmt.Sprintf("%v", statusRaw)

			unexpectedStatus := []string{
				"Abnormal", "Error", "ScalingFailed", "CreationFailed",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return getResourcePoolRespBody, status, nil
			}

			// check if the resource pool is in the process of changing scope
			if rawArray, ok := d.GetOk("scope"); ok {
				for _, v := range rawArray.(*schema.Set).List() {
					scopeStatus := fmt.Sprintf("status.scope[?scopeType=='%s']|[0].state", v)
					log.Println("scopeStatus: ", scopeStatus)
					if utils.PathSearch(scopeStatus, getResourcePoolRespBody, "").(string) != "Enabled" {
						return getResourcePoolRespBody, "PENDING", nil
					}
				}
			}

			creating := utils.PathSearch("status.resources.creating", getResourcePoolRespBody, make([]interface{}, 0)).([]interface{})
			deleting := utils.PathSearch("status.resources.deleting", getResourcePoolRespBody, make([]interface{}, 0)).([]interface{})
			creationFailed := utils.PathSearch("status.resources.creationFailed", getResourcePoolRespBody, make([]interface{}, 0)).([]interface{})

			// check if the resource pool is in the process of expanding capacity
			if len(creating) == 0 && len(deleting) == 0 && len(creationFailed) == 0 {
				return getResourcePoolRespBody, "COMPLETED", nil
			}

			return getResourcePoolRespBody, "PENDING", nil
		},
		Timeout:                   t,
		Delay:                     20 * time.Second,
		PollInterval:              20 * time.Second,
		ContinuousTargetOccurence: 2,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
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

	bssClient, err := cfg.NewServiceClient("bssv2", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating BSS client: %s", err)
	}

	// If there are nodes in the prepaid billing mode under the resource pool (pre-paid or post-paid), we must unsubscribe the nodes first.
	if err := unsubscribePrePaidBillingNodes(ctx, client, bssClient, resourcePoolName, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error unsubscribing nodes under specified resource pool (%s): %s", resourcePoolName, err)
	}

	if d.Get("charging_mode").(string) == "prePaid" {
		resourcePoolId := d.Get("resource_pool_id")
		if resourcePoolId == nil {
			return diag.Errorf("error getting resource ID from the resource pool (%s)", resourcePoolName)
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

	err = deleteResourcePoolWaitingForStateCompleted(ctx, cfg, region, d)
	if err != nil {
		return diag.Errorf("error waiting for the Modelarts resource pool (%s) deletion to complete: %s", resourcePoolName, err)
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
		fmt.Sprintf(`[?metadata.annotations."os.modelarts/billing.mode"=='%s'].metadata.labels."os.modelarts/resource.id"`,
			billingModePrePaid),
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

func deleteResourcePool(client *golangsdk.ServiceClient, resourcePoolName string) error {
	deleteResourcePoolHttpUrl := "v2/{project_id}/pools/{id}"
	deleteResourcePoolPath := client.Endpoint + deleteResourcePoolHttpUrl
	deleteResourcePoolPath = strings.ReplaceAll(deleteResourcePoolPath, "{project_id}", client.ProjectID)
	deleteResourcePoolPath = strings.ReplaceAll(deleteResourcePoolPath, "{id}", resourcePoolName)

	deleteResourcePoolOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	_, err := client.Request("DELETE", deleteResourcePoolPath, &deleteResourcePoolOpt)
	if err != nil {
		return fmt.Errorf("error deleting Modelarts resource pool: %s", err)
	}
	return nil
}

func deleteResourcePoolWaitingForStateCompleted(ctx context.Context, cfg *config.Config, region string, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
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
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        30 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
