package hss

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

func getContainerKubernetesClusterDaemonsetResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "hss"
		epsId   = state.Primary.Attributes["enterprise_project_id"]
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating HSS client: %s", err)
	}

	queryPath := client.Endpoint + "v5/{project_id}/container/kubernetes/clusters/{cluster_id}/daemonsets"
	queryPath = strings.ReplaceAll(queryPath, "{project_id}", client.ProjectID)
	queryPath = strings.ReplaceAll(queryPath, "{cluster_id}", state.Primary.ID)
	if epsId != "" {
		queryPath = fmt.Sprintf("%s?enterprise_project_id=%s", queryPath, epsId)
	}

	queryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", queryPath, &queryOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving HSS container kubernetes cluster daemonset: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	runtimeInfoResp := utils.PathSearch("runtime_info", respBody, nil)
	yamlContentResp := utils.PathSearch("yaml_content", respBody, nil)
	if runtimeInfoResp == nil || yamlContentResp == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return respBody, nil
}

func TestAccContainerKubernetesClusterDaemonset_basic(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_hss_container_kubernetes_cluster_daemonset.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getContainerKubernetesClusterDaemonsetResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires the preparation of a CCE cluster under the default enterprise project.
			acceptance.TestAccPreCheckHSSCCEProtection(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testContainerKubernetesClusterDaemonset_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "cluster_id", acceptance.HW_CCE_CLUSTER_ID),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", acceptance.HW_CCE_CLUSTER_NAME),
					resource.TestCheckResourceAttr(resourceName, "auto_upgrade", "true"),
					resource.TestCheckResourceAttr(resourceName, "runtime_info.0.runtime_name", "crio_endpoint"),
					resource.TestCheckResourceAttr(resourceName, "runtime_info.0.runtime_path", "user/test"),
					resource.TestCheckResourceAttr(resourceName, "schedule_info.0.node_selector.0", "test=test"),
				),
			},
			{
				Config: testContainerKubernetesClusterDaemonset_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "schedule_info.0.node_selector.0", "test_update=test_update"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"cluster_id", "cluster_name", "auto_upgrade", "agent_version", "invoked_service", "charging_mode",
					"cce_protection_type", "prefer_packet_cycle", "enterprise_project_id",
				},
			},
		},
	})
}

func testContainerKubernetesClusterDaemonset_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_container_kubernetes_cluster_daemonset" "test" {
  cluster_id   = "%[1]s"
  cluster_name = "%[2]s"
  auto_upgrade = true

  runtime_info {
    runtime_name = "crio_endpoint"
    runtime_path = "user/test"
  }

  schedule_info {
    node_selector = ["test=test"]
  }

  enterprise_project_id = "0"
}
`, acceptance.HW_CCE_CLUSTER_ID, acceptance.HW_CCE_CLUSTER_NAME)
}

func testContainerKubernetesClusterDaemonset_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_container_kubernetes_cluster_daemonset" "test" {
  cluster_id   = "%[1]s"
  cluster_name = "%[2]s"
  auto_upgrade = true

  runtime_info {
    runtime_name = "crio_endpoint"
    runtime_path = "user/test"
  }

  schedule_info {
    node_selector = ["test_update=test_update"]
  }

  enterprise_project_id = "0"
}
`, acceptance.HW_CCE_CLUSTER_ID, acceptance.HW_CCE_CLUSTER_NAME)
}
