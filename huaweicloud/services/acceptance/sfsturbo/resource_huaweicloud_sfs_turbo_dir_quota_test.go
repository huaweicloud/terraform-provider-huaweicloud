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

func sfsTurboDirQuotaReadfunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v1/{project_id}/sfs-turbo/shares/{share_id}/fs/dir-quota"
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
		return nil, fmt.Errorf("error retrieving SFS Turbo directory quota %s", err)
	}

	return utils.FlattenResponse(getResp)
}

func TestAccSfsTurboDirQuota_basic(t *testing.T) {
	var obj interface{}
	name := acceptance.RandomAccResourceName()
	path := "/tmp" + acctest.RandString(10)
	resourceName := "huaweicloud_sfs_turbo_dir_quota.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		sfsTurboDirQuotaReadfunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSfsTurboDirQuotaBasic(name, path),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "path", path),
					resource.TestCheckResourceAttr(resourceName, "capacity", "100"),
					resource.TestCheckResourceAttr(resourceName, "inode", "100"),
				),
			},
			{
				Config: testSfsTurboDirQuotaUpdateBasic(name, path),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "path", path),
					resource.TestCheckResourceAttr(resourceName, "capacity", "50"),
					resource.TestCheckResourceAttr(resourceName, "inode", "30"),
				),
			},
		},
	})
}

func testSfsTurboDirQuotaBasic(rName string, path string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_sfs_turbo_dir_quota" "test" {
  path     = "%[2]s"
  share_id = huaweicloud_sfs_turbo_dir.test.share_id
  capacity = 100
  inode    = 100
}
`, testSfsTurboDirBasic(rName, path), path)
}

func testSfsTurboDirQuotaUpdateBasic(rName string, path string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_sfs_turbo_dir_quota" "test" {
  path     = "%[2]s"
  share_id = huaweicloud_sfs_turbo_dir.test.share_id
  capacity = 50
  inode    = 30
}
`, testSfsTurboDirBasic(rName, path), path)
}
