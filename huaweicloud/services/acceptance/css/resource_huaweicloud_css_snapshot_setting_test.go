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

func getSnapshotSettingResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("css", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CSS v1 client: %s", err)
	}

	getClusterDetailsHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}"
	getClusterDetailsPath := client.Endpoint + getClusterDetailsHttpUrl
	getClusterDetailsPath = strings.ReplaceAll(getClusterDetailsPath, "{project_id}", client.ProjectID)
	getClusterDetailsPath = strings.ReplaceAll(getClusterDetailsPath, "{cluster_id}", state.Primary.ID)

	getClusterDetailsPathOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getClusterDetailsResp, err := client.Request("GET", getClusterDetailsPath, &getClusterDetailsPathOpt)
	if err != nil {
		return false, fmt.Errorf("error retrieving CSS cluster details: %s", err)
	}

	getClusterDetailsRespBody, err := utils.FlattenResponse(getClusterDetailsResp)
	if err != nil {
		return false, fmt.Errorf("error flattening CSS cluster details response: %s", err)
	}

	// The snapshot function is closed when the backupAvailable is false.
	backupAvailable := utils.PathSearch("backupAvailable", getClusterDetailsRespBody, false).(bool)
	if !backupAvailable {
		return nil, golangsdk.ErrDefault404{}
	}

	getSnapshotSettingHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/index_snapshot/policy"
	getSnapshotSettingPath := client.Endpoint + getSnapshotSettingHttpUrl
	getSnapshotSettingPath = strings.ReplaceAll(getSnapshotSettingPath, "{project_id}", client.ProjectID)
	getSnapshotSettingPath = strings.ReplaceAll(getSnapshotSettingPath, "{cluster_id}", state.Primary.ID)

	getSnapshotSettingPathOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getSnapshotSettingResp, err := client.Request("GET", getSnapshotSettingPath, &getSnapshotSettingPathOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CSS snapshot setting: %s", err)
	}

	getSnapshotSettingRespBody, err := utils.FlattenResponse(getSnapshotSettingResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening CSS snapshot setting response: %s", err)
	}
	return getSnapshotSettingRespBody, nil
}

func TestAccSnapshotSetting_basic(t *testing.T) {
	var obj interface{}
	rName := "huaweicloud_css_snapshot_setting.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getSnapshotSettingResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCSSClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSnapshotSetting_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_CSS_CLUSTER_ID),
					resource.TestCheckResourceAttrPair(rName, "bucket", "huaweicloud_obs_bucket.bucket1", "bucket"),
					resource.TestCheckResourceAttr(rName, "base_path", "css_repository/css-snapshot"),
					resource.TestCheckResourceAttr(rName, "agency", "css_obs_agency"),
					resource.TestCheckResourceAttr(rName, "max_snapshot_bytes_per_seconds", "100mb"),
					resource.TestCheckResourceAttr(rName, "max_restore_bytes_per_seconds", "100mb"),
					resource.TestCheckResourceAttr(rName, "enable", "true"),
					resource.TestCheckResourceAttr(rName, "indices", "index1"),
					resource.TestCheckResourceAttr(rName, "prefix", "snapshot"),
					resource.TestCheckResourceAttr(rName, "period", "01:00 GMT+08:00"),
					resource.TestCheckResourceAttr(rName, "keepday", "7"),
					resource.TestCheckResourceAttr(rName, "frequency", "DAY"),
					resource.TestCheckResourceAttr(rName, "delete_auto", "false"),
					resource.TestCheckResourceAttr(rName, "snapshot_cmk_id", ""),
				),
			},
			{
				Config: testSnapshotSetting_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_CSS_CLUSTER_ID),
					resource.TestCheckResourceAttrPair(rName, "bucket", "huaweicloud_obs_bucket.bucket2", "bucket"),
					resource.TestCheckResourceAttr(rName, "base_path", "css_repository/css-snapshot-update"),
					resource.TestCheckResourceAttr(rName, "agency", "css_obs_agency"),
					resource.TestCheckResourceAttr(rName, "max_snapshot_bytes_per_seconds", "0"),
					resource.TestCheckResourceAttr(rName, "max_restore_bytes_per_seconds", "0"),
					resource.TestCheckResourceAttr(rName, "enable", "true"),
					resource.TestCheckResourceAttr(rName, "indices", "index1,index2"),
					resource.TestCheckResourceAttr(rName, "prefix", "snapshot-update"),
					resource.TestCheckResourceAttr(rName, "period", "02:00 GMT+08:00"),
					resource.TestCheckResourceAttr(rName, "keepday", "10"),
					resource.TestCheckResourceAttr(rName, "frequency", "SUN"),
					resource.TestCheckResourceAttr(rName, "delete_auto", "false"),
					resource.TestCheckResourceAttr(rName, "snapshot_cmk_id", ""),
				),
			},
			{
				Config: testSnapshotSetting_update_closeAuto(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_CSS_CLUSTER_ID),
					resource.TestCheckResourceAttrPair(rName, "bucket", "huaweicloud_obs_bucket.bucket2", "bucket"),
					resource.TestCheckResourceAttr(rName, "base_path", "css_repository/css-snapshot-update"),
					resource.TestCheckResourceAttr(rName, "agency", "css_obs_agency"),
					resource.TestCheckResourceAttr(rName, "max_snapshot_bytes_per_seconds", "0"),
					resource.TestCheckResourceAttr(rName, "max_restore_bytes_per_seconds", "0"),
					resource.TestCheckResourceAttr(rName, "enable", "false"),
					resource.TestCheckResourceAttr(rName, "indices", "index1,index2"),
					resource.TestCheckResourceAttr(rName, "prefix", "snapshot-update"),
					resource.TestCheckResourceAttr(rName, "period", "02:00 GMT+08:00"),
					resource.TestCheckResourceAttr(rName, "keepday", "10"),
					resource.TestCheckResourceAttr(rName, "frequency", "SUN"),
					resource.TestCheckResourceAttr(rName, "delete_auto", "true"),
					resource.TestCheckResourceAttr(rName, "snapshot_cmk_id", ""),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_auto"},
			},
		},
	})
}

func testSnapshotSetting_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket1" {
  bucket        = "tf-test-css-1"
  storage_class = "STANDARD"
  acl           = "private"
}

resource "huaweicloud_obs_bucket" "bucket2" {
  bucket        = "tf-test-css-2"
  storage_class = "STANDARD"
  acl           = "private"
}

resource "huaweicloud_css_snapshot_setting" "test" {
  cluster_id                     = "%s"
  bucket                         = huaweicloud_obs_bucket.bucket1.bucket
  agency                         = "css_obs_agency"
  base_path                      = "css_repository/css-snapshot"
  max_snapshot_bytes_per_seconds = "100mb"
  max_restore_bytes_per_seconds  = "100mb"
  enable                         = "true"
  indices                        = "index1"
  prefix                         = "snapshot"
  period                         = "01:00 GMT+08:00"
  keepday                        = 7
  frequency                      = "DAY"
  delete_auto                    = "false"
}
`, acceptance.HW_CSS_CLUSTER_ID)
}

func testSnapshotSetting_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket1" {
  bucket        = "tf-test-css-1"
  storage_class = "STANDARD"
  acl           = "private"
}

resource "huaweicloud_obs_bucket" "bucket2" {
  bucket        = "tf-test-css-2"
  storage_class = "STANDARD"
  acl           = "private"
}

resource "huaweicloud_css_snapshot_setting" "test" {
  cluster_id                     = "%s"
  bucket                         = huaweicloud_obs_bucket.bucket2.bucket
  agency                         = "css_obs_agency"
  base_path                      = "css_repository/css-snapshot-update"
  max_snapshot_bytes_per_seconds = "0"
  max_restore_bytes_per_seconds  = "0"
  enable                         = "true"
  indices                        = "index1,index2"
  prefix                         = "snapshot-update"
  period                         = "02:00 GMT+08:00"
  keepday                        = 10
  frequency                      = "SUN"
  delete_auto                    = "false"
}
`, acceptance.HW_CSS_CLUSTER_ID)
}

func testSnapshotSetting_update_closeAuto() string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket1" {
  bucket        = "tf-test-css-1"
  storage_class = "STANDARD"
  acl           = "private"
}

resource "huaweicloud_obs_bucket" "bucket2" {
  bucket        = "tf-test-css-2"
  storage_class = "STANDARD"
  acl           = "private"
}

resource "huaweicloud_css_snapshot_setting" "test" {
  cluster_id  = "%s"
  bucket      = huaweicloud_obs_bucket.bucket2.bucket
  agency      = "css_obs_agency"
  base_path   = "css_repository/css-snapshot-update"
  enable      = "false"
  delete_auto = "true"
}
`, acceptance.HW_CSS_CLUSTER_ID)
}
