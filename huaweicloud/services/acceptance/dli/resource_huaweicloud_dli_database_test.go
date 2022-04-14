package dli

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/dli/v1/databases"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dli"
)

func getDatabaseResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.DliV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud DLI v1 client: %s", err)
	}

	return dli.GetDliSqlDatabaseByName(c, state.Primary.ID)
}

func TestAccDliDatabase_basic(t *testing.T) {
	var database databases.Database

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dli_database.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&database,
		getDatabaseResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDliDatabase_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "For terraform acc test"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttrSet(resourceName, "owner"),
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

func testAccDliDatabase_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_database" "test" {
  name                  = "%s"
  description           = "For terraform acc test"
  enterprise_project_id = "%s"
}
`, rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
