package geminidb

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API GaussDBforNoSQL GET /v3/{project_id}/instances
// @API GaussDBforNoSQL GET /v3/{project_id}/dedicated-resources
// @API GaussDBforNoSQL POST /v3/{project_id}/instances
// @API GaussDBforNoSQL POST /v3/{project_id}/instances/{instance_id}/tags/action
// @API GaussDBforNoSQL GET /v3/{project_id}/instances/{instance_id}/tags
// @API GaussDBforNoSQL PUT /v3/{project_id}/instances/{instance_id}/name
// @API GaussDBforNoSQL PUT /v3/{project_id}/instances/{instance_id}/password
// @API GaussDBforNoSQL PUT /v3/{project_id}/configurations/{config_id}/apply
// @API GaussDBforNoSQL GET /v3/{project_id}/configurations/{config_id}
// @API GaussDBforNoSQL GET /v3/{project_id}/instances/{instance_id}/configurations
// @API GaussDBforNoSQL POST /v3/{project_id}/instances/{instance_id}/extend-volume
// @API GaussDBforNoSQL POST /v3/{project_id}/instances/{instance_id}/enlarge-node
// @API GaussDBforNoSQL POST /v3/{project_id}/instances/{instance_id}/reduce-node
// @API GaussDBforNoSQL PUT /v3/{project_id}/instances/{instance_id}/resize
// @API GaussDBforNoSQL PUT /v3/{project_id}/instances/{instance_id}/security-group
// @API GaussDBforNoSQL PUT /v3/{project_id}/instances/{instance_id}/backups/policy
// @API GaussDBforNoSQL DELETE /v3/{project_id}/instances/{instance_id}
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
func ResourceGaussDBMongoInstanceV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDBMongoInstanceCreate,
		ReadContext:   resourceGeminiDBInstanceV3Read,
		UpdateContext: resourceGaussDBMongoInstanceUpdate,
		DeleteContext: resourceGeminiDBInstanceV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"node_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{3}),
			},
			"volume_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"configuration_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dedicated_resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"dedicated_resource_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"datastore": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"engine": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"storage_engine": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"backup_strategy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"keep_days": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"force_import": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"private_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lb_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lb_port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nodes": {
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
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"support_reduce": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},

			// charge info: charging_mode, period_unit, period, auto_renew
			// make ForceNew false here but do nothing in update method!
			"charging_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"prePaid", "postPaid",
				}, false),
			},
			"period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"period"},
				ValidateFunc: validation.StringInSlice([]string{
					"month", "year",
				}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"period_unit"},
			},
			"auto_renew": common.SchemaAutoRenewUpdatable(nil),

			"tags": common.TagsSchema(),
		},
	}
}

func resourceGaussDBMongoInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaults := defaultValues{
		Mode:      "ReplicaSet",
		dbType:    "mongodb",
		dbVersion: "4.0",
		logName:   "mongo",
	}
	return resourceGeminiDBInstanceV3Create(ctx, d, meta, defaults)
}

func resourceGaussDBMongoInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaults := defaultValues{
		Mode:      "ReplicaSet",
		dbType:    "mongodb",
		dbVersion: "4.0",
		logName:   "mongo",
	}
	return resourceGeminiDBInstanceV3Update(ctx, d, meta, defaults)
}
