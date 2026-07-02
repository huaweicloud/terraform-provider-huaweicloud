package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDatabaseInstances_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dsc_database_instances.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDatabaseInstances_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instance_type"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.#"),
				),
			},
		},
	})
}

const testDataSourceDatabaseInstances_basic = `
data "huaweicloud_dsc_database_instances" "test" {
  instance_type = "RDS"
}
`
