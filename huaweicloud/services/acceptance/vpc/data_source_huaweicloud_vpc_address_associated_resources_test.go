package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpcAddressAssociatedResources_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpc_address_group_associated_resources.test"
	dc := acceptance.InitDataSourceCheck(dataSource)
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			// 测试用例1: 无过滤参数
			{
				Config: testAccDataSourceVpcAddressAssociatedResources_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "address_groups.#"),
				),
			},
			// 测试用例2: 按ID过滤
			{
				Config: testAccDataSourceVpcAddressAssociatedResources_filterById(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "address_groups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "address_groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "address_groups.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "address_groups.0.dependency.#"),
					resource.TestCheckResourceAttrSet(dataSource, "address_groups.0.dependency.0.instance_id"),
					resource.TestCheckResourceAttr(dataSource, "address_groups.0.dependency.0.instance_name", rName+"-secgroup"),
					resource.TestCheckResourceAttr(dataSource, "address_groups.0.dependency.0.type", "sg"),
				),
			},
			// 测试用例3: 按企业项目ID过滤
			{
				Config: testAccDataSourceVpcAddressAssociatedResources_filterByEnterpriseProject(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "address_groups.#"),
				),
			},
			// 测试用例4: 同时按ID和企业项目ID过滤
			{
				Config: testAccDataSourceVpcAddressAssociatedResources_filterByAll(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "address_groups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "address_groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "address_groups.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "address_groups.0.dependency.#"),
					resource.TestCheckResourceAttrSet(dataSource, "address_groups.0.dependency.0.instance_id"),
					resource.TestCheckResourceAttr(dataSource, "address_groups.0.dependency.0.instance_name", rName+"-secgroup"),
					resource.TestCheckResourceAttr(dataSource, "address_groups.0.dependency.0.type", "sg"),
				),
			},
		},
	})
}

func testAccDataSourceVpcAddressAssociatedResources_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpc_address_group_associated_resources" "test" {
  ip_address_group_id  = ""
  enterprise_project_id =""

  depends_on = [huaweicloud_networking_secgroup_rule.test]
}
`, testAccNetworkingSecGroupWithAddressGroup(name))
}

func testAccDataSourceVpcAddressAssociatedResources_filterByEnterpriseProject(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_vpc_address_group_associated_resources" "test" {
  enterprise_project_id = "%[2]s"
  
  depends_on = [huaweicloud_networking_secgroup_rule.test]
}
`, testAccNetworkingSecGroupWithAddressGroup(name), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccDataSourceVpcAddressAssociatedResources_filterById(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpc_address_group_associated_resources" "test" {
  ip_address_group_id = huaweicloud_vpc_address_group.test.id

  depends_on = [huaweicloud_networking_secgroup_rule.test]
}
`, testAccNetworkingSecGroupWithAddressGroup(name))
}

func testAccDataSourceVpcAddressAssociatedResources_filterByAll(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_vpc_address_group_associated_resources" "test" {
  ip_address_group_id   = huaweicloud_vpc_address_group.test.id
  enterprise_project_id = "%[2]s"

  depends_on = [huaweicloud_networking_secgroup_rule.test]
}
`, testAccNetworkingSecGroupWithAddressGroup(name), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccNetworkingSecGroupWithAddressGroup(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "test" {
  name                  = "%[1]s-secgroup"
  description           = "terraform security group rule acceptance test"
  enterprise_project_id = "%[2]s"
}

resource "huaweicloud_vpc_address_group" "test" {
  name                  = "%[1]s-address-group"
  enterprise_project_id = "%[2]s"
  addresses             = [
    "192.168.10.12",
    "192.168.11.0-192.168.11.240",
  ]
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  security_group_id       = huaweicloud_networking_secgroup.test.id
  direction               = "ingress"
  ethertype               = "IPv4"
  ports                   = 80
  protocol                = "tcp"
  remote_address_group_id = huaweicloud_vpc_address_group.test.id
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
