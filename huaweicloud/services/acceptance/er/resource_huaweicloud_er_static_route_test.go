package er

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/er/v3/routes"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getStaticRouteFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ErV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ER v3 client: %s", err)
	}

	return routes.Get(client, state.Primary.Attributes["route_table_id"], state.Primary.ID)
}

func TestAccStaticRoute_basic(t *testing.T) {
	var (
		obj routes.Route

		sourceSelfResName = "huaweicloud_er_static_route.source_self"
		destSelfResName   = "huaweicloud_er_static_route.destination_self"
		crossVpcResName   = "huaweicloud_er_static_route.cross_vpc"
		name              = acceptance.RandomAccResourceName()
		bgpAsNum          = acctest.RandIntRange(64512, 65534)

		sourceSelfRes = acceptance.InitResourceCheck(sourceSelfResName, &obj, getStaticRouteFunc)
		destSelfRes   = acceptance.InitResourceCheck(destSelfResName, &obj, getStaticRouteFunc)
		crossVpcRes   = acceptance.InitResourceCheck(crossVpcResName, &obj, getStaticRouteFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      sourceSelfRes.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccStaticRoute_basic_step1(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					sourceSelfRes.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(sourceSelfResName, "route_table_id",
						"huaweicloud_er_route_table.source", "id"),
					resource.TestCheckResourceAttrPair(sourceSelfResName, "destination",
						"huaweicloud_vpc.source", "cidr"),
					resource.TestCheckResourceAttrPair(sourceSelfResName, "attachment_id",
						"huaweicloud_er_vpc_attachment.source", "id"),
					resource.TestCheckResourceAttrSet(sourceSelfResName, "type"),
					resource.TestCheckResourceAttrSet(sourceSelfResName, "status"),
					resource.TestMatchResourceAttr(sourceSelfResName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(sourceSelfResName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					destSelfRes.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(destSelfResName, "route_table_id",
						"huaweicloud_er_route_table.destination", "id"),
					resource.TestCheckResourceAttrPair(destSelfResName, "destination",
						"huaweicloud_vpc.destination", "cidr"),
					resource.TestCheckResourceAttrPair(destSelfResName, "attachment_id",
						"huaweicloud_er_vpc_attachment.destination", "id"),
					crossVpcRes.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(crossVpcResName, "route_table_id",
						"huaweicloud_er_route_table.source", "id"),
					resource.TestCheckResourceAttrPair(crossVpcResName, "destination",
						"huaweicloud_vpc.destination", "cidr"),
					resource.TestCheckResourceAttrPair(crossVpcResName, "attachment_id",
						"huaweicloud_er_vpc_attachment.source", "id"),
				),
			},
			{
				Config: testAccStaticRoute_basic_step2(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					sourceSelfRes.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(sourceSelfResName, "attachment_id",
						"huaweicloud_er_vpc_attachment.destination", "id"),
					destSelfRes.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(destSelfResName, "attachment_id",
						"huaweicloud_er_vpc_attachment.source", "id"),
					crossVpcRes.CheckResourceExists(),
					resource.TestCheckResourceAttr(crossVpcResName, "is_blackhole", "true"),
				),
			},
			{
				ResourceName:      sourceSelfResName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStaticRouteImportStateFunc(sourceSelfResName),
			},
			{
				ResourceName:      destSelfResName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStaticRouteImportStateFunc(destSelfResName),
			},
			{
				ResourceName:      crossVpcResName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStaticRouteImportStateFunc(crossVpcResName),
			},
		},
	})
}

func testAccStaticRouteImportStateFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var routeTableId, staticRouteId string
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) of ER static route is not found in the tfstate", rsName)
		}
		routeTableId = rs.Primary.Attributes["route_table_id"]
		staticRouteId = rs.Primary.ID
		if routeTableId == "" || staticRouteId == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<route_table_id>/<id>', but got '%s/%s'",
				routeTableId, staticRouteId)
		}
		return fmt.Sprintf("%s/%s", routeTableId, staticRouteId), nil
	}
}

func testAccStaticRoute_base(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
data "huaweicloud_er_availability_zones" "test" {}

variable "base_vpc_cidr" {
  type    = string
  default = "192.168.0.0/16"
}

resource "huaweicloud_vpc" "source" {
  name = "%[1]s_source"
  cidr = cidrsubnet(var.base_vpc_cidr, 2, 1)
}

resource "huaweicloud_vpc" "destination" {
  name = "%[1]s_destination"
  cidr = cidrsubnet(var.base_vpc_cidr, 2, 2)
}

resource "huaweicloud_vpc_subnet" "source" {
  vpc_id = huaweicloud_vpc.source.id

  name       = "%[1]s_source"
  cidr       = cidrsubnet(huaweicloud_vpc.source.cidr, 2, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.source.cidr, 2, 1), 1)
}

resource "huaweicloud_vpc_subnet" "destination" {
  vpc_id = huaweicloud_vpc.destination.id

  name       = "%[1]s_destination"
  cidr       = cidrsubnet(huaweicloud_vpc.destination.cidr, 2, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.destination.cidr, 2, 1), 1)
}

resource "huaweicloud_er_instance" "test" {
  availability_zones = slice(data.huaweicloud_er_availability_zones.test.names, 0, 1)
  name               = "%[1]s"
  asn                = %[2]d
}

resource "huaweicloud_er_route_table" "source" {
  instance_id = huaweicloud_er_instance.test.id
  name        = "%[1]s_source"
}

resource "huaweicloud_er_route_table" "destination" {
  instance_id = huaweicloud_er_instance.test.id
  name        = "%[1]s_destination"
}

resource "huaweicloud_er_vpc_attachment" "source" {
  instance_id = huaweicloud_er_instance.test.id
  vpc_id      = huaweicloud_vpc.source.id
  subnet_id   = huaweicloud_vpc_subnet.source.id
  name        = "%[1]s_source"
}

resource "huaweicloud_er_vpc_attachment" "destination" {
  instance_id = huaweicloud_er_instance.test.id
  vpc_id      = huaweicloud_vpc.destination.id
  subnet_id   = huaweicloud_vpc_subnet.destination.id
  name        = "%[1]s_destination"
}
`, name, bgpAsNum)
}

func testAccStaticRoute_basic_step1(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_er_static_route" "source_self" {
  route_table_id = huaweicloud_er_route_table.source.id
  destination    = huaweicloud_vpc.source.cidr
  attachment_id  = huaweicloud_er_vpc_attachment.source.id
}

resource "huaweicloud_er_static_route" "destination_self" {
  route_table_id = huaweicloud_er_route_table.destination.id
  destination    = huaweicloud_vpc.destination.cidr
  attachment_id  = huaweicloud_er_vpc_attachment.destination.id
}

resource "huaweicloud_er_static_route" "cross_vpc" {
  route_table_id = huaweicloud_er_route_table.source.id
  destination    = huaweicloud_vpc.destination.cidr
  attachment_id  = huaweicloud_er_vpc_attachment.source.id
}
`, testAccStaticRoute_base(name, bgpAsNum))
}

func testAccStaticRoute_basic_step2(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%[1]s

// Update the VPC attachment ID.
resource "huaweicloud_er_static_route" "source_self" {
  route_table_id = huaweicloud_er_route_table.source.id
  destination    = huaweicloud_vpc.source.cidr
  attachment_id  = huaweicloud_er_vpc_attachment.destination.id
}

// Update the route destination CIDR.
resource "huaweicloud_er_static_route" "destination_self" {
  route_table_id = huaweicloud_er_route_table.destination.id
  destination    = huaweicloud_vpc.destination.cidr
  attachment_id  = huaweicloud_er_vpc_attachment.source.id
}

// Change the static route to the black hole route.
resource "huaweicloud_er_static_route" "cross_vpc" {
  route_table_id = huaweicloud_er_route_table.source.id
  destination    = huaweicloud_vpc.destination.cidr
  is_blackhole   = true
}
`, testAccStaticRoute_base(name, bgpAsNum))
}
