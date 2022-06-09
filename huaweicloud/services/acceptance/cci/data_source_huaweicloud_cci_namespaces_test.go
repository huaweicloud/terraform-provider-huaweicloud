package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataCciNamespaces_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_cci_namespaces.test"
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataCciNamespaces_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "namespaces.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "namespaces.0.auto_expend_enabled", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "namespaces.0.container_network_enabled", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "namespaces.0.rbac_enabled", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "namespaces.0.status", "Active"),
					resource.TestCheckResourceAttr(dataSourceName, "namespaces.0.type", "general-computing"),
					resource.TestCheckResourceAttrSet(dataSourceName, "namespaces.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "namespaces.0.recycling_interval"),
					resource.TestCheckResourceAttrSet(dataSourceName, "namespaces.0.warmup_pool_size"),
					resource.TestCheckResourceAttrPair(dataSourceName, "namespaces.0.network.0.name",
						"huaweicloud_cci_network.test", "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "namespaces.0.network.0.security_group_id",
						"huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "namespaces.0.network.0.vpc.0.id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "namespaces.0.network.0.vpc.0.subnet_id",
						"huaweicloud_vpc_subnet.test", "subnet_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "namespaces.0.network.0.vpc.0.subnet_cidr",
						"huaweicloud_vpc_subnet.test", "cidr"),
					resource.TestCheckResourceAttrPair(dataSourceName, "namespaces.0.network.0.vpc.0.network_id",
						"huaweicloud_vpc_subnet.test", "id"),
				),
			},
		},
	})
}

func TestAccDataCciNamespaces_noNetwork(t *testing.T) {
	dataSourceName := "data.huaweicloud_cci_namespaces.test"
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataCciNamespaces_noNetwork(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "namespaces.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "namespaces.0.name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "namespaces.0.type", "general-computing"),
				),
			},
		},
	})
}

func testAccDataCciNamespaces_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc.test.id
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%s"
}

resource "huaweicloud_cci_namespace" "test" {
  name = "%s"
  type = "general-computing"
}

resource "huaweicloud_cci_network" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  namespace         = huaweicloud_cci_namespace.test.name
  name              = "%s"
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
}
`, rName, rName, rName, rName, rName)
}

func testAccDataCciNamespaces_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cci_namespaces" "test" {
  depends_on = [huaweicloud_cci_network.test]

  name = "%s"
  type = "general-computing"
}
`, testAccDataCciNamespaces_base(rName), rName)
}

func testAccDataCciNamespaces_noNetwork(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cci_namespace" "test" {
  name = "%[1]s"
  type = "general-computing"
}

data "huaweicloud_cci_namespaces" "test" {
  name = huaweicloud_cci_namespace.test.name
}
`, rName)
}
