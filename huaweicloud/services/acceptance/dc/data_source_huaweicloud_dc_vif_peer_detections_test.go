package dc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcVifPeerDetections_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSource := "data.huaweicloud_dc_vif_peer_detections.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDcDirectConnection(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDcVifPeerDetections_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "vif_peer_detections.#"),
					resource.TestCheckResourceAttrSet(dataSource, "vif_peer_detections.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "vif_peer_detections.0.vif_peer_id"),
					resource.TestCheckResourceAttrSet(dataSource, "vif_peer_detections.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "vif_peer_detections.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "vif_peer_detections.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "vif_peer_detections.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "vif_peer_detections.0.loss_ratio"),

					resource.TestCheckOutput("sort_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDcVifPeerDetections_basic(name string) string {
	return fmt.Sprintf(`
%s


data "huaweicloud_dc_vif_peer_detections" "test" {
  depends_on = [huaweicloud_dc_vif_peer_detection.test]

  vif_peer_id = huaweicloud_dc_vif_peer_detection.test.id
}

data "huaweicloud_dc_vif_peer_detections" "sort_filter" {
  depends_on = [huaweicloud_dc_vif_peer_detection.test]

  vif_peer_id = huaweicloud_dc_vif_peer_detection.test.id
  sort_key    = "id"
  sort_dir    = ["asc"]
}

output "sort_filter_is_useful" {
  value = length(data.huaweicloud_dc_vif_peer_detections.sort_filter.vif_peer_detections) > 0
}
`, testResourceDcVifPeerDetection_basic(name))
}
