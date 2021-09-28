package iam

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/v3/groups"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccIdentityV3Group_basic(t *testing.T) {
	var group groups.Group
	var groupName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_identity_group.group_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckIdentityV3GroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3Group_basic(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3GroupExists(resourceName, &group),
					resource.TestCheckResourceAttrPtr(resourceName, "name", &group.Name),
					resource.TestCheckResourceAttrPtr(resourceName, "description", &group.Description),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIdentityV3Group_update(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3GroupExists(resourceName, &group),
					resource.TestCheckResourceAttrPtr(resourceName, "name", &group.Name),
					resource.TestCheckResourceAttrPtr(resourceName, "description", &group.Description),
				),
			},
		},
	})
}

func testAccCheckIdentityV3GroupDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	identityClient, err := config.IdentityV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud identity client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_identity_group" {
			continue
		}

		_, err := groups.Get(identityClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Group still exists")
		}
	}

	return nil
}

func testAccCheckIdentityV3GroupExists(n string, group *groups.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		identityClient, err := config.IdentityV3Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud identity client: %s", err)
		}

		found, err := groups.Get(identityClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("Group not found")
		}

		*group = *found

		return nil
	}
}

func testAccIdentityV3Group_basic(groupName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_group" "group_1" {
  name        = "%s"
  description = "A ACC test group"
}
`, groupName)
}

func testAccIdentityV3Group_update(groupName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_group" "group_1" {
  name        = "%s"
  description = "Some Group"
}
`, groupName)
}
