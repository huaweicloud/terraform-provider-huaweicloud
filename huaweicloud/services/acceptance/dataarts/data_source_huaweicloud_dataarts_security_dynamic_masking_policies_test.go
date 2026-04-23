package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// HW_DATAARTS_CONNECTION_ID must be the connection ID corresponding to a DWS type data connection.
func TestAccDataSecurityDynamicMaskingPolicies_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dataarts_security_dynamic_masking_policies.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_dataarts_security_dynamic_masking_policies.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byClusterName   = "data.huaweicloud_dataarts_security_dynamic_masking_policies.filter_by_cluster_name"
		dcByClusterName = acceptance.InitDataSourceCheck(byClusterName)

		byDatabaseName   = "data.huaweicloud_dataarts_security_dynamic_masking_policies.filter_by_database_name"
		dcByDatabaseName = acceptance.InitDataSourceCheck(byDatabaseName)

		byTableName   = "data.huaweicloud_dataarts_security_dynamic_masking_policies.filter_by_table_name"
		dcByTableName = acceptance.InitDataSourceCheck(byTableName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsConnectionID(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
			acceptance.TestAccPreCheckUserName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSecurityDynamicMaskingPolicies_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "policies.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					// Filter by 'name' parameter.
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(byName, "policies.0.id"),
					resource.TestCheckResourceAttrSet(byName, "policies.0.name"),
					resource.TestCheckResourceAttrSet(byName, "policies.0.datasource_type"),
					resource.TestCheckResourceAttrSet(byName, "policies.0.cluster_id"),
					resource.TestCheckResourceAttrSet(byName, "policies.0.cluster_name"),
					resource.TestCheckResourceAttrSet(byName, "policies.0.database_name"),
					resource.TestCheckResourceAttrSet(byName, "policies.0.table_name"),
					resource.TestCheckResourceAttrSet(byName, "policies.0.user_groups"),
					resource.TestCheckResourceAttrSet(byName, "policies.0.users"),
					resource.TestCheckResourceAttrSet(byName, "policies.0.sync_status"),
					resource.TestMatchResourceAttr(byName, "policies.0.create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(byName, "policies.0.create_user"),
					resource.TestMatchResourceAttr(byName, "policies.0.update_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(byName, "policies.0.update_user"),
					// Filter by 'cluster_name' parameter.
					dcByClusterName.CheckResourceExists(),
					resource.TestCheckOutput("is_cluster_name_filter_useful", "true"),
					// Filter by 'database_name' parameter.
					dcByDatabaseName.CheckResourceExists(),
					resource.TestCheckOutput("is_database_name_filter_useful", "true"),
					// Filter by 'table_name' parameter.
					dcByTableName.CheckResourceExists(),
					resource.TestCheckOutput("is_table_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSecurityDynamicMaskingPolicies_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_dataarts_studio_data_connections" "test" {
  workspace_id  = "%[1]s"
  connection_id = "%[2]s"
}

data "huaweicloud_dws_clusters" "test" {}

resource "huaweicloud_identity_group" "test" {
  name = "%[3]s"
}

locals {
  connection_name = try(data.huaweicloud_dataarts_studio_data_connections.test.connections[0].name, "NOT_FOUND")
  cluster_name    = try([for v in data.huaweicloud_dws_clusters.test.clusters : v.name if v.id == "%[4]s"][0], "NOT_FOUND")
}

resource "huaweicloud_dataarts_security_dynamic_masking_policy" "test" {
  workspace_id    = "%[1]s"
  name            = "%[3]s"
  datasource_type = "DWS"
  cluster_id      = "%[4]s"
  cluster_name    = local.cluster_name
  database_name   = "gaussdb"
  table_name      = "%[3]s"
  users           = "%[5]s"
  user_groups     = huaweicloud_identity_group.test.name
  conn_id         = "%[2]s"
  conn_name       = local.connection_name
  table_id        = "NativeTable-%[4]s-gaussdb-public-%[3]s"
  schema_name     = "public"

  policy_list {
    column_name          = "total"
    column_type          = "int8"
    algorithm_type       = "DWS_SELF_CONFIG"
    algorithm_detail_dto = jsonencode({
      start      = 1
      end        = 2
      int_target = 0
    })
  }
  policy_list {
    column_name    = "enabled"
    column_type    = "bool"
    algorithm_type = "MASK"
  }
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_CONNECTION_ID, name,
		acceptance.HW_DWS_CLUSTER_ID, acceptance.HW_USER_NAME)
}

func testAccDataSecurityDynamicMaskingPolicies_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameters.
data "huaweicloud_dataarts_security_dynamic_masking_policies" "test" {
  workspace_id = "%[2]s"
}

# Filter by 'name' parameter.
locals {
  policy_name = huaweicloud_dataarts_security_dynamic_masking_policy.test.name
}

data "huaweicloud_dataarts_security_dynamic_masking_policies" "filter_by_name" {
  workspace_id = "%[2]s"
  name         = local.policy_name

  depends_on = [huaweicloud_dataarts_security_dynamic_masking_policy.test]
}

locals {
  name_filter_result = [for v in data.huaweicloud_dataarts_security_dynamic_masking_policies.filter_by_name.policies[*].name :
  strcontains(v, local.policy_name)]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by 'cluster_name' parameter.
data "huaweicloud_dataarts_security_dynamic_masking_policies" "filter_by_cluster_name" {
  workspace_id = "%[2]s"
  cluster_name = local.cluster_name

  depends_on = [huaweicloud_dataarts_security_dynamic_masking_policy.test]
}

locals {
  cluster_name_filter_result = [
    for v in data.huaweicloud_dataarts_security_dynamic_masking_policies.filter_by_cluster_name.policies[*].cluster_name :
    strcontains(v, local.cluster_name)
  ]
}

output "is_cluster_name_filter_useful" {
  value = length(local.cluster_name_filter_result) > 0 && alltrue(local.cluster_name_filter_result)
}

# Filter by 'database_name' parameter.
locals {
  database_name = huaweicloud_dataarts_security_dynamic_masking_policy.test.database_name
}

data "huaweicloud_dataarts_security_dynamic_masking_policies" "filter_by_database_name" {
  workspace_id  = "%[2]s"
  database_name = local.database_name

  depends_on = [huaweicloud_dataarts_security_dynamic_masking_policy.test]
}

locals {
  database_name_filter_result = [
    for v in data.huaweicloud_dataarts_security_dynamic_masking_policies.filter_by_database_name.policies[*].database_name :
    strcontains(v, local.database_name)
  ]
}

output "is_database_name_filter_useful" {
  value = length(local.database_name_filter_result) > 0 && alltrue(local.database_name_filter_result)
}

# Filter by 'table_name' parameter.
locals {
  table_name = huaweicloud_dataarts_security_dynamic_masking_policy.test.table_name
}

data "huaweicloud_dataarts_security_dynamic_masking_policies" "filter_by_table_name" {
  workspace_id = "%[2]s"
  table_name   = local.table_name

  depends_on = [huaweicloud_dataarts_security_dynamic_masking_policy.test]
}

locals {
  table_name_filter_result = [
    for v in data.huaweicloud_dataarts_security_dynamic_masking_policies.filter_by_table_name.policies[*].table_name :
    strcontains(v, local.table_name)
  ]
}

output "is_table_name_filter_useful" {
  value = length(local.table_name_filter_result) > 0 && alltrue(local.table_name_filter_result)
}
`, testAccDataSecurityDynamicMaskingPolicies_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID)
}
