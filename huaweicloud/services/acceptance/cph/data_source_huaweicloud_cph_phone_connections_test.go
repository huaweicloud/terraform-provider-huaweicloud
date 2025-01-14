package cph

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCphPhoneConnections_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cph_phone_connections.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCphPhoneConnections_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "connect_infos.#"),
					resource.TestCheckResourceAttrSet(dataSource, "connect_infos.0.phone_id"),
					resource.TestCheckResourceAttrSet(dataSource, "connect_infos.0.access_info.#"),
					resource.TestCheckResourceAttrSet(dataSource, "connect_infos.0.access_info.0.access_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "connect_infos.0.access_info.0.intranet_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "connect_infos.0.access_info.0.access_port"),
					resource.TestCheckResourceAttrSet(dataSource, "connect_infos.0.access_info.0.session_id"),
					resource.TestCheckResourceAttrSet(dataSource, "connect_infos.0.access_info.0.access_time"),
					resource.TestCheckResourceAttrSet(dataSource, "connect_infos.0.access_info.0.ticket"),
				),
			},
			{
				Config: testCphServerBase(rName),
				Check: resource.ComposeTestCheckFunc(
					waitForDeletionCooldownComplete(),
				),
			},
		},
	})
}

func testDataSourceDataSourceCphPhoneConnections_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cph_phones" "test" {
  depends_on = [huaweicloud_cph_server.test]
}

data "huaweicloud_cph_phone_connections" "test" {
  phone_ids   = [data.huaweicloud_cph_phones.test.phones[0].phone_id]
  client_type = "ANDROID"
}
`, testCphServer_basic(name))
}
