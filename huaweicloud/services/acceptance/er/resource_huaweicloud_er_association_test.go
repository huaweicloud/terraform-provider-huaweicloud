package er

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/er/v3/associations"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/er"
)

func getAssociationResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.ErV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ER v3 client: %s", err)
	}

	return er.QueryAssociationById(client, state.Primary.Attributes["instance_id"],
		state.Primary.Attributes["route_table_id"], state.Primary.ID)
}

func TestAccAssociation_basic(t *testing.T) {
	var (
		obj associations.Association

		rName    = "huaweicloud_er_association.test"
		name     = acceptance.RandomAccResourceName()
		bgpAsNum = acctest.RandIntRange(64512, 65534)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAssociationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAssociation_basic(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_er_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "route_table_id",
						"huaweicloud_er_route_table.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "attachment_id",
						"huaweicloud_er_vpc_attachment.test", "id"),
					resource.TestCheckResourceAttr(rName, "attachment_type", "vpc"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAssociationImportStateFunc(),
			},
		},
	})
}

func testAccAssociationImportStateFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var instanceId, routeTableId, associationId string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "huaweicloud_er_association" {
				instanceId = rs.Primary.Attributes["instance_id"]
				routeTableId = rs.Primary.Attributes["route_table_id"]
				associationId = rs.Primary.ID
			}
		}
		if instanceId == "" || routeTableId == "" || associationId == "" {
			return "", fmt.Errorf("some import IDs are missing, want "+
				"'<instance_id>/<route_table_id>/<association_id>', but '%s/%s/%s'",
				instanceId, routeTableId, associationId)
		}
		return fmt.Sprintf("%s/%s/%s", instanceId, routeTableId, associationId), nil
	}
}

func testAccAssociation_base(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id = huaweicloud_vpc.test.id

  name       = "%[1]s"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1), 1)
}

resource "huaweicloud_er_instance" "test" {
  availability_zones = ["%[2]s"]

  name = "%[1]s"
  asn  = %[3]d
}

resource "huaweicloud_er_vpc_attachment" "test" {
  instance_id = huaweicloud_er_instance.test.id
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id

  name                   = "%[1]s"
  auto_create_vpc_routes = true
}

resource "huaweicloud_er_route_table" "test" {
  instance_id = huaweicloud_er_instance.test.id

  name = "%[1]s"
}
`, name, acceptance.HW_AVAILABILITY_ZONE, bgpAsNum)
}

func testAccAssociation_basic(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_er_association" "test" {
  instance_id    = huaweicloud_er_instance.test.id
  route_table_id = huaweicloud_er_route_table.test.id
  attachment_id  = huaweicloud_er_vpc_attachment.test.id
}
`, testAccAssociation_base(name, bgpAsNum))
}
