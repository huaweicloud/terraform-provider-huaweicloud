package das

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDatabaseInstanceConnections_basic(t *testing.T) {
	var (
		name      = acceptance.RandomAccResourceName()
		password  = acceptance.RandomPassword()
		extraName = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_das_database_instance_connections.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByCondition   = "data.huaweicloud_das_database_instance_connections.filter_by_condition"
		dcFilterByCondition = acceptance.InitDataSourceCheck(filterByCondition)

		filterByInstanceId   = "data.huaweicloud_das_database_instance_connections.filter_by_instance_id"
		dcFilterByInstanceId = acceptance.InitDataSourceCheck(filterByInstanceId)

		filterByNetworkType   = "data.huaweicloud_das_database_instance_connections.filter_by_network_type"
		dcFilterByNetworkType = acceptance.InitDataSourceCheck(filterByNetworkType)

		filterByDatastoreType   = "data.huaweicloud_das_database_instance_connections.filter_by_datastore_type"
		dcFilterByDatastoreType = acceptance.InitDataSourceCheck(filterByDatastoreType)

		filterByConnectionType   = "data.huaweicloud_das_database_instance_connections.filter_by_connection_type"
		dcFilterByConnectionType = acceptance.InitDataSourceCheck(filterByConnectionType)

		filterByUsernameFuzzy   = "data.huaweicloud_das_database_instance_connections.filter_by_username_fuzzy"
		dcFilterByUsernameFuzzy = acceptance.InitDataSourceCheck(filterByUsernameFuzzy)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataDatabaseInstanceConnections_basic(name, password, extraName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "connections.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),

					// filter by condition
					dcFilterByCondition.CheckResourceExists(),
					resource.TestCheckOutput("is_condition_filter_useful", "true"),
					resource.TestCheckResourceAttrPair(filterByCondition, "connections.0.id",
						"huaweicloud_das_database_instance_connection.test", "id"),
					resource.TestCheckResourceAttrPair(filterByCondition, "connections.0.instance_id",
						"huaweicloud_das_database_instance_connection.test", "instance_id"),
					resource.TestCheckResourceAttrPair(filterByCondition, "connections.0.engine_type",
						"huaweicloud_das_database_instance_connection.test", "engine_type"),
					resource.TestCheckResourceAttrPair(filterByCondition, "connections.0.network_type",
						"huaweicloud_das_database_instance_connection.test", "network_type"),
					resource.TestCheckResourceAttrPair(filterByCondition, "connections.0.username",
						"huaweicloud_das_database_instance_connection.test", "username"),
					resource.TestCheckResourceAttrPair(filterByCondition, "connections.0.is_save_password",
						"huaweicloud_das_database_instance_connection.test", "is_save_password"),
					resource.TestCheckResourceAttrPair(filterByCondition, "connections.0.description",
						"huaweicloud_das_database_instance_connection.test", "description"),
					resource.TestCheckResourceAttrPair(filterByCondition, "connections.0.port",
						"huaweicloud_das_database_instance_connection.test", "port"),
					resource.TestCheckResourceAttrPair(filterByCondition, "connections.0.instance_name",
						"huaweicloud_das_database_instance_connection.test", "instance_name"),
					resource.TestCheckResourceAttrPair(filterByCondition, "connections.0.datastore_version",
						"huaweicloud_das_database_instance_connection.test", "datastore_version"),
					resource.TestCheckResourceAttrPair(filterByCondition, "connections.0.ip_address",
						"huaweicloud_das_database_instance_connection.test", "ip_address"),
					resource.TestCheckResourceAttrPair(filterByCondition, "connections.0.created_at",
						"huaweicloud_das_database_instance_connection.test", "created_at"),
					resource.TestCheckResourceAttrPair(filterByCondition, "connections.0.status",
						"huaweicloud_das_database_instance_connection.test", "status"),
					resource.TestCheckResourceAttrPair(filterByCondition, "connections.0.conn_share_type",
						"huaweicloud_das_database_instance_connection.test", "conn_share_type"),

					// filter by instance ID
					dcFilterByInstanceId.CheckResourceExists(),
					resource.TestCheckOutput("is_instance_id_filter_useful", "true"),

					// filter by network type
					dcFilterByNetworkType.CheckResourceExists(),
					resource.TestCheckOutput("is_network_type_filter_useful", "true"),

					// filter by datastore type
					dcFilterByDatastoreType.CheckResourceExists(),
					resource.TestCheckOutput("is_datastore_type_filter_useful", "true"),

					// filter by connection type
					dcFilterByConnectionType.CheckResourceExists(),
					resource.TestCheckOutput("is_connection_type_filter_useful", "true"),

					// filter by username fuzzy
					dcFilterByUsernameFuzzy.CheckResourceExists(),
					resource.TestCheckOutput("is_username_fuzzy_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataDatabaseInstanceConnections_basic(name, password, extraName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_mysql_account" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  password    = "%[3]s"

  hosts = [
    "%%"
  ]
}

resource "huaweicloud_das_database_instance_connection" "test" {
  depends_on = [
    huaweicloud_rds_mysql_account.test,
  ]

  instance_id      = "%[1]s"
  engine_type      = "mysql"
  network_type     = "rds"
  username         = "%[2]s"
  password         = "%[3]s"
  is_save_password = true
  node_ids         = ["%[1]s"]
  description      = "Created by terraform script"
  database_name    = "%[2]s"
  sql_record_flag  = true
}

resource "huaweicloud_rds_mysql_account" "extra_test" {
  instance_id = "%[1]s"
  name        = "%[4]s"
  password    = "%[3]s"

  hosts = [
    "%%"
  ]
}

resource "huaweicloud_das_database_instance_connection" "extra_test" {
  depends_on = [
    huaweicloud_rds_mysql_account.extra_test,
  ]

  instance_id      = "%[1]s"
  engine_type      = "mysql"
  network_type     = "rds"
  username         = "%[4]s"
  password         = "%[3]s"
  is_save_password = true
  node_ids         = ["%[1]s"]
  description      = "Created by terraform script"
  database_name    = "%[4]s"
  sql_record_flag  = true
}

data "huaweicloud_das_database_instance_connections" "all" {
  depends_on = [
    huaweicloud_das_database_instance_connection.test,
  ]
}

# Filter by condition
locals {
  condition = huaweicloud_das_database_instance_connection.test.username
}

data "huaweicloud_das_database_instance_connections" "filter_by_condition" {
  depends_on = [
    huaweicloud_das_database_instance_connection.test,
  ]

  condition = local.condition
}

locals {
  condition_filter_result = [
    for v in data.huaweicloud_das_database_instance_connections.filter_by_condition.connections : v.username == local.condition
  ]
}

output "is_condition_filter_useful" {
  value = length(local.condition_filter_result) > 0 && alltrue(local.condition_filter_result)
}

# Filter by instance ID
locals {
  instance_id = huaweicloud_das_database_instance_connection.test.instance_id
}

data "huaweicloud_das_database_instance_connections" "filter_by_instance_id" {
  depends_on = [
    huaweicloud_das_database_instance_connection.test,
  ]

  instance_id = local.instance_id
}

locals {
  instance_id_filter_result = [
    for v in data.huaweicloud_das_database_instance_connections.filter_by_instance_id.connections : v.instance_id == local.instance_id
  ]
}

output "is_instance_id_filter_useful" {
  value = length(local.instance_id_filter_result) > 0 && alltrue(local.instance_id_filter_result)
}

# Filter by network type
locals {
  network_type = huaweicloud_das_database_instance_connection.test.network_type
}

data "huaweicloud_das_database_instance_connections" "filter_by_network_type" {
  depends_on = [
    huaweicloud_das_database_instance_connection.test,
  ]

  network_type = local.network_type
}

locals {
  network_type_filter_result = [
    for v in data.huaweicloud_das_database_instance_connections.filter_by_network_type.connections : v.network_type == local.network_type
  ]
}

output "is_network_type_filter_useful" {
  value = length(local.network_type_filter_result) > 0 && alltrue(local.network_type_filter_result)
}

# Filter by datastore type
locals {
  datastore_type = huaweicloud_das_database_instance_connection.test.engine_type
}

data "huaweicloud_das_database_instance_connections" "filter_by_datastore_type" {
  depends_on = [
    huaweicloud_das_database_instance_connection.test,
  ]

  datastore_type = local.datastore_type
}

locals {
  datastore_type_filter_result = [
    for v in data.huaweicloud_das_database_instance_connections.filter_by_datastore_type.connections : v.engine_type == local.datastore_type
  ]
}

output "is_datastore_type_filter_useful" {
  value = length(local.datastore_type_filter_result) > 0 && alltrue(local.datastore_type_filter_result)
}

# Filter by connection type
locals {
  connection_type = huaweicloud_das_database_instance_connection.test.conn_share_type
}

data "huaweicloud_das_database_instance_connections" "filter_by_connection_type" {
  depends_on = [
    huaweicloud_das_database_instance_connection.test,
  ]

  connection_type = local.connection_type
}

locals {
  connection_type_filter_result = [
    for v in data.huaweicloud_das_database_instance_connections.filter_by_connection_type.connections : v.conn_share_type == local.connection_type
  ]
}

output "is_connection_type_filter_useful" {
  value = length(local.connection_type_filter_result) > 0 && alltrue(local.connection_type_filter_result)
}

# Filter by username fuzzy
locals {
  username = huaweicloud_das_database_instance_connection.test.username
  username_prefix = split("_", local.username)[0]
}

data "huaweicloud_das_database_instance_connections" "filter_by_username_fuzzy" {
  depends_on = [
    huaweicloud_das_database_instance_connection.test,
    huaweicloud_das_database_instance_connection.extra_test,
  ]

  condition = local.username_prefix
}

locals {
  username_fuzzy_filter_result = [
    for v in data.huaweicloud_das_database_instance_connections.filter_by_username_fuzzy.connections : strcontains(v.username, local.username_prefix)
  ]
}

output "is_username_fuzzy_filter_useful" {
  value = length(local.username_fuzzy_filter_result) > 1 && alltrue(local.username_fuzzy_filter_result)
}
`, acceptance.HW_RDS_INSTANCE_ID, name, password, extraName)
}
