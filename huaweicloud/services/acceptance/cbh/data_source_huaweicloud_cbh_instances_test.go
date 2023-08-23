package cbh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceCbhInstances_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_cbh_instances.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceCbhInstances_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "instances.0.name", name),
					resource.TestCheckResourceAttr(dataSourceName, "instances.0.flavor_id", "cbh.basic.10"),
					resource.TestCheckResourceAttr(dataSourceName, "instances.0.status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(dataSourceName, "instances.0.vpc_id",
						"data.huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "instances.0.subnet_id",
						"data.huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "instances.0.security_group_id",
						"data.huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "instances.0.public_ip_id",
						"huaweicloud_vpc_eip.test", "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "instances.0.public_ip",
						"huaweicloud_vpc_eip.test", "address"),
				),
			},
		},
	})
}

func testAccDatasourceCbhInstances_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cbh_instances" "test" {
  name              = huaweicloud_cbh_instance.test.name
  vpc_id            = data.huaweicloud_vpc.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  flavor_id         = "cbh.basic.10"
}
`, testCBHInstance_basic(name))
}
