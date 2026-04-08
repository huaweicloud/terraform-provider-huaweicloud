package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceWafInstanceTags_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_waf_instance_tags.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
		rName          = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceWafInstanceTags_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.0.value"),
				),
			},
		},
	})
}

func testDataSourceWafInstanceTags_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_dedicated_instance" "test" {
  name               = "%[2]s"
  available_zone     = data.huaweicloud_availability_zones.test.names[1]
  specification_code = "waf.instance.enterprise"
  ecs_flavor         = data.huaweicloud_compute_flavors.test.ids[0]
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  res_tenant         = true
  anti_affinity      = true

  tags = {
    foo = "bar"
    key = "value"
  }

  security_group = [
    huaweicloud_networking_secgroup.test.id
  ]
}

data "huaweicloud_waf_instance_tags" "test" {
  depends_on    = [huaweicloud_waf_dedicated_instance.test]
  resource_type = "waf-instance"
  resource_id   = huaweicloud_waf_dedicated_instance.test.id
}
`, common.TestBaseComputeResources(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
