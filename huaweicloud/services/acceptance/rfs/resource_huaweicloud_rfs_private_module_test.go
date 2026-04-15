package rfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rfs"
)

func getPrivateModuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("rfs", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating RFS client: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("unable to generate UUID: %s", err)
	}

	return rfs.QueryPrivateModule(client, state.Primary.ID, uuid)
}

func TestAccPrivateModule_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_rfs_private_module.test"
		name  = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPrivateModuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateModule_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "module_name", name),
					resource.TestCheckResourceAttr(rName, "module_description", "test module description"),
					resource.TestCheckResourceAttrSet(rName, "module_id"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
				),
			},
			{
				Config: testAccPrivateModule_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "module_name", name),
					resource.TestCheckResourceAttr(rName, "module_description", "test module description update"),
					resource.TestCheckResourceAttrSet(rName, "module_id"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"module_version",
				},
			},
		},
	})
}

func testAccPrivateModule_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rfs_private_module" "test" {
  module_name        = "%s"
  module_version     = "1.2.2"
  module_description = "test module description"
}
`, name)
}

func testAccPrivateModule_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rfs_private_module" "test" {
  module_name        = "%s"
  module_version     = "1.2.2"
  module_description = "test module description update"
}
`, name)
}
