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

func sfsTurboDirReadfunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v1/{project_id}/sfs-turbo/shares/{share_id}/fs/dir"
	)
	sfsClient, err := cfg.SfsV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating SFS v1 client: %s", err)
	}

	getPath := sfsClient.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", sfsClient.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{share_id}", state.Primary.Attributes["share_id"])
	getPath += fmt.Sprintf("?path=%s", state.Primary.ID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := sfsClient.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving SFS Turbo directory %s", err)
	}

	return utils.FlattenResponse(getResp)
}

func TestAccSfsTurboDir_basic(t *testing.T) {
	var obj interface{}
	name := acceptance.RandomAccResourceName()
	path := "/tmp" + acctest.RandString(10)
	resourceName := "huaweicloud_sfs_turbo_dir.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		sfsTurboDirReadfunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSfsTurboDirBasic(name, path),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "path", path),
				),
			},
		},
	})
}

func testSfsTurboDirBasic(rName string, path string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_sfs_turbo_dir" "test" {
  path     = "%[2]s"
  share_id = huaweicloud_sfs_turbo.test.id
  mode     = 777
  gid      = 100
  uid      = 100
}
`, testAccSFSTurbo_basic(rName), path)
}
