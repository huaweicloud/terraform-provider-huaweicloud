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

func getLogSettingResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.CssV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CSS v1 client: %s", err)
	}

	getLogSettingHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/logs/settings"
	getLogSettingPath := client.Endpoint + getLogSettingHttpUrl
	getLogSettingPath = strings.ReplaceAll(getLogSettingPath, "{project_id}", client.ProjectID)
	getLogSettingPath = strings.ReplaceAll(getLogSettingPath, "{cluster_id}", state.Primary.ID)

	getLogSettingPathOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getLogSettingResp, err := client.Request("GET", getLogSettingPath, &getLogSettingPathOpt)
	if err != nil {
		return getLogSettingResp, err
	}

	getLogSettingRespBody, err := utils.FlattenResponse(getLogSettingResp)
	if err != nil {
		return nil, fmt.Errorf("erorr retrieving CSS cluster log setting: %s", err)
	}

	return getLogSettingRespBody, nil
}

func TestAccLogSetting_elastic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_css_log_setting.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getLogSettingResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testLogSetting_elastic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "cluster_id", "huaweicloud_css_cluster.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "bucket", "huaweicloud_obs_bucket.cssObs", "bucket"),
					resource.TestCheckResourceAttr(rName, "base_path", "css_repository/css-log"),
					resource.TestCheckResourceAttr(rName, "agency", "css_obs_agency"),
					resource.TestCheckResourceAttr(rName, "auto_enabled", "false"),
					resource.TestCheckResourceAttr(rName, "log_switch", "true"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testLogSetting_elastic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "cluster_id", "huaweicloud_css_cluster.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "bucket", "huaweicloud_obs_bucket.cssObs", "bucket"),
					resource.TestCheckResourceAttr(rName, "base_path", "css_repository/css-log-update"),
					resource.TestCheckResourceAttr(rName, "agency", "css_obs_agency"),
					resource.TestCheckResourceAttr(rName, "auto_enabled", "true"),
					resource.TestCheckResourceAttr(rName, "period", "00:00 GMT+08:00"),
					resource.TestCheckResourceAttr(rName, "log_switch", "true"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testLogSetting_elastic_updateNull(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "auto_enabled", "false"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogSetting_logstash(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_css_log_setting.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getLogSettingResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testLogSetting_logstash(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "cluster_id", "huaweicloud_css_logstash_cluster.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "bucket", "huaweicloud_obs_bucket.bucket1", "bucket"),
					resource.TestCheckResourceAttr(rName, "base_path", "css_repository/logstash-log"),
					resource.TestCheckResourceAttr(rName, "agency", "css_obs_agency"),
					resource.TestCheckResourceAttr(rName, "auto_enabled", "false"),
					resource.TestCheckResourceAttr(rName, "log_switch", "true"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testLogSetting_logstash_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "cluster_id", "huaweicloud_css_logstash_cluster.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "bucket", "huaweicloud_obs_bucket.bucket2", "bucket"),
					resource.TestCheckResourceAttr(rName, "base_path", "css_repository/logstash-log-2"),
					resource.TestCheckResourceAttr(rName, "agency", "css_obs_agency"),
					resource.TestCheckResourceAttr(rName, "auto_enabled", "true"),
					resource.TestCheckResourceAttr(rName, "period", "00:00 GMT+08:00"),
					resource.TestCheckResourceAttr(rName, "log_switch", "true"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testLogSetting_logstash_updateNull(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "auto_enabled", "false"),
				),
			},
		},
	})
}

func testLogSetting_elastic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_css_log_setting" "test" {
  cluster_id = huaweicloud_css_cluster.test.id
  agency     = "css_obs_agency"
  base_path  = "css_repository/css-log"
  bucket     = huaweicloud_obs_bucket.cssObs.bucket
}
`, testAccCssCluster_basic(name, "Test@passw0rd", 1, "tag"))
}

func testLogSetting_elastic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_css_log_setting" "test" {
  cluster_id = huaweicloud_css_cluster.test.id
  agency     = "css_obs_agency"
  base_path  = "css_repository/css-log-update"
  bucket     = huaweicloud_obs_bucket.cssObs.bucket
  period     = "00:00 GMT+08:00"
}
`, testAccCssCluster_basic(name, "Test@passw0rd", 1, "tag"))
}

func testLogSetting_elastic_updateNull(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_css_log_setting" "test" {
  cluster_id = huaweicloud_css_cluster.test.id
  agency     = "css_obs_agency"
  base_path  = "css_repository/css-log-update"
  bucket     = huaweicloud_obs_bucket.cssObs.bucket
  period     = ""
}
`, testAccCssCluster_basic(name, "Test@passw0rd", 1, "tag"))
}

func testLogSetting_logstash(name string) string {
	return fmt.Sprintf(`
%s

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

resource "huaweicloud_css_log_setting" "test" {
  cluster_id = huaweicloud_css_logstash_cluster.test.id
  agency     = "css_obs_agency"
  base_path  = "css_repository/logstash-log"
  bucket     = huaweicloud_obs_bucket.bucket1.bucket
}
`, testAccLogstashCluster_basic(name, 1, "bar"))
}

func testLogSetting_logstash_update(name string) string {
	return fmt.Sprintf(`
%s

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

resource "huaweicloud_css_log_setting" "test" {
  cluster_id = huaweicloud_css_logstash_cluster.test.id
  agency     = "css_obs_agency"
  base_path  = "css_repository/logstash-log-2"
  bucket     = huaweicloud_obs_bucket.bucket2.bucket
  period     = "00:00 GMT+08:00"
}
`, testAccLogstashCluster_basic(name, 1, "bar"))
}

func testLogSetting_logstash_updateNull(name string) string {
	return fmt.Sprintf(`
%s

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

resource "huaweicloud_css_log_setting" "test" {
  cluster_id = huaweicloud_css_logstash_cluster.test.id
  agency     = "css_obs_agency"
  base_path  = "css_repository/logstash-log-2"
  bucket     = huaweicloud_obs_bucket.bucket2.bucket
  period     = ""
}
`, testAccLogstashCluster_basic(name, 1, "bar"))
}
