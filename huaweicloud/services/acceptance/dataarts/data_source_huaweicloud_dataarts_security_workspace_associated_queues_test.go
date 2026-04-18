package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSecurityWorkspaceAssociatedQueues_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dataarts_security_workspace_associated_queues.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byType   = "data.huaweicloud_dataarts_security_workspace_associated_queues.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byClusterId   = "data.huaweicloud_dataarts_security_workspace_associated_queues.filter_by_cluster_id"
		dcByClusterId = acceptance.InitDataSourceCheck(byClusterId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckMrsClusterID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSecurityWorkspaceAssociatedQueues_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttrSet(all, "region"),
					resource.TestMatchResourceAttr(all, "queues.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(all, "queues.0.id"),
					resource.TestCheckResourceAttrSet(all, "queues.0.source_type"),
					resource.TestCheckResourceAttrSet(all, "queues.0.name"),
					resource.TestCheckResourceAttrSet(all, "queues.0.type"),
					resource.TestCheckResourceAttrSet(all, "queues.0.connection_id"),
					resource.TestCheckResourceAttrSet(all, "queues.0.conn_name"),
					resource.TestCheckResourceAttrSet(all, "queues.0.created_at"),
					resource.TestCheckResourceAttrSet(all, "queues.0.create_user"),
					resource.TestCheckResourceAttrSet(all, "queues.0.project_id"),

					// filter by type
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),

					// filter by cluster_id
					dcByClusterId.CheckResourceExists(),
					resource.TestCheckOutput("is_cluster_id_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(byClusterId, "queues.0.cluster_id"),
					resource.TestCheckResourceAttrSet(byClusterId, "queues.0.cluster_name"),
				),
			},
		},
	})
}

func testDataSourceSecurityWorkspaceAssociatedQueues_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dataarts_security_workspace_associated_queues" "all" {
  workspace_id = "%[1]s"
}

// filter by type
data "huaweicloud_dataarts_security_workspace_associated_queues" "filter_by_type" {
  workspace_id = "%[1]s"
  type         = "dli"
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_dataarts_security_workspace_associated_queues.filter_by_type.queues[*].source_type : 
    v == "dli"
  ]
}

output "is_type_filter_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}

// filter by cluster_id.
data "huaweicloud_dataarts_security_workspace_associated_queues" "filter_by_cluster_id" {
  workspace_id = "%[1]s"
  type         = "mrs"
  cluster_id   = "%[2]s"
}

locals {
  cluster_id_filter_result = [
    for v in data.huaweicloud_dataarts_security_workspace_associated_queues.filter_by_cluster_id.queues[*].cluster_id :
    v == "%[2]s"
  ]
}

output "is_cluster_id_filter_useful" {
  value = length(local.cluster_id_filter_result) > 0 && alltrue(local.cluster_id_filter_result)
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_MRS_CLUSTER_ID)
}
