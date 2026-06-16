package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataClusterDatabaseUsers_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dws_cluster_database_users.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byType   = "data.huaweicloud_dws_cluster_database_users.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byUserType   = "data.huaweicloud_dws_cluster_database_users.filter_by_user_type"
		dcByUserType = acceptance.InitDataSourceCheck(byUserType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataClusterDatabaseUsers_clusterNotFound(),
				ExpectError: regexp.MustCompile("Cluster does not exist or has been deleted"),
			},
			{
				Config: testAccDataClusterDatabaseUsers_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "users.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "users.0.name"),
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					dcByUserType.CheckResourceExists(),
					resource.TestCheckOutput("user_type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataClusterDatabaseUsers_clusterNotFound() string {
	clusterId, _ := uuid.NewRandom()
	return fmt.Sprintf(`
data "huaweicloud_dws_cluster_database_users" "all" {
  cluster_id = "%s"
}
`, clusterId.String())
}

func testAccDataClusterDatabaseUsers_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dws_cluster_database_users" "all" {
  cluster_id = "%[1]s"
}

# Filter by type
locals {
  type = "USER"
}

data "huaweicloud_dws_cluster_database_users" "filter_by_type" {
  cluster_id = "%[1]s"
  type       = local.type
}

output "type_filter_is_useful" {
  value = length(data.huaweicloud_dws_cluster_database_users.filter_by_type.users) > 0
}

# Filter by user type
locals {
  user_type = "COMMON"
}

data "huaweicloud_dws_cluster_database_users" "filter_by_user_type" {
  cluster_id = "%[1]s"
  type       = "USER"
  user_type  = local.user_type
}

locals {
  user_type_filter_result = [
    for v in data.huaweicloud_dws_cluster_database_users.filter_by_user_type.users[*].user_type : v == local.user_type
  ]
}

output "user_type_filter_is_useful" {
  value = length(local.user_type_filter_result) > 0 && alltrue(local.user_type_filter_result)
}
`, acceptance.HW_DWS_CLUSTER_ID)
}
