package gaussdb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGaussDBKmsDetailDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_gaussdb_kms_detail.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBKMSKeyId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBKmsDetailDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "kms_project_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "key_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "key_alias"),
					resource.TestCheckResourceAttrSet(dataSourceName, "key_state"),
				),
			},
		},
	})
}

func testAccGaussDBKmsDetailDataSource_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_kms_detail" "test" {
  kms_project_name = "cn-north-4"
  key_id           = "%s"
}
`, acceptance.HW_GAUSSDB_KMS_KEY_ID)
}
