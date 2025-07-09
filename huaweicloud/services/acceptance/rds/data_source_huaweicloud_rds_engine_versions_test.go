package rds

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRdsEngineVersionsV3DataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_rds_engine_versions.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsEngineVersions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "versions.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "versions.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "versions.0.name"),
				),
			},
		},
	})
}

func testDataSourceRdsEngineVersions_basic() string {
	return `
data "huaweicloud_rds_engine_versions" "test" {
  type = "MySQL"
}
`
}
