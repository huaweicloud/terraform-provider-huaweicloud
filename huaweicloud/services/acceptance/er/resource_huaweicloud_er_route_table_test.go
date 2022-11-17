package er

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/er/v3/routetables"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getRouteTableResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.ErV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ER v3 client: %s", err)
	}

	return routetables.Get(client, state.Primary.Attributes["instance_id"], state.Primary.ID)
}

func TestAccRouteTable_basic(t *testing.T) {
	var (
		obj routetables.RouteTable

		rName      = "huaweicloud_er_route_table.test"
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
		bgpAsNum   = acctest.RandIntRange(64512, 65534)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRouteTableResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckER(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRouteTable_basic(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Create by acc test"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testRouteTable_basic_update(updateName, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", ""),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccRouteTableImportStateFunc(),
			},
		},
	})
}

func testAccRouteTableImportStateFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var instanceId, routeTableId string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "huaweicloud_er_route_table" {
				instanceId = rs.Primary.Attributes["instance_id"]
				routeTableId = rs.Primary.ID
			}
		}
		if instanceId == "" || routeTableId == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<instance_id>/<route_table_id>', but '%s/%s'",
				instanceId, routeTableId)
		}
		return fmt.Sprintf("%s/%s", instanceId, routeTableId), nil
	}
}

func testRouteTable_base(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
resource "huaweicloud_er_instance" "test" {
  availability_zones = ["%[2]s"]
  name               = "%[1]s"
  asn                = %[3]d
}
`, name, acceptance.HW_AVAILABILITY_ZONE, bgpAsNum)
}

func testRouteTable_basic(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_er_route_table" "test" {
  instance_id = huaweicloud_er_instance.test.id
  name        = "%[2]s"
  description = "Create by acc test"

  tags = {
    foo = "bar"
  }
}
`, testRouteTable_base(name, bgpAsNum), name)
}

func testRouteTable_basic_update(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_er_route_table" "test" {
  instance_id = huaweicloud_er_instance.test.id
  name        = "%[2]s"

  tags = {
    foo = "bar"
  }
}
`, testRouteTable_base(name, bgpAsNum), name)
}
