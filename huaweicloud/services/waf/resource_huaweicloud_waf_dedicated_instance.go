package waf

import (
	"context"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	instances "github.com/chnsz/golangsdk/openstack/waf_hw/v1/premium_instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

const (
	// runStatusCreating the instance is creating.
	runStatusCreating = 0
	// runStatusRunning the instance has been created.
	runStatusRunning = 1
	// runStatusDeleting the instance deleting.
	runStatusDeleting = 2
	// runStatusDeleting the instance has be deleted.
	runStatusDeleted = 3
)

const (
	// defaultCount the number of instances created.
	defaultCount = 1
	// Billing mode, payPerUseMode: pay pre use mode
	payPerUseMode = 30
)

// ResourceWafDedicatedInstance the resource of managing a dedicated mode instance within HuaweiCloud.
// @API WAF DELETE /v1/{project_id}/premium-waf/instance/{instance_id}
// @API WAF GET /v1/{project_id}/premium-waf/instance/{instance_id}
// @API WAF PUT /v1/{project_id}/premium-waf/instance/{instance_id}
// @API WAF POST /v1/{project_id}/premium-waf/instance
// @API EPS POST /v1.0/enterprise-projects/{enterprise_project_id}/resources-migrate
func ResourceWafDedicatedInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDedicatedInstanceCreate,
		ReadContext:   resourceDedicatedInstanceRead,
		UpdateContext: resourceDedicatedInstanceUpdate,
		DeleteContext: resourceDedicatedInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceWAFImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"available_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"specification_code": {
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
			"security_group": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"cpu_architecture": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "x86",
				ForceNew: true,
			},
			"ecs_flavor": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"res_tenant": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"anti_affinity": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The following are the attributes
			"server_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"run_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"access_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"upgradable": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildCreateOpts(d *schema.ResourceData, region string) *instances.CreateInstanceOpts {
	sg := d.Get("security_group").([]interface{})
	groups := make([]string, 0, len(sg))
	for _, v := range sg {
		groups = append(groups, v.(string))
	}
	logp.Printf("[DEBUG] The security_group parameters are: %+v.", groups)

	createOpts := instances.CreateInstanceOpts{
		Region:        region,
		ChargeMode:    payPerUseMode,
		AvailableZone: d.Get("available_zone").(string),
		Arch:          d.Get("cpu_architecture").(string),
		NamePrefix:    d.Get("name").(string),
		Specification: d.Get("specification_code").(string),
		CpuFlavor:     d.Get("ecs_flavor").(string),
		VpcId:         d.Get("vpc_id").(string),
		SubnetId:      d.Get("subnet_id").(string),
		SecurityGroup: groups,
		Count:         defaultCount,
		PoolId:        d.Get("group_id").(string),
		ResTenant:     utils.Bool(d.Get("res_tenant").(bool)),
	}
	if d.Get("res_tenant").(bool) {
		// `anti_affinity` is valid only when `res_tenant` is true
		createOpts.AntiAffinity = utils.Bool(d.Get("anti_affinity").(bool))
	}

	return &createOpts
}

func waitForInstanceCreated(c *golangsdk.ServiceClient, id string, epsId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := instances.GetWithEpsId(c, id, epsId)
		if err != nil {
			return nil, "Error", err
		}

		switch r.RunStatus {
		case runStatusCreating:
			return r, "Creating", nil
		case runStatusRunning:
			return r, "Created", nil
		default:
			err = fmtp.Errorf("Error in create WAF dedicated instance[%s]. "+
				"Unexpected run_status: %v.", r.Id, r.RunStatus)
			return r, "Error", err
		}
	}
}

func resourceDedicatedInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.WafDedicatedV1Client(cfg.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud WAF dedicated client : %s", err)
	}

	createOpts := buildCreateOpts(d, cfg.GetRegion(d))
	epsId := cfg.GetEnterpriseProjectID(d)

	r, err := instances.CreateWithEpsId(client, *createOpts, epsId)
	if err != nil {
		return fmtp.DiagErrorf("error creating WAF dedicated : %w", err)
	}
	d.SetId(r.Instances[0].Id)

	logp.Printf("[DEBUG] Waiting for WAF dedicated instance[%s] to be created.", r.Instances[0].Id)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Creating"},
		Target:       []string{"Created"},
		Refresh:      waitForInstanceCreated(client, r.Instances[0].Id, epsId),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 15 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err == nil {
		err = updateInstanceName(client, r.Instances[0].Id, d.Get("name").(string), epsId)
	}
	if err != nil {
		logp.Printf("[DEBUG] Error while waiting to create  Waf dedicated instance. %s : %v", d.Id(), err)
		return diag.FromErr(err)
	}

	return resourceDedicatedInstanceRead(ctx, d, meta)
}

func resourceDedicatedInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.WafDedicatedV1Client(cfg.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud WAF dedicated client: %s", err)
	}

	epsId := cfg.GetEnterpriseProjectID(d)
	r, err := instances.GetWithEpsId(client, d.Id(), epsId)
	if err != nil {
		// If the dedicated instance does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving WAF dedicated instance")
	}
	logp.Printf("[DEBUG] Get a WAF dedicated instance :%#v", r)

	mErr := multierror.Append(nil,
		d.Set("region", r.Region),
		d.Set("name", r.InstanceName),
		d.Set("available_zone", r.Zone),
		d.Set("cpu_architecture", r.Arch),
		d.Set("ecs_flavor", r.CupFlavor),
		d.Set("vpc_id", r.VpcId),
		d.Set("subnet_id", r.SubnetId),
		d.Set("security_group", r.SecurityGroupIds),
		d.Set("server_id", r.ServerId),
		d.Set("service_ip", r.ServiceIp),
		d.Set("run_status", r.RunStatus),
		d.Set("access_status", r.AccessStatus),
		d.Set("upgradable", r.Upgradable),
		d.Set("specification_code", r.ResourceSpecCode),
	)
	// Only ELB mode uses this field
	d.Set("group_id", r.PoolId)

	if mErr.ErrorOrNil() != nil {
		return fmtp.DiagErrorf("error setting WAF dedicated instance fields: %s", err)
	}
	return nil
}

// updateInstanceName call API to change the instance name.
func updateInstanceName(c *golangsdk.ServiceClient, id, name, epsId string) error {
	opt := instances.UpdateInstanceOpts{
		InstanceName: name,
	}

	_, err := instances.UpdateWithEpsId(c, id, opt, epsId)
	if err != nil {
		return fmtp.Errorf("error update name of WAF dedicate instance %s: %s", id, err)
	}
	return nil
}

func resourceDedicatedInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.WafDedicatedV1Client(region)
	if err != nil {
		return diag.Errorf("error creating WAF dedicated client: %s", err)
	}
	epsId := cfg.GetEnterpriseProjectID(d)
	instanceId := d.Id()
	if d.HasChanges("name") {
		err = updateInstanceName(client, instanceId, d.Get("name").(string), epsId)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   instanceId,
			ResourceType: "waf-instance",
			RegionId:     region,
			ProjectId:    client.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
		// check waf with enterprise_project_id
		_, err = instances.GetWithEpsId(client, instanceId, d.Get("enterprise_project_id").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		return nil
	}
	return resourceDedicatedInstanceRead(ctx, d, meta)
}

func waitForInstanceDeleted(c *golangsdk.ServiceClient, id string, epsId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := instances.GetWithEpsId(c, id, epsId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] The Waf dedicated instance has been deleted(ID:%s).", id)
				return &(instances.DedicatedInstance{}), "Deleted", nil
			}
			return nil, "Error", err
		}

		switch r.RunStatus {
		case runStatusDeleting:
			return r, "Deleting", nil
		case runStatusDeleted:
			return r, "Deleted", nil
		default:
			err = fmtp.Errorf("Error in delete WAF dedicated instance[%s]. "+
				"Unexpected run_status: %s.", r.Id, r.RunStatus)
			return r, "Error", err
		}
	}
}

func resourceDedicatedInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.WafDedicatedV1Client(cfg.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud WAF dedicated client: %s", err)
	}

	epsId := cfg.GetEnterpriseProjectID(d)
	_, err = instances.DeleteWithEpsId(client, d.Id(), epsId)
	if err != nil {
		// If the dedicated instance does not exist, the response HTTP status code of the deletion API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting WAF dedicated instance")
	}

	logp.Printf("[DEBUG] Waiting for WAF dedicated instance to be deleted(ID:%s).", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Deleting"},
		Target:       []string{"Deleted"},
		Refresh:      waitForInstanceDeleted(client, d.Id(), epsId),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 15 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		logp.Printf("[DEBUG] Error while waiting to delete Waf dedicated instance. \n%s : %v", d.Id(), err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
