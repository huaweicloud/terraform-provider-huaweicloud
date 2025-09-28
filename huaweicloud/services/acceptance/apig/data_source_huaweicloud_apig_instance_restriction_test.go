package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceInstanceRestriction_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_apig_instance_restriction.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceInstanceRestriction_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "restrict_cidrs.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestCheckResourceAttrSet(dcName, "resource_subnet_cidr"),
				),
			},
		},
	})
}

func testAccDataSourceInstanceRestriction_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_apig_instances" "test" {
  instance_id = "%[1]s"
}

locals {
  instance_id = data.huaweicloud_apig_instances.test.instances[0].id
}

data "huaweicloud_apig_instance_restriction" "test" {
  instance_id = local.instance_id
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}
