package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDdsDatabaseUsers_basic(t *testing.T) {
	rName := "data.huaweicloud_dds_database_users.all"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDdsDatabaseUsers_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "users.#"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_db_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceDdsDatabaseUsers_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dds_database_users" "all" {
  depends_on = [huaweicloud_dds_database_user.test]

  instance_id = huaweicloud_dds_instance.test.id
}

// filter by name
data "huaweicloud_dds_database_users" "filter_by_name" {
  depends_on = [huaweicloud_dds_database_user.test]

  instance_id = huaweicloud_dds_instance.test.id
  name        = huaweicloud_dds_database_user.test.name
}

locals {
  filter_result_by_name = [for v in data.huaweicloud_dds_database_users.filter_by_name.users[*].name : 
    v == huaweicloud_dds_database_user.test.name]
}

output "is_name_filter_useful" {
  value = length(local.filter_result_by_name) > 0 && alltrue(local.filter_result_by_name) 
}

// filter by db_name
data "huaweicloud_dds_database_users" "filter_by_db_name" {
  depends_on = [huaweicloud_dds_database_user.test]

  instance_id = huaweicloud_dds_instance.test.id
  db_name     = huaweicloud_dds_database_user.test.db_name
}

locals {
  filter_result_by_db_name = [for v in data.huaweicloud_dds_database_users.filter_by_db_name.users[*].db_name : 
    v == huaweicloud_dds_database_user.test.db_name]
}

output "is_db_name_filter_useful" {
  value = length(local.filter_result_by_db_name) > 0 && alltrue(local.filter_result_by_db_name) 
}
`, testAccDatabaseUser_basic(name))
}
