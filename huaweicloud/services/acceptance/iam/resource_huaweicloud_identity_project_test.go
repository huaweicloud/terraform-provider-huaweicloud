package iam

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/v3/projects"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getIdentityProjectResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.IAMV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("Error creating HuaweiCloud IAM client: %s", err)
	}
	return projects.Get(client, state.Primary.ID).Extract()
}

func TestAccIdentityV3Project_basic(t *testing.T) {
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
				Config: testAccIdentityV3Project_basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPtr(resourceName, "name", &project.Name),
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
				Config: testAccIdentityV3Project_update(projectName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPtr(resourceName, "name", &project.Name),
					resource.TestCheckResourceAttr(resourceName, "description", "An updated project"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "parent_id"),
				),
			},
		},
	})
}

func testAccIdentityV3Project_basic(projectName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_project" "project_1" {
  name        = "%s_%s"
  description = "A project"
}
`, acceptance.HW_REGION_NAME, projectName)
}

func testAccIdentityV3Project_update(projectName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_project" "project_1" {
  name        = "%s_%s"
  description = "An updated project"
}
`, acceptance.HW_REGION_NAME, projectName)
}
