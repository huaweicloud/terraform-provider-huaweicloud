package dds

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDDSFlavorV3DataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dds_flavors.flavor"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDDSFlavorV3DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.engine_name", "DDS-Community"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.engine_versions.0", "3.4"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.vcpus", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.memory", "4"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.type", "mongos"),

					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.az_status.%"),
				),
			},
		},
	})
}

var testAccDDSFlavorV3DataSource_basic = `
data "huaweicloud_dds_flavors" "flavor" {
  engine_name    = "DDS-Community"
  engine_version = "3.4"
  vcpus          = 2
  memory         = 4
  type           = "mongos"
}
`
