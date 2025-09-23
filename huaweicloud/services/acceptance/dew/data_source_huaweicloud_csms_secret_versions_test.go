package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCsmsSecretVersions_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_csms_secret_versions.test"
		rName      = acceptance.RandomAccResourceName()
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCsmsSecretVersions_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "versions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.kms_key_id"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.secret_name"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.version_stages.#"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.created_at"),
				),
			},
		},
	})
}

func testDataSourceCsmsSecretVersions_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_csms_secret" "test" {
  name        = "%s"
  description = "desc"
  secret_text = "terraform"
}

data "huaweicloud_csms_secret_versions" "test" {
  secret_name = huaweicloud_csms_secret.test.name
}
`, name)
}
