package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataHostNets_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dws_host_nets.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byInstanceName   = "data.huaweicloud_dws_host_nets.filter_by_instance_name"
		dcByInstanceName = acceptance.InitDataSourceCheck(byInstanceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataHostNets_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "host_nets.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "host_nets.0.virtual_cluster_id"),
					resource.TestCheckResourceAttrSet(all, "host_nets.0.ctime"),
					resource.TestCheckResourceAttrSet(all, "host_nets.0.host_id"),
					resource.TestCheckResourceAttrSet(all, "host_nets.0.host_name"),
					resource.TestCheckResourceAttrSet(all, "host_nets.0.instance_name"),
					resource.TestCheckResourceAttrSet(all, "host_nets.0.interface_name"),
					resource.TestCheckResourceAttrSet(all, "host_nets.0.up"),
					resource.TestCheckResourceAttrSet(all, "host_nets.0.speed"),
					resource.TestCheckResourceAttrSet(all, "host_nets.0.recv_packets"),
					resource.TestCheckResourceAttrSet(all, "host_nets.0.send_packets"),
					resource.TestCheckResourceAttrSet(all, "host_nets.0.recv_drop"),
					resource.TestCheckResourceAttrSet(all, "host_nets.0.recv_rate"),
					resource.TestCheckResourceAttrSet(all, "host_nets.0.send_rate"),
					resource.TestCheckResourceAttrSet(all, "host_nets.0.io_rate"),
					dcByInstanceName.CheckResourceExists(),
					resource.TestCheckOutput("instance_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataHostNets_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dws_host_nets" "test" {
  cluster_id = "%[1]s"
}

locals {
  instance_name = data.huaweicloud_dws_host_nets.test.host_nets[0].instance_name
}

data "huaweicloud_dws_host_nets" "filter_by_instance_name" {
  cluster_id    = "%[1]s"
  instance_name = local.instance_name
}

locals {
  instance_names = [
    for v in data.huaweicloud_dws_host_nets.filter_by_instance_name.host_nets[*].instance_name : v
  ]
}

output "instance_name_filter_is_useful" {
  value = length(local.instance_names) > 0 && alltrue([
    for v in local.instance_names : v == local.instance_name
  ])
}
`, acceptance.HW_DWS_CLUSTER_ID)
}
