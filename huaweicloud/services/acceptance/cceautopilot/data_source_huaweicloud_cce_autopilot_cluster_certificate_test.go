package cceautopilot

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCceAutopilotClusterCertificate_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cce_autopilot_cluster_certificate.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCceAutopilotClusterCertificate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "duration", "30"),
					resource.TestCheckResourceAttr(dataSource, "clusters.#", "2"),
					resource.TestCheckResourceAttr(dataSource, "users.#", "1"),
					resource.TestCheckResourceAttr(dataSource, "contexts.#", "2"),
					resource.TestCheckResourceAttrSet(dataSource, "current_context"),
					resource.TestCheckResourceAttrSet(dataSource, "kube_config_raw"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCceAutopilotClusterCertificate_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cce_autopilot_cluster_certificate" "test" {
  cluster_id = huaweicloud_cce_autopilot_cluster.test.id
  duration   = 30
}
`, testAccCluster_basic(name))
}
