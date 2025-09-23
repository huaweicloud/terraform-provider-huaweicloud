package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceV2Pods_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_cciv2_pods.test"
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceV2Pods_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "pods.0.namespace", rName),
					resource.TestCheckResourceAttr(dataSourceName, "pods.0.name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "pods.0.containers.0.image", "nginx:stable-alpine-perl"),
					resource.TestCheckResourceAttr(dataSourceName, "pods.0.containers.0.name", "c1"),
					resource.TestCheckResourceAttr(dataSourceName, "pods.0.containers.0.resources.0.limits.cpu", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "pods.0.containers.0.resources.0.limits.memory", "2G"),
					resource.TestCheckResourceAttr(dataSourceName, "pods.0.containers.0.resources.0.requests.cpu", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "pods.0.containers.0.resources.0.requests.memory", "2G"),
					resource.TestCheckResourceAttr(dataSourceName, "pods.0.image_pull_secrets.0.name", "imagepull-secret"),
					resource.TestCheckResourceAttrSet(dataSourceName, "pods.0.annotations.%"),
				),
			},
		},
	})
}

func testAccDataSourceV2Pods_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cciv2_pods" "test" {
  depends_on = [huaweicloud_cciv2_pod.test]

  namespace = huaweicloud_cciv2_namespace.test.name
}
`, testAccV2Pod_basic(rName))
}
