package rfs

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rfs"
)

func getPrivateModuleVersionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region        = acceptance.HW_REGION_NAME
		moduleName    = state.Primary.Attributes["module_name"]
		moduleVersion = state.Primary.Attributes["module_version"]
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return nil, fmt.Errorf("error creating RFS client: %s", err)
	}

	requestId, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("error creating RFS client: %s", err)
	}

	return rfs.QueryPrivateModuleVersion(client, moduleName, moduleVersion, requestId.String())
}
func TestAccPrivateModuleVersion_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_rfs_private_module_version.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPrivateModuleVersionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRfsModuleURI(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPrivateModuleVersion_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "module_name", "huaweicloud_rfs_private_module.test", "module_name"),
					resource.TestCheckResourceAttrSet(rName, "module_version"),
					resource.TestCheckResourceAttrSet(rName, "module_uri"),
					resource.TestCheckResourceAttrSet(rName, "version_description"),
					resource.TestCheckResourceAttrSet(rName, "module_id"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"module_uri"},
				ImportStateIdFunc:       testPrivateModuleVersionImportState(rName),
			},
		},
	})
}

func testPrivateModuleVersion_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rfs_private_module" "test" {
  module_name        = "%[1]s"
  module_description = "acc private module for version test"
}
`, name)
}

func testPrivateModuleVersion_basic(name string) string {
	return fmt.Sprintf(`
%[1]s
resource "huaweicloud_rfs_private_module_version" "test" {
  module_name         = huaweicloud_rfs_private_module.test.module_name
  module_version      = "3.0.0"
  module_uri          = "%[2]s"
  version_description = "private module version test"
}
`, testPrivateModuleVersion_base(name), acceptance.HW_RFS_MODULE_URI)
}

func testPrivateModuleVersionImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		moduleName := rs.Primary.Attributes["module_name"]
		moduleVersion := rs.Primary.Attributes["module_version"]
		if moduleName == "" || moduleVersion == "" {
			return "", fmt.Errorf("the module_name (%s) or module_version (%s) is nil",
				moduleName, moduleVersion)
		}

		return moduleName + "/" + moduleVersion, nil
	}
}
