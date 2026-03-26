package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccConfigurationCompare_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_rds_parametergroup_compare.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testConfigurationCompare_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "source_name", fmt.Sprintf("%s_source", name)),
					resource.TestCheckResourceAttr(rName, "target_name", fmt.Sprintf("%s_target", name)),
					resource.TestCheckResourceAttr(rName, "parameters.#", "1"),
					resource.TestCheckResourceAttr(rName, "parameters.0.name", "auto_increment_increment"),
					resource.TestCheckResourceAttr(rName, "parameters.0.source_value", "2"),
					resource.TestCheckResourceAttr(rName, "parameters.0.target_value", "4"),
				),
			},
		},
	})
}

func testConfigurationCompare_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_parametergroup" "source" {
  name = "%[1]s_source"

  values = {
    auto_increment_increment = "2"
  }

  datastore {
    type    = "mysql"
    version = "8.0"
  }
}

resource "huaweicloud_rds_parametergroup" "target" {
  name = "%[1]s_target"

  values = {
    auto_increment_increment = "4"
  }

  datastore {
    type    = "mysql"
    version = "8.0"
  }
}

resource "huaweicloud_rds_parametergroup_compare" "test" {
  source_id = huaweicloud_rds_parametergroup.source.id
  target_id = huaweicloud_rds_parametergroup.target.id
}
`, name)
}
