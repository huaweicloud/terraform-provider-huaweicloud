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

func getDuTaskResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.SfsV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating SFS v1 client: %s", err)
	}

	getDuTaskHttpUrl := "sfs-turbo/shares/{share_id}/fs/{feature}/tasks/{task_id}"
	getDuTaskPath := client.ResourceBaseURL() + getDuTaskHttpUrl
	getDuTaskPath = strings.ReplaceAll(getDuTaskPath, "{share_id}", state.Primary.Attributes["share_id"])
	getDuTaskPath = strings.ReplaceAll(getDuTaskPath, "{feature}", "dir-usage")
	getDuTaskPath = strings.ReplaceAll(getDuTaskPath, "{task_id}", state.Primary.ID)

	getDuTaskOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getDuTaskResp, err := client.Request("GET", getDuTaskPath, &getDuTaskOpts)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DU task: %s", err)
	}

	return utils.FlattenResponse(getDuTaskResp)
}

func TestAccDuTask_basic(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		path  = "/temp" + acctest.RandString(10)
		rName = "huaweicloud_sfs_turbo_du_task.test"
		rc    = acceptance.InitResourceCheck(
			rName,
			&obj,
			getDuTaskResourceFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDuTask_basic(name, path),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "share_id", "huaweicloud_sfs_turbo.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "path", "huaweicloud_sfs_turbo_dir.test", "path"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "begin_time"),
					resource.TestCheckResourceAttrSet(rName, "end_time"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccDuTaskImportStateFunc(rName),
			},
		},
	})
}

func testAccDuTask_basic(name, path string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_sfs_turbo_dir" "test" {
  share_id = huaweicloud_sfs_turbo.test.id
  path     = "%[2]s"
  mode     = 777
  gid      = 100
  uid      = 100
}

resource "huaweicloud_sfs_turbo_du_task" "test" {
  share_id = huaweicloud_sfs_turbo.test.id
  path     = huaweicloud_sfs_turbo_dir.test.id
}
`, testAccSFSTurbo_shareTypeHpc(name), path)
}

func testAccDuTaskImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var shareId, taskId string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
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
