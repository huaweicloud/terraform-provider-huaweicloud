package cci

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/cci/v1/namespaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cci"
)

func getNamespaceResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.CciV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud CCI v1 client: %s", err)
	}
	return cci.GetCciNamespaceInfoById(c, state.Primary.ID)
}

func TestAccCciNamespace_basic(t *testing.T) {
	var ns namespaces.Namespace
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_cci_namespace.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&ns,
		getNamespaceResourceFunc,
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
				Config: testAccCciNamespace_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "gpu-accelerated"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "auto_expend_enabled", "true"),
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
				ImportStateIdFunc: testAccCciNamespaceImportStateFunc(resourceName),
			},
		},
	})
}

func TestAccCciNamespace_network(t *testing.T) {
	var ns namespaces.Namespace
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_cci_namespace.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&ns,
		getNamespaceResourceFunc,
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
				Config: testAccCciNamespace_network(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "gpu-accelerated"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "auto_expend_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "warmup_pool_size", "15"),
					resource.TestCheckResourceAttr(resourceName, "recycling_interval", "30"),
					resource.TestCheckResourceAttr(resourceName, "container_network_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "rbac_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "Active"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCciNamespaceImportStateFunc(resourceName),
			},
		},
	})
}

func testAccCciNamespaceImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", rName, rs)
		}
		return rs.Primary.Attributes["name"], nil
	}
}

func testAccCciNamespace_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cci_namespace" "test" {
  name                      = "%s"
  type                      = "gpu-accelerated"
  auto_expend_enabled       = true
  rbac_enabled              = true
  enterprise_project_id     = "%s"
}
`, rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

// The container network of namespace is only supported in cn-north-4.
func testAccCciNamespace_network(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cci_namespace" "test" {
  name                      = "%s"
  type                      = "gpu-accelerated"
  auto_expend_enabled       = true
  warmup_pool_size          = 15
  recycling_interval        = 30
  container_network_enabled = true
  rbac_enabled              = true
  enterprise_project_id     = "%s"
}
`, rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
