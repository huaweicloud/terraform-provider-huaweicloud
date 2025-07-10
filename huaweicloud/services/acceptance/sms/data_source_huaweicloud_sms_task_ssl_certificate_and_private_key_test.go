package sms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSmsTaskSslCertificateAndPrivateKey_basic(t *testing.T) {
	dataSource := "data.huaweicloud_sms_task_ssl_certificate_and_private_key.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceSmsTaskSslCertificateAndPrivateKey_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "cert"),
					resource.TestCheckResourceAttrSet(dataSource, "private_key"),
					resource.TestCheckResourceAttrSet(dataSource, "ca"),
					resource.TestCheckResourceAttrSet(dataSource, "target_mgmt_cert"),
					resource.TestCheckResourceAttrSet(dataSource, "target_mgmt_private_key"),
					resource.TestCheckResourceAttrSet(dataSource, "target_data_cert"),
					resource.TestCheckResourceAttrSet(dataSource, "target_data_private_key"),
				),
			},
		},
	})
}

func testDataSourceDataSourceSmsTaskSslCertificateAndPrivateKey_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_sms_task_ssl_certificate_and_private_key" "test" {
  task_id        = huaweicloud_sms_task.migration.id
  enable_ca_cert = true
}
`, testAccMigrationTask_basic(name))
}
