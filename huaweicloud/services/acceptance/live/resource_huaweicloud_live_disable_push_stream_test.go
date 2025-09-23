package live

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/live"
)

func getDisablePushStreamFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return nil, fmt.Errorf("error creating Live client: %s", err)
	}

	getRespBody, err := live.GetDisablePushStream(client, state.Primary.Attributes["domain_name"],
		state.Primary.Attributes["app_name"], state.Primary.Attributes["stream_name"])
	if err != nil {
		return nil, fmt.Errorf("error retrieving disabled push stream: %s", err)
	}

	return getRespBody, nil
}

func TestAccDisablePushStream_basic(t *testing.T) {
	var (
		disablePushStreamObj interface{}
		rName                = "huaweicloud_live_disable_push_stream.test"
		domainName           = fmt.Sprintf("%s.huaweicloud.com", acceptance.RandomAccResourceNameWithDash())
		createTime           = time.Now().UTC().Add(24 * time.Hour).Format("2006-01-02T15:04:05Z")
		updateTime           = time.Now().UTC().Add(48 * time.Hour).Format("2006-01-02T15:04:05Z")
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&disablePushStreamObj,
		getDisablePushStreamFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDisablePushStream_basic(domainName, createTime),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "domain_name", "huaweicloud_live_domain.test", "name"),
					resource.TestCheckResourceAttr(rName, "app_name", "live"),
					resource.TestCheckResourceAttr(rName, "stream_name", "tf-test"),
					resource.TestCheckResourceAttr(rName, "resume_time", createTime),
				),
			},
			{
				Config: testAccDisablePushStream_update(domainName, updateTime),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "resume_time", updateTime),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateIdFunc: testAccDisablePushStreamImportState(rName),
			},
		},
	})
}

func testAccDisablePushStream_basic(name, nowTime string) string {
	return fmt.Sprintf(`
resource "huaweicloud_live_domain" "test" {
  name = "%[1]s"
  type = "push"
}

resource "huaweicloud_live_disable_push_stream" "test" {
  domain_name = huaweicloud_live_domain.test.name
  app_name    = "live"
  stream_name = "tf-test"
  resume_time = "%[2]s"
}
`, name, nowTime)
}

func testAccDisablePushStream_update(name, updateTime string) string {
	return fmt.Sprintf(`
resource "huaweicloud_live_domain" "test" {
  name = "%[1]s"
  type = "push"
}

resource "huaweicloud_live_disable_push_stream" "test" {
  domain_name = huaweicloud_live_domain.test.name
  app_name    = "live"
  stream_name = "tf-test"
  resume_time = "%[2]s"
}
`, name, updateTime)
}

func testAccDisablePushStreamImportState(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var domainName, appName, streamName string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", rName)
		}

		domainName = rs.Primary.Attributes["domain_name"]
		appName = rs.Primary.Attributes["app_name"]
		streamName = rs.Primary.Attributes["stream_name"]
		if domainName == "" || appName == "" || streamName == "" {
			return "", fmt.Errorf("invalid format specified for import ID, "+
				"want '<domain_name>/<app_name>/<stream_name>',but got '%s/%s/%s'", domainName, appName, streamName)
		}
		return fmt.Sprintf("%s/%s/%s", domainName, appName, streamName), nil
	}
}
