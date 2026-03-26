package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccConfigurationReset_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_rds_parametergroup_reset.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testConfigurationReset_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "config_name"),
					resource.TestCheckResourceAttrSet(rName, "need_restart"),
				),
			},
		},
	})
}

func testConfigurationReset_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_parametergroup_reset" "test" {
  config_id = huaweicloud_rds_parametergroup.test.id
}
`, testAccRdsConfig_basic(name))
}
