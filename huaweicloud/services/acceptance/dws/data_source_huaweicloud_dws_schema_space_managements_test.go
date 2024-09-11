package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSchemaSpaceManagements_basic(t *testing.T) {
	var (
		notFoundDatabase   = "data.huaweicloud_dws_schema_space_managements.test"
		dcNotFoundDatabase = acceptance.InitDataSourceCheck(notFoundDatabase)
		dataSource         = "data.huaweicloud_dws_schema_space_managements.test"
		dc                 = acceptance.InitDataSourceCheck(dataSource)
		bySchemaName       = "data.huaweicloud_dws_schema_space_managements.filter_by_schema_name"
		dcBySchemaName     = acceptance.InitDataSourceCheck(bySchemaName)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testDataSourceSchemaSpaceManagements_clusterIdNotExist(),
				ExpectError: regexp.MustCompile("Cluster does not exist or has been deleted"),
			},
			{
				Config: testDataSourceSchemaSpaceManagements_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					dcNotFoundDatabase.CheckResourceExists(),
					resource.TestCheckOutput("not_found_database", "true"),
					resource.TestMatchResourceAttr(dataSource, "schemas.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("assert_space_limit", "true"),
					dcBySchemaName.CheckResourceExists(),
					resource.TestCheckResourceAttr(bySchemaName, "schemas.0.database_name", "gaussdb"),
					resource.TestCheckResourceAttrSet(bySchemaName, "schemas.0.schema_name"),
					resource.TestCheckResourceAttrSet(bySchemaName, "schemas.0.used"),
					resource.TestCheckResourceAttrSet(bySchemaName, "schemas.0.space_limit"),
					resource.TestCheckResourceAttrSet(bySchemaName, "schemas.0.skew_percent"),
					resource.TestCheckResourceAttrSet(bySchemaName, "schemas.0.dn_num"),
				),
			},
		},
	})
}

func testDataSourceSchemaSpaceManagements_clusterIdNotExist() string {
	clusterId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dws_schema_space_managements" "test" {
  cluster_id    = "%s"
  database_name = "gaussdb"
}
`, clusterId)
}

func testDataSourceSchemaSpaceManagements_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dws_schema_space_managements" "not_found_database" {
  cluster_id    = "%[1]s"
  database_name = "not_found_database"
}

output "not_found_database" {
  value = length(data.huaweicloud_dws_schema_space_managements.not_found_database.schemas) == 0
}

data "huaweicloud_dws_schema_space_managements" "test" {
  depends_on = [
    huaweicloud_dws_schema_space_management.test
  ]

  cluster_id    = "%[1]s"
  database_name = huaweicloud_dws_schema_space_management.test.database_name
}

# Modify space quota for scheduler to 2MB (2048 Byte).
resource "huaweicloud_dws_schema_space_management" "test" {
  cluster_id    = "%[1]s"
  database_name = "gaussdb"
  schema_name   = "scheduler"
  space_limit   = "2048"
}

# Filter by schema name.
data "huaweicloud_dws_schema_space_managements" "filter_by_schema_name" {
  depends_on = [
    huaweicloud_dws_schema_space_management.test
  ]
  
  cluster_id    = "%[1]s"
  database_name = huaweicloud_dws_schema_space_management.test.database_name
  schema_name   = local.schema_name
}

locals {
  schema_name = huaweicloud_dws_schema_space_management.test.schema_name

  # Convert the obtained value from Byte to MB.
  space_limit = try([for v in data.huaweicloud_dws_schema_space_managements.filter_by_schema_name.schemas : ceil(v.space_limit / 1024 / 1024)
  if v.schema_name == local.schema_name][0], null)
}

output "assert_space_limit" {
  value = local.space_limit == 2
}
`, acceptance.HW_DWS_CLUSTER_ID)
}
