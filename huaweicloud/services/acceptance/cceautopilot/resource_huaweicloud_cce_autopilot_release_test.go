package cceautopilot

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getAutopilotReleaseFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	var (
		getAutopilotReleaseHttpUrl = "autopilot/cam/v3/clusters/{cluster_id}/namespace/{namespace}/releases/{name}"
		getAutopilotReleaseProduct = "cce"
	)
	getAutopilotReleaseClient, err := cfg.NewServiceClient(getAutopilotReleaseProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CCE client: %s", err)
	}

	getAutopilotReleaseHttpPath := getAutopilotReleaseClient.Endpoint + getAutopilotReleaseHttpUrl
	getAutopilotReleaseHttpPath = strings.ReplaceAll(getAutopilotReleaseHttpPath, "{cluster_id}", state.Primary.Attributes["cluster_id"])
	getAutopilotReleaseHttpPath = strings.ReplaceAll(getAutopilotReleaseHttpPath, "{namespace}", state.Primary.Attributes["namespace"])
	getAutopilotReleaseHttpPath = strings.ReplaceAll(getAutopilotReleaseHttpPath, "{name}", state.Primary.ID)

	getAutopilotReleaseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getAutopilotReleaseResp, err := getAutopilotReleaseClient.Request("GET", getAutopilotReleaseHttpPath, &getAutopilotReleaseOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CCE Autopilot release: %s", err)
	}

	return utils.FlattenResponse(getAutopilotReleaseResp)
}

func TestAccAutopilotRelease_basic(t *testing.T) {
	var (
		release      interface{}
		resourceName = "huaweicloud_cce_autopilot_release.test"
		name         = acceptance.RandomAccResourceNameWithDash()

		rc = acceptance.InitResourceCheck(
			resourceName,
			&release,
			getAutopilotReleaseFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCceChartPath(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAutopilotRelease_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "cluster_id",
						"huaweicloud_cce_autopilot_cluster.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "chart_id",
						"huaweicloud_cce_autopilot_chart.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "namespace", "default"),
					resource.TestCheckResourceAttrSet(resourceName, "chart_name"),
					resource.TestCheckResourceAttrSet(resourceName, "chart_version"),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_name"),
					resource.TestCheckResourceAttr(resourceName, "status", "DEPLOYED"),
					resource.TestCheckResourceAttrSet(resourceName, "status_description"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttr(resourceName, "release_version", "1"),
				),
			},
			{
				Config: testAccAutopilotRelease_upgrade(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "release_version", "2"),
				),
			},
			{
				Config: testAccAutopilotRelease_rollback(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "release_version", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCEAutopilotReleaseImportStateIdFunc("default", name),
				ImportStateVerifyIgnore: []string{
					"action", "version", "values", "chart_id", "description", "parameters",
				},
			},
		},
	})
}

func testAccCCEAutopilotReleaseImportStateIdFunc(namespace, name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		cluster, ok := s.RootModule().Resources["huaweicloud_cce_autopilot_cluster.test"]
		if !ok {
			return "", fmt.Errorf("cluster not found: %s", cluster)
		}
		if cluster.Primary.ID == "" {
			return "", fmt.Errorf("resource not found: %s/%s/%s", cluster.Primary.ID, namespace, name)
		}
		return fmt.Sprintf("%s/%s/%s", cluster.Primary.ID, namespace, name), nil
	}
}

func testAccAutopilotRelease_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_cce_autopilot_release" "test" {
  cluster_id = huaweicloud_cce_autopilot_cluster.test.id
  chart_id   = huaweicloud_cce_autopilot_chart.test.id
  name       = "%[3]s"
  namespace  = "default"
  version    = "1.0.0"

  values {
    image_pull_policy = "IfNotPresent"
    image_tag         = "5.0"
  }

  description = "created by terraform"

  parameters {
    dry_run = false
  }
}
`, testAccAutopilotRelease_base(name), testAccAutopilotChart_basic(), name)
}

func testAccAutopilotRelease_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_autopilot_cluster" "test" {
  name        = "%[2]s"
  flavor      = "cce.autopilot.cluster"
  description = "created by terraform"

  host_network {
    vpc    = huaweicloud_vpc.test.id
    subnet = huaweicloud_vpc_subnet.test.id
  }

  container_network {
    mode = "eni"
  }

  eni_network {
    subnets {
      subnet_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id
    }
  }

  # for pull image from SWR
  enable_swr_image_access = true

  lifecycle {
    ignore_changes = [
      tags
    ]
  }
}
`, common.TestVpc(rName), rName)
}

func testAccAutopilotRelease_upgrade(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_cce_autopilot_release" "test" {
  cluster_id = huaweicloud_cce_autopilot_cluster.test.id
  chart_id   = huaweicloud_cce_autopilot_chart.test.id
  name       = "%[3]s"
  namespace  = "default"
  version    = "1.0.0"

  values {
    image_pull_policy = "IfNotPresent"
    image_tag         = "5.0.6"
  }

  description = "created by terraform"

  action = "upgrade"

  parameters {
    dry_run = false
  }
}
`, testAccAutopilotRelease_base(name), testAccAutopilotChart_basic(), name)
}

func testAccAutopilotRelease_rollback(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_cce_autopilot_release" "test" {
  cluster_id = huaweicloud_cce_autopilot_cluster.test.id
  chart_id   = huaweicloud_cce_autopilot_chart.test.id
  name       = "%[3]s"
  namespace  = "default"
  version    = "1.0.0"

  values {
    image_pull_policy = "IfNotPresent"
    image_tag         = "5.0.6"
  }

  description = "created by terraform"

  action = "rollback"

  parameters {
    dry_run         = false
    release_version = 1
  }
}
`, testAccAutopilotRelease_base(name), testAccAutopilotChart_basic(), name)
}
