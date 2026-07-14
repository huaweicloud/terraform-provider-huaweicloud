package gaussdb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGaussDBKmsKeysDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_gaussdb_kms_keys.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBKmsKeysDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "kms_project_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "authorized_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "authorized_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "key_details.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "key_details.0.key_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "key_details.0.default_key_flag"),
					resource.TestCheckResourceAttrSet(dataSourceName, "key_details.0.key_alias"),
					resource.TestCheckResourceAttrSet(dataSourceName, "key_details.0.key_spec"),
					resource.TestCheckResourceAttrSet(dataSourceName, "key_details.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "key_details.0.key_state"),
				),
			},
		},
	})
}

const testAccGaussDBKmsKeysDataSource_basic = `
data "huaweicloud_gaussdb_kms_keys" "test" {
  kms_project_name = "cn-north-4"
}
`
