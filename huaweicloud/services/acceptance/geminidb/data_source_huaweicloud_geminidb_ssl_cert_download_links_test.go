package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSslCertDownloadLinks_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_geminidb_ssl_cert_download_links.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccCheckGeminidbInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSslCertDownloadLinks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "certs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "certs.0.category"),
					resource.TestCheckResourceAttrSet(dataSource, "certs.0.download_link"),
				),
			},
		},
	})
}

func testAccDataSourceSslCertDownloadLinks_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_geminidb_ssl_cert_download_links" "test" {
  instance_id = "%s"
}
`, acceptance.HW_GEMINIDB_INSATNCE_ID)
}
