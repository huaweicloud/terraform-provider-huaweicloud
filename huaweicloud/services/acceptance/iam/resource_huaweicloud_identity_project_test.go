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

func getIdentityProjectResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.IdentityV3ExtClient(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}
	return projects.Get(client, state.Primary.ID).Extract()
}

func TestAccIdentityProject_basic(t *testing.T) {
	var project projects.Project
	var projectName = acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_identity_project.project_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&project,
		getIdentityProjectResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
			acceptance.TestAccPreCheckProject(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProject_basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPtr(resourceName, "name", &project.Name),
					resource.TestCheckResourceAttr(resourceName, "status", "normal"),
					resource.TestCheckResourceAttr(resourceName, "description", "A project"),
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
				Config: testAccIdentityProject_update(projectName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPtr(resourceName, "name", &project.Name),
					resource.TestCheckResourceAttr(resourceName, "status", "suspended"),
					resource.TestCheckResourceAttr(resourceName, "description", "An updated project"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "parent_id"),
				),
			},
		},
	})
}

func testAccIdentityProject_basic(projectName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_project" "project_1" {
  name        = "%s_%s"
  description = "A project"
}
`, acceptance.HW_REGION_NAME, projectName)
}

func testAccIdentityProject_update(projectName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_project" "project_1" {
  name        = "%s_%s"
  status      = "suspended"
  description = "An updated project"
}
`, acceptance.HW_REGION_NAME, projectName)
}
