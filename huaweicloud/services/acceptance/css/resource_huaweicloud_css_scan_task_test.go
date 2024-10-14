package css

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

func getScanTaskFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		getScanTaskHttpUrl = "v1.0/{project_id}/clusters/{cluster_id}/ai-ops"
		getScanTaskProduct = "css"
	)

	getScanTaskClient, err := conf.NewServiceClient(getScanTaskProduct, acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CSS client: %s", err)
	}

	getScanTaskPath := getScanTaskClient.Endpoint + getScanTaskHttpUrl
	getScanTaskPath = strings.ReplaceAll(getScanTaskPath, "{project_id}", getScanTaskClient.ProjectID)
	getScanTaskPath = strings.ReplaceAll(getScanTaskPath, "{cluster_id}", state.Primary.Attributes["cluster_id"])

	getScanTaskPathOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	currentTotal := 0
	for {
		currentPath := fmt.Sprintf("%s?limit=10&start=%d", getScanTaskPath, currentTotal)
		getScanTaskResp, err := getScanTaskClient.Request("GET", currentPath, &getScanTaskPathOpt)
		if err != nil {
			return nil, err
		}
		getScanTaskRespBody, err := utils.FlattenResponse(getScanTaskResp)
		if err != nil {
			return nil, err
		}
		scanTasks := utils.PathSearch("aiops_list", getScanTaskRespBody, make([]interface{}, 0)).([]interface{})
		findAiopsList := fmt.Sprintf("aiops_list | [?name=='%s'] | [0]", state.Primary.Attributes["name"])
		scanTask := utils.PathSearch(findAiopsList, getScanTaskRespBody, nil)
		if scanTask != nil {
			return scanTask, nil
		}
		total := utils.PathSearch("total_size", getScanTaskRespBody, float64(0)).(float64)
		currentTotal += len(scanTasks)
		if float64(currentTotal) == total {
			break
		}
	}

	return nil, golangsdk.ErrDefault404{}
}

func TestAccScanTask_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_css_scan_task.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getScanTaskFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccScanTask_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "scan task test"),
					resource.TestCheckResourceAttr(resourceName, "alarm.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "alarm.0.level", "suggestion"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "summary.#"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alarm"},
				ImportStateIdFunc:       testAccResourceScanTaskImportStateIDFunc(resourceName),
			},
		},
	})
}

func testAccResourceScanTaskImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		clusterID := rs.Primary.Attributes["cluster_id"]
		scanTaskName := rs.Primary.Attributes["name"]
		if clusterID == "" || scanTaskName == "" {
			return "", fmt.Errorf("invalid format specified for import ID, " +
				"must be '<cluster_id>/<name>'")
		}
		return fmt.Sprintf("%s/%s", clusterID, scanTaskName), nil
	}
}

func testAccScanTask_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_smn_topic" "test" {
  name = "%[2]s"
}

resource "huaweicloud_css_scan_task" "test" {
  cluster_id  = huaweicloud_css_cluster.test.id
  name        = "%[2]s"
  description = "scan task test"
  
  alarm {
    level     = "suggestion"
    smn_topic = huaweicloud_smn_topic.test.name
  }
}
`, testAccCssCluster_basic(rName, "Test@passw0rd", 7, "bar"), rName)
}
