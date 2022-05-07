package servicestage

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/servicestage/v2/metadata"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceComponentRuntimes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComponentRuntimesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"runtimes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceComponentRuntimesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.ServiceStageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage v2 client: %s", err)
	}

	resp, err := metadata.ListRuntimes(client)
	if err != nil {
		return diag.Errorf("error retrieving the runtimes list: %s", err)
	}

	filter := map[string]interface{}{
		"Type":                 d.Get("name"),
		"ContainerDefaultPort": d.Get("default_port"),
	}
	filtResult, err := utils.FilterSliceWithField(resp, filter)
	if err != nil {
		return diag.Errorf("filter component runtimes failed: %s", err)
	}
	log.Printf("filter %d component runtimes from %d through option %v", len(filtResult), len(resp), filter)

	types := make([]string, len(filtResult))
	runtimes := make([]map[string]interface{}, len(filtResult))
	for i, val := range filtResult {
		runtime := val.(metadata.Runtime)
		types[i] = runtime.Type
		runtimes[i] = map[string]interface{}{
			"name":         runtime.Type,
			"default_port": runtime.ContainerDefaultPort,
			"description":  runtime.TypeDesc,
		}
	}
	d.SetId(hashcode.Strings(types))

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("runtimes", runtimes),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
