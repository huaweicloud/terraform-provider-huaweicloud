package waf

import (
	"time"

	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/valuelists"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

// DataSourceWafReferenceTablesV1 the function is used for data source 'huaweicloud_waf_reference_tables'.
func DataSourceWafReferenceTablesV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceWafReferenceTablesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

func dataSourceWafReferenceTablesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF client: %s", err)
	}

	r, err := valuelists.List(client, valuelists.ListValueListOpts{})
	if err != nil {
		return common.CheckDeleted(d, err, "Error obtain WAF reference table information")
	}

	if len(r.Items) == 0 {
		return fmtp.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	ids := make([]string, 0, len(r.Items))
	tables := make([]map[string]interface{}, 0, len(r.Items))
	for _, t := range r.Items {
		tab := map[string]interface{}{
			"id":            t.Id,
			"name":          t.Name,
			"type":          t.Type,
			"conditions":    t.Values,
			"description":   t.Description,
			"creation_time": time.Unix(t.CreationTime/1000, 0).Format("2006-01-02 15:04:05"),
		}
		tables = append(tables, tab)
		ids = append(ids, t.Id)
	}

	d.SetId(hashcode.Strings(ids))
	if err = d.Set("tables", tables); err != nil {
		return fmtp.Errorf("error setting WAF reference table fields: %s", err)
	}

	return nil
}
