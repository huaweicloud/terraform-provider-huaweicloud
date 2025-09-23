package fgs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/fgs"
)

func getDependencyResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.FgsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating FunctionGraph client: %s", err)
	}
	return fgs.GetDependencyById(client, state.Primary.ID)
}

func TestAccDependency_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_fgs_dependency.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getDependencyResourceFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckFgsDependencyLink(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDependency_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "runtime", "Python2.7"),
					resource.TestCheckResourceAttr(resourceName, "link", acceptance.HW_FGS_DEPENDENCY_OBS_LINK),
				),
			},
			{
				Config: testAccDependency_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name+"_update"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "runtime", "Python3.6"),
					resource.TestCheckResourceAttr(resourceName, "link", acceptance.HW_FGS_DEPENDENCY_OBS_LINK),
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

func testAccDependency_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_fgs_dependency" "test" {
  name        = "%[1]s"
  description = "Created by terraform script"
  runtime     = "Python2.7"
  link        = "%[2]s"
}
`, name, acceptance.HW_FGS_DEPENDENCY_OBS_LINK)
}

func testAccDependency_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_fgs_dependency" "test" {
  name        = "%[1]s_update"
  description = "Updated by terraform script" # Does not support empty.
  runtime     = "Python3.6"
  link        = "%[2]s"
}
`, name, acceptance.HW_FGS_DEPENDENCY_OBS_LINK)
}
