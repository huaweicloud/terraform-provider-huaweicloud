package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCCEAddonTemplateV3DataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEAddonTemplateV3DataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.huaweicloud_cce_addon_template.spark_operator_test", "spec"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_cce_addon_template.nginx_ingress_test", "spec"),
				),
			},
		},
	})
}

func testAccCCEAddonTemplateV3DataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cce_addon_template" "spark_operator_test" {
  cluster_id = huaweicloud_cce_cluster.test.id
  name       = "spark-operator"
  version    = "1.0.1"
}

data "huaweicloud_cce_addon_template" "nginx_ingress_test" {
  cluster_id = huaweicloud_cce_cluster.test.id
  name       = "nginx-ingress"
  version    = "1.2.2"
}
`, testAccCCEClusterV3_basic(rName))
}
