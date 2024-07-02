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
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
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
%[1]s

data "huaweicloud_apig_instances" "test" {
  instance_id = "%[2]s"
}

locals {
  instance_id = data.huaweicloud_apig_instances.test.instances[0].id
}

resource "huaweicloud_apig_application" "test" {
  count = 2

  instance_id = local.instance_id
  name        = format("%[3]s_%%d", count.index)
}

resource "huaweicloud_apig_appcode" "test" {
  instance_id    = local.instance_id
  application_id = huaweicloud_apig_application.test[0].id
}
`, common.TestBaseNetwork(name), acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccDataSourceAppcodes_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_apig_appcodes" "test" {
  depends_on = [
    huaweicloud_apig_appcode.test
  ]

  instance_id    = local.instance_id
  application_id = huaweicloud_apig_application.test[0].id
}

data "huaweicloud_apig_appcodes" "not_found" {
  instance_id    = local.instance_id
  application_id = huaweicloud_apig_application.test[1].id
}
`, testAccDataSourceAppcodes_basic_base())
}
