package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKMSAliasesDataSource_basic(t *testing.T) {
	var (
		dataSourceName1 = "data.huaweicloud_kms_aliases.basic"
		dc1             = acceptance.InitDataSourceCheck(dataSourceName1)
		dataSourceName2 = "data.huaweicloud_kms_aliases.filter_by_keyId"
		dc2             = acceptance.InitDataSourceCheck(dataSourceName2)
		name            = acceptance.RandomAccResourceName()
		aliasName       = fmt.Sprintf("alias/%s", name)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckKmsKeyID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKmsAliases_basic(aliasName),
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
					resource.TestCheckOutput("is_keyId_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceKmsAliases_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_kms_aliases" "basic" {
  depends_on = [huaweicloud_kms_alias.test]
}

data "huaweicloud_kms_aliases" "filter_by_keyId" {
  key_id = "%[2]s"
  depends_on = [huaweicloud_kms_alias.test]
}

locals {
  keyId_filter_result = [for v in data.huaweicloud_kms_aliases.filter_by_keyId.aliases[*].key_id : v == "%[2]s"]
}

output "is_keyId_filter_useful" {
  value = alltrue(local.keyId_filter_result) && length(local.keyId_filter_result) > 0
}
`, testKmsAlias_basic(name), acceptance.HW_KMS_KEY_ID)
}
