package dataarts

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dayu/v1/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio POST /v1/{project_id}/instances/onekey-purchase
// @API DataArtsStudio GET /v1/{project_id}/instances

// ResourceStudioInstance is the impl of huaweicloud_dataarts_studio_instance
func ResourceStudioInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStudioInstanceCreate,
		ReadContext:   resourceStudioInstanceRead,
		UpdateContext: resourceStudioInstanceupdate,
		DeleteContext: resourceStudioInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

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
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				Required: true,
				ForceNew: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"period_unit": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"month", "year",
				}, false),
			},
			"period": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"auto_renew": common.SchemaAutoRenew(nil),
			"tags":       common.TagsForceNewSchema(),
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"charging_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"order_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expire_days": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

type PeriodType int

const (
	Monthly PeriodType = 2
	Yearly  PeriodType = 3
)

func formatPeriodType(unit string) int {
	if unit == "year" {
		return int(Yearly)
	}
	return int(Monthly)
}

func formatAutoRenew(renew string) *int {
	var auto int
	if renew == "true" {
		auto = 1
	}
	return &auto
}

func resourceStudioInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.DataArtsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio v1 client: %s", err)
	}
	bssClient, err := conf.BssV2Client(region)
	if err != nil {
		return diag.Errorf("error creating BSS v2 client: %s", err)
	}

	createOpts := instances.CreateOpts{
		Region:              region,
		Name:                d.Get("name").(string),
		SpecCode:            d.Get("version").(string),
		VpcID:               d.Get("vpc_id").(string),
		SubnetID:            d.Get("subnet_id").(string),
		SecurityGroupID:     d.Get("security_group_id").(string),
		AvailabilityZone:    d.Get("availability_zone").(string),
		EnterpriseProjectID: conf.GetEnterpriseProjectID(d),

		PeriodNum:   d.Get("period").(int),
		PeriodType:  formatPeriodType(d.Get("period_unit").(string)),
		IsAutoRenew: formatAutoRenew(d.Get("auto_renew").(string)),
	}

	if v, ok := d.GetOk("tags"); ok {
		createOpts.Tags = utils.ExpandResourceTags(v.(map[string]interface{}))
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	resp, err := instances.Create(client, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating DataArts Studio instance: %s", err)
	}

	if resp.OrderID == "" || resp.ID == "" {
		return diag.Errorf("error creating DataArts Studio instance: the instance ID is not exist")
	}

	err = common.WaitOrderComplete(ctx, bssClient, resp.OrderID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, resp.OrderID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resourceId)
	return resourceStudioInstanceRead(ctx, d, meta)
}

func resourceStudioInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.DataArtsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio v1 client: %s", err)
	}

	instanceID := d.Id()
	object, err := findStudioInstanceByID(client, instanceID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DataArts Studio instance")
	}

	log.Printf("[DEBUG] fetching DataArts Studio instance %s: %#v", instanceID, object)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", object.Name),
		d.Set("version", object.SpecCode),
		d.Set("vpc_id", object.VpcID),
		d.Set("subnet_id", object.SubnetID),
		d.Set("security_group_id", object.SecurityGroupID),
		d.Set("enterprise_project_id", object.EnterpriseProjectID),
		d.Set("availability_zone", object.AvailabilityZone),
		d.Set("charging_mode", "prePaid"),
		d.Set("order_id", object.OrderID),
		d.Set("expire_days", object.ExpireDays),
		d.Set("status", object.Status),
		d.Set("tags", d.Get("tags")),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting DataArts Studio instance fields: %s", err)
	}
	return nil
}

func resourceStudioInstanceupdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	instanceID := d.Id()

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   instanceID,
			ResourceType: "dayu-instance",
			RegionId:     region,
			ProjectId:    cfg.GetProjectID(region),
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceStudioInstanceRead(ctx, d, meta)
}

func resourceStudioInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.DataArtsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio v1 client: %s", err)
	}

	instanceID := d.Id()
	if err := common.UnsubscribePrePaidResource(d, conf, []string{instanceID}); err != nil {
		return diag.Errorf("Error unsubscribing DataArts Studio instance %s: %s", instanceID, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"deleting"},
		Target:       []string{"deleted"},
		Refresh:      refreshInstanceStatusFunc(client, instanceID),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        30 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error deleting DataArts Studio instance: %s", err)
	}

	return nil
}

func refreshInstanceStatusFunc(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := findStudioInstanceByID(client, id)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return resp, "deleted", nil
			}
			return resp, "error", err
		}
		return resp, "deleting", nil
	}
}

func findStudioInstanceByID(client *golangsdk.ServiceClient, id string) (*instances.Instance, error) {
	resp, err := instances.List(client, nil)
	if err != nil {
		return nil, err
	}

	for _, item := range resp {
		if item.ID == id {
			return &item, nil
		}
	}
	return nil, golangsdk.ErrDefault404{}
}
