package er

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceAvailableRoutes_basic(t *testing.T) {
	var (
		baseConfig = testAccDataSourceAvailableRoutes_base()

		all = "data.huaweicloud_er_available_routes.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byDestination   = "data.huaweicloud_er_available_routes.filter_by_destination"
		dcByDestination = acceptance.InitDataSourceCheck(byDestination)

		byResourceType   = "data.huaweicloud_er_available_routes.filter_by_resource_type"
		dcByResourceType = acceptance.InitDataSourceCheck(byResourceType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAvailableRoutes_basic(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "routes.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),

					dcByDestination.CheckResourceExists(),
					resource.TestCheckOutput("is_destination_filter_useful", "true"),
					resource.TestMatchResourceAttr(byDestination, "routes.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "routes.0.id"),
					resource.TestCheckResourceAttrSet(all, "routes.0.destination"),
					resource.TestMatchResourceAttr(all, "routes.0.next_hops.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrPair(all, "routes.0.next_hops.0.attachment_id",
						"huaweicloud_er_vpc_attachment.test", "id"),
					resource.TestCheckResourceAttrPair(all, "routes.0.next_hops.0.resource_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttr(all, "routes.0.next_hops.0.resource_type", "vpc"),
					resource.TestCheckResourceAttrSet(all, "routes.0.is_blackhole"),
					resource.TestCheckResourceAttrSet(all, "routes.0.type"),

					dcByResourceType.CheckResourceExists(),
					resource.TestCheckOutput("is_resource_type_filter_useful", "true"),
				),
			},
			{
				Config:      testAccDataSourceAvailableRoutes_routeTableNotFound(baseConfig),
				ExpectError: regexp.MustCompile(`route table [a-f0-9-]+ not found`),
			},
		},
	})
}

func testAccDataSourceAvailableRoutes_base() string {
	var (
		name     = acceptance.RandomAccResourceName()
		bgpAsNum = acctest.RandIntRange(64512, 65534)
	)

	return fmt.Sprintf(`
data "huaweicloud_er_availability_zones" "test" {}

%[1]s

resource "huaweicloud_er_instance" "test" {
  availability_zones = slice(data.huaweicloud_er_availability_zones.test.names, 0, 1)

  name = "%[2]s"
  asn  = %[3]d
}

resource "huaweicloud_er_route_table" "test" {
  instance_id = huaweicloud_er_instance.test.id
  name        = "%[2]s"
}

resource "huaweicloud_er_vpc_attachment" "test" {
  instance_id = huaweicloud_er_instance.test.id
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id

  name                   = "%[2]s"
  auto_create_vpc_routes = true
}

resource "huaweicloud_er_static_route" "test" {
  route_table_id = huaweicloud_er_route_table.test.id
  destination    = huaweicloud_vpc.test.cidr
  attachment_id  = huaweicloud_er_vpc_attachment.test.id
}
`, common.TestVpc(name), name, bgpAsNum)
}

func testAccDataSourceAvailableRoutes_basic(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_er_available_routes" "test" {
  depends_on = [huaweicloud_er_static_route.test]

  route_table_id = huaweicloud_er_route_table.test.id
}

# Filter by destination
locals {
  dest_cidr = huaweicloud_er_static_route.test.destination
}

data "huaweicloud_er_available_routes" "filter_by_destination" {
  depends_on = [huaweicloud_er_static_route.test]

  route_table_id = huaweicloud_er_route_table.test.id
  destination    = local.dest_cidr
}

locals {
  destination_filter_result = [for _, v in data.huaweicloud_er_available_routes.filter_by_destination.routes[*].destination: v == local.dest_cidr]
}

output "is_destination_filter_useful" {
  value = length(local.destination_filter_result) >0 && alltrue(local.destination_filter_result)
}

# Filter by resource type
locals {
  resource_type = "vpc"
}

data "huaweicloud_er_available_routes" "filter_by_resource_type" {
  depends_on = [huaweicloud_er_static_route.test]

  route_table_id = huaweicloud_er_route_table.test.id
  resource_type  = local.resource_type
}

locals {
  resource_type_filter_result = [for _, v in data.huaweicloud_er_available_routes.filter_by_destination.routes[*].next_hops:
    contains(v[*].resource_type, local.resource_type)]
}

output "is_resource_type_filter_useful" {
  value = length(local.resource_type_filter_result) >0 && alltrue(local.resource_type_filter_result)
}
`, baseConfig)
}

func testAccDataSourceAvailableRoutes_routeTableNotFound(baseConfig string) string {
	randUUID, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_er_available_routes" "test" {
  route_table_id = "%[2]s"
}  
`, baseConfig, randUUID)
}
