package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceWorkloadQueueAssociatedUsers_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dws_workload_queue_associated_users.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
			acceptance.TestAccPreCheckDwsClusterUserNames(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testDataSourceWorkloadQueueAssociatedUsers_clusterNotExist(),
				ExpectError: regexp.MustCompile("Cluster does not exist or has been deleted"),
			},
			{
				Config:      testDataSourceWorkloadQueueAssociatedUsers_queueNotExist(),
				ExpectError: regexp.MustCompile("resource pool not exist"),
			},
			{
				Config: testDataSourceWorkloadQueueAssociatedUsers_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_exist_users", "true"),
					resource.TestCheckResourceAttrSet(dataSource, "users.0.name"),
					resource.TestMatchResourceAttr(dataSource, "users.0.occupy_resource_list.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "users.0.occupy_resource_list.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "users.0.occupy_resource_list.0.resource_value"),
					resource.TestCheckResourceAttrSet(dataSource, "users.0.occupy_resource_list.0.value_unit"),
				),
			},
		},
	})
}

func testDataSourceWorkloadQueueAssociatedUsers_clusterNotExist() string {
	clusterId, _ := uuid.GenerateUUID()
	name := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dws_workload_queue_associated_users" "test" {
  cluster_id = "%[2]s"
  queue_name = huaweicloud_dws_workload_queue.test.id
}
`, testAccWorkloadQueue_basic(name), clusterId)
}

func testDataSourceWorkloadQueueAssociatedUsers_queueNotExist() string {
	return fmt.Sprintf(`
data "huaweicloud_dws_workload_queue_associated_users" "test" {
  cluster_id = "%s"
  queue_name = "not_found"
}
`, acceptance.HW_DWS_CLUSTER_ID)
}

func testDataSourceWorkloadQueueAssociatedUsers_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dws_workload_queue_user_associate" "test" {
  cluster_id = "%[2]s"
  queue_name = huaweicloud_dws_workload_queue.test.id
  user_names = split(",", "%[3]s")
}

data "huaweicloud_dws_workload_queue_associated_users" "test" {
  depends_on = [huaweicloud_dws_workload_queue_user_associate.test]

  cluster_id = "%[2]s"
  queue_name = huaweicloud_dws_workload_queue.test.id
}

output "is_exist_users" {
  value = length(data.huaweicloud_dws_workload_queue_associated_users.test.users) >= 2
}
`, testAccWorkloadQueue_basic(name), acceptance.HW_DWS_CLUSTER_ID, acceptance.HW_DWS_ASSOCIATE_USER_NAMES)
}
