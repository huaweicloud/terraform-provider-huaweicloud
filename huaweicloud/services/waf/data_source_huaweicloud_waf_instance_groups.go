package waf

import (
	"context"

	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/pools"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func DataSourceWafInstanceGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceWafInstanceGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"body_limit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"header_limit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"connection_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"write_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"read_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"load_balancers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"dedicated_instances": {
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
						"domain_names": {
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
				},
			},
		},
	}
}

func DataSourceWafInstanceGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WafDedicatedV1Client(conf.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud WAF dedicated client : %s", err)
	}

	opts := pools.ListPoolOpts{
		Name:   d.Get("name").(string),
		VpcID:  d.Get("vpc_id").(string),
		Detail: true,
	}
	page, err := pools.List(client, opts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error querying WAF instance group.")
	}

	p, err := page.AllPages()
	if err != nil {
		return fmtp.DiagErrorf("error querying WAF instance group: %s", err)
	}
	groups, err := pools.ExtractGroups(p)
	if err != nil {
		return fmtp.DiagErrorf("error querying WAF instance group: %s", err)
	}

	if len(groups) == 0 {
		return fmtp.DiagErrorf("Your query returned no results.  " +
			"Please change your search criteria and try again.")
	}

	ids := make([]string, len(groups))
	grps := make([]map[string]interface{}, 0, len(groups))
	for i, g := range groups {
		ids[i] = g.ID

		loadBalances := make([]interface{}, 0, len(g.Bindings))
		for _, v := range g.Bindings {
			loadBalances = append(loadBalances, v.Name)
		}

		domainNames := make([]interface{}, 0, len(g.Hosts))
		for _, v := range g.Hosts {
			dName := map[string]interface{}{
				"id":   v.ID,
				"name": v.Name,
			}
			domainNames = append(domainNames, dName)
		}
		group := map[string]interface{}{
			"region":             g.Region,
			"name":               g.Name,
			"vpc_id":             g.VpcID,
			"description":        g.Description,
			"body_limit":         g.Option.BodyLimit,
			"header_limit":       g.Option.HeaderLimit,
			"connection_timeout": g.Option.ConnectTimeout,
			"write_timeout":      g.Option.SendTimeout,
			"read_timeout":       g.Option.ReadTimeout,
			"load_balancers":     loadBalances,
			"domain_names":       domainNames,
		}
		grps = append(grps, group)
	}

	d.SetId(hashcode.Strings(ids))
	mErr := multierror.Append(nil,
		d.Set("region", conf.GetRegion(d)),
		d.Set("groups", grps),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("error setting WAF instance group attributes: %s", err)
	}

	return nil
}
