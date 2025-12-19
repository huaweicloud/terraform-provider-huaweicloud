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
		resourceName = "huaweicloud_cce_autopiot_release.test"
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
					resource.TestCheckResourceAttr(resourceName, "version", "4.9.0"),
					resource.TestCheckResourceAttrSet(resourceName, "chart_name"),
					resource.TestCheckResourceAttrSet(resourceName, "chart_version"),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_name"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "status_description"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCEAutopilotReleaseImportStateIdFunc("default", name),
				ImportStateVerifyIgnore: []string{
					"version", "values", "chart_id", "description", "parameters",
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

resource "huaweicloud_cce_autopiot_release" "test" {
  cluster_id = huaweicloud_cce_autopilot_cluster.test.id
  chart_id   = huaweicloud_cce_autopilot_chart.test.id
  name       = "%[3]s"
  namespace  = "default"
  version    = "4.9.0"

  values {
    image_tag         = "v1"
    image_pull_policy = "IfNotPresent"
  }

  description = "created by terraform"

  parameters {
    dry_run = false
  }
}
`, testAccCluster_basic(name), testAccAutopilotChart_basic(), name)
}
