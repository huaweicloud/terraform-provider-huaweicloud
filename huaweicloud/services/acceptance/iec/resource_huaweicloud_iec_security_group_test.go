package iec

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/iec/v1/security/groups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSecurityGroupResource_basic(t *testing.T) {
	resourceName := "huaweicloud_iec_security_group.my_group"
	rName := fmt.Sprintf("iec-secgroup-%s", acctest.RandString(5))
	description := "This is a test of iec security group"

	var group groups.RespSecurityGroupEntity

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroup_Basic(rName, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "security_group_rules.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSecurityGroupExists(n string, group *groups.RespSecurityGroupEntity) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID has been seted")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		iecClient, err := config.IECV1Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating IEC client: %s", err)
		}

		found, err := groups.Get(iecClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("IEC security group not found")
		}
		*group = *found
		return nil
	}
}

func testAccCheckSecurityGroupDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	iecClient, err := cfg.IECV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating IEC client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_iec_security_group" {
			continue
		}
		_, err := groups.Get(iecClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("IEC security group still exists")
		}
	}

	return nil
}

func testAccSecurityGroup_Basic(rName, description string) string {
	return fmt.Sprintf(`
resource "huaweicloud_iec_security_group" "my_group" {
  name        = "%s"
  description = "%s"
}
`, rName, description)
}
