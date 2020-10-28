package huaweicloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/huaweicloud/golangsdk/openstack/taurusdb/v3/configurations"
)

func dataSourceGaussdbMysqlConfigurations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGaussdbMysqlConfigurationsRead,

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

func dataSourceGaussdbMysqlConfigurationsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	client, err := config.gaussdbV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud GaussDB client: %s", err)
	}

	configsList, err := configurations.List(client).Extract()
	if err != nil {
		return fmt.Errorf("Unable to retrieve configurations: %s", err)
	}
	if len(configsList) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if hasFilledOpt(d, "name") {
		var filteredConfigs []configurations.Configuration
		for _, conf := range configsList {
			if conf.Name == d.Get("name").(string) {
				filteredConfigs = append(filteredConfigs, conf)
			}
		}
		configsList = filteredConfigs
	}

	if len(configsList) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}
	Configuration := configsList[0]

	d.SetId(Configuration.ID)

	d.Set("name", Configuration.Name)
	d.Set("description", Configuration.Description)
	d.Set("datastore_version", Configuration.DataStoreVer)
	d.Set("datastore_name", Configuration.DataStoreName)

	return nil
}
