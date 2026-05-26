package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataStudioInstances_basic(t *testing.T) {
	var (
		beforeInstanceCreation   = "data.huaweicloud_dataarts_studio_instances.before_instance_creation"
		dcBeforeInstanceCreation = acceptance.InitDataSourceCheck(beforeInstanceCreation)

		all = "data.huaweicloud_dataarts_studio_instances.test"
		dc  = acceptance.InitDataSourceCheck(all)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataArtsStudioInstances_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dcBeforeInstanceCreation.CheckResourceExists(),
					resource.TestCheckResourceAttr(beforeInstanceCreation, "instances.#", "0"),
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "instances.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrPair(all, "instances.0.id", "huaweicloud_dataarts_studio_instance.test", "id"),
					resource.TestCheckResourceAttrPair(all, "instances.0.name", "huaweicloud_dataarts_studio_instance.test", "name"),
					resource.TestCheckResourceAttrPair(all, "instances.0.version", "huaweicloud_dataarts_studio_instance.test", "version"),
					resource.TestCheckResourceAttrPair(all, "instances.0.order_id", "huaweicloud_dataarts_studio_instance.test", "order_id"),
					resource.TestCheckResourceAttrSet(all, "instances.0.product_id"),
					resource.TestCheckResourceAttrPair(all, "instances.0.auto_renew", "huaweicloud_dataarts_studio_instance.test", "auto_renew"),
					resource.TestCheckResourceAttrPair(all, "instances.0.enterprise_project_id",
						"huaweicloud_dataarts_studio_instance.test", "enterprise_project_id"),
					resource.TestCheckResourceAttrPair(all, "instances.0.status", "huaweicloud_dataarts_studio_instance.test", "status"),
					resource.TestCheckResourceAttrPair(all, "instances.0.vpc_id", "huaweicloud_dataarts_studio_instance.test", "vpc_id"),
					resource.TestCheckResourceAttrPair(all, "instances.0.subnet_id", "huaweicloud_dataarts_studio_instance.test", "subnet_id"),
					resource.TestCheckResourceAttrPair(all, "instances.0.availability_zone",
						"huaweicloud_dataarts_studio_instance.test", "availability_zone"),
					resource.TestMatchResourceAttr(all, "instances.0.effective_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(all, "instances.0.created_by"),
					resource.TestMatchResourceAttr(all, "instances.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(all, "instances.0.workspace_mode"),
				),
			},
		},
	})
}

func testAccDataSourceDataArtsStudioInstances_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dataarts_studio_instances" "before_instance_creation" {}

resource "huaweicloud_dataarts_studio_instance" "test" {
  name                  = "%[2]s"
  version               = "dayu.free"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  period_unit = "month"
  period      = 1
  auto_renew  = "true"
}

data "huaweicloud_dataarts_studio_instances" "test" {
  depends_on = [huaweicloud_dataarts_studio_instance.test]
}
`, common.TestBaseNetwork(name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST), name)
}
