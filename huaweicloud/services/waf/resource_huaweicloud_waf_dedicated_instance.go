package waf

import (
	"context"
	"errors"
	"fmt"
	"log"
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
)

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

		CustomizeDiff: config.MergeDefaultTags(),

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
			// This field has no response return value.
			"tags": common.TagsForceNewSchema(),

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

			// Deprecated; Reasons for abandonment are as follows:
			// `group_id`: Legacy fields are no longer supported.
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `schema: Deprecated;`,
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

	createOpts := instances.CreateInstanceOpts{
		Region:        region,
		ChargeMode:    30,
		AvailableZone: d.Get("available_zone").(string),
		Arch:          d.Get("cpu_architecture").(string),
		NamePrefix:    d.Get("name").(string),
		Specification: d.Get("specification_code").(string),
		CpuFlavor:     d.Get("ecs_flavor").(string),
		VpcId:         d.Get("vpc_id").(string),
		SubnetId:      d.Get("subnet_id").(string),
		SecurityGroup: groups,
		Count:         1,
		PoolId:        d.Get("group_id").(string),
		ResTenant:     utils.Bool(d.Get("res_tenant").(bool)),
		Tags:          utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
	}
	if d.Get("res_tenant").(bool) {
		// `anti_affinity` is valid only when `res_tenant` is true
		createOpts.AntiAffinity = utils.Bool(d.Get("anti_affinity").(bool))
	}

	return &createOpts
}

func waitForInstanceCreated(c *golangsdk.ServiceClient, id string, epsId string) resource.StateRefreshFunc {
	unexpectedRunStatus := []interface{}{2, 3, 4, 5, 6, 7, 8}
	return func() (interface{}, string, error) {
		r, err := instances.GetWithEpsId(c, id, epsId)
		if err != nil {
			return nil, "ERROR", err
		}

		runStatus := r.RunStatus
		if runStatus == 1 {
			return r, "COMPLETED", nil
		}

		if utils.SliceContains(unexpectedRunStatus, runStatus) {
			return r, "ERROR", fmt.Errorf("got unexpected run status %d", runStatus)
		}
		return r, "PENDING", nil
	}
}

func resourceDedicatedInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.WafDedicatedV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF dedicated client: %s", err)
	}

	createOpts := buildCreateOpts(d, cfg.GetRegion(d))
	epsId := cfg.GetEnterpriseProjectID(d)

	r, err := instances.CreateWithEpsId(client, *createOpts, epsId)
	if err != nil {
		return diag.Errorf("error creating WAF dedicated instance: %s", err)
	}

	if r == nil || len(r.Instances) == 0 || r.Instances[0].Id == "" {
		return diag.Errorf("error creating WAF dedicated instance: ID is not found in API response")
	}

	d.SetId(r.Instances[0].Id)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      waitForInstanceCreated(client, d.Id(), epsId),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 15 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for WAF dedicated instance (%s) creation to completed: %s", d.Id(), err)
	}

	err = updateInstanceName(client, d.Id(), d.Get("name").(string), epsId)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDedicatedInstanceRead(ctx, d, meta)
}

func resourceDedicatedInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.WafDedicatedV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF dedicated client: %s", err)
	}

	epsId := cfg.GetEnterpriseProjectID(d)
	r, err := instances.GetWithEpsId(client, d.Id(), epsId)
	if err != nil {
		// If the dedicated instance does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving WAF dedicated instance")
	}

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
		// Only ELB mode uses this field
		d.Set("group_id", r.PoolId),
		d.Set("tags", d.Get("tags")),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting WAF dedicated instance fields: %s", err)
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
		return fmt.Errorf("error updating WAF dedicated instance (%s) name: %s", id, err)
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

	// Prioritize changes to enterprise project ID to avoid failure of other changes.
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
	}

	if d.HasChanges("name") {
		err = updateInstanceName(client, instanceId, d.Get("name").(string), epsId)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceDedicatedInstanceRead(ctx, d, meta)
}

func waitForInstanceDeleted(c *golangsdk.ServiceClient, id string, epsId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := instances.GetWithEpsId(c, id, epsId)
		if err != nil {
			var errDefault404 golangsdk.ErrDefault404
			if errors.As(err, &errDefault404) {
				log.Printf("[DEBUG] The WAF dedicated instance (%s) has been deleted.", id)
				return "success deleted", "COMPLETED", nil
			}
			return nil, "ERROR", err
		}

		runStatus := r.RunStatus
		if runStatus == 3 {
			return r, "COMPLETED", nil
		}

		return r, "PENDING", nil
	}
}

func resourceDedicatedInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.WafDedicatedV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF dedicated client: %s", err)
	}

	epsId := cfg.GetEnterpriseProjectID(d)
	_, err = instances.DeleteWithEpsId(client, d.Id(), epsId)
	if err != nil {
		// If the dedicated instance does not exist, the response HTTP status code of the deletion API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting WAF dedicated instance")
	}

	log.Printf("[DEBUG] Waiting for WAF dedicated instance (%s) to be deleted.", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      waitForInstanceDeleted(client, d.Id(), epsId),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 15 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for WAF dedicated instance (%s) deletion to completed: %s", d.Id(), err)
	}
	return nil
}
