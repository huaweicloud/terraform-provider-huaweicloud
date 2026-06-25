package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussDBInstanceDatabaseRoles_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_instance_database_roles.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGaussDBInstanceDatabaseRoles_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "roles.#"),
					resource.TestCheckResourceAttrSet(dataSource, "roles.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "roles.0.memberof"),
					resource.TestCheckResourceAttrSet(dataSource, "roles.0.lock_status"),
					resource.TestCheckResourceAttrSet(dataSource, "roles.0.attribute.#"),
					resource.TestCheckResourceAttrSet(dataSource, "roles.0.attribute.0.rolsuper"),
					resource.TestCheckResourceAttrSet(dataSource, "roles.0.attribute.0.rolinherit"),
					resource.TestCheckResourceAttrSet(dataSource, "roles.0.attribute.0.rolcreaterole"),
					resource.TestCheckResourceAttrSet(dataSource, "roles.0.attribute.0.rolcreatedb"),
					resource.TestCheckResourceAttrSet(dataSource, "roles.0.attribute.0.rolcanlogin"),
					resource.TestCheckResourceAttrSet(dataSource, "roles.0.attribute.0.rolconnlimit"),
					resource.TestCheckResourceAttrSet(dataSource, "roles.0.attribute.0.rolreplication"),
					resource.TestCheckResourceAttrSet(dataSource, "roles.0.attribute.0.rolbypassrls"),
				),
			},
		},
	})
}

func testAccDataSourceGaussDBInstanceDatabaseRoles_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_gaussdb_instance_database_role" "test" {
  instance_id = huaweicloud_gaussdb_instance.test.id
  name        = "%[2]s"
  password    = "Test@74123698"
}

data "huaweicloud_gaussdb_instance_database_roles" "test" {
  instance_id = huaweicloud_gaussdb_instance.test.id
}
`, testDataSourceGaussdbInstanceMetrics_base(name), name)
}
