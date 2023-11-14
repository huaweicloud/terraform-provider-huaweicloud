package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceLogtanks_basic(t *testing.T) {
	rName := "data.huaweicloud_elb_logtanks.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceLogtanks_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "logtanks.#"),
					resource.TestCheckResourceAttrSet(rName, "logtanks.0.id"),
					resource.TestCheckResourceAttrSet(rName, "logtanks.0.log_topic_id"),
					resource.TestCheckResourceAttrSet(rName, "logtanks.0.log_group_id"),
					resource.TestCheckResourceAttrSet(rName, "logtanks.0.loadbalancer_id"),
					resource.TestCheckOutput("logtank_id_filter_is_useful", "true"),
					resource.TestCheckOutput("log_topic_id_filter_is_useful", "true"),
					resource.TestCheckOutput("log_group_id_filter_is_useful", "true"),
					resource.TestCheckOutput("loadbalancer_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceLogtanks_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_elb_logtanks" "test" {
  depends_on = [huaweicloud_elb_logtank.test]
}

locals {
  logtank_id = huaweicloud_elb_logtank.test.id
}
data "huaweicloud_elb_logtanks" "logtank_id_filter" {
  logtank_id = local.logtank_id
}
output "logtank_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_logtanks.logtank_id_filter.logtanks) > 0 && alltrue(
  [for v in data.huaweicloud_elb_logtanks.logtank_id_filter.logtanks[*].id : v == local.logtank_id]
)
}

locals {
  log_topic_id = huaweicloud_elb_logtank.test.log_topic_id
}
data "huaweicloud_elb_logtanks" "log_topic_id_filter" {
  log_topic_id = local.log_topic_id
}
output "log_topic_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_logtanks.log_topic_id_filter.logtanks) > 0 && alltrue(
  [for v in data.huaweicloud_elb_logtanks.log_topic_id_filter.logtanks[*].log_topic_id : v == local.log_topic_id]
)
}

locals {
  log_group_id = huaweicloud_elb_logtank.test.log_group_id
}
data "huaweicloud_elb_logtanks" "log_group_id_filter" {
  log_group_id = local.log_group_id
}
output "log_group_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_logtanks.log_group_id_filter.logtanks) > 0 && alltrue(
  [for v in data.huaweicloud_elb_logtanks.log_group_id_filter.logtanks[*].log_group_id : v == local.log_group_id]
)
}

locals {
  loadbalancer_id = huaweicloud_elb_logtank.test.loadbalancer_id
}
data "huaweicloud_elb_logtanks" "loadbalancer_id_filter" {
  loadbalancer_id = local.loadbalancer_id
}
output "loadbalancer_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_logtanks.loadbalancer_id_filter.logtanks) > 0 && alltrue(
  [for v in data.huaweicloud_elb_logtanks.loadbalancer_id_filter.logtanks[*].loadbalancer_id : v == local.loadbalancer_id]
  )  
}

`, testAccElbLogTankConfig_basic(name), name)
}
