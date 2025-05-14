package sfsturbo

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDataTaskResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.SfsV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating SFS v1 client: %s", err)
	}

	getDataTaskHttpUrl := "sfs-turbo/{share_id}/hpc-cache/task/{task_id}"
	getDataTaskPath := client.ResourceBaseURL() + getDataTaskHttpUrl
	getDataTaskPath = strings.ReplaceAll(getDataTaskPath, "{share_id}", state.Primary.Attributes["share_id"])
	getDataTaskPath = strings.ReplaceAll(getDataTaskPath, "{task_id}", state.Primary.ID)

	getDataTaskOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getDataTaskResp, err := client.Request("GET", getDataTaskPath, &getDataTaskOpts)
	if err != nil {
		return nil, fmt.Errorf("error retrieving data task: %s", err)
	}

	return utils.FlattenResponse(getDataTaskResp)
}

func TestAccDataTask_basic(t *testing.T) {
	var (
		obj     interface{}
		name    = acceptance.RandomAccResourceName()
		randInt = acctest.RandInt()
		rName   = "huaweicloud_sfs_turbo_data_task.test"
		rc      = acceptance.InitResourceCheck(
			rName,
			&obj,
			getDataTaskResourceFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBSEndpoint(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataTask_basic(name, randInt),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "share_id", "huaweicloud_sfs_turbo.test", "id"),
					resource.TestCheckResourceAttr(rName, "type", "import_metadata"),
					resource.TestCheckResourceAttr(rName, "src_prefix", "test"),
					resource.TestCheckResourceAttr(rName, "dest_prefix", "test"),
					resource.TestCheckResourceAttrPair(rName, "src_target", "huaweicloud_sfs_turbo_obs_target.test",
						"file_system_path"),
					resource.TestCheckResourceAttrPair(rName, "dest_target", "huaweicloud_sfs_turbo_obs_target.test",
						"file_system_path"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccDataTaskImportStateFunc(rName),
			},
		},
	})
}

func testAccDataTask_basic(name string, randInt int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_sfs_turbo_data_task" "test" {
  share_id    = huaweicloud_sfs_turbo.test.id
  type        = "import_metadata"
  src_target  = huaweicloud_sfs_turbo_obs_target.test.file_system_path
  dest_target = huaweicloud_sfs_turbo_obs_target.test.file_system_path
  src_prefix  = "test"
  dest_prefix = "test"
}
`, testAccOBSTarget_basic(name, randInt))
}

func testAccDataTaskImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var shareId, taskId string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", rName)
		}

		shareId = rs.Primary.Attributes["share_id"]
		taskId = rs.Primary.ID
		if shareId == "" || taskId == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<share_id>/<id>', but got '%s/%s'",
				shareId, taskId)
		}
		return fmt.Sprintf("%s/%s", shareId, taskId), nil
	}
}
