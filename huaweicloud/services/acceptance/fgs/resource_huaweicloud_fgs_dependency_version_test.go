package fgs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/fgs/v2/dependencies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/fgs"
)

func getDependencyVersionFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("fgs", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating FunctionGraph client: %s", err)
	}

	return fgs.GetDependencyVersionById(client, state.Primary.ID)
}

func TestAccDependencyVersion_basic(t *testing.T) {
	var (
		obj dependencies.Dependency

		rName        = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_fgs_dependency_version.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getDependencyVersionFunc)
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
				Config: testAccDependencyVersion_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "runtime", "Python2.7"),
					resource.TestCheckResourceAttr(resourceName, "link", acceptance.HW_FGS_DEPENDENCY_OBS_LINK),
					resource.TestCheckResourceAttrSet(resourceName, "dependency_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"link",
				},
			},
			// Test the ID format: <depend_name>/<version>
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccDependencyVersionImportStateFunc_withDependName(resourceName),
				ImportStateVerifyIgnore: []string{
					"link",
				},
			},
			// Test the ID format: <depend_name>/<version_id>
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccDependencyVersionImportStateFunc_withVersionId(resourceName),
				ImportStateVerifyIgnore: []string{
					"link",
				},
			},
		},
	})
}

func testAccDependencyVersionImportStateFunc_withDependName(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var dependName, version string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		dependName = rs.Primary.Attributes["name"]
		version = rs.Primary.Attributes["version"]
		if dependName == "" || version == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<depend_name>/<version>', but got '%s/%s'",
				dependName, version)
		}
		return fmt.Sprintf("%s/%s", dependName, version), nil
	}
}

func testAccDependencyVersionImportStateFunc_withVersionId(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var dependName, versionId string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		dependName = rs.Primary.Attributes["name"]
		versionId = rs.Primary.Attributes["version_id"]
		if dependName == "" || versionId == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<depend_name>/<version_id>', but got '%s/%s'",
				dependName, versionId)
		}
		return fmt.Sprintf("%s/%s", dependName, versionId), nil
	}
}

func testAccDependencyVersion_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_fgs_dependency_version" "test" {
  name        = "%s"
  description = "Created by terraform script"
  runtime     = "Python2.7"
  link        = "%s"
}
`, rName, acceptance.HW_FGS_DEPENDENCY_OBS_LINK)
}
