package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDdsParameterTemplates_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dds_parameter_templates.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDdsParameterTemplates_noArgs(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.node_type"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.datastore_version"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.datastore_name"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.user_defined"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.updated_at"),
				),
			},
			{
				Config: testDdsParameterTemplates_args(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "configurations.#", "1"),
					resource.TestCheckResourceAttr(dataSource, "configurations.0.name", rName),
					resource.TestCheckResourceAttr(dataSource, "configurations.0.node_type", "single"),
					resource.TestCheckResourceAttr(dataSource, "configurations.0.datastore_version", "4.0"),
					resource.TestCheckResourceAttr(dataSource, "configurations.0.description", "created by acc"),

					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.datastore_name"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.user_defined"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.updated_at"),
				),
			},
		},
	})
}

func testDdsParameterTemplates_noArgs(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dds_parameter_template" "test"{
  name         = "%s"
  node_type    = "single"
  node_version = "4.0"
  description  = "created by acc"

  parameter_values = {
    connPoolMaxConnsPerHost        = 200
    connPoolMaxShardedConnsPerHost = 200
  }
}

data "huaweicloud_dds_parameter_templates" "test" {
  depends_on = [huaweicloud_dds_parameter_template.test]
}`, name)
}

func testDdsParameterTemplates_args(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dds_parameter_template" "test"{
  name         = "%s"
  node_type    = "single"
  node_version = "4.0"
  description  = "created by acc"

  parameter_values = {
    connPoolMaxConnsPerHost        = 200
    connPoolMaxShardedConnsPerHost = 200
  }
}

data "huaweicloud_dds_parameter_templates" "test" {
  name              = "%s"
  node_type         = "single"
  datastore_version = "4.0"

  depends_on = [huaweicloud_dds_parameter_template.test]
}
`, name, name)
}
