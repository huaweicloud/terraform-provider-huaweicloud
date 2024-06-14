package css

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/css/v1/cluster"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getCssEsCoreUpgradeDetailFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.CssV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CSS v1 client: %s", err)
	}

	getEsUpgradeDetailHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/upgrade/detail"
	getEsUpgradeDetailPath := client.Endpoint + getEsUpgradeDetailHttpUrl
	getEsUpgradeDetailPath = strings.ReplaceAll(getEsUpgradeDetailPath, "{project_id}", client.ProjectID)
	getEsUpgradeDetailPath = strings.ReplaceAll(
		getEsUpgradeDetailPath, "{cluster_id}", state.Primary.Attributes["cluster_id"])

	getEsUpgradeDetailPathOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getEsUpgradeDetailResp, err := client.Request("GET", getEsUpgradeDetailPath, &getEsUpgradeDetailPathOpt)
	if err != nil {
		return getEsUpgradeDetailResp, err
	}
	getEsUpgradeDetailRespBody, err := utils.FlattenResponse(getEsUpgradeDetailResp)
	if err != nil {
		return getEsUpgradeDetailRespBody, err
	}

	expression := fmt.Sprintf("detailList | [?id=='%s'] | [0]", state.Primary.ID)
	upgradeDetail := utils.PathSearch(expression, getEsUpgradeDetailRespBody, nil)
	if upgradeDetail == nil {
		return upgradeDetail, golangsdk.ErrDefault404{}
	}

	return upgradeDetail, nil
}

func TestAccCssEsCoreUpgrade_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_css_es_core_upgrade.test"

	var obj cluster.ClusterDetailResponse
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCssEsCoreUpgradeDetailFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCssLowEngineVersion(t)
			acceptance.TestAccPreCheckCssTargetImageId(t)
			acceptance.TestAccPreCheckCSSUpgradeAgency(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccCssEsCoreUpgrade_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "upgrade_detail.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "upgrade_detail.0.start_time"),
					resource.TestCheckResourceAttrSet(resourceName, "upgrade_detail.0.end_time"),
					resource.TestCheckResourceAttrSet(resourceName, "upgrade_detail.0.status"),
					resource.TestCheckResourceAttrSet(resourceName, "upgrade_detail.0.datastore.0.type"),
					resource.TestCheckResourceAttrSet(resourceName, "upgrade_detail.0.datastore.0.version"),
				),
			},
		},
	})
}

func testAccCssEsCoreUpgrade_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_es_core_upgrade" "test" {
  cluster_id           = huaweicloud_css_cluster.test.id
  target_image_id      = "%[2]s"
  upgrade_type         = "cross"
  agency               = "%[3]s"
  indices_backup_check = true
  cluster_load_check   = true
}
`, testAccCssCluster_upgrade(rName), acceptance.HW_CSS_TARGET_IMAGE_ID, acceptance.HW_CSS_UPGRADE_AGENCY)
}

func testAccCssCluster_upgrade(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_cluster" "test" {
  name           = "%[2]s"
  engine_version = "%[3]s"

  ess_node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = 3
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id

  lifecycle {
    ignore_changes = [
      engine_version
    ]
  }
}
`, testAccCssBase(name), name, acceptance.HW_CSS_LOW_ENGINE_VERSION)
}
