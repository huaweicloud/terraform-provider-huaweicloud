package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceGaussdbInstanceDatabaseRole_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_gaussdb_instance_database_role.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGaussdbInstanceDatabaseRole_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "instance_id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrSet(rName, "attribute.#"),
					resource.TestCheckResourceAttrSet(rName, "attribute.0.rolsuper"),
					resource.TestCheckResourceAttrSet(rName, "attribute.0.rolinherit"),
					resource.TestCheckResourceAttrSet(rName, "attribute.0.rolcreaterole"),
					resource.TestCheckResourceAttrSet(rName, "attribute.0.rolcreatedb"),
					resource.TestCheckResourceAttrSet(rName, "attribute.0.rolcanlogin"),
					resource.TestCheckResourceAttrSet(rName, "attribute.0.rolconnlimit"),
					resource.TestCheckResourceAttrSet(rName, "attribute.0.rolreplication"),
					resource.TestCheckResourceAttrSet(rName, "attribute.0.rolbypassrls"),
					resource.TestCheckResourceAttrSet(rName, "memberof"),
					resource.TestCheckResourceAttrSet(rName, "lock_status"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
				ImportStateIdFunc:       testAccResourceGaussdbInstanceDatabaseRoleImportState(rName),
			},
		},
	})
}

func testAccResourceGaussdbInstanceDatabaseRole_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_gaussdb_instance_database_role" "test" {
  instance_id = huaweicloud_gaussdb_instance.test.id
  name        = "%[2]s"
  password    = "Test@74123698"
}
`, testDataSourceGaussdbInstanceMetrics_base(name), name)
}

func testAccResourceGaussdbInstanceDatabaseRoleImportState(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", rName)
		}
		instanceId := rs.Primary.Attributes["instance_id"]
		name := rs.Primary.Attributes["name"]
		return fmt.Sprintf("%s/%s", instanceId, name), nil
	}
}
