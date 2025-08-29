package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKMSAliasesDataSource_basic(t *testing.T) {
	var (
		dataSourceName1 = "data.huaweicloud_kms_aliases.test"
		dc1             = acceptance.InitDataSourceCheck(dataSourceName1)
		dataSourceName2 = "data.huaweicloud_kms_aliases.filter_by_key_id"
		dc2             = acceptance.InitDataSourceCheck(dataSourceName2)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare a import KMS key ID with an alias and config it to the environment variable.
			acceptance.TestAccPreCheckKmsKeyID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKmsAliases_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName1, "aliases.#"),
					resource.TestCheckResourceAttrSet(dataSourceName1, "aliases.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSourceName1, "aliases.0.key_id"),
					resource.TestCheckResourceAttrSet(dataSourceName1, "aliases.0.alias"),
					resource.TestCheckResourceAttrSet(dataSourceName1, "aliases.0.alias_urn"),
					resource.TestCheckResourceAttrSet(dataSourceName1, "aliases.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSourceName1, "aliases.0.update_time"),
					resource.TestCheckOutput("is_key_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceKmsAliases_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_kms_aliases" "test" {}

data "huaweicloud_kms_aliases" "filter_by_key_id" {
  key_id = "%[1]s"
}

locals {
  key_id_filter_result = [for v in data.huaweicloud_kms_aliases.filter_by_key_id.aliases[*].key_id : v == "%[1]s"]
}

output "is_key_id_filter_useful" {
  value = alltrue(local.key_id_filter_result) && length(local.key_id_filter_result) > 0
}
`, acceptance.HW_KMS_KEY_ID)
}
