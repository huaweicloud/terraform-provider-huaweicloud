package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDdsBackupDownloadLinks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dds_backup_download_links.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDdsBackupDownloadLinks_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "files.#"),
					resource.TestCheckResourceAttrSet(dataSource, "files.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "files.0.download_link"),
					resource.TestCheckResourceAttrSet(dataSource, "files.0.size"),
					resource.TestCheckResourceAttrSet(dataSource, "files.0.link_expired_time"),
					resource.TestCheckResourceAttrSet(dataSource, "bucket"),
				),
			},
		},
	})
}

func testDataSourceDdsBackupDownloadLinks_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dds_backup_download_links" "test" {
  instance_id = huaweicloud_dds_instance.instance.id
  backup_id   = huaweicloud_dds_backup.test.id
}
`, testDdsBackup_basic(name))
}
