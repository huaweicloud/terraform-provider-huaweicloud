package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/iec/v1/security/groups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccIecSecurityGroupResource_basic(t *testing.T) {

	resourceName := "huaweicloud_iec_security_group.my_group"
	rName := fmt.Sprintf("iec-secgroup-%s", acctest.RandString(5))
	description := "This is a test of iec security group"

	var group groups.RespSecurityGroupEntity

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIecSecurityGroupV1Destory,
		Steps: []resource.TestStep{
			{
				Config: testAccIecSecurityGroupV1_Basic(rName, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIecSecurityGroupV1Exists(resourceName, &group),
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

func testAccCheckIecSecurityGroupV1RuleCount(group *groups.RespSecurityGroupEntity, count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(group.SecurityGroupRules) == count {
			return nil
		}

		return fmtp.Errorf("Unexpected number of rules in group %s. Expected %d, got %d",
			group.ID, count, len(group.SecurityGroupRules))
	}
}

func testAccCheckIecSecurityGroupV1Exists(n string, group *groups.RespSecurityGroupEntity) resource.TestCheckFunc {

	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID has been seted")
		}

		config := testAccProvider.Meta().(*config.Config)
		iecClient, err := config.IECV1Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating Huaweicloud IEC client: %s", err)
		}

		found, err := groups.Get(iecClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("IEC Security group not found")
		}
		*group = *found
		return nil
	}
}

func testAccCheckIecSecurityGroupV1Destory(s *terraform.State) error {

	config := testAccProvider.Meta().(*config.Config)
	iecClient, err := config.IECV1Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud IEC client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_iec_security_group" {
			continue
		}
		_, err := groups.Get(iecClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("IEC Security group still exists")
		}
	}

	return nil
}

func testAccIecSecurityGroupV1_Basic(rName, description string) string {
	return fmt.Sprintf(`
resource "huaweicloud_iec_security_group" "my_group" {
  name        = "%s"
  description = "%s"
}
`, rName, description)
}
