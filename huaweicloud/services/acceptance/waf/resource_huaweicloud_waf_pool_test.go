package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/waf"
)

func getPoolResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("waf", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating WAF client: %s", err)
	}

	return waf.GetWafPool(client, state.Primary.ID, state.Primary.Attributes["enterprise_project_id"])
}

// Before running the test case, please ensure that there is at least one WAF instance in the current region.
func TestAccPool_basic(t *testing.T) {
	var (
		obj          interface{}
		randName     = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_waf_pool.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getPoolResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWafPool_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "type", "detector-cloud"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform"),
					resource.TestCheckResourceAttrSet(resourceName, "hosts.#"),
					resource.TestCheckResourceAttrSet(resourceName, "instances.#"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testWAFPoolImportState(resourceName),
			},
		},
	})
}

func testAccWafPool_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}
`, name)
}

func testAccWafPool_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_pool" "test" {
  name                  = "%[2]s"
  type                  = "detector-cloud"
  vpc_id                = huaweicloud_vpc.test.id
  description           = "created by terraform"
  enterprise_project_id = "%[3]s"
}
`, testAccWafPool_base(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testWAFPoolImportState(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		epsId := rs.Primary.Attributes["enterprise_project_id"]
		id := rs.Primary.ID
		if id == "" || epsId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, "+
				"want '<id>/<enterprise_project_id>', but got '%s/%s'",
				id, epsId)
		}

		return fmt.Sprintf("%s/%s", id, epsId), nil
	}
}
