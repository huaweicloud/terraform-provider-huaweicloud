package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccClusterOpenIDJWKSDataSource_basic(t *testing.T) {
	datasourceName := "data.huaweicloud_cce_cluster_openid_jwks.test"
	dc := acceptance.InitDataSourceCheck(datasourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCceClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testClousterOpenIDJWKS_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(datasourceName, "keys.0.use"),
					resource.TestCheckResourceAttrSet(datasourceName, "keys.0.kty"),
					resource.TestCheckResourceAttrSet(datasourceName, "keys.0.kid"),
					resource.TestCheckResourceAttrSet(datasourceName, "keys.0.alg"),
					resource.TestCheckResourceAttrSet(datasourceName, "keys.0.n"),
					resource.TestCheckResourceAttrSet(datasourceName, "keys.0.e"),
				),
			},
		},
	})
}

func testClousterOpenIDJWKS_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cce_cluster_openid_jwks" "test" {
  cluster_id = "%s"
}`, acceptance.HW_CCE_CLUSTER_ID)
}
