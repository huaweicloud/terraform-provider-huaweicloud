package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceKmsPublicKey_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_kms_public_key.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckKmsKeyID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceKmsPublicKey_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "public_key"),
				),
			},
		},
	})
}

func testDataSourceKmsPublicKey_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_kms_public_key" "test" {
  key_id = "%s"
}
`, acceptance.HW_KMS_KEY_ID)
}
