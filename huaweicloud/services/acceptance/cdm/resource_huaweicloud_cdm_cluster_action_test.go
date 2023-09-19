package cdm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cdm/v1/clusters"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getClusterActionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.CdmV11Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CDM v1 client: %s", err)
	}
	return clusters.Get(client, state.Primary.ID)
}

func TestAccResourceClusterAction_basic(t *testing.T) {
	var obj clusters.ClusterCreateOpts
	resourceName := "huaweicloud_cdm_cluster_action.test"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		"huaweicloud_cdm_cluster.test",
		&obj,
		getClusterActionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterAction_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "restart"),
					resource.TestCheckResourceAttr(resourceName, "restart.0.level", "SERVICE"),
					resource.TestCheckResourceAttr(resourceName, "restart.0.mode", "IMMEDIATELY"),
				),
			},
			{
				Config: testAccClusterAction_restart(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "restart"),
					resource.TestCheckResourceAttr(resourceName, "restart.0.level", "VM"),
					resource.TestCheckResourceAttr(resourceName, "restart.0.mode", "FORCIBLY"),
				),
			},
			{
				Config:      testAccClusterAction_start(name),
				ExpectError: regexp.MustCompile(`CDM.5090`), // 5090: Start action is not allowed in the current status
			},
		},
	})
}

func testAccClusterAction_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cdm_cluster_action" "test" {
  cluster_id = huaweicloud_cdm_cluster.test.id
  type       = "restart"

  restart {
    level      = "SERVICE"
    mode       = "IMMEDIATELY"
    delay_time = 0
  }
}
`, testAccCdmCluster_basic(name))
}

func testAccClusterAction_restart(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cdm_cluster_action" "test" {
  cluster_id = huaweicloud_cdm_cluster.test.id
  type       = "restart"

  restart {
    level      = "VM"
    mode       = "FORCIBLY"
    delay_time = 0
  }
}
`, testAccCdmCluster_basic(name))
}

func testAccClusterAction_start(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cdm_cluster_action" "test" {
  cluster_id = huaweicloud_cdm_cluster.test.id
  type       = "start"
}
`, testAccCdmCluster_basic(name))
}
