package taurusdb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbc"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var htapStarrocksInstanceNoneUpdatableParams = []string{
	"instance", "name", "ha", "ha.*.mode", "az_mode", "az_code", "time_zone",
	"engine", "engine.*.type", "engine.*.version", "tags_info", "tags_info.*.sys_tags",
	"fe_volume", "fe_volume.*.io_type", "fe_volume.*.capacity_in_gb", "fe_count",
	"be_volume", "be_volume.*.io_type", "be_volume.*.capacity_in_gb", "be_count",
	"charging_mode", "period_unit", "period",
}

var htapStarrocksInstanceSchema = map[string]*schema.Schema{
	"region": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		ForceNew: true,
	},
	// URI parameters
	"instance_id": {
		Type:     schema.TypeString,
		Required: true,
	},
	// Request parameters
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"engine": {
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem:     starrocksInstanceEngineSchema(),
	},
	"ha": {
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem:     starrocksInstanceHaSchema(),
	},
	"fe_flavor_id": {
		Type:     schema.TypeString,
		Required: true,
	},
	"be_flavor_id": {
		Type:     schema.TypeString,
		Required: true,
	},
	"db_root_pwd": {
		Type:      schema.TypeString,
		Required:  true,
		Sensitive: true,
	},
	"fe_count": {
		Type:     schema.TypeInt,
		Required: true,
	},
	"be_count": {
		Type:     schema.TypeInt,
		Required: true,
	},
	"az_mode": {
		Type:     schema.TypeString,
		Required: true,
	},
	"fe_volume": {
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem:     starrocksInstanceFeVolumeSchema(),
	},
	"be_volume": {
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem:     starrocksInstanceBeVolumeSchema(),
	},
	"az_code": {
		Type:     schema.TypeString,
		Required: true,
	},
	"time_zone": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		DiffSuppressFunc: func(_, old, newVal string, _ *schema.ResourceData) bool {
			return newVal == "" && old != ""
		},
	},
	"tags_info": {
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem:     starrocksInstanceTagsInfoSchema(),
	},
	"security_group_id": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"enable_users_sync": {
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
	},
	"open_slow_log_switch": {
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
	},
	"query_queue_rule": {
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"query_queue_max_queued_queries": {
					Type:     schema.TypeInt,
					Required: true,
				},
				"query_queue_pending_timeout_second": {
					Type:     schema.TypeInt,
					Required: true,
				},
				"query_queue_concurrency_limit": {
					Type:     schema.TypeInt,
					Required: true,
				},
				"query_queue_mem_used_pct_limit": {
					Type:     schema.TypeInt,
					Required: true,
				},
				"query_queue_cpu_used_pct_limit": {
					Type:     schema.TypeInt,
					Required: true,
				},
				"enable_query_queue_select": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				},
			},
		},
	},
	"be_parameter_values": {
		Type:     schema.TypeMap,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	"fe_parameter_values": {
		Type:     schema.TypeMap,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	// charge info: charging_mode, period_unit, period, auto_renew
	"charging_mode": {
		Type:     schema.TypeString,
		Optional: true,
		ValidateFunc: validation.StringInSlice([]string{
			"prePaid", "postPaid",
		}, false),
		DiffSuppressFunc: func(_, old, newVal string, _ *schema.ResourceData) bool {
			return newVal == "" && old != ""
		},
	},
	"period_unit": {
		Type:         schema.TypeString,
		Optional:     true,
		RequiredWith: []string{"period"},
		ValidateFunc: validation.StringInSlice([]string{"month", "year"}, false),
	},
	"period": {
		Type:         schema.TypeInt,
		Optional:     true,
		RequiredWith: []string{"period_unit"},
	},
	"auto_renew": common.SchemaAutoRenewUpdatable(nil),
	"enable_force_new": {
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
		Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
	},
	// Computed attributes
	"project_id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"public_ip": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"data_vip": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"can_enable_public_access": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"current_backup_node_id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"cluster_mode": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"status": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"is_frozen": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"frozen_time": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"db_user": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"bak_period": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"bak_keep_day": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"bak_expected_start_time": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"data_store_version_id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"data_store_version": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"data_store_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"create_at": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"update_at": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"delete_at": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"db_port": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"param_group": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"actions": {
		Type:     schema.TypeList,
		Computed: true,
		Elem:     starrocksInstanceActionsSchema(),
	},
	"create_fail_error_code": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"users_sync_switch_on": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"groups": {
		Type:     schema.TypeList,
		Computed: true,
		Elem:     starrocksInstanceGroupsSchema(),
	},
	"ops_window": {
		Type:     schema.TypeList,
		Computed: true,
		Elem:     starrocksInstanceOpsWindowsSchema(),
	},
	"backup_used_space": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"enterprise_project_id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"port_info": {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"mysql_port": {
					Type:     schema.TypeInt,
					Computed: true,
				},
			},
		},
	},
	"fe_node_volume_code": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"be_node_volume_code": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"fe_node_volume_size": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"be_node_volume_size": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"support_data_replication": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"new_version_available": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"ssl_option": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"dedicated_resource_id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"vpc_id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"subnet_id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"be_configurations": {
		Type:     schema.TypeList,
		Computed: true,
		Elem:     starrocksParametersConfigurationsSchema(),
	},
	"be_parameters": {
		Type:     schema.TypeList,
		Computed: true,
		Elem:     starrocksParametersParameterValuesSchema(),
	},
	"fe_configurations": {
		Type:     schema.TypeList,
		Computed: true,
		Elem:     starrocksParametersConfigurationsSchema(),
	},
	"fe_parameters": {
		Type:     schema.TypeList,
		Computed: true,
		Elem:     starrocksParametersParameterValuesSchema(),
	},
}

// @API BSS POST /v3/orders/customer-orders/pay
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/{project_id}/orders/auto-renew
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
// @API TaurusDB POST /v3/{project_id}/instances/{instance_id}/starrocks
// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/starrocks/configurations
// @API TaurusDB PUT /v3/{project_id}/instances/{instance_id}/starrocks/configurations
// @API TaurusDB PUT /v3/{project_id}/instances/{starrocks_instance_id}/starrocks/restart
// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/starrocks/{starrocks_instance_id}
// @API TaurusDB DELETE /v3/{project_id}/instances/{instance_id}/starrocks/{starrocks_instance_id}
// @API TaurusDB POST /v3/{project_id}/instances/{instance_id}/starrocks/resize-flavor
// @API TaurusDB PUT /v3/{project_id}/instances/{instance_id}/starrocks/users/password
// @API TaurusDB PUT /v3/{project_id}/instances/{instance_id}/starrocks/security-group
// @API TaurusDB POST /v3/{project_id}/instances/{instance_id}/starrocks/users/sync
// @API TaurusDB PUT /v3/{project_id}/instances/{instance_id}/starrocks/slowlog-sensitive
// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/starrocks/slowlog-sensitive
// @API TaurusDB PUT /v3/{project_id}/instances/{instance_id}/query-queue/rules
// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/query-queue/rules
// @API TaurusDB POST /v3/{project_id}/instances/{instance_id}/htap/query-queue/switch
// @API TaurusDB GET /v3/{project_id}/jobs
func ResourceTaurusDBHtapStarrocksInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaurusDBHtapStarrocksInstanceCreate,
		UpdateContext: resourceTaurusDBHtapStarrocksInstanceUpdate,
		ReadContext:   resourceTaurusDBHtapStarrocksInstanceRead,
		DeleteContext: resourceTaurusDBHtapStarrocksInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceTaurusDBHtapStarrocksInstanceImportState,
		},

		CustomizeDiff: customdiff.All(config.FlexibleForceNew(htapStarrocksInstanceNoneUpdatableParams, htapStarrocksInstanceSchema)),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: htapStarrocksInstanceSchema,
	}
}

func starrocksInstanceEngineSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func starrocksInstanceHaSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"mode": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func starrocksInstanceFeVolumeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"io_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"capacity_in_gb": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func starrocksInstanceBeVolumeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"io_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"capacity_in_gb": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func starrocksInstanceActionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"object_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func starrocksInstanceGroupsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_type_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_node_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     starrocksInstanceNodesSchema(),
			},
		},
	}
}

func starrocksInstanceNodesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"period": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"cpu": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mem": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datastore": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"actions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     starrocksInstanceActionsSchema(),
			},
			"priority": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"frozen_flag": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"db_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"pay_model": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"order_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"traffic_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"traffic_ipv6": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"az_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"az_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"az_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flavor_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"iass_flavor_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_connections": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"need_restart": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"sg_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"param_group": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func starrocksInstanceOpsWindowsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"period": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func starrocksInstanceTagsInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"sys_tags": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

// handlePrePaidInstanceCreation handles the creation process for prePaid charging mode
func handlePrePaidInstanceCreation(ctx context.Context, cfg *config.Config, d *schema.ResourceData,
	taurusdbInstanceId string, createRespBody interface{}) error {
	orderId := utils.PathSearch("order_id", createRespBody, "").(string)
	if orderId == "" {
		return fmt.Errorf("error creating HTAP StarRocks instance for TaurusDB instance (%s): "+
			"order_id is not found in API response", taurusdbInstanceId)
	}

	region := cfg.GetRegion(d)
	bssClient, err := cfg.BssV2Client(region)
	if err != nil {
		return fmt.Errorf("error creating BSS client: %s", err)
	}

	// The expand nodes interface does not support auto-payment, so the CBC side interface is called to pay for the order.
	if err := cbc.PaySubscriptionOrder(bssClient, orderId); err != nil {
		return fmt.Errorf("error paying for expansion order (%s): %s", orderId, err)
	}

	// 1. If charging mode is PrePaid, wait for the order to be completed.
	err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	// 2. get the instance ID, must be after order success
	htapInstanceId, err := common.WaitOrderResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	d.SetId(htapInstanceId)
	return nil
}

// handlePostPaidInstanceCreation handles the creation process for postPaid charging mode
func handlePostPaidInstanceCreation(ctx context.Context, client *golangsdk.ServiceClient,
	d *schema.ResourceData, taurusdbInstanceId string, createRespBody interface{}) error {
	htapInstanceId := utils.PathSearch("instance.id", createRespBody, "").(string)
	if htapInstanceId == "" {
		return fmt.Errorf("error creating HTAP StarRocks instance for TaurusDB instance (%s): "+
			"instance ID is not found in API response", taurusdbInstanceId)
	}

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return fmt.Errorf("error creating HTAP StarRocks instance for TaurusDB instance (%s): "+
			"job_id is not found in API response", taurusdbInstanceId)
	}

	err := checkGaussDBMySQLJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error waiting for creating TaurusDB HTAP StarRocks instance (%s) job to complete: %s",
			htapInstanceId, err)
	}

	d.SetId(htapInstanceId)
	return nil
}

// configureStarRocksParameters configures BE and FE parameters after instance creation
func configureStarRocksParameters(ctx context.Context, client *golangsdk.ServiceClient,
	d *schema.ResourceData, taurusdbInstanceId string) error {
	restartRequired := false

	beParamValues := d.Get("be_parameter_values").(map[string]interface{})
	if len(beParamValues) > 0 {
		var err error
		restartRequired, err = modifyStarrocksParameters(ctx, client, d, "be", beParamValues)
		if err != nil {
			return err
		}
	}

	feParamValues := d.Get("fe_parameter_values").(map[string]interface{})
	if len(feParamValues) > 0 {
		var err error
		restartRequired, err = modifyStarrocksParameters(ctx, client, d, "fe", feParamValues)
		if err != nil {
			return err
		}
	}

	if restartRequired {
		err := restartStarrocksInstance(ctx, client, taurusdbInstanceId, d.Id(), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceTaurusDBHtapStarrocksInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/instances/{instance_id}/starrocks"
		product = "gaussdb"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	taurusdbInstanceId := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", taurusdbInstanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildCreateStarrocksInstanceBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating HTAP StarRocks instance for TaurusDB instance (%s): %s", taurusdbInstanceId, err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// Handle instance creation based on charging mode
	if d.Get("charging_mode") == "prePaid" {
		if err := handlePrePaidInstanceCreation(ctx, cfg, d, taurusdbInstanceId, createRespBody); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := handlePostPaidInstanceCreation(ctx, client, d, taurusdbInstanceId, createRespBody); err != nil {
			return diag.FromErr(err)
		}
	}

	err = waitHtapInstanceStatusNormal(ctx, client, taurusdbInstanceId, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for HTAP StarRocks instance (%s) status to be normal: %s", d.Id(), err)
	}

	// Open the users sync
	if v, ok := d.GetOk("enable_users_sync"); ok && v.(string) == "true" {
		err = updateStarRocksInstanceUsersSyncSwitch(ctx, client, taurusdbInstanceId, d.Id(), true, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Configure StarRocks parameters (BE and FE)
	if err := configureStarRocksParameters(ctx, client, d, taurusdbInstanceId); err != nil {
		return diag.FromErr(err)
	}

	// Open slow log switch
	if v, ok := d.GetOk("open_slow_log_switch"); ok && v.(string) == "true" {
		if err := updateStarRocksSlowLogSwitch(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	// Update query queue rule and switch
	enableQueryQueueRule := d.Get("query_queue_rule.0.enable_query_queue_select").(string)
	if enableQueryQueueRule == "true" {
		if err := updateHtapQueryQueueRuleSwitch(client, d, enableQueryQueueRule); err != nil {
			return diag.FromErr(err)
		}
		if err := updateHtapQueryQueueRule(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceTaurusDBHtapStarrocksInstanceRead(ctx, d, meta)
}

func buildCreateStarrocksInstanceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":              d.Get("name"),
		"engine":            buildStarrocksInstanceEngineBodyParams(d),
		"ha":                buildStarrocksInstanceHaBodyParams(d),
		"fe_flavor_id":      d.Get("fe_flavor_id"),
		"be_flavor_id":      d.Get("be_flavor_id"),
		"db_root_pwd":       d.Get("db_root_pwd"),
		"fe_count":          d.Get("fe_count"),
		"be_count":          d.Get("be_count"),
		"az_mode":           d.Get("az_mode"),
		"fe_volume":         buildStarrocksInstanceFeVolumeBodyParams(d),
		"be_volume":         buildStarrocksInstanceBeVolumeBodyParams(d),
		"az_code":           d.Get("az_code"),
		"time_zone":         utils.ValueIgnoreEmpty(d.Get("time_zone")),
		"tags_info":         buildStarrocksInstanceTagsInfoBodyParams(d),
		"security_group_id": utils.ValueIgnoreEmpty(d.Get("security_group_id")),
		"pay_info":          buildStarrocksInstancePayInfoBodyParams(d),
	}
	return bodyParams
}

func buildStarrocksInstanceEngineBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"type":    d.Get("engine.0.type"),
		"version": d.Get("engine.0.version"),
	}
}

func buildStarrocksInstanceHaBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"mode": d.Get("ha.0.mode"),
	}
}

func buildStarrocksInstanceBeVolumeBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"io_type":        d.Get("be_volume.0.io_type"),
		"capacity_in_gb": d.Get("be_volume.0.capacity_in_gb"),
	}
}

func buildStarrocksInstanceFeVolumeBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"io_type":        d.Get("fe_volume.0.io_type"),
		"capacity_in_gb": d.Get("fe_volume.0.capacity_in_gb"),
	}
}

func buildStarrocksInstanceTagsInfoBodyParams(d *schema.ResourceData) map[string]interface{} {
	tagsInfoRaw := d.Get("tags_info.0").(map[string]interface{})
	if tagsInfoRaw == nil {
		return nil
	}

	result := make(map[string]interface{})

	// Only build sys_tags list for API request (tags is output-only)
	if sysTagsRaw, ok := tagsInfoRaw["sys_tags"].(*schema.Set); ok && sysTagsRaw.Len() > 0 {
		sysTagsList := make([]map[string]interface{}, 0, sysTagsRaw.Len())
		for _, tag := range sysTagsRaw.List() {
			tagMap := tag.(map[string]interface{})
			sysTagsList = append(sysTagsList, map[string]interface{}{
				"key":   tagMap["key"],
				"value": tagMap["value"],
			})
		}
		result["sys_tags"] = sysTagsList
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

func buildStarrocksInstancePayInfoBodyParams(d *schema.ResourceData) map[string]interface{} {
	// Only build pay info for prePaid charging mode
	if chargingMode := d.Get("charging_mode").(string); chargingMode != "prePaid" {
		return nil
	}

	payModel := "1"
	periodType := "2"
	if periodUnit := d.Get("period_unit").(string); periodUnit == "year" {
		periodType = "3"
	}

	period := d.Get("period").(int)

	isAutoRenew := "0"
	if autoRenew, ok := d.Get("auto_renew").(string); ok && autoRenew == "true" {
		isAutoRenew = "1"
	}

	return map[string]interface{}{
		"pay_model":     payModel,
		"period_type":   periodType,
		"period":        period,
		"is_auto_renew": isAutoRenew,
	}
}

func updateStarRocksInstanceUsersSyncSwitch(ctx context.Context, client *golangsdk.ServiceClient,
	taurusdbInstanceId, htapInstanceId string, enableUsersSync bool, timeout time.Duration) error {
	enableUsersSyncPath := client.Endpoint + "v3/{project_id}/instances/{instance_id}/starrocks/users/sync"
	enableUsersSyncPath = strings.ReplaceAll(enableUsersSyncPath, "{project_id}", client.ProjectID)
	enableUsersSyncPath = strings.ReplaceAll(enableUsersSyncPath, "{instance_id}", htapInstanceId)

	enableUsersSyncOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildStarRocksInstanceUsersSyncBodyParams(enableUsersSync),
	}

	retryFunc := func() (interface{}, bool, error) {
		resp, err := client.Request("POST", enableUsersSyncPath, &enableUsersSyncOpt)
		shouldRetry, err := handleMultiOperationsError(err)
		return resp, shouldRetry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     htapInstanceStateRefreshFunc(client, taurusdbInstanceId, htapInstanceId),
		WaitTarget:   []string{"normal"},
		Timeout:      timeout,
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating users sync switch of HTAP StarRocks instance(%s): %s", htapInstanceId, err)
	}
	respBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return err
	}
	success := utils.PathSearch("success", respBody, false).(bool)
	if !success {
		return fmt.Errorf("error updating StarRocks instance users sync switch on: %v", enableUsersSync)
	}

	if err = waitHtapInstanceStatusNormal(ctx, client, taurusdbInstanceId, htapInstanceId, timeout); err != nil {
		return fmt.Errorf("error waiting for HTAP StarRocks instance(%s) actions to be empty "+
			"after updating users sync switch: %s", htapInstanceId, err)
	}
	return nil
}

func updateStarRocksSlowLogSwitch(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var httpUrl = "v3/{project_id}/instances/{instance_id}/starrocks/slowlog-sensitive"

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Id())

	openSlowLogSwitch, err := strconv.ParseBool(d.Get("open_slow_log_switch").(string))
	if err != nil {
		return err
	}
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"open_slow_log_switch": openSlowLogSwitch,
		},
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating slow log switch of HTAP StarRocks instance(%s): %s", d.Id(), err)
	}

	return nil
}

func getStarRocksSlowLogSwitch(client *golangsdk.ServiceClient, htapInstanceId string) (string, error) {
	var httpUrl = "v3/{project_id}/instances/{instance_id}/starrocks/slowlog-sensitive"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", htapInstanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return "", fmt.Errorf("error querying slow log switch of HTAP StarRocks instance(%s): %s", htapInstanceId, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return "", err
	}

	enabled := utils.PathSearch("open_slow_log_switch", respBody, "").(string)
	return enabled, nil
}

func updateHtapQueryQueueRule(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var httpUrl = "v3/{project_id}/instances/{instance_id}/query-queue/rules"

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Id())

	ruleRaw := d.Get("query_queue_rule").([]interface{})
	ruleMap := ruleRaw[0].(map[string]interface{})

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"query_queue_rule": map[string]interface{}{
				"query_queue_max_queued_queries":     ruleMap["query_queue_max_queued_queries"],
				"query_queue_pending_timeout_second": ruleMap["query_queue_pending_timeout_second"],
				"query_queue_concurrency_limit":      ruleMap["query_queue_concurrency_limit"],
				"query_queue_mem_used_pct_limit":     ruleMap["query_queue_mem_used_pct_limit"],
				"query_queue_cpu_used_pct_limit":     ruleMap["query_queue_cpu_used_pct_limit"],
			},
		},
	}

	resp, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating query queue rule of HTAP StarRocks instance(%s): %s", d.Id(), err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	status := utils.PathSearch("status", respBody, "").(string)
	if status != "success" {
		msg := utils.PathSearch("msg", respBody, "").(string)
		return fmt.Errorf("error updating query queue rule of HTAP StarRocks instance(%s): status=%s, msg=%s", d.Id(), status, msg)
	}

	return nil
}

func updateHtapQueryQueueRuleSwitch(client *golangsdk.ServiceClient, d *schema.ResourceData, enableSwitch string) error {
	var httpUrl = "v3/{project_id}/instances/{instance_id}/htap/query-queue/switch"

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"enable_query_queue_select": enableSwitch,
		},
	}

	resp, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating query queue switch of HTAP StarRocks instance(%s): %s", d.Id(), err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	status := utils.PathSearch("status", respBody, "").(string)
	if status != "success" {
		msg := utils.PathSearch("msg", respBody, "").(string)
		return fmt.Errorf("error updating query queue switch of HTAP StarRocks instance(%s): status=%s, msg=%s", d.Id(), status, msg)
	}

	return nil
}

func getHtapQueryQueueRule(client *golangsdk.ServiceClient, htapInstanceId string) (interface{}, error) {
	var httpUrl = "v3/{project_id}/instances/{instance_id}/query-queue/rules"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", htapInstanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error querying query queue rule of HTAP StarRocks instance(%s): %s", htapInstanceId, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("query_queue_rule", respBody, nil), nil
}

func modifyStarrocksParameters(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	nodeType string, paramValues map[string]interface{}) (bool, error) {
	var (
		taurusdbInstanceId  = d.Get("instance_id").(string)
		starrocksInstanceId = d.Id()
		timeout             = d.Timeout(schema.TimeoutCreate)
	)
	modifyPath := client.Endpoint + "v3/{project_id}/instances/{instance_id}/starrocks/configurations"
	modifyPath = strings.ReplaceAll(modifyPath, "{project_id}", client.ProjectID)
	modifyPath = strings.ReplaceAll(modifyPath, "{instance_id}", starrocksInstanceId)

	modifyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"node_type":        nodeType,
			"parameter_values": utils.ExpandToStringMap(paramValues),
		},
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", modifyPath, &modifyOpt)
		shouldRetry, err := handleMultiOperationsError(err)
		return res, shouldRetry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     htapInstanceStateRefreshFunc(client, taurusdbInstanceId, starrocksInstanceId),
		WaitTarget:   []string{"normal"},
		Timeout:      timeout,
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return false, fmt.Errorf("error modifying TaurusDB Htap StarRocks instance(%s) parameters: %s",
			starrocksInstanceId, err)
	}
	modifyRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return false, err
	}

	jobId := utils.PathSearch("job_id", modifyRespBody, "").(string)
	if jobId == "" {
		return false, fmt.Errorf("error modifying TaurusDB Htap StarRocks instance(%s) parameters, job_id is not found in the response",
			starrocksInstanceId)
	}

	// Wait for the modify job to complete
	err = checkGaussDBMySQLJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return false, fmt.Errorf("error waiting for modifying TaurusDB Htap StarRocks parameters job (%s) to complete: %s", jobId, err)
	}

	restartRequired := utils.PathSearch("restart_required", modifyRespBody, false).(bool)

	return restartRequired, nil
}

func waitHtapInstanceStatusNormal(ctx context.Context, client *golangsdk.ServiceClient,
	taurusdbInstanceId, htapInstanceId string, timeout time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending:      []string{"pending"},
		Target:       []string{"normal"},
		Refresh:      htapInstanceStateRefreshFunc(client, taurusdbInstanceId, htapInstanceId),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for HTAP StarRocks instance(%s) to be normal: %s", htapInstanceId, err)
	}
	return nil
}

func buildStarRocksInstanceUsersSyncBodyParams(enableUsersSync bool) map[string]interface{} {
	action := "stopSyncTaurusAccount"
	if enableUsersSync {
		action = "startSyncTaurusAccount"
	}
	return map[string]interface{}{
		"action": action,
	}
}

func resourceTaurusDBHtapStarrocksInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                = meta.(*config.Config)
		region             = cfg.GetRegion(d)
		taurusdbInstanceId = d.Get("instance_id").(string)
		htapInstanceId     = d.Id()
		product            = "gaussdb"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	htapInstanceDetail, err := GetHtapInstanceDetail(client, taurusdbInstanceId, htapInstanceId)
	if err != nil {
		err = common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.05000044")
		return common.CheckDeletedDiag(d, err, fmt.Sprintf(
			"error retrieving HTAP StarRocks instance(%s) for TaurusDB instance(%s)",
			htapInstanceId, taurusdbInstanceId))
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instance_id", d.Get("instance_id")),
		d.Set("engine", flattenStarrocksInstanceEngineWithState(d, htapInstanceDetail)),
		d.Set("fe_volume", flattenStarrocksInstanceFeVolumeWithState(d, htapInstanceDetail)),
		d.Set("be_volume", flattenStarrocksInstanceBeVolumeWithState(d, htapInstanceDetail)),
		d.Set("ha", flattenStarrocksInstanceHa(htapInstanceDetail)),
		d.Set("az_mode", utils.PathSearch("az_mode", htapInstanceDetail, nil)),
		d.Set("tags_info", flattenStarrocksInstanceTagsInfoWithState(d, htapInstanceDetail)),
		d.Set("time_zone", utils.PathSearch("time_zone", htapInstanceDetail, nil)),
		d.Set("name", utils.PathSearch("name", htapInstanceDetail, nil)),
		d.Set("project_id", utils.PathSearch("project_id", htapInstanceDetail, nil)),
		d.Set("status", utils.PathSearch("status", htapInstanceDetail, nil)),
		d.Set("actions", flattenStarrocksInstanceActions(htapInstanceDetail)),
		d.Set("data_vip", utils.PathSearch("data_vip", htapInstanceDetail, nil)),
		d.Set("az_mode", utils.PathSearch("az_mode", htapInstanceDetail, nil)),
		d.Set("fe_node_volume_code", utils.PathSearch("fe_node_volume_code", htapInstanceDetail, nil)),
		d.Set("be_node_volume_code", utils.PathSearch("be_node_volume_code", htapInstanceDetail, nil)),
		d.Set("fe_node_volume_size", utils.PathSearch("fe_node_volume_size", htapInstanceDetail, nil)),
		d.Set("be_node_volume_size", utils.PathSearch("be_node_volume_size", htapInstanceDetail, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", htapInstanceDetail, nil)),
		d.Set("ssl_option", utils.PathSearch("ssl_option", htapInstanceDetail, nil)),
		d.Set("groups", flattenStarrocksInstanceGroups(htapInstanceDetail)),
		d.Set("ops_window", flattenStarrocksInstanceOpsWindow(htapInstanceDetail)),
		d.Set("users_sync_switch_on", utils.PathSearch("users_sync_switch_on", htapInstanceDetail, nil)),
		d.Set("can_enable_public_access", utils.PathSearch("can_enable_public_access", htapInstanceDetail, nil)),
		d.Set("public_ip", utils.PathSearch("public_ip", htapInstanceDetail, nil)),
		d.Set("current_backup_node_id", utils.PathSearch("current_backup_node_id", htapInstanceDetail, nil)),
		d.Set("cluster_mode", utils.PathSearch("cluster_mode", htapInstanceDetail, nil)),
		d.Set("is_frozen", utils.PathSearch("is_frozen", htapInstanceDetail, nil)),
		d.Set("frozen_time", utils.PathSearch("frozen_time", htapInstanceDetail, nil)),
		d.Set("db_user", utils.PathSearch("db_user", htapInstanceDetail, nil)),
		d.Set("bak_period", utils.PathSearch("bak_period", htapInstanceDetail, nil)),
		d.Set("bak_keep_day", utils.PathSearch("bak_keep_day", htapInstanceDetail, nil)),
		d.Set("bak_expected_start_time", utils.PathSearch("bak_expected_start_time", htapInstanceDetail, nil)),
		d.Set("data_store_type", utils.PathSearch("data_store_type", htapInstanceDetail, nil)),
		d.Set("data_store_version", utils.PathSearch("data_store_version", htapInstanceDetail, nil)),
		d.Set("data_store_version_id", utils.PathSearch("data_store_version_id", htapInstanceDetail, nil)),
		d.Set("create_at", utils.PathSearch("create_at", htapInstanceDetail, nil)),
		d.Set("update_at", utils.PathSearch("update_at", htapInstanceDetail, nil)),
		d.Set("delete_at", utils.PathSearch("delete_at", htapInstanceDetail, nil)),
		d.Set("db_port", utils.PathSearch("db_port", htapInstanceDetail, nil)),
		d.Set("param_group", utils.PathSearch("param_group", htapInstanceDetail, nil)),
		d.Set("create_fail_error_code", utils.PathSearch("create_fail_error_code", htapInstanceDetail, nil)),
		d.Set("backup_used_space", utils.PathSearch("backup_used_space", htapInstanceDetail, nil)),
		d.Set("port_info", flattenStarrocksInstancePortInfo(htapInstanceDetail)),
		d.Set("support_data_replication", utils.PathSearch("support_data_replication", htapInstanceDetail, nil)),
		d.Set("dedicated_resource_id", utils.PathSearch("dedicated_resource_id", htapInstanceDetail, nil)),
		d.Set("new_version_available", utils.PathSearch("new_version_available", htapInstanceDetail, nil)),
	)

	// Set enable_users_sync based on users_sync_switch_on status
	usersSyncSwitchOn := utils.PathSearch("users_sync_switch_on", htapInstanceDetail, false).(bool)
	enableUsersSync := "false"
	if usersSyncSwitchOn {
		enableUsersSync = "true"
	}
	mErr = multierror.Append(mErr, d.Set("enable_users_sync", enableUsersSync))

	// Get charging_mode from instance detail
	payModel := utils.PathSearch("pay_model", htapInstanceDetail, "").(string)
	chargingMode := "postPaid"
	if payModel == "1" {
		chargingMode = "prePaid"
	}
	mErr = multierror.Append(mErr,
		d.Set("charging_mode", chargingMode),
		d.Set("period", d.Get("period")),
		d.Set("period_unit", d.Get("period_unit")),
		d.Set("auto_renew", d.Get("auto_renew")),
	)

	// Get network, flavor, nodes, volume from groups/nodes
	groups := utils.PathSearch("groups", htapInstanceDetail, make([]interface{}, 0)).([]interface{})
	if len(groups) > 0 {
		// Get network info from first node of first group
		firstGroup := groups[0]
		nodes := utils.PathSearch("nodes", firstGroup, make([]interface{}, 0)).([]interface{})
		if len(nodes) > 0 {
			firstNode := nodes[0]
			mErr = multierror.Append(mErr,
				d.Set("vpc_id", utils.PathSearch("vpc_id", firstNode, nil)),
				d.Set("subnet_id", utils.PathSearch("subnet_id", firstNode, nil)),
				d.Set("security_group_id", utils.PathSearch("sg_id", firstNode, nil)),
				d.Set("az_code", utils.PathSearch("az_code", firstNode, nil)),
			)
		}

		for _, group := range groups {
			groupTypeName := utils.PathSearch("group_node_type", group, "").(string)
			if groupTypeName == "fe" {
				feNodes := utils.PathSearch("nodes", group, make([]interface{}, 0)).([]interface{})
				if feNodesCount := len(feNodes); feNodesCount > 0 {
					feFlavorId := utils.PathSearch("flavor_id", feNodes[0], "").(string)
					mErr = multierror.Append(mErr,
						d.Set("fe_count", feNodesCount),
						d.Set("fe_flavor_id", feFlavorId),
					)
				}
			} else if groupTypeName == "be" {
				beNodes := utils.PathSearch("nodes", group, make([]interface{}, 0)).([]interface{})
				if beNodesCount := len(beNodes); beNodesCount > 0 {
					beFlavorId := utils.PathSearch("flavor_id", beNodes[0], "").(string)
					mErr = multierror.Append(mErr,
						d.Set("be_count", beNodesCount),
						d.Set("be_flavor_id", beFlavorId),
					)
				}
			}
		}
	}
	// Set configurations and parameters of be nodes
	beConfiguration, beParameterValues, err := getStarrocksParameters(client, htapInstanceId, "be")
	if err != nil {
		return diag.FromErr(err)
	}
	mErr = multierror.Append(mErr,
		d.Set("be_configurations", flattenStarrocksParametersConfigurationsBody(beConfiguration)),
		d.Set("be_parameters", flattenStarrocksParametersParameterValuesBody(beParameterValues)),
	)

	// Set configurations and parameters of fe nodes
	feConfiguration, feParameterValues, err := getStarrocksParameters(client, htapInstanceId, "fe")
	if err != nil {
		return diag.FromErr(err)
	}
	mErr = multierror.Append(mErr,
		d.Set("fe_configurations", flattenStarrocksParametersConfigurationsBody(feConfiguration)),
		d.Set("fe_parameters", flattenStarrocksParametersParameterValuesBody(feParameterValues)),
	)

	// set slow log switch
	mErr = multierror.Append(mErr, setOpenSlowLogSwitch(d, client, htapInstanceId))

	// Set query queue rule
	mErr = multierror.Append(mErr, setQueryQueueRule(d, client, htapInstanceId))

	return diag.FromErr(mErr.ErrorOrNil())
}

func setOpenSlowLogSwitch(d *schema.ResourceData, client *golangsdk.ServiceClient, htapInstanceId string) error {
	slowLogSwitch, err := getStarRocksSlowLogSwitch(client, htapInstanceId)
	if err != nil {
		log.Printf("[WARN] query HTAP instance %s open slow log switch failed: %s", htapInstanceId, err)
		return nil
	}
	return d.Set("open_slow_log_switch", slowLogSwitch)
}

func setQueryQueueRule(d *schema.ResourceData, client *golangsdk.ServiceClient, htapInstanceId string) error {
	// Get query queue rule
	queryQueueRule, err := getHtapQueryQueueRule(client, htapInstanceId)

	if err != nil {
		log.Printf("[WARN] query HTAP instance %s open slow log switch failed: %s", htapInstanceId, err)
		return nil
	}
	return d.Set("query_queue_rule", flattenQueryQueueRule(queryQueueRule))
}

func flattenQueryQueueRule(rule interface{}) []interface{} {
	if rule == nil {
		return []interface{}{}
	}
	result := map[string]interface{}{
		"query_queue_max_queued_queries":     int(utils.PathSearch("query_queue_max_queued_queries", rule, float64(0)).(float64)),
		"query_queue_pending_timeout_second": int(utils.PathSearch("query_queue_pending_timeout_second", rule, float64(0)).(float64)),
		"query_queue_concurrency_limit":      int(utils.PathSearch("query_queue_concurrency_limit", rule, float64(0)).(float64)),
		"query_queue_mem_used_pct_limit":     int(utils.PathSearch("query_queue_mem_used_pct_limit", rule, float64(0)).(float64)),
		"query_queue_cpu_used_pct_limit":     int(utils.PathSearch("query_queue_cpu_used_pct_limit", rule, float64(0)).(float64)),
		"enable_query_queue_select":          utils.PathSearch("enable_query_queue_select", rule, "false").(string),
	}
	return []interface{}{result}
}

func flattenStarrocksInstanceActions(instance interface{}) []interface{} {
	actions := utils.PathSearch("actions", instance, make([]interface{}, 0)).([]interface{})
	res := make([]interface{}, 0, len(actions))
	for _, a := range actions {
		res = append(res, map[string]interface{}{
			"id":         utils.PathSearch("id", a, nil),
			"action":     utils.PathSearch("action", a, nil),
			"object_id":  utils.PathSearch("object_id", a, nil),
			"type":       utils.PathSearch("type", a, nil),
			"job_id":     utils.PathSearch("job_id", a, nil),
			"status":     utils.PathSearch("status", a, nil),
			"created_at": utils.PathSearch("created_at", a, nil),
			"updated_at": utils.PathSearch("updated_at", a, nil),
		})
	}
	return res
}

func flattenStarrocksInstanceGroups(instance interface{}) []interface{} {
	groups := utils.PathSearch("groups", instance, make([]interface{}, 0)).([]interface{})
	res := make([]interface{}, 0, len(groups))
	for _, g := range groups {
		res = append(res, map[string]interface{}{
			"id":              utils.PathSearch("id", g, nil),
			"name":            utils.PathSearch("name", g, nil),
			"group_type_name": utils.PathSearch("group_type_name", g, nil),
			"group_node_type": utils.PathSearch("group_node_type", g, nil),
			"nodes":           flattenStarrocksInstanceNodes(g),
		})
	}
	return res
}

func flattenStarrocksInstanceNodes(group interface{}) []interface{} {
	nodes := utils.PathSearch("nodes", group, make([]interface{}, 0)).([]interface{})
	res := make([]interface{}, 0, len(nodes))
	for _, n := range nodes {
		res = append(res, map[string]interface{}{
			"id":              utils.PathSearch("id", n, nil),
			"name":            utils.PathSearch("name", n, nil),
			"type":            utils.PathSearch("type", n, nil),
			"status":          utils.PathSearch("status", n, nil),
			"period":          utils.PathSearch("period", n, nil),
			"volume":          flattenStarrocksInstanceNodeVolume(n),
			"cpu":             utils.PathSearch("cpu", n, nil),
			"mem":             utils.PathSearch("mem", n, nil),
			"datastore":       flattenStarrocksInstanceNodeDatastore(n),
			"actions":         flattenStarrocksInstanceNodeActions(n),
			"priority":        utils.PathSearch("priority", n, nil),
			"frozen_flag":     utils.PathSearch("frozen_flag", n, nil),
			"db_port":         utils.PathSearch("db_port", n, nil),
			"pay_model":       utils.PathSearch("pay_model", n, nil),
			"order_id":        utils.PathSearch("order_id", n, nil),
			"traffic_ip":      utils.PathSearch("traffic_ip", n, nil),
			"traffic_ipv6":    utils.PathSearch("traffic_ipv6", n, nil),
			"az_code":         utils.PathSearch("az_code", n, nil),
			"az_description":  utils.PathSearch("az_description", n, nil),
			"az_type":         utils.PathSearch("az_type", n, nil),
			"region_code":     utils.PathSearch("region_code", n, nil),
			"create_at":       utils.PathSearch("create_at", n, nil),
			"update_at":       utils.PathSearch("update_at", n, nil),
			"flavor_id":       utils.PathSearch("flavor_id", n, nil),
			"flavor_ref":      utils.PathSearch("flavor_ref", n, nil),
			"iass_flavor_ref": utils.PathSearch("iass_flavor_ref", n, nil),
			"max_connections": utils.PathSearch("max_connections", n, nil),
			"vpc_id":          utils.PathSearch("vpc_id", n, nil),
			"subnet_id":       utils.PathSearch("subnet_id", n, nil),
			"need_restart":    utils.PathSearch("need_restart", n, nil),
			"sg_id":           utils.PathSearch("sg_id", n, nil),
			"param_group":     flattenStarrocksInstanceNodeParamGroup(n),
		})
	}
	return res
}

func flattenStarrocksInstanceNodeVolume(node interface{}) []interface{} {
	volume := utils.PathSearch("volume", node, nil)
	if volume == nil {
		return nil
	}
	res := make([]interface{}, 0)
	res = append(res, map[string]interface{}{
		"type": utils.PathSearch("type", volume, nil),
		"size": utils.PathSearch("size", volume, nil),
	})
	return res
}

func flattenStarrocksInstanceNodeDatastore(node interface{}) []interface{} {
	datastore := utils.PathSearch("datastore", node, nil)
	if datastore == nil {
		return nil
	}
	res := make([]interface{}, 0)
	res = append(res, map[string]interface{}{
		"id":      utils.PathSearch("id", datastore, nil),
		"type":    utils.PathSearch("type", datastore, nil),
		"version": utils.PathSearch("version", datastore, nil),
	})
	return res
}

func flattenStarrocksInstanceNodeActions(node interface{}) []interface{} {
	actions := utils.PathSearch("actions", node, make([]interface{}, 0)).([]interface{})
	res := make([]interface{}, 0)
	for _, a := range actions {
		res = append(res, map[string]interface{}{
			"id":         utils.PathSearch("id", a, nil),
			"action":     utils.PathSearch("action", a, nil),
			"object_id":  utils.PathSearch("object_id", a, nil),
			"type":       utils.PathSearch("type", a, nil),
			"job_id":     utils.PathSearch("job_id", a, nil),
			"status":     utils.PathSearch("status", a, nil),
			"created_at": utils.PathSearch("created_at", a, nil),
			"updated_at": utils.PathSearch("updated_at", a, nil),
		})
	}
	return res
}

func flattenStarrocksInstanceNodeParamGroup(node interface{}) []interface{} {
	paramGroup := utils.PathSearch("param_group", node, nil)
	if paramGroup == nil {
		return nil
	}
	res := make([]interface{}, 0)
	res = append(res, map[string]interface{}{
		"id":   utils.PathSearch("id", paramGroup, nil),
		"name": utils.PathSearch("name", paramGroup, nil),
	})
	return res
}

func flattenStarrocksInstanceOpsWindow(instance interface{}) []interface{} {
	opsWindow := utils.PathSearch("ops_window", instance, nil)
	res := make([]interface{}, 0)
	if opsWindow != nil {
		res = append(res, map[string]interface{}{
			"period":     utils.PathSearch("period", opsWindow, nil),
			"start_time": utils.PathSearch("start_time", opsWindow, nil),
			"end_time":   utils.PathSearch("end_time", opsWindow, nil),
		})
	}
	return res
}

func flattenStarrocksInstancePortInfo(instance interface{}) []interface{} {
	portInfo := utils.PathSearch("port_info", instance, nil)
	res := make([]interface{}, 0)
	if portInfo != nil {
		res = append(res, map[string]interface{}{
			"mysql_port": utils.PathSearch("mysql_port", portInfo, nil),
		})
	}
	return res
}

func flattenStarrocksInstanceHa(instance interface{}) []interface{} {
	haModel := utils.PathSearch("cluster_mode", instance, "").(string)
	res := make([]interface{}, 0)
	if haModel != "" {
		res = append(res, map[string]interface{}{
			"mode": haModel,
		})
	}
	return res
}

func flattenStarrocksInstanceEngineWithState(d *schema.ResourceData, instance interface{}) []interface{} {
	engineType := utils.PathSearch("data_store_type", instance, "").(string)
	engineVersion := utils.PathSearch("data_store_version", instance, "").(string)
	userEngineVersion := d.Get("engine.0.version").(string)
	if userEngineVersion != "" && strings.HasPrefix(engineVersion, userEngineVersion) {
		engineVersion = userEngineVersion
	}
	res := make([]interface{}, 0, 1)
	res = append(res, map[string]interface{}{
		"type":    engineType,
		"version": engineVersion,
	})
	return res
}

func flattenStarrocksInstanceVolumeTypeFromGroup(instance interface{}, groupNodeType string) string {
	groups := utils.PathSearch("groups", instance, make([]interface{}, 0)).([]interface{})
	for _, g := range groups {
		groupTypeName := utils.PathSearch("group_node_type", g, "").(string)
		if groupTypeName == groupNodeType {
			volume := utils.PathSearch("nodes[0].volume", g, nil)
			if volume != nil {
				return utils.PathSearch("type", volume, "").(string)
			}
		}
	}
	return ""
}

// flattenStarrocksInstanceFeVolumeWithState preserves user input from state when available,
// and only falls back to API response during import (when state is empty).
// This avoids state drift caused by API returning calculated values that differ from input.
func flattenStarrocksInstanceFeVolumeWithState(d *schema.ResourceData, instance interface{}) []interface{} {
	if v := d.Get("fe_volume").([]interface{}); len(v) > 0 {
		return v
	}
	feVolumeSizeRaw := utils.PathSearch("fe_node_volume_size", instance, "").(string)
	feVolumeSize, _ := strconv.Atoi(feVolumeSizeRaw)
	res := make([]interface{}, 0)
	return append(res, map[string]interface{}{
		"io_type":        flattenStarrocksInstanceVolumeTypeFromGroup(instance, "fe"),
		"capacity_in_gb": feVolumeSize / 1000000000,
	})
}

// flattenStarrocksInstanceBeVolumeWithState preserves user input from state when available,
// and only falls back to API response during import (when state is empty).
// This avoids state drift caused by API returning calculated values that differ from input.
func flattenStarrocksInstanceBeVolumeWithState(d *schema.ResourceData, instance interface{}) []interface{} {
	if v := d.Get("be_volume").([]interface{}); len(v) > 0 {
		return v
	}
	beVolumeSizeRaw := utils.PathSearch("be_node_volume_size", instance, "").(string)
	beVolumeSize, _ := strconv.Atoi(beVolumeSizeRaw)
	res := make([]interface{}, 0)
	return append(res, map[string]interface{}{
		"io_type":        flattenStarrocksInstanceVolumeTypeFromGroup(instance, "be"),
		"capacity_in_gb": beVolumeSize / 1000000000,
	})
}

// flattenStarrocksInstanceTagsInfoWithState preserves user input from state when available,
// and only falls back to API response during import (when state is empty).
func flattenStarrocksInstanceTagsInfoWithState(d *schema.ResourceData, instance interface{}) []interface{} {
	if v := d.Get("tags_info").([]interface{}); len(v) > 0 {
		return v
	}
	return flattenStarrocksInstanceTagsInfo(instance)
}

func flattenStarrocksInstanceTagsInfo(instance interface{}) []interface{} {
	tagsInfo := utils.PathSearch("tags_info", instance, nil)
	if tagsInfo == nil {
		return nil
	}

	tags := utils.PathSearch("tags", tagsInfo, make([]interface{}, 0)).([]interface{})
	sysTags := utils.PathSearch("sys_tags", tagsInfo, make([]interface{}, 0)).([]interface{})

	flattenedTags := make([]interface{}, 0, len(tags))
	for _, t := range tags {
		flattenedTags = append(flattenedTags, map[string]interface{}{
			"key":   utils.PathSearch("key", t, nil),
			"value": utils.PathSearch("value", t, nil),
		})
	}

	flattenedSysTags := make([]interface{}, 0, len(sysTags))
	for _, t := range sysTags {
		flattenedSysTags = append(flattenedSysTags, map[string]interface{}{
			"key":   utils.PathSearch("key", t, nil),
			"value": utils.PathSearch("value", t, nil),
		})
	}

	res := make([]interface{}, 0, 1)
	res = append(res, map[string]interface{}{
		"tags":     flattenedTags,
		"sys_tags": flattenedSysTags,
	})
	return res
}

// handleUpdateAutoRenew handles auto_renew update
func handleUpdateAutoRenew(bssClient *golangsdk.ServiceClient, d *schema.ResourceData) error {
	if err := cbc.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), d.Id()); err != nil {
		return fmt.Errorf("error updating the auto-renew of the HTAP StarRocks instance (%s): %s", d.Id(), err)
	}
	return nil
}

// handleUpdatePassword handles password update
func handleUpdatePassword(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	return updateStarRocksPassword(d, client)
}

// handleUpdateSecurityGroup handles security group update
func handleUpdateSecurityGroup(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	return updateStarRocksInstanceSecurityGroup(ctx, d, client)
}

// handleUpdateUsersSync handles users sync switch update
func handleUpdateUsersSync(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	taurusDBInstanceId, starrocksInstanceId string) error {
	enableUsersSyncRaw := d.Get("enable_users_sync").(string)
	return updateStarRocksInstanceUsersSyncSwitch(ctx, client, taurusDBInstanceId, starrocksInstanceId,
		enableUsersSyncRaw == "true", d.Timeout(schema.TimeoutUpdate))
}

// handleUpdateFlavor handles flavor update
func handleUpdateFlavor(ctx context.Context, d *schema.ResourceData, client, bssClient *golangsdk.ServiceClient) error {
	// For prePaid instance, API only allows changing fe_flavor_id or be_flavor_id at the same time
	if d.Get("charging_mode").(string) == "prePaid" {
		if d.HasChange("fe_flavor_id") {
			if err := resizeStarRocksFlavor(ctx, d, client, bssClient, "fe_flavor_id"); err != nil {
				return err
			}
		}
		if d.HasChange("be_flavor_id") {
			if err := resizeStarRocksFlavor(ctx, d, client, bssClient, "be_flavor_id"); err != nil {
				return err
			}
		}
	} else {
		// For postPaid instance, API allows changing both fe_flavor_id and be_flavor_id at the same time
		if err := resizeStarRocksFlavor(ctx, d, client, bssClient, "all"); err != nil {
			return err
		}
	}
	return nil
}

// handleUpdateParameters handles BE and FE parameters update
func handleUpdateParameters(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) (context.Context, error) {
	restartRequired := false

	beParamValues := d.Get("be_parameter_values").(map[string]interface{})
	if len(beParamValues) > 0 {
		var err error
		restartRequired, err = modifyStarrocksParameters(ctx, client, d, "be", beParamValues)
		if err != nil {
			return ctx, err
		}
	}

	feParamValues := d.Get("fe_parameter_values").(map[string]interface{})
	if len(feParamValues) > 0 {
		var err error
		restartRequired, err = modifyStarrocksParameters(ctx, client, d, "fe", feParamValues)
		if err != nil {
			return ctx, err
		}
	}

	if restartRequired {
		// Sending needRestartByParametersChanged to Read to warn users the instance needs a reboot.
		ctx = context.WithValue(ctx, ctxType("needRestartByParametersChanged"), "true")
	}

	return ctx, nil
}

func resourceTaurusDBHtapStarrocksInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "gaussdb"
	)

	taurusDBInstanceId := d.Get("instance_id").(string)
	starrocksInstanceId := d.Id()
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	bssClient, err := cfg.BssV2Client(region)
	if err != nil {
		return diag.Errorf("error creating BSS V2 client: %s", err)
	}

	// Handle auto_renew update
	if d.HasChange("auto_renew") {
		if err := handleUpdateAutoRenew(bssClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	// Update password
	if d.HasChange("db_root_pwd") {
		if err := handleUpdatePassword(d, client); err != nil {
			return diag.FromErr(err)
		}
	}

	// Update security group
	if d.HasChange("security_group_id") {
		if err := handleUpdateSecurityGroup(ctx, d, client); err != nil {
			return diag.FromErr(err)
		}
	}

	// Update users sync switch
	if d.HasChange("enable_users_sync") {
		if err := handleUpdateUsersSync(ctx, d, client, taurusDBInstanceId, starrocksInstanceId); err != nil {
			return diag.FromErr(err)
		}
	}

	// Update flavor
	if d.HasChanges("fe_flavor_id", "be_flavor_id") {
		if err := handleUpdateFlavor(ctx, d, client, bssClient); err != nil {
			return diag.FromErr(err)
		}
	}

	// Update parameters
	if d.HasChanges("be_parameter_values", "fe_parameter_values") {
		var err error
		ctx, err = handleUpdateParameters(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Update slow log switch
	if d.HasChange("open_slow_log_switch") {
		if err := updateStarRocksSlowLogSwitch(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	// Update query queue rule
	if d.HasChange("query_queue_rule") {
		// Update query queue switch if enable_query_queue_select changed
		if d.HasChange("query_queue_rule.0.enable_query_queue_select") {
			enableSwitch := d.Get("query_queue_rule.0.enable_query_queue_select").(string)
			if err := updateHtapQueryQueueRuleSwitch(client, d, enableSwitch); err != nil {
				return diag.FromErr(err)
			}
		}
		if err := updateHtapQueryQueueRule(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceTaurusDBHtapStarrocksInstanceRead(ctx, d, meta)
}

func resizeStarRocksFlavor(ctx context.Context, d *schema.ResourceData, client, bssClient *golangsdk.ServiceClient, flavorField string) error {
	var (
		taurusDBInstanceId  = d.Get("instance_id").(string)
		starrocksInstanceId = d.Id()
		httpUrl             = "v3/{project_id}/instances/{instance_id}/starrocks/resize-flavor"
	)

	resizePath := client.Endpoint + httpUrl
	resizePath = strings.ReplaceAll(resizePath, "{project_id}", client.ProjectID)
	resizePath = strings.ReplaceAll(resizePath, "{instance_id}", d.Id())

	resizeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildInstanceResizeFlavorParams(d, flavorField),
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", resizePath, &resizeOpt)
		shouldRetry, err := handleMultiOperationsError(err)
		return res, shouldRetry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     htapInstanceStateRefreshFunc(client, taurusDBInstanceId, starrocksInstanceId),
		WaitTarget:   []string{"normal"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error resizing flavor of HTAP StarRocks instance(%s): %s", starrocksInstanceId, err)
	}

	respBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return err
	}

	if d.Get("charging_mode").(string) == "prePaid" {
		orderId := utils.PathSearch("orderId", respBody, "").(string)
		if orderId == "" {
			return errors.New("error resizing TaurusDB HTAP StarRocks instance flavor: orderId is empty")
		}
		if err = cbc.PaySubscriptionOrder(bssClient, orderId); err != nil {
			return fmt.Errorf("error paying for resizing order (%s): %s", orderId, err)
		}
		if err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return fmt.Errorf("error waiting for resizing order (%s) to complete: %s", orderId, err)
		}
	} else {
		jobId := utils.PathSearch("job_id", respBody, "").(string)
		if jobId == "" {
			return errors.New("error resizing TaurusDB HTAP StarRocks instance flavor: job_id is empty")
		}
		// No API to query job status, so wait for HTAP StarRocks instance to be normal and flavor to be updated
		if err = checkHtapInstanceResizeFlavorJobFinish(ctx, client, d); err != nil {
			return fmt.Errorf("error waiting for HTAP StarRocks instance(%s) resize flavor to be completed: %s", d.Id(), err)
		}
	}

	return nil
}

func buildInstanceResizeFlavorParams(d *schema.ResourceData, flavorField string) map[string]interface{} {
	if flavorField == "all" {
		return map[string]interface{}{
			"fe_flavor_id": d.Get("fe_flavor_id").(string),
			"be_flavor_id": d.Get("be_flavor_id").(string),
		}
	}
	return map[string]interface{}{
		flavorField: d.Get(flavorField).(string),
	}
}

func checkHtapInstanceResizeFlavorJobFinish(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &retry.StateChangeConf{
		Pending:      []string{"Running"},
		Target:       []string{"Completed"},
		Refresh:      htapInstanceFlavorRefreshFunc(client, d),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for TaurusDB HTAP StarRocks instance flavor resize to be completed: %s", err)
	}
	return nil
}

func checkStarRocksNodeStatus(nodes []interface{}) error {
	for _, node := range nodes {
		nodeStatus := utils.PathSearch("status", node, "").(string)
		if nodeStatus == "abnormal" || nodeStatus == "createfail" {
			return fmt.Errorf("StarRocks node is in %s state", nodeStatus)
		}
	}
	return nil
}

func countFlavorMatchedNodes(nodes []interface{}, targetFlavorId string) int {
	matched := 0
	for _, node := range nodes {
		nodeFlavorId := utils.PathSearch("flavor_id", node, "").(string)
		if nodeFlavorId == targetFlavorId {
			matched++
		}
	}
	return matched
}

func htapInstanceFlavorRefreshFunc(client *golangsdk.ServiceClient, d *schema.ResourceData) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			taurusdbInstanceId = d.Get("instance_id").(string)
			htapInstanceId     = d.Id()
			beFlavorId         = d.Get("be_flavor_id").(string)
			feFlavorId         = d.Get("fe_flavor_id").(string)
		)

		htapInstanceDetail, err := GetHtapInstanceDetail(client, taurusdbInstanceId, htapInstanceId)
		if err != nil {
			return nil, "", err
		}

		htapInstanceStatus := utils.PathSearch("status", htapInstanceDetail, "").(string)
		if htapInstanceStatus == "abnormal" || htapInstanceStatus == "createfail" {
			return nil, "", fmt.Errorf("TaurusDB HTAP instance(%s) is in %s state",
				htapInstanceId, htapInstanceStatus)
		}

		htapActions := utils.PathSearch("actions", htapInstanceDetail, make([]interface{}, 0)).([]interface{})
		if len(htapActions) > 0 {
			htapAction := utils.PathSearch("actions[0].action", htapInstanceDetail, "").(string)
			if htapAction == "resize_flavor" {
				return htapInstanceDetail, "Running", nil
			}
		}

		groups := utils.PathSearch("groups", htapInstanceDetail, make([]interface{}, 0)).([]interface{})
		beNodeCount := 0
		feNodeCount := 0
		beFlavorMatchedNum := 0
		feFlavorMatchedNum := 0

		for _, group := range groups {
			groupNodeType := utils.PathSearch("group_node_type", group, "").(string)
			nodes := utils.PathSearch("nodes", group, make([]interface{}, 0)).([]interface{})
			if len(nodes) == 0 {
				continue
			}

			if err := checkStarRocksNodeStatus(nodes); err != nil {
				return nil, "", err
			}

			if groupNodeType == "be" {
				beNodeCount = len(nodes)
				beFlavorMatchedNum += countFlavorMatchedNodes(nodes, beFlavorId)
			} else if groupNodeType == "fe" {
				feNodeCount = len(nodes)
				feFlavorMatchedNum += countFlavorMatchedNodes(nodes, feFlavorId)
			}
		}

		if feNodeCount > 0 && beNodeCount > 0 &&
			feFlavorMatchedNum == feNodeCount && beFlavorMatchedNum == beNodeCount {
			return htapInstanceDetail, "Completed", nil
		}

		return htapInstanceDetail, "Running", nil
	}
}

func updateStarRocksPassword(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/starrocks/users/password"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	updateOpt.JSONBody = map[string]interface{}{
		"user_name": "root",
		"password":  d.Get("db_root_pwd"),
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating TaurusDB HTAP StarRocks instance password: %s", err)
	}

	return nil
}

func updateStarRocksInstanceSecurityGroup(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		taurusDBInstanceId  = d.Get("instance_id").(string)
		starrocksInstanceId = d.Id()
		httpUrl             = "v3/{project_id}/instances/{instance_id}/starrocks/security-group"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", starrocksInstanceId)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	updateOpt.JSONBody = map[string]interface{}{
		"security_group_id": d.Get("security_group_id").(string),
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", updatePath, &updateOpt)
		shouldRetry, err := handleMultiOperationsError(err)
		return res, shouldRetry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     htapInstanceStateRefreshFunc(client, taurusDBInstanceId, starrocksInstanceId),
		WaitTarget:   []string{"normal"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating security group of HTAP StarRocks instance(%s): %s", starrocksInstanceId, err)
	}

	respBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return err
	}
	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return fmt.Errorf("error updating security group of HTAP StarRocks instance(%s). job_id is not found in the response", starrocksInstanceId)
	}

	if err := checkGaussDBMySQLJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return fmt.Errorf("error waiting for HTAP StarRocks instance security group update job (%s) to be completed: %s", jobId, err)
	}

	if err = waitHtapInstanceStatusNormal(ctx, client, taurusDBInstanceId, starrocksInstanceId, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return fmt.Errorf("error waiting for HTAP StarRocks instance(%s) to be normal after security group update: %s", starrocksInstanceId, err)
	}

	return nil
}

func resourceTaurusDBHtapStarrocksInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                = meta.(*config.Config)
		region             = cfg.GetRegion(d)
		httpUrl            = "v3/{project_id}/instances/{instance_id}/starrocks/{starrocks_instance_id}"
		product            = "gaussdb"
		taurusdbInstanceId = d.Get("instance_id").(string)
		htapInstanceId     = d.Id()
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}
	timeout := d.Timeout(schema.TimeoutDelete)

	_, err = GetHtapInstanceDetail(client, taurusdbInstanceId, htapInstanceId)
	if err != nil {
		err = common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.05000044")
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting StarRocks instance(%s) for TaurusDB instance(%s)",
			htapInstanceId, taurusdbInstanceId))
	}

	if chargingMode := d.Get("charging_mode").(string); chargingMode == "prePaid" {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		deleteResourceIds := []interface{}{htapInstanceId}

		retryFunc := func() (interface{}, bool, error) {
			err := cbc.UnsubscribePrePaidResources(bssClient, deleteResourceIds)
			if err != nil {
				shouldRetry, err := handleMultiOperationsError(err)
				return nil, shouldRetry, err
			}
			return nil, false, nil
		}
		_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     htapInstanceStateRefreshFunc(client, taurusdbInstanceId, htapInstanceId),
			WaitTarget:   []string{"normal"},
			Timeout:      timeout,
			DelayTimeout: 10 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			err = common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.05000044")
			err = common.ConvertExpected403ErrInto404Err(err, "error_code", "DBS.2000044")
			return diag.Errorf("error unsubscribing HTAP StarRocks instance(%s) for TaurusDB instance(%s): %s",
				htapInstanceId, taurusdbInstanceId, err)
		}

		err = cbc.WaitForResourcesUnsubscribed(ctx, bssClient, deleteResourceIds, timeout)
		if err != nil {
			return diag.Errorf("error waiting for HTAP StarRocks instance(%s) for TaurusDB instance(%s) "+
				"to be unsubscribed: %s ", htapInstanceId, taurusdbInstanceId, err)
		}
	} else {
		deletePath := client.Endpoint + httpUrl
		deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
		deletePath = strings.ReplaceAll(deletePath, "{instance_id}", taurusdbInstanceId)
		deletePath = strings.ReplaceAll(deletePath, "{starrocks_instance_id}", htapInstanceId)

		deleteOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		}

		retryFunc := func() (interface{}, bool, error) {
			resp, err := client.Request("DELETE", deletePath, &deleteOpt)
			shouldRetry, err := handleMultiOperationsError(err)
			return resp, shouldRetry, err
		}
		r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     htapInstanceStateRefreshFunc(client, taurusdbInstanceId, htapInstanceId),
			WaitTarget:   []string{"normal"},
			Timeout:      d.Timeout(schema.TimeoutDelete),
			DelayTimeout: 10 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			err = common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.05000044")
			err = common.ConvertExpected403ErrInto404Err(err, "error_code", "DBS.2000044")
			return common.CheckDeletedDiag(d, err,
				fmt.Sprintf("error deleting HTAP StarRocks instance(%s) for TaurusDB instance(%s)",
					htapInstanceId, taurusdbInstanceId))
		}

		respBody, err := utils.FlattenResponse(r.(*http.Response))
		if err != nil {
			return diag.FromErr(err)
		}

		workflowId := utils.PathSearch("workflow_id", respBody, "").(string)
		if workflowId != "" {
			if err = checkGaussDBMySQLJobFinish(ctx, client, workflowId, d.Timeout(schema.TimeoutDelete)); err != nil {
				return diag.Errorf("error waiting for deleting TaurusDB HTAP StarRocks instance (%s) job to complete: %s",
					d.Id(), err)
			}
		}
	}
	return nil
}

func resourceTaurusDBHtapStarrocksInstanceImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <taurusdb_instance_id>/<htap_instance_id>")
	}
	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("instance_id", parts[0])
}
