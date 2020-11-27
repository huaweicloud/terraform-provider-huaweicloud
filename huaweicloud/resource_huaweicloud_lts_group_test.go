package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/lts/huawei/loggroups"
)

func TestAccLogTankGroupV2_basic(t *testing.T) {
	var group loggroups.LogGroup

	resourceName := "huaweicloud_lts_group.testacc_group"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLogTankGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLogTankGroupV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogTankGroupV2Exists(
						resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "group_name", "testacc_group"),
					resource.TestCheckResourceAttr(
						resourceName, "ttl_in_days", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLogTankGroupV2_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "group_name", "testacc_group"),
					resource.TestCheckResourceAttr(
						resourceName, "ttl_in_days", "7"),
				),
			},
		},
	})
}

func testAccCheckLogTankGroupV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	ltsclient, err := config.ltsV2Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud LTS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_lts_group" {
			continue
		}

		groups, err := loggroups.List(ltsclient).Extract()
		if err != nil {
			return fmt.Errorf("Log group get list err: %s", err.Error())
		}
		for _, group := range groups.LogGroups {
			if group.ID == rs.Primary.ID {
				return fmt.Errorf("Log group (%s) still exists.", rs.Primary.ID)
			}
		}

	}
	return nil
}

func testAccCheckLogTankGroupV2Exists(n string, group *loggroups.LogGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		ltsclient, err := config.ltsV2Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud LTS client: %s", err)
		}

		var founds *loggroups.LogGroups
		founds, err = loggroups.List(ltsclient).Extract()
		if err != nil {
			return err
		}
		for _, loggroup := range founds.LogGroups {
			if rs.Primary.ID == loggroup.ID {
				*group = loggroup
				return nil
			}
		}

		return fmt.Errorf("Error HuaweiCloud log group %s: No Found", rs.Primary.ID)
	}
}

const testAccLogTankGroupV2_basic = `
resource "huaweicloud_lts_group" "testacc_group" {
	group_name  = "testacc_group"
	ttl_in_days = 1
}
`

const testAccLogTankGroupV2_update = `
resource "huaweicloud_lts_group" "testacc_group" {
	group_name  = "testacc_group"
	ttl_in_days = 7
}
`
