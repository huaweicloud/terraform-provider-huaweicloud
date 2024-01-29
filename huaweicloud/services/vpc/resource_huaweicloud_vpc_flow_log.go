package vpc

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/networking/v1/flowlogs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API VPC DELETE /v1/{project_id}/fl/flow_logs/{id}
// @API VPC GET /v1/{project_id}/fl/flow_logs/{id}
// @API VPC PUT /v1/{project_id}/fl/flow_logs/{id}
// @API VPC POST /v1/{project_id}/fl/flow_logs
func ResourceVpcFlowLog() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcFlowLogCreate,
		ReadContext:   resourceVpcFlowLogRead,
		UpdateContext: resourceVpcFlowLogUpdate,
		DeleteContext: resourceVpcFlowLogDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"port", "network", "vpc",
				}, true),
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_stream_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"traffic_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "all",
				ValidateFunc: validation.StringInSlice([]string{
					"all", "accept", "reject",
				}, true),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVpcFlowLogCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	vpcClient, err := cfg.NetworkingV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	createOpts := flowlogs.CreateOpts{
		Name:         d.Get("name").(string),
		ResourceType: d.Get("resource_type").(string),
		ResourceID:   d.Get("resource_id").(string),
		TrafficType:  d.Get("traffic_type").(string),
		LogGroupID:   d.Get("log_group_id").(string),
		LogTopicID:   d.Get("log_stream_id").(string),
		Description:  d.Get("description").(string),
	}

	log.Printf("[DEBUG] Create VPC flow Log Options: %#v", createOpts)
	fl, err := flowlogs.Create(vpcClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating VPC flow log: %s", err)
	}

	d.SetId(fl.ID)

	// disable the flow log function if `enabled = false`
	enabled := d.Get("enabled").(bool)
	if !enabled {
		updateOpts := flowlogs.UpdateOpts{
			AdminState: false,
		}

		_, err = flowlogs.Update(vpcClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error disable VPC flow log: %s", err)
		}
	}
	return resourceVpcFlowLogRead(ctx, d, meta)
}

func resourceVpcFlowLogRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	vpcClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	fl, err := flowlogs.Get(vpcClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving VPC flow log")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", fl.Name),
		d.Set("description", fl.Description),
		d.Set("resource_type", fl.ResourceType),
		d.Set("resource_id", fl.ResourceID),
		d.Set("traffic_type", fl.TrafficType),
		d.Set("log_group_id", fl.LogGroupID),
		d.Set("log_stream_id", fl.LogTopicID),
		d.Set("enabled", fl.AdminState),
		d.Set("status", fl.Status),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceVpcFlowLogUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	vpcClient, err := cfg.NetworkingV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	if d.HasChanges("name", "description", "enabled") {
		updateOpts := flowlogs.UpdateOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			AdminState:  d.Get("enabled").(bool),
		}

		_, err = flowlogs.Update(vpcClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating VPC flow log: %s", err)
		}
	}

	return resourceVpcFlowLogRead(ctx, d, meta)
}

func resourceVpcFlowLogDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	vpcClient, err := cfg.NetworkingV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	err = flowlogs.Delete(vpcClient, d.Id()).ExtractErr()
	if err != nil {
		// ignore ErrDefault404
		return common.CheckDeletedDiag(d, err, "error deleting VPC flow log")
	}

	return nil
}
