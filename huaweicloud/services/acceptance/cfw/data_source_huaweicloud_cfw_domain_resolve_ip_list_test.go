package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCfwDomainResolveIpList_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cfw_domain_resolve_ip_list.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCfwDomainResolveIpList_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.excess_ip.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.parsed_success_ip.#"),
				),
			},
		},
	})
}

func testDataSourceCfwDomainResolveIpList_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_firewalls" "test" {
  fw_instance_id = "%s"
}

data "huaweicloud_cfw_domain_name_groups" "test" {
  depends_on     = [data.huaweicloud_cfw_firewalls.test]
  fw_instance_id = data.huaweicloud_cfw_firewalls.test.records[0].fw_instance_id
  object_id      = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
}

data "huaweicloud_cfw_domain_resolve_ip_list" "test" {
  depends_on        = [data.huaweicloud_cfw_firewalls.test, data.huaweicloud_cfw_domain_name_groups.test]
  fw_instance_id    = data.huaweicloud_cfw_firewalls.test.records[0].fw_instance_id
  domain_address_id = data.huaweicloud_cfw_domain_name_groups.test.records[0].domain_names[0].domain_address_id
}
`, acceptance.HW_CFW_INSTANCE_ID)
}
