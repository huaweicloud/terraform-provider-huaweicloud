package rds

import (
	"context"
	"sort"

	"github.com/chnsz/golangsdk/openstack/rds/v3/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func DataSourceRdsEngineVersionsV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsEngineVersionsV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"MySQL", "PostgreSQL", "SQLServer",
				}, false),
				Default: "MySQL",
			},
			"versions": {
				Type:     schema.TypeList,
				Optional: true,
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
	}
}

func dataSourceRdsEngineVersionsV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.RdsV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud RDS v3 client: %s", err)
	}

	engineType := d.Get("type").(string)
	engine, err := instances.ListEngine(client, engineType)
	if err != nil {
		return fmtp.DiagErrorf("Error getting version list of specific database engine: %s", err)
	}

	versions := engine.Versions
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Name < versions[j].Name
	})
	logp.Printf("After sorting, the engine version list is: %+v", versions)

	ids := make([]string, len(versions))
	result := make([]map[string]interface{}, len(versions))

	for i, engine := range versions {
		vMap := map[string]interface{}{
			"id":   engine.ID,
			"name": engine.Name,
		}
		result[i] = vMap
		ids[i] = engine.ID
	}

	d.SetId(hashcode.Strings(ids))

	return diag.FromErr(d.Set("versions", result))
}
