package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDdsSslCertDownloadLinks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dds_ssl_cert_download_links.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDdsSslCertDownloadLinks_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "certs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "certs.0.download_link"),
					resource.TestCheckResourceAttrSet(dataSource, "certs.0.category"),
				),
			},
		},
	})
}

func testDataSourceDdsSslCertDownloadLinks_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dds_ssl_cert_download_links" "test" {
  instance_id = huaweicloud_dds_instance.instance.id
}
`, testAccDDSInstanceV3Config_basic(name, 8800))
}
