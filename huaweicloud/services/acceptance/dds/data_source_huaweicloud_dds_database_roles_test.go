package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDdsDatabaseRoles_basic(t *testing.T) {
	rName := "data.huaweicloud_dds_database_roles.all"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDdsDatabaseRoles_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "roles.#"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_db_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceDdsDatabaseRoles_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dds_database_roles" "all" {
  depends_on = [huaweicloud_dds_database_role.test]

  instance_id = huaweicloud_dds_instance.test.id
}

// filter by name
data "huaweicloud_dds_database_roles" "filter_by_name" {
  depends_on = [huaweicloud_dds_database_role.test]

  instance_id = huaweicloud_dds_instance.test.id
  name        = huaweicloud_dds_database_role.test.name
}

locals {
  filter_result_by_name = [for v in data.huaweicloud_dds_database_roles.filter_by_name.roles[*].name : 
    v == huaweicloud_dds_database_role.test.name]
}

output "is_name_filter_useful" {
  value = length(local.filter_result_by_name) > 0 && alltrue(local.filter_result_by_name) 
}

// filter by db_name
data "huaweicloud_dds_database_roles" "filter_by_db_name" {
  depends_on = [huaweicloud_dds_database_role.test]

  instance_id = huaweicloud_dds_instance.test.id
  db_name     = huaweicloud_dds_database_role.test.db_name
}

locals {
  filter_result_by_db_name = [for v in data.huaweicloud_dds_database_roles.filter_by_db_name.roles[*].db_name : 
    v == huaweicloud_dds_database_role.test.db_name]
}

output "is_db_name_filter_useful" {
  value = length(local.filter_result_by_db_name) > 0 && alltrue(local.filter_result_by_db_name) 
}
`, testAccDatabaseRole_basic(name))
}
