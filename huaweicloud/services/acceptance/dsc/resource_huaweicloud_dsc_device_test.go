package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dsc"
)

func getDscDeviceResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dsc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DSC client: %s", err)
	}

	return dsc.GetDeviceById(client, state.Primary.ID)
}

func TestAccDscDevice_basic(t *testing.T) {
	var (
		device interface{}
		rName  = "huaweicloud_dsc_device.test"
		name   = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&device,
		getDscDeviceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDevice_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "0"),
					resource.TestCheckResourceAttr(rName, "mode", "SINGLE"),
					resource.TestCheckResourceAttrPair(
						rName, "vpc_id",
						"huaweicloud_vpc.test", "id",
					),
					resource.TestCheckResourceAttrPair(
						rName, "subnet_id",
						"huaweicloud_vpc_subnet.test", "id",
					),
					resource.TestCheckResourceAttr(rName, "description", "demo description"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
				),
			},
			{
				Config: testAccDevice_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s-update", name)),
					resource.TestCheckResourceAttr(rName, "type", "0"),
					resource.TestCheckResourceAttr(rName, "mode", "SINGLE"),
					resource.TestCheckResourceAttrPair(
						rName, "vpc_id",
						"huaweicloud_vpc.test", "id",
					),
					resource.TestCheckResourceAttrPair(
						rName, "subnet_id",
						"huaweicloud_vpc_subnet.test", "id",
					),
					resource.TestCheckResourceAttr(rName, "description", "updated description"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
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

func testAccDevice_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dsc_device" "test" {
  name        = "%s"
  type        = 0
  mode        = "SINGLE"
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id
  description = "demo description"
}
`, common.TestVpc(name), name)
}

func testAccDevice_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dsc_device" "test" {
  name        = "%s-update"
  type        = 0
  mode        = "SINGLE"
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id
  description = "updated description"
}
`, common.TestVpc(name), name)
}
