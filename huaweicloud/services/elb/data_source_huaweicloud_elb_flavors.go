package elb

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/elb/v3/flavors"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
)

// @API ELB GET /v3/{project_id}/elb/flavors
func DataSourceElbFlavorsV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceElbFlavorsV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_connections": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cps": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"qps": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// Computed values.
			"flavors": {
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
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_connections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cps": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"qps": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceElbFlavorsV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	listOpts := flavors.ListOpts{}
	if v, ok := d.GetOk("type"); ok {
		listOpts.Type = []string{v.(string)}
	}

	pages, err := flavors.List(elbClient, listOpts).AllPages()
	if err != nil {
		return diag.FromErr(err)
	}

	allFlavors, err := flavors.ExtractFlavors(pages)
	if err != nil {
		return diag.Errorf("unable to retrieve flavors: %s", err)
	}

	maxConnections := d.Get("max_connections").(int)
	cps := d.Get("cps").(int)
	qps := d.Get("qps").(int)
	bandwidth := d.Get("bandwidth").(int)

	var ids []string
	var flavorInfos []map[string]interface{}
	for _, flavor := range allFlavors {
		if flavor.SoldOut {
			continue
		}

		if maxConnections > 0 && flavor.Info.Connection != maxConnections {
			continue
		}

		if cps > 0 && flavor.Info.Cps != cps {
			continue
		}

		if qps > 0 && flavor.Info.Qps != qps {
			continue
		}

		if bandwidth > 0 && flavor.Info.Bandwidth != bandwidth*1000 {
			continue
		}

		ids = append(ids, flavor.ID)
		flavorInfo := map[string]interface{}{
			"id":              flavor.ID,
			"name":            flavor.Name,
			"type":            flavor.Type,
			"max_connections": flavor.Info.Connection,
			"cps":             flavor.Info.Cps,
			"qps":             flavor.Info.Qps,
			"bandwidth":       flavor.Info.Bandwidth / 1000,
		}
		flavorInfos = append(flavorInfos, flavorInfo)
	}

	if len(ids) < 1 {
		return diag.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	d.SetId(hashcode.Strings(ids))

	mErr := multierror.Append(
		d.Set("region", cfg.GetRegion(d)),
		d.Set("ids", ids),
		d.Set("flavors", flavorInfos),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
