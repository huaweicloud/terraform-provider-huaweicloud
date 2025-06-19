package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceKmsParametersForImport_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_kms_parameters_for_import.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare a import KMS key ID and config it to the environment variable.
			acceptance.TestAccPreCheckKmsKeyID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceKmsParametersForImport_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "import_token"),
					resource.TestCheckResourceAttrSet(dataSource, "expiration_time"),
					resource.TestCheckResourceAttrSet(dataSource, "public_key"),
				),
			},
		},
	})
}

func testDataSourceKmsParametersForImport_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_kms_parameters_for_import" "test" {
  key_id             = "%s"
  wrapping_algorithm = "RSAES_OAEP_SHA_256"
}
`, acceptance.HW_KMS_KEY_ID)
}
