package iec

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/iec/v1/security/groups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IEC POST /v1/security-groups
// @API IEC DELETE /v1/security-groups/{security_group_id}
// @API IEC GET /v1/security-groups/{security_group_id}
func ResourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityGroupCreate,
		ReadContext:   resourceSecurityGroupRead,
		DeleteContext: resourceSecurityGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"security_group_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"direction": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ethertype": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_range_max": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"port_range_min": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"remote_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_ip_prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceSecurityGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iecClient, err := cfg.IECV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	createOpts := groups.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	group, err := groups.Create(iecClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating IEC security group: %s", err)
	}

	d.SetId(group.ID)
	return resourceSecurityGroupRead(ctx, d, meta)
}

func resourceSecurityGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iecClient, err := cfg.IECV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	group, err := groups.Get(iecClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "iec security group")
	}

	secRules := make([]map[string]interface{}, len(group.SecurityGroupRules))
	for index, rule := range group.SecurityGroupRules {
		secRules[index] = map[string]interface{}{
			"id":                rule.ID,
			"security_group_id": rule.SecurityGroupID,
			"description":       rule.Description,
			"direction":         rule.Direction,
			"ethertype":         rule.EtherType,
			"protocol":          rule.Protocol,
			"remote_group_id":   rule.RemoteGroupID,
			"remote_ip_prefix":  rule.RemoteIPPrefix,
		}

		if ret, err := strconv.Atoi(rule.PortRangeMax.(string)); err == nil {
			secRules[index]["port_range_max"] = ret
		}
		if ret, err := strconv.Atoi(rule.PortRangeMin.(string)); err == nil {
			secRules[index]["port_range_min"] = ret
		}
	}

	mErr := multierror.Append(
		d.Set("name", group.Name),
		d.Set("description", group.Description),
		d.Set("security_group_rules", secRules),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting fields: %s", err)
	}

	return nil
}

func resourceSecurityGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iecClient, err := cfg.IECV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating networking client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForSecurityGroupDelete(iecClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error deleting IEC security group: %s", err)
	}

	return nil
}

func waitForSecurityGroupDelete(iecClient *golangsdk.ServiceClient, groupID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] attempting to delete security group %s.\n", groupID)
		sg, err := groups.Get(iecClient, groupID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] successfully deleted IEC security group %s", groupID)
				return sg, "DELETED", nil
			}
			return sg, "ACTIVE", err
		}

		err = groups.Delete(iecClient, groupID).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] successfully deleted IEC security group %s", groupID)
				return sg, "DELETED", nil
			}
			if _, ok := err.(golangsdk.ErrDefault409); ok {
				return sg, "ACTIVE", nil
			}
			return sg, "ACTIVE", err
		}

		log.Printf("[DEBUG] IEC security group %s still active.\n", groupID)
		return sg, "ACTIVE", nil
	}
}
