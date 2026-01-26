package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/v3/projects"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getProjectResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.IdentityV3ExtClient(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}
	return projects.Get(client, state.Primary.ID).Extract()
}

func TestAccProject_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_identity_project.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getProjectResourceFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// Currently, the DELETE method is not publicly available.
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProject_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", acceptance.HW_REGION_NAME+"_"+name),
					resource.TestCheckResourceAttr(resourceName, "status", "normal"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "parent_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIdentityProject_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", acceptance.HW_REGION_NAME+"_"+name),
					resource.TestCheckResourceAttr(resourceName, "status", "suspended"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "parent_id"),
				),
			},
		},
	})
}

func testAccIdentityProject_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_project" "test" {
  name        = "%[1]s_%[2]s"
  description = "Created by terraform script"
}
`, acceptance.HW_REGION_NAME, name)
}

func testAccIdentityProject_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_project" "test" {
  name        = "%[1]s_%[2]s"
  status      = "suspended"
  description = "Updated by terraform script"
}
`, acceptance.HW_REGION_NAME, name)
}
