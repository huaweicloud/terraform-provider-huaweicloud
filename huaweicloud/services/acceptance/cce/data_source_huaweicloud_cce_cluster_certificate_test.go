package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccClusterCertificateDataSource_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()
	datasourceName := "data.huaweicloud_cce_cluster_certificate.test"
	dc := acceptance.InitDataSourceCheck(datasourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testClousterCertificate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(datasourceName, "duration", "30"),
					resource.TestCheckResourceAttr(datasourceName, "clusters.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "users.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "contexts.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceName, "current_context"),
					resource.TestCheckResourceAttrSet(datasourceName, "kube_config_raw"),
				),
			},
		},
	})
}

func testClousterCertificate_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cce_cluster_certificate" "test" {
  cluster_id = huaweicloud_cce_cluster.test.id
  duration   = 30
}`, testAccCluster_basic(name))
}
