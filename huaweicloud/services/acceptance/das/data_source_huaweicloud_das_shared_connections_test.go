package das

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSharedConnections_basic(t *testing.T) {
	var (
		name        = acceptance.RandomAccResourceName()
		password    = acceptance.RandomPassword()
		currentTime = time.Now().Local().Format(time.RFC3339)

		all = "data.huaweicloud_das_shared_connections.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByUserId   = "data.huaweicloud_das_shared_connections.filter_by_user_id"
		dcFilterByUserId = acceptance.InitDataSourceCheck(filterByUserId)

		filterByUserNameFuzzy   = "data.huaweicloud_das_shared_connections.filter_by_user_name_fuzzy"
		dcFilterByUserNameFuzzy = acceptance.InitDataSourceCheck(filterByUserNameFuzzy)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.12.1",
			},
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSharedConnections_nonExistentSharedConnections(),
				ExpectError: regexp.MustCompile("error querying DAS shared connections"),
			},
			{
				Config: testAccDataSharedConnections_basic(name, password, currentTime),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "shared_connections.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),

					// filter by user ID
					dcFilterByUserId.CheckResourceExists(),
					resource.TestCheckOutput("is_user_id_filter_useful", "true"),
					resource.TestCheckResourceAttrPair(filterByUserId, "shared_connections.0.user_id",
						"huaweicloud_das_shared_connection.test", "user_id"),
					resource.TestCheckResourceAttrPair(filterByUserId, "shared_connections.0.user_name",
						"huaweicloud_das_shared_connection.test", "user_name"),
					resource.TestCheckResourceAttrPair(filterByUserId, "shared_connections.0.expired_at",
						"huaweicloud_das_shared_connection.test", "expired_at"),
					resource.TestCheckResourceAttrPair(filterByUserId, "shared_connections.0.shared_at",
						"huaweicloud_das_shared_connection.test", "shared_at"),

					// filter by user name fuzzy
					dcFilterByUserNameFuzzy.CheckResourceExists(),
					resource.TestCheckOutput("is_user_name_fuzzy_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSharedConnections_nonExistentSharedConnections() string {
	randUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
data "huaweicloud_das_shared_connections" "non_existent_shared_connections" {
  connection_id = "%[1]s"
}
`, randUUID.String())
}

func testAccDataSharedConnections_basic_base(name, password, currentTime string) string {
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

data "huaweicloud_identity_users" "test" {
  depends_on = [
    huaweicloud_das_database_instance_connection.test,
  ]
}

resource "huaweicloud_das_shared_connection" "test" {
  depends_on = [
    data.huaweicloud_identity_users.test,
  ]

  connection_id = huaweicloud_das_database_instance_connection.test.id
  user_id       = data.huaweicloud_identity_users.test.users[0].id
  user_name     = data.huaweicloud_identity_users.test.users[0].name
  expired_at    = format("%%s+08:00", split("+", timeadd("%[4]s", "1h"))[0])
}
`, acceptance.HW_RDS_INSTANCE_ID, name, password, currentTime)
}

func testAccDataSharedConnections_basic(name, password, currentTime string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_das_shared_connections" "all" {
  depends_on = [
    huaweicloud_das_shared_connection.test,
  ]

  connection_id = huaweicloud_das_database_instance_connection.test.id
}

# Filter by user ID
locals {
  user_id = huaweicloud_das_shared_connection.test.user_id
}

data "huaweicloud_das_shared_connections" "filter_by_user_id" {
  depends_on = [
    data.huaweicloud_das_shared_connections.all,
  ]

  connection_id = huaweicloud_das_database_instance_connection.test.id
  keywords      = local.user_id
}

locals {
  user_id_filter_result = [
    for v in data.huaweicloud_das_shared_connections.filter_by_user_id.shared_connections : v.user_id == local.user_id
  ]
}

output "is_user_id_filter_useful" {
  value = length(local.user_id_filter_result) > 0 && alltrue(local.user_id_filter_result)
}

# Filter by user name fuzzy
locals {
  user_name        = data.huaweicloud_identity_users.test.users[0].name
  user_name_prefix = substr(local.user_name, 1, 2)
}

data "huaweicloud_das_shared_connections" "filter_by_user_name_fuzzy" {
  depends_on = [
    data.huaweicloud_das_shared_connections.all,
  ]

  connection_id = huaweicloud_das_database_instance_connection.test.id
  keywords      = local.user_name_prefix
}

locals {
	user_name_fuzzy_filter_result = [
    for v in data.huaweicloud_das_shared_connections.filter_by_user_name_fuzzy.shared_connections : strcontains(v.user_name, local.user_name_prefix)
  ]
}

output "is_user_name_fuzzy_filter_useful" {
  value = length(local.user_name_fuzzy_filter_result) > 0 && alltrue(local.user_name_fuzzy_filter_result)
}
`, testAccDataSharedConnections_basic_base(name, password, currentTime))
}
