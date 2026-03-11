package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataClusterDatabaseSchemas_basic(t *testing.T) {
	var (
		all          = "data.huaweicloud_dws_cluster_database_schemas.test"
		dc           = acceptance.InitDataSourceCheck(all)

		byKeywords   = "data.huaweicloud_dws_cluster_database_schemas.filter_by_keywords"
		dcByKeywords = acceptance.InitDataSourceCheck(byKeywords)

		emptyResult  = "data.huaweicloud_dws_cluster_database_schemas.not_found_database"
		dcEmpty      = acceptance.InitDataSourceCheck(emptyResult)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataClusterDatabaseSchemas_clusterNotFound(),
				ExpectError: regexp.MustCompile("Cluster does not exist or has been deleted"),
			},
			{
				Config: testAccDataClusterDatabaseSchemas_databaseNotFound(),
				Check: resource.ComposeTestCheckFunc(
					dcEmpty.CheckResourceExists(),
					resource.TestCheckResourceAttr(emptyResult, "schemas.#", "0"),
				),
			},
			{
				Config: testAccDataClusterDatabaseSchemas_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "schemas.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "schemas.0.schema_name"),
					resource.TestCheckResourceAttrSet(all, "schemas.0.database_name"),
					resource.TestCheckResourceAttrSet(all, "schemas.0.total_value"),
					resource.TestCheckResourceAttrSet(all, "schemas.0.perm_space"),
					resource.TestCheckResourceAttrSet(all, "schemas.0.skew_percent"),
					resource.TestCheckResourceAttrSet(all, "schemas.0.min_value"),
					resource.TestCheckResourceAttrSet(all, "schemas.0.max_value"),
					resource.TestCheckResourceAttrSet(all, "schemas.0.min_dn"),
					resource.TestCheckResourceAttrSet(all, "schemas.0.max_dn"),
					resource.TestCheckResourceAttrSet(all, "schemas.0.dn_num"),
					dcByKeywords.CheckResourceExists(),
					resource.TestCheckOutput("keywords_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataClusterDatabaseSchemas_clusterNotFound() string {
	randomUUID, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dws_cluster_database_schemas" "test" {
  cluster_id    = "%s"
  database_name = "gaussdb"
}
`, randomUUID)
}

func testAccDataClusterDatabaseSchemas_databaseNotFound() string {
	return fmt.Sprintf(`
data "huaweicloud_dws_cluster_database_schemas" "not_found_database" {
  cluster_id    = "%s"
  database_name = "database_name_not_exist_for_acc_test"
}
`, acceptance.HW_DWS_CLUSTER_ID)
}

func testAccDataClusterDatabaseSchemas_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dws_cluster_database_schemas" "test" {
  cluster_id    = "%[1]s"
  database_name = "gaussdb"
}

locals {
  schema_name = data.huaweicloud_dws_cluster_database_schemas.test.schemas[0].schema_name
}

data "huaweicloud_dws_cluster_database_schemas" "filter_by_keywords" {
  cluster_id    = "%[1]s"
  database_name = "gaussdb"
  keywords      = local.schema_name
}

locals {
  schema_names = data.huaweicloud_dws_cluster_database_schemas.filter_by_keywords.schemas[*].schema_name
}

output "keywords_filter_is_useful" {
  value = length(local.schema_names) > 0 && contains(local.schema_names, local.schema_name)
}
`, acceptance.HW_DWS_CLUSTER_ID)
}
