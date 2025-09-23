package deprecated

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/deprecated"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getCCEProtectionFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region      = acceptance.HW_REGION_NAME
		product     = "hss"
		clusterID   = state.Primary.Attributes["cluster_id"]
		clusterName = state.Primary.Attributes["cluster_name"]
		epsID       = state.Primary.Attributes["enterprise_project_id"]
	)

	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating HSS client: %s", err)
	}

	respBody, err := deprecated.ReadCCEProtection(client, clusterID, clusterName, epsID, region)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CCE protection configuration: %s", err)
	}

	clusterRespBody := utils.PathSearch("data_list|[0]", respBody, nil)
	if clusterRespBody == nil {
		return nil, fmt.Errorf("error retrieving HSS protection configuration for CCE:" +
			" configuration is empty in API response")
	}

	protectStatus := utils.PathSearch("protect_status", clusterRespBody, "").(string)
	if protectStatus == "" {
		return nil, fmt.Errorf("error retrieving HSS protection configuration for CCE:" +
			" protect_status is not found in API response")
	}

	if protectStatus == "unprotect" {
		return nil, golangsdk.ErrDefault404{}
	}
	return clusterRespBody, nil
}

// Due to the limitations of the test environment, only some scenarios and abnormal situations can be tested.
// The current test cases cannot be executed 100% successfully. We need to re-verify after the open API is launched on
// the official website.
// Currently, the test case can only be called in the `cn-north-7` region.
func TestAccCCEProtection_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_hss_cce_protection.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCCEProtectionFunc,
	)

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSCCEProtection(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEProtection_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cluster_type", "existing"),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_CCE_CLUSTER_ID),
					resource.TestCheckResourceAttr(rName, "cluster_name", acceptance.HW_CCE_CLUSTER_NAME),
					resource.TestCheckResourceAttr(rName, "charging_mode", "on_demand"),
					resource.TestCheckResourceAttr(rName, "cce_protection_type", "cluster_level"),
					resource.TestCheckResourceAttr(rName, "prefer_packet_cycle", "false"),
					resource.TestCheckResourceAttr(rName, "protect_status", "protecting"),
				),
			},
			{
				Config: testAccCCEProtection_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_CCE_CLUSTER_ID),
					resource.TestCheckResourceAttr(rName, "node_total_num", "1"),
					resource.TestCheckResourceAttr(rName, "cluster_name", acceptance.HW_CCE_CLUSTER_NAME),
					resource.TestCheckResourceAttr(rName, "charging_mode", "on_demand"),
					resource.TestCheckResourceAttr(rName, "prefer_packet_cycle", "true"),
					resource.TestCheckResourceAttr(rName, "protect_status", "protecting"),
					resource.TestCheckResourceAttr(rName, "cluster_type", "existing"),
					resource.TestCheckResourceAttr(rName, "fail_reason", ""),
					resource.TestCheckResourceAttrSet(rName, "protect_type"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"cce_protection_type",
					"enterprise_project_id",
					"agent_version",
					"runtime_info",
					"schedule_info",
					"invoked_service",
				},
				ImportStateIdFunc: testCCEProtectionImportState(rName),
			},
		},
	})
}

func testCCEProtectionImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		clusterID := rs.Primary.Attributes["cluster_id"]
		clusterName := rs.Primary.Attributes["cluster_name"]

		if clusterID == "" {
			return "", fmt.Errorf("attribute (cluster_id) of resource (%s) not found", name)
		}

		if clusterName == "" {
			return "", fmt.Errorf("attribute (cluster_name) of resource (%s) not found", name)
		}

		return fmt.Sprintf("%s/%s", clusterID, clusterName), nil
	}
}

func testAccCCEProtection_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_cce_protection" "test" {
  cluster_type          = "existing"
  cluster_id            = "%[1]s"
  cluster_name          = "%[2]s"
  charging_mode         = "on_demand"
  cce_protection_type   = "cluster_level"
  prefer_packet_cycle   = false
  enterprise_project_id = "%[3]s"
}
`, acceptance.HW_CCE_CLUSTER_ID, acceptance.HW_CCE_CLUSTER_NAME, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccCCEProtection_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_cce_protection" "test" {
  cluster_type          = "existing"
  cluster_id            = "%[1]s"
  cluster_name          = "%[2]s"
  charging_mode         = "on_demand"
  cce_protection_type   = "node_level"
  prefer_packet_cycle   = true
  enterprise_project_id = "%[3]s"
  agent_version         = "1.23.1"
  auto_upgrade          = false
  invoked_service       = "cce"

  runtime_info {
    runtime_name = "ciro_endpoint"
    runtime_path = "/var/run/crio/crio.sock"
  }
  runtime_info {
    runtime_name = "containerd_endpoint"
    runtime_path = "/var/run/containerd/containerd.sock"
  }
  runtime_info {
    runtime_name = "docker_endpoint"
    runtime_path = "/var/run/docker.sock"
  }
  runtime_info {
    runtime_name = "isulad_endpoint"
    runtime_path = "/var/run/isulad.sock"
  }
}
`, acceptance.HW_CCE_CLUSTER_ID, acceptance.HW_CCE_CLUSTER_NAME, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
