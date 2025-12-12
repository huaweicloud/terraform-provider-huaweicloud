package eps

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/eps/v1/enterpriseprojects"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getResourceEnterpriseProject(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	epsClient, err := conf.EnterpriseProjectClient(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("unable to create EPS client: %s", err)
	}

	return enterpriseprojects.Get(epsClient, state.Primary.ID).Extract()
}

func TestAccEnterpriseProject_basic(t *testing.T) {
	var project enterpriseprojects.Project
	rName := acceptance.RandomAccResourceName()
	updateName := rName + "update"
	resourceName := "huaweicloud_enterprise_project.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&project,
		getResourceEnterpriseProject,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckEnterpriseProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnterpriseProject_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform test"),
					resource.TestCheckResourceAttr(resourceName, "status", "1"),
				),
			},
			{
				Config: testAccEnterpriseProject_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform test update"),
					resource.TestCheckResourceAttr(resourceName, "status", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_flag"},
			},
		},
	})
}

func TestAccEnterpriseProject_delete(t *testing.T) {
	var project enterpriseprojects.Project
	deleteName := acceptance.RandomAccResourceName() + "delete"
	resourceDeleteName := "huaweicloud_enterprise_project.test_delete"

	rc_delete := acceptance.InitResourceCheck(
		resourceDeleteName,
		&project,
		getResourceEnterpriseProject,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckEnterpriseProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnterpriseProject_delete(deleteName),
				Check: resource.ComposeTestCheckFunc(
					rc_delete.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceDeleteName, "name", deleteName),
					resource.TestCheckResourceAttr(resourceDeleteName, "description", "terraform test delete"),
					resource.TestCheckResourceAttr(resourceDeleteName, "status", "1"),
				),
			},
			{
				ResourceName:            resourceDeleteName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_flag"},
			},
		},
	})
}

func testAccCheckEnterpriseProjectDestroy(s *terraform.State) error {
	conf := acceptance.TestAccProvider.Meta().(*config.Config)
	epsClient, err := conf.EnterpriseProjectClient(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("unable to create EPS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_enterprise_project" {
			continue
		}

		project, err := enterpriseprojects.Get(epsClient, rs.Primary.ID).Extract()
		if err == nil {
			if project.Status != 2 {
				return fmt.Errorf("project still active")
			}
		}
	}

	return nil
}

func testAccEnterpriseProject_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_enterprise_project" "test" {
  name        = "%s"
  description = "terraform test"
}`, rName)
}

func testAccEnterpriseProject_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_enterprise_project" "test" {
  name        = "%s"
  description = "terraform test update"
}`, rName)
}

func testAccEnterpriseProject_delete(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_enterprise_project" "test_delete" {
  name        = "%s"
  description = "terraform test delete"
  delete_flag = true
}`, rName)
}
