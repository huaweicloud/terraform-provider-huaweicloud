package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccHSSContainerKubernetesSyncMCCS_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testContainerKubernetesSyncMCCS_basic,
			},
		},
	})
}

const testContainerKubernetesSyncMCCS_basic = `
resource "huaweicloud_hss_container_kubernetes_sync_mccs" "test" {
  total_num             = 2
  enterprise_project_id = "0"

  data_list {
    cluster_id = "cluster1-id"
  }

  data_list {
    cluster_id = "cluster2-id"
  }
}
`
