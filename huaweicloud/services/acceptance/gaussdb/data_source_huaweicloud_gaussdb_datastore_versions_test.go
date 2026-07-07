package gaussdb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDatastoreVersions_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_gaussdb_datastore_versions.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDatastoreVersions_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "database_versions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "database_versions.0.software_version"),
					resource.TestCheckResourceAttrSet(dataSource, "database_versions.0.hotfixes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "database_versions.0.hotfixes.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "database_versions.0.hotfixes.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "database_versions.0.hotfixes.0.deploy_modes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "database_versions.0.hotfixes.0.deploy_modes.0.default_upgrade"),
					resource.TestCheckResourceAttrSet(dataSource, "database_versions.0.hotfixes.0.deploy_modes.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "database_versions.0.hotfixes.0.deploy_modes.0.mode"),
				),
			},
		},
	})
}

const testDataSourceDatastoreVersions_basic = `data "huaweicloud_gaussdb_datastore_versions" "test" {}`
