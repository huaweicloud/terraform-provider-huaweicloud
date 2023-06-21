package sdrs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/sdrs/v1/drill"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getDrillResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.SdrsV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating SDRS Client: %s", err)
	}
	return drill.Get(client, state.Primary.ID).Extract()
}

func TestAccDrill_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_sdrs_drill.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDrillResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDrill_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "group_id", "huaweicloud_sdrs_protection_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "drill_vpc_id", "huaweicloud_vpc.drill_vpc", "id"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				Config: testDrill_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testDrill_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc" "drill_vpc" {
  name = "%[2]s_drill_vpc"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "drill_vpc_subnet" {
  name       = "%[2]s_drill_vpc_subnet"
  vpc_id     = huaweicloud_vpc.drill_vpc.id
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.2"
}
`, testProtectedInstance_basic(name), name)
}

func testDrill_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_sdrs_drill" "test" {
  name         = "%[2]s"
  group_id     = huaweicloud_sdrs_protection_group.test.id
  drill_vpc_id = huaweicloud_vpc.drill_vpc.id

  depends_on = [
    huaweicloud_sdrs_protected_instance.test,
    huaweicloud_vpc_subnet.drill_vpc_subnet,
  ]
}
`, testDrill_base(name), name)
}

func testDrill_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_sdrs_drill" "test" {
  name         = "%[2]s_update"
  group_id     = huaweicloud_sdrs_protection_group.test.id
  drill_vpc_id = huaweicloud_vpc.drill_vpc.id

  depends_on = [
    huaweicloud_sdrs_protected_instance.test,
    huaweicloud_vpc_subnet.drill_vpc_subnet,
  ]
}
`, testDrill_base(name), name)
}
