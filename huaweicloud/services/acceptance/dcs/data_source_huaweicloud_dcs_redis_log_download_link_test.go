package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcsRedisLogDownloadLink_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dcs_redis_log_download_link.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDcsRedisLogDownloadLink_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.huaweicloud_dcs_redis_run_logs.test", "file_list.0.id"),
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "backup_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "link"),
				),
			},
		},
	})
}

func testAccDataSourceDcsRedisLogDownloadLink_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dcs_redis_log_download_link" "test" {
  depends_on  = [data.huaweicloud_dcs_redis_run_logs.test]

  instance_id = huaweicloud_dcs_instance.test.id
  log_id      = data.huaweicloud_dcs_redis_run_logs.test.file_list.0.id
}
`, testDataSourceRedisRunLogs_basic(name))
}
