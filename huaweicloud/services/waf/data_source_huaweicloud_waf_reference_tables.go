package waf

import (
	"context"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/valuelists"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// DataSourceWafReferenceTablesV1 the function is used for data source 'huaweicloud_waf_reference_tables'.
// @API WAF GET /v1/{project_id}/waf/valuelist
func DataSourceWafReferenceTablesV1() *schema.Resource {
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
	client, err := conf.WafV1Client(conf.GetRegion(d))
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

	if len(r.Items) == 0 {
		return nil
	}
	// filter data by name
	filterData, err := utils.FilterSliceWithField(r.Items, map[string]interface{}{
		"Name": d.Get("name").(string),
	})
	if err != nil {
		return diag.Errorf("error filtering WAF reference tables: %s", err)
	}
	tables := make([]map[string]interface{}, 0, len(filterData))
	ids := make([]string, 0, len(r.Items))
	for _, t := range filterData {
		v := t.(valuelists.WafValueList)
		tab := map[string]interface{}{
			"id":            v.Id,
			"name":          v.Name,
			"type":          v.Type,
			"conditions":    v.Values,
			"description":   v.Description,
			"creation_time": time.Unix(v.CreationTime/1000, 0).Format("2006-01-02 15:04:05"),
		}
		tables = append(tables, tab)
		ids = append(ids, v.Id)
	}

	d.SetId(hashcode.Strings(ids))
	mErr := multierror.Append(nil, d.Set("tables", tables))

	return diag.FromErr(mErr.ErrorOrNil())
}
