package waf

import (
	"context"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/pools"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func ResourceWafInstGroupAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafInsGroupAssociateCreate,
		ReadContext:   resourceWafInsGroupAssociateRead,
		UpdateContext: resourceWafInsGroupAssociateUpdate,
		DeleteContext: resourceWafInsGroupAssociateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"load_balancers": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceWafInsGroupAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WafDedicatedV1Client(conf.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud WAF dedicated client : %s", err)
	}

	groupID := d.Get("group_id").(string)
	group, err := pools.Get(client, groupID)
	if err != nil {
		return fmtp.DiagErrorf("Error querying WAF instance group: %s", err)
	}

	if len(group.Bindings) > 0 {
		// Remove the bound instance
		for _, v := range group.Bindings {
			err = pools.RemoveELB(client, groupID, v.ID)
			if err != nil {
				return fmtp.DiagErrorf("Error removing load balance[%s] from the group[%s]: %s", v, groupID, err)
			}
		}
	}

	elbIDs := d.Get("load_balancers").(*schema.Set).List()
	mErr := addELBInstances(client, groupID, elbIDs)
	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}
	d.SetId(groupID)

	return resourceWafInsGroupAssociateRead(ctx, d, meta)
}

func addELBInstances(c *golangsdk.ServiceClient, groupID string, ids []interface{}) *multierror.Error {
	var mErr *multierror.Error
	for _, v := range ids {
		lbID := v.(string)
		_, e := pools.AddELB(c, groupID, lbID)
		if e != nil {
			err := fmtp.Errorf("Error in binding load balance[%s] to the group[%s]: %s", lbID, groupID, e)
			mErr = multierror.Append(mErr, err)
		}
	}
	return mErr
}

func resourceWafInsGroupAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WafDedicatedV1Client(conf.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud WAF dedicated client : %s", err)
	}

	group, err := pools.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error querying WAF instance group")
	}

	loadBalances := make([]interface{}, 0, len(group.Bindings))
	for _, v := range group.Bindings {
		loadBalances = append(loadBalances, v.Name)
	}
	mErr := multierror.Append(nil,
		d.Set("group_id", group.ID),
		d.Set("load_balancers", loadBalances),
	)

	if mErr.ErrorOrNil() != nil {
		return fmtp.DiagErrorf("error setting WAF dedicated group attributes: %s", err)
	}

	return nil
}

func resourceWafInsGroupAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WafDedicatedV1Client(conf.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud WAF dedicated client : %s", err)
	}

	mErr := &multierror.Error{}
	oldVal, newVal := d.GetChange("load_balancers")
	oldValSet := oldVal.(*schema.Set)
	newValSet := newVal.(*schema.Set)

	addBindings := newValSet.Difference(oldValSet)
	removeBindings := oldValSet.Difference(newValSet)

	if addBindings.Len() > 0 {
		errs := addELBInstances(client, d.Id(), addBindings.List())
		mErr = multierror.Append(mErr, errs.Errors...)
	}
	if removeBindings.Len() > 0 {
		errs := batchRemoveELBInstances(client, d.Id(), removeBindings.List())
		mErr = multierror.Append(mErr, errs.Errors...)
	}
	if mErr.ErrorOrNil() != nil {
		return fmtp.DiagErrorf("error setting WAF dedicated group attributes: %s", err)
	}

	return resourceWafInsGroupAssociateRead(ctx, d, meta)
}

func resourceWafInsGroupAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WafDedicatedV1Client(conf.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error in creating HuaweiCloud WAF dedicated client : %s", err)
	}
	// remove the bound ELB instances before deleting the group
	elbIDs := d.Get("load_balancers").(*schema.Set)
	if elbIDs.Len() > 0 {
		mErr := batchRemoveELBInstances(client, d.Id(), elbIDs.List())
		if mErr.ErrorOrNil() != nil {
			return fmtp.DiagErrorf("error in removing ELB instances from group: %s", err)
		}
	}

	d.SetId("")
	return nil
}
