package waf

import (
	"context"
	"regexp"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/pools"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceWafInstanceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafInstanceGroupCreate,
		ReadContext:   resourceWafInstanceGroupRead,
		UpdateContext: resourceWafInstanceGroupUpdate,
		DeleteContext: resourceWafInstanceGroupDelete,
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[\w]{1,64}$`),
					"The maximum length is 64 characters. "+
						"Only letters, digits and underscores (_) are allowed"),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 256),
			},
			"body_limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      4000,
				ValidateFunc: validation.IntBetween(2000, 8000),
			},
			"header_limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      8000,
				ValidateFunc: validation.IntBetween(4000, 20000),
			},
			"connection_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      10,
				ValidateFunc: validation.IntBetween(1, 20),
			},
			"write_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      10,
				ValidateFunc: validation.IntBetween(1, 20),
			},
			"read_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      10,
				ValidateFunc: validation.IntBetween(1, 20),
			},
			"load_balances": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceWafInstanceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WafDedicatedV1Client(conf.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud WAF dedicated client : %s", err)
	}

	createOpts := pools.CreateOpts{
		Name:        d.Get("name").(string),
		Region:      conf.GetRegion(d),
		Type:        "elb",
		VpcID:       d.Get("vpc_id").(string),
		Description: d.Get("description").(string),
	}
	logp.Printf("[DEBUG] Create WAF instance group options: %#v", createOpts)

	pool, err := pools.Create(client, createOpts)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(pool.ID)

	return resourceWafInstanceGroupRead(ctx, d, meta)
}

func resourceWafInstanceGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WafDedicatedV1Client(conf.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud WAF dedicated client : %s", err)
	}

	pool, err := pools.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error querying WAF instance group")
	}

	mErr := multierror.Append(nil,
		d.Set("region", pool.Region),
		d.Set("name", pool.Name),
		d.Set("vpc_id", pool.VpcID),
		d.Set("body_limit", pool.Option.BodyLimit),
		d.Set("header_limit", pool.Option.HeaderLimit),
		d.Set("connection_timeout", pool.Option.ConnectTimeout),
		d.Set("write_timeout", pool.Option.SendTimeout),
		d.Set("read_timeout", pool.Option.ReadTimeout),
	)

	loadBalances := make([]interface{}, 0, len(pool.Bindings))
	for _, v := range pool.Bindings {
		loadBalances = append(loadBalances, v.Name)
	}
	d.Set("load_balances", loadBalances)
	d.Set("description", pool.Description)

	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("error setting WAF dedicated group attributes: %s", err)
	}

	return nil
}

func batchAddELBInstances(c *golangsdk.ServiceClient, poolID string, ids []interface{}) *multierror.Error {
	mErr := &multierror.Error{}
	for _, v := range ids {
		lbID := v.(string)
		_, e := pools.AddELB(c, poolID, lbID)
		if e != nil {
			err := fmtp.Errorf("Error in binding load balance[%s] to the group[%s]: %s", lbID, poolID, e)
			mErr = multierror.Append(mErr, err)
		}
	}
	return mErr
}

func batchRemoveELBInstances(c *golangsdk.ServiceClient, poolID string, ids []interface{}) *multierror.Error {
	mErr := &multierror.Error{}

	page, err := pools.ListELB(c, poolID).AllPages()
	if err != nil {
		return multierror.Append(mErr, err)
	}
	bindELBs, err := pools.ExtractBindELBs(page)
	if err != nil {
		return multierror.Append(mErr, err)
	}

	idMapping := make(map[string]string)
	for _, v := range bindELBs {
		idMapping[v.LoadBalancerID] = v.ID
	}
	for _, v := range ids {
		if bindingID, ok := idMapping[v.(string)]; ok {
			err = pools.RemoveELB(c, poolID, bindingID)
			if err != nil {
				err = fmtp.Errorf("Error in removing load balance[%s] from the group[%s]: %s", v, poolID, err)
				mErr = multierror.Append(mErr, err)
			}
		}
	}
	return mErr
}

func resourceWafInstanceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WafDedicatedV1Client(conf.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud WAF dedicated client : %s", err)
	}

	mErr := &multierror.Error{}
	if d.HasChanges("name", "description",
		"body_limit", "header_limit", "connection_timeout", "write_timeout", "read_timeout") {
		desc := d.Get("description").(string)
		updateOpts := pools.UpdatePoolOpts{
			Name:        d.Get("name").(string),
			Description: &desc,
		}
		if d.HasChanges("body_limit", "header_limit", "connection_timeout", "write_timeout", "read_timeout") {
			option := pools.PoolOption{
				BodyLimit:      d.Get("body_limit").(int),
				HeaderLimit:    d.Get("header_limit").(int),
				ConnectTimeout: d.Get("connection_timeout").(int),
				SendTimeout:    d.Get("write_timeout").(int),
				ReadTimeout:    d.Get("read_timeout").(int),
			}
			updateOpts.Option = &option
		}
		logp.Printf("[DEBUG] Create WAF instance group options: %#v", updateOpts)

		_, err = pools.Update(client, d.Id(), updateOpts)
		if err != nil {
			mErr = multierror.Append(mErr, err)
		}
	}

	if d.HasChange("load_balances") {
		oldVal, newVal := d.GetChange("load_balances")
		oldValSet := oldVal.(*schema.Set)
		newValSet := newVal.(*schema.Set)

		addBindings := newValSet.Difference(oldValSet)
		removeBindings := oldValSet.Difference(newValSet)

		if addBindings.Len() > 0 {
			errs := batchAddELBInstances(client, d.Id(), addBindings.List())
			mErr = multierror.Append(mErr, errs.Errors...)
		}
		if removeBindings.Len() > 0 {
			errs := batchRemoveELBInstances(client, d.Id(), removeBindings.List())
			mErr = multierror.Append(mErr, errs.Errors...)
		}
	}
	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("Error in update the groups[%s] : %s", d.Id(), err)
	}

	return resourceWafInstanceGroupRead(ctx, d, meta)
}

func resourceWafInstanceGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WafDedicatedV1Client(conf.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error in creating HuaweiCloud WAF dedicated client : %s", err)
	}
	// remove the bound ELB instances before deleting the group
	elbs := d.Get("load_balances").(*schema.Set)
	if elbs.Len() > 0 {
		mErr := batchRemoveELBInstances(client, d.Id(), elbs.List())
		if err = mErr.ErrorOrNil(); err != nil {
			return fmtp.DiagErrorf("error in removing ELB instances from group: %s", err)
		}
	}

	_, err = pools.Delete(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error in deleting WAF instance group: %s")
	}

	d.SetId("")
	return nil
}
