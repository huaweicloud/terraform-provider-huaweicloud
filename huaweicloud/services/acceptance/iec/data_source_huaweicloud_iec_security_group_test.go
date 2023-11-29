package iec

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSecurityGroupDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("iec-secgroup-%s", acctest.RandString(5))
	description := "This is a test of iec security group"

	dataSourceName := "data.huaweicloud_iec_security_group.by_name"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecurityGroup_basic(rName, description),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "description", description),
				),
			},
		},
	})
}

func testAccDataSourceSecurityGroup_basic(rName, description string) string {
	return fmt.Sprintf(`
resource "huaweicloud_iec_security_group" "my_group" {
  name        = "%s"
  description = "%s"
}

data "huaweicloud_iec_security_group" "by_name" {
  name = huaweicloud_iec_security_group.my_group.name
}
`, rName, description)
}
