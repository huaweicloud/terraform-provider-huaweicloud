package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccClusterCertificateRotatecredentials_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckCceChartPath(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterCertificateRotatecredentials_basic(rName),
			},
		},
	})
}

func testAccClusterCertificateRotatecredentials_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_cluster_certificate_rotatecredentials" "test" {
  cluster_id = huaweicloud_cce_cluster.test.id	
  component  = "service-account-controller"
}
`, testAccNode_Base(rName))
}
