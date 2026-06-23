package geminidb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceConfigurationDatastores_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_geminidb_configuration_datastores.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConfigurationDatastores_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "datastores.#"),
					resource.TestCheckResourceAttrSet(dataSource, "datastores.0.datastore_name"),
					resource.TestCheckResourceAttrSet(dataSource, "datastores.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "datastores.0.mode"),
				),
			},
		},
	})
}

const testAccDataSourceConfigurationDatastores_basic = `data "huaweicloud_geminidb_configuration_datastores" "test" {}`
