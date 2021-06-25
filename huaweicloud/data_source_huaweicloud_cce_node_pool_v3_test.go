package huaweicloud

import (
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccCCENodePoolV3DataSource_basic(t *testing.T) {
	rName := fmtp.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "data.huaweicloud_cce_node_pool.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodePoolV3DataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodePoolV3DataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func testAccCheckCCENodePoolV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find node pools data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("Node pool data source ID not set ")
		}

		return nil
	}
}

func testAccCCENodePoolV3DataSource_basic(rName string) string {
	return fmtp.Sprintf(`
%s

data "huaweicloud_cce_node_pool" "test" {
  cluster_id = huaweicloud_cce_cluster.test.id
  name       = huaweicloud_cce_node_pool.test.name
}
`, testAccCCENodePool_basic(rName))
}
