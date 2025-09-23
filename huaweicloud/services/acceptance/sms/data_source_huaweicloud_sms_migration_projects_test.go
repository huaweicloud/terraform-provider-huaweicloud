package sms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSmsMigrationProjects_basic(t *testing.T) {
	dataSource := "data.huaweicloud_sms_migration_projects.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceSmsMigrationProjects_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "migprojects.#"),
					resource.TestCheckResourceAttrSet(dataSource, "migprojects.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "migprojects.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "migprojects.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "migprojects.0.enterprise_project"),
					resource.TestCheckResourceAttrSet(dataSource, "migprojects.0.exist_server"),
					resource.TestCheckResourceAttrSet(dataSource, "migprojects.0.is_default"),
					resource.TestCheckResourceAttrSet(dataSource, "migprojects.0.region"),
					resource.TestCheckResourceAttrSet(dataSource, "migprojects.0.speed_limit"),
					resource.TestCheckResourceAttrSet(dataSource, "migprojects.0.start_network_check"),
					resource.TestCheckResourceAttrSet(dataSource, "migprojects.0.start_target_server"),
					resource.TestCheckResourceAttrSet(dataSource, "migprojects.0.syncing"),
					resource.TestCheckResourceAttrSet(dataSource, "migprojects.0.use_public_ip"),
				),
			},
		},
	})
}

func testDataSourceDataSourceSmsMigrationProjects_basic() string {
	return `
data "huaweicloud_sms_migration_projects" "test" {}
`
}
