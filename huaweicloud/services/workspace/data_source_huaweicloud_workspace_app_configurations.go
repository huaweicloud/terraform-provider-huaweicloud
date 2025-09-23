package workspace

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace POST /v1/{project_id}/bundles/batch-query-config-info
func DataSourceAppConfigurations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppConfigurationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the Workspace APP are located.",
			},
			"items": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of configuration keys to be queried.",
			},
			"configurations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The key of the configuration.",
						},
						"config_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The value corresponding to the configuration key.",
						},
					},
				},
				Description: "The list of configuration information.",
			},
		},
	}
}

func dataSourceAppConfigurationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/bundles/batch-query-config-info"
	)
	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"items": d.Get("items"),
		},
	}

	requestResp, err := client.Request("POST", getPath, &opt)
	if err != nil {
		return diag.Errorf("error querying Workspace APP configurations: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("configurations", flattenAppConfigurations(utils.PathSearch("config_infos", respBody, make([]interface{}, 0)).([]interface{}))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAppConfigurations(configurations []interface{}) []map[string]interface{} {
	if len(configurations) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(configurations))
	for i, v := range configurations {
		result[i] = map[string]interface{}{
			"config_key":   utils.PathSearch("config_key", v, nil),
			"config_value": utils.PathSearch("config_value", v, nil),
		}
	}

	return result
}
