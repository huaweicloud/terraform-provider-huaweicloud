package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/v3/groups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getIdentityGroupResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.IdentityV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}
	return groups.Get(client, state.Primary.ID).Extract()
}

func TestAccIdentityGroup_basic(t *testing.T) {
	var group groups.Group
	groupName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_identity_group.group_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&group,
		getIdentityGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityGroup_basic(groupName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", groupName),
					resource.TestCheckResourceAttr(resourceName, "description", "An ACC test group"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIdentityGroup_update(groupName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", groupName),
					resource.TestCheckResourceAttr(resourceName, "description", "An ACC update group"),
				),
			},
		},
	})
}

func testAccIdentityGroup_basic(groupName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_group" "group_1" {
  name        = "%s"
  description = "An ACC test group"
}
`, groupName)
}

func testAccIdentityGroup_update(groupName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_group" "group_1" {
  name        = "%s"
  description = "An ACC update group"
}
`, groupName)
}
