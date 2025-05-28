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

func getOBSTargetResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.SfsV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating SFS v1 client: %s", err)
	}

	getObsTargetHttpUrl := "sfs-turbo/shares/{share_id}/targets/{target_id}"
	getObsTargetPath := client.ResourceBaseURL() + getObsTargetHttpUrl
	getObsTargetPath = strings.ReplaceAll(getObsTargetPath, "{share_id}", state.Primary.Attributes["share_id"])
	getObsTargetPath = strings.ReplaceAll(getObsTargetPath, "{target_id}", state.Primary.ID)

	getObsTargetOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getObsTargetResp, err := client.Request("GET", getObsTargetPath, &getObsTargetOpts)
	if err != nil {
		return nil, fmt.Errorf("error retrieving OBS target: %s", err)
	}

	return utils.FlattenResponse(getObsTargetResp)
}

func TestAccOBSTarget_basic(t *testing.T) {
	var (
		obj     interface{}
		name    = acceptance.RandomAccResourceName()
		randInt = acctest.RandInt()
		rName   = "huaweicloud_sfs_turbo_obs_target.test"
		rc      = acceptance.InitResourceCheck(
			rName,
			&obj,
			getOBSTargetResourceFunc,
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
				Config: testAccOBSTarget_basic(name, randInt),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "share_id", "huaweicloud_sfs_turbo.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "obs.0.bucket", "huaweicloud_obs_bucket.test", "id"),
					resource.TestCheckResourceAttr(rName, "file_system_path", "obsdir"),
					resource.TestCheckResourceAttr(rName, "obs.0.endpoint", acceptance.HW_OBS_ENDPOINT),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_data_in_file_system"},
				ImportStateIdFunc:       testAccOBSTargetImportStateFunc(rName),
			},
		},
	})
}

func testAccOBSTarget_basic_base(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "tf-test-bucket-%d"
  storage_class = "STANDARD"
  acl           = "private"
}
`, randInt)
}

func testAccOBSTarget_basic(name string, randInt int) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "huaweicloud_sfs_turbo_obs_target" "test" {
  share_id         = huaweicloud_sfs_turbo.test.id
  file_system_path = "obsdir"

  obs {
    bucket   = huaweicloud_obs_bucket.test.id
    endpoint = "%[3]s"

    policy {
      auto_export_policy {
        events = ["NEW","DELETED"]
        prefix = "pre"
        suffix = "suf"
      }
    }

    attributes {
      file_mode = "421"
      dir_mode  = "750"
      uid       = 101
      gid       = 234
    }
  }
}
`, testAccOBSTarget_basic_base(randInt), testAccSFSTurbo_shareTypeHpc(name), acceptance.HW_OBS_ENDPOINT)
}

func testAccOBSTargetImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var shareId, targetId string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		shareId = rs.Primary.Attributes["share_id"]
		targetId = rs.Primary.ID
		if shareId == "" || targetId == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<instance_id>/<id>', but got '%s/%s'",
				shareId, targetId)
		}
		return fmt.Sprintf("%s/%s", shareId, targetId), nil
	}
}
