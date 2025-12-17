package bms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBmsVolumeAttachments_basic(t *testing.T) {
	dataSource := "data.huaweicloud_bms_volume_attachments.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceBmsVolumeAttachments_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "volume_attachments.#"),
					resource.TestCheckResourceAttrSet(dataSource, "volume_attachments.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "volume_attachments.0.server_id"),
					resource.TestCheckResourceAttrSet(dataSource, "volume_attachments.0.volume_id"),
					resource.TestCheckResourceAttrSet(dataSource, "volume_attachments.0.device"),
				),
			},
		},
	})
}

func testDataSourceBmsVolumeAttachments_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_bms_volume_attachments" "test" {
  depends_on = [huaweicloud_bms_volume_attach.test]

  server_id = huaweicloud_bms_instance.test.id
}
`, testAccBmsVolumeAttach_basic(name))
}
