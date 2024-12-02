package taurusdb

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/taurusdb/v3/configurations"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API GaussDBforMySQL GET /v3/{project_id}/configurations
func DataSourceGaussdbMysqlConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussdbMysqlConfigurationRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datastore_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datastore_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGaussdbMysqlConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.GaussdbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	configsList, err := configurations.List(client).Extract()
	if err != nil {
		return diag.Errorf("unable to retrieve configurations: %s", err)
	}
	if len(configsList) < 1 {
		return diag.Errorf("your query returned no results. " +
			"please change your search criteria and try again.")
	}

	if common.HasFilledOpt(d, "name") {
		var filteredConfigs []configurations.Configuration
		for _, conf := range configsList {
			if conf.Name == d.Get("name").(string) {
				filteredConfigs = append(filteredConfigs, conf)
			}
		}
		configsList = filteredConfigs
	}

	if len(configsList) < 1 {
		return diag.Errorf("your query returned no results. " +
			"please change your search criteria and try again.")
	}
	configuration := configsList[0]

	d.SetId(configuration.ID)

	mErr := multierror.Append(
		d.Set("name", configuration.Name),
		d.Set("description", configuration.Description),
		d.Set("datastore_version", configuration.DataStoreVer),
		d.Set("datastore_name", configuration.DataStoreName),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
