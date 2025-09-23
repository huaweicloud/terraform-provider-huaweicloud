package waf

import (
	"context"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/valuelists"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API WAF GET /v1/{project_id}/waf/valuelist
func DataSourceWafReferenceTables() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWafReferenceTablesRead,

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
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tables": {
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
						"conditions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceWafReferenceTablesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.WafV1Client(region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	opts := valuelists.ListValueListOpts{
		EnterpriseProjectId: conf.GetEnterpriseProjectID(d),
	}

	r, err := valuelists.List(client, opts)
	if err != nil {
		return diag.Errorf("error retrieving WAF reference tables %s", err)
	}

	tables := make([]map[string]interface{}, 0, len(r.Items))
	name := d.Get("name").(string)
	for _, v := range r.Items {
		if name != "" && name != v.Name {
			continue
		}

		tab := map[string]interface{}{
			"id":            v.Id,
			"name":          v.Name,
			"type":          v.Type,
			"conditions":    v.Values,
			"description":   v.Description,
			"creation_time": time.Unix(v.CreationTime/1000, 0).Format("2006-01-02 15:04:05"),
		}
		tables = append(tables, tab)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("tables", tables),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
