package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpcTrafficMirrorSessions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpc_traffic_mirror_sessions.test1"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceVpcTrafficMirrorSessions_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "name", rName),
					resource.TestCheckResourceAttr(dataSource, "traffic_mirror_target_type", "eni"),
					resource.TestCheckResourceAttr(dataSource, "priority", "10"),
					resource.TestCheckResourceAttrPair(dataSource, "traffic_mirror_session_id",
						"huaweicloud_vpc_traffic_mirror_session.test", "id"),
					resource.TestCheckResourceAttrPair(dataSource, "traffic_mirror_filter_id",
						"huaweicloud_vpc_traffic_mirror_filter.test", "id"),
				),
			},
		},
	})
}

func testDataSourceDataSourceVpcTrafficMirrorSessions_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpc_traffic_mirror_sessions" "test1" {
  traffic_mirror_session_id  = huaweicloud_vpc_traffic_mirror_session.test.id
  name                       = "%s"
  traffic_mirror_filter_id   = huaweicloud_vpc_traffic_mirror_filter.test.id
  traffic_mirror_target_id   = huaweicloud_compute_instance.test[0].network[0].port
  traffic_mirror_target_type = "eni"
  priority                   = 10
}
`, testAccTrafficMirrorSession_basic(name), name)
}
