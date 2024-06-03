package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceAppcodes_basic(t *testing.T) {
	var (
		rName = "data.huaweicloud_apig_appcodes.test"
		dc    = acceptance.InitDataSourceCheck(rName)

		rNameNotFound = "data.huaweicloud_apig_appcodes.not_found"
		dcNotFound    = acceptance.InitDataSourceCheck(rNameNotFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppcodes_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "appcodes.#", "1"),
					dcNotFound.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameNotFound, "appcodes.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceAppcodes_basic_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

%[1]s

resource "huaweicloud_apig_instance" "test" {
  name                  = "%[2]s"
  edition               = "BASIC"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "0"

  availability_zones = try(slice(data.huaweicloud_availability_zones.test.names, 0, 1), null)
}

resource "huaweicloud_apig_application" "test" {
  count = 2

  instance_id = huaweicloud_apig_instance.test.id
  name        = format("%[2]s_%%d", count.index)
}

resource "huaweicloud_apig_appcode" "test" {
  instance_id    = huaweicloud_apig_instance.test.id
  application_id = huaweicloud_apig_application.test[0].id
}
`, common.TestBaseNetwork(name), name)
}

func testAccDataSourceAppcodes_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_apig_appcodes" "test" {
  depends_on = [
    huaweicloud_apig_appcode.test
  ]

  instance_id    = huaweicloud_apig_instance.test.id
  application_id = huaweicloud_apig_application.test[0].id
}

data "huaweicloud_apig_appcodes" "not_found" {
  instance_id    = huaweicloud_apig_instance.test.id
  application_id = huaweicloud_apig_application.test[1].id
}
`, testAccDataSourceAppcodes_basic_base())
}
