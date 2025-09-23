package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGaussDBMysqlTemplateCompare_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_gaussdb_mysql_parameter_template_compare.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testParameterTemplateCompare_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "differences.#", "1"),
					resource.TestCheckResourceAttr(rName, "differences.0.parameter_name", "auto_increment_increment"),
					resource.TestCheckResourceAttr(rName, "differences.0.source_value", "1"),
					resource.TestCheckResourceAttr(rName, "differences.0.target_value", "4"),
				),
			},
		},
	})
}

func testParameterTemplateCompare_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_mysql_parameter_template" "source" {
  name              = "%[1]s_source"
  description       = "test gaussdb mysql parameter template"
  datastore_engine  = "gaussdb-mysql"
  datastore_version = "8.0"

  parameter_values = {
    auto_increment_increment = "1"
  }
}

resource "huaweicloud_gaussdb_mysql_parameter_template" "target" {
  name              = "%[1]s_target"
  description       = "test gaussdb mysql parameter template"
  datastore_engine  = "gaussdb-mysql"
  datastore_version = "8.0"

  parameter_values = {
    auto_increment_increment = "4"
  }
}

resource "huaweicloud_gaussdb_mysql_parameter_template_compare" "test" {
  source_configuration_id = huaweicloud_gaussdb_mysql_parameter_template.source.id
  target_configuration_id = huaweicloud_gaussdb_mysql_parameter_template.target.id
}
`, name)
}
