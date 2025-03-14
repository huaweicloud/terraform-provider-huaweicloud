package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cci/v1/namespaces"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cci"
)

func getV2NamespaceResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.CciV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud CCI v1 client: %s", err)
	}
	return cci.GetNamespaceDetail(c, state.Primary.ID)
}

func TestAccNamespace_basic(t *testing.T) {
	var ns namespaces.Namespace
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_cciv2_namespace.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&ns,
		getV2NamespaceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNamespace_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "flavor", "general-computing"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "warmup_pool_size", "0"),
					resource.TestCheckResourceAttr(resourceName, "recycling_interval", "0"),
					resource.TestCheckResourceAttr(resourceName, "container_network_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "rbac_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "Active"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNamespace_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cciv2_namespace" "test" {
  name = %s
	
  annotations = {
    namespace.kubernetes.io/flavor            = "gpu-accelerated"
    network.cci.io/warm-pool-size             = "10"
    network.cci.io/warm-pool-recycle-interval = "24"
    network.cci.io/ready-before-pod-run       = "vpc-network-ready"
  }
	
  labels = {
    rbac.authorization.cci.io/enable-k8s-rbac = "true",
    sys_enterprise_project_id                 = "0"
  }
}
`, rName)
}
