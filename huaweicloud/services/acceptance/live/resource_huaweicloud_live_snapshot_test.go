package live

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

func getLiveSnapshotResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v1/{project_id}/stream/snapshot"
		product = "live"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Live client: %s", err)
	}

	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ID format, want '<domain_name>/<app_name>', but got '%s'", state.Primary.ID)
	}
	domainName := parts[0]
	appName := parts[1]

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildGetLiveSnapshotQueryParams(domainName, appName)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Live snapshot: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	snapshot := utils.PathSearch("snapshot_config_list|[0]", respBody, nil)
	if snapshot == nil {
		return nil, fmt.Errorf("error get live snapshot")
	}

	return snapshot, nil
}

func buildGetLiveSnapshotQueryParams(domainName, appName string) string {
	return fmt.Sprintf("?domain=%s&app_name=%s", domainName, appName)
}

func TestAccLiveSnapshot_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceNameWithDash()
	domainName := fmt.Sprintf("%s.huaweicloud.com", name)
	rName := "huaweicloud_live_snapshot.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getLiveSnapshotResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testLiveSnapshot_basic(domainName, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "domain_name", domainName),
					resource.TestCheckResourceAttr(rName, "app_name", "live"),
					resource.TestCheckResourceAttr(rName, "frequency", "10"),
					resource.TestCheckResourceAttr(rName, "storage_mode", "0"),
					resource.TestCheckResourceAttr(rName, "storage_bucket", name),
					resource.TestCheckResourceAttr(rName, "storage_path", "input/"),
					resource.TestCheckResourceAttr(rName, "call_back_enabled", "off"),
				),
			},
			{
				Config: testLiveSnapshot_basic_update(domainName, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "domain_name", domainName),
					resource.TestCheckResourceAttr(rName, "app_name", "live"),
					resource.TestCheckResourceAttr(rName, "frequency", "20"),
					resource.TestCheckResourceAttr(rName, "storage_mode", "1"),
					resource.TestCheckResourceAttr(rName, "storage_bucket", name+"-update"),
					resource.TestCheckResourceAttr(rName, "storage_path", "output/"),
					resource.TestCheckResourceAttr(rName, "call_back_enabled", "on"),
					resource.TestCheckResourceAttr(rName, "call_back_url", "https://test.com"),
					resource.TestCheckResourceAttr(rName, "call_back_auth_key",
						"ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"),
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

func testLiveSnapshot_basic(domainName, name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_live_domain" "test" {
  name = "%[1]s"
  type = "push"
}

resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[2]s"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_live_bucket_authorization" "test" {
  depends_on = [huaweicloud_obs_bucket.test]

  bucket = "%[2]s"
}

resource "huaweicloud_live_snapshot" "test" {
  depends_on = [huaweicloud_live_bucket_authorization.test]

  domain_name       = huaweicloud_live_domain.test.name
  app_name          = "live"
  frequency         = 10
  storage_mode      = 0
  storage_bucket    = huaweicloud_obs_bucket.test.bucket
  storage_path      = "input/"
  call_back_enabled = "off"
}
`, domainName, name)
}

func testLiveSnapshot_basic_update(domainName, name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_live_domain" "test" {
  name = "%[1]s"
  type = "push"
}

resource "huaweicloud_obs_bucket" "test_update" {
  bucket        = "%[2]s-update"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_live_bucket_authorization" "test_update" {
  depends_on = [huaweicloud_obs_bucket.test_update]

  bucket = "%[2]s-update"
}

resource "huaweicloud_live_snapshot" "test" {
  depends_on = [huaweicloud_live_bucket_authorization.test_update]

  domain_name        = huaweicloud_live_domain.test.name
  app_name           = "live"
  frequency          = 20
  storage_mode       = 1
  storage_bucket     = huaweicloud_obs_bucket.test_update.bucket
  storage_path       = "output/"
  call_back_enabled  = "on"
  call_back_url      = "https://test.com"
  call_back_auth_key = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
}
`, domainName, name)
}
