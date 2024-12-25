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

func getNotificationConfigFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return nil, fmt.Errorf("error creating Live client: %s", err)
	}

	getHttpUrl := "v1/{project_id}/notifications/publish"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = fmt.Sprintf("%s?domain=%v", getPath, state.Primary.Attributes["domain_name"])

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, fmt.Errorf("error retrieving notification configuration: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	key := utils.PathSearch("url", getRespBody, "").(string)
	if key == "" {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func TestAccNotificationConfiguration_basic(t *testing.T) {
	var (
		notifyConfigObj interface{}
		rName           = "huaweicloud_live_notification_configuration.test"
		domainName      = fmt.Sprintf("%s.huaweicloud.com", acceptance.RandomAccResourceNameWithDash())
		authKey1        = "8DDEtfv0ZE3Z3EAUokupq2Zf2NAgyL5vSZAmlgmX5jzBd2fHohrA9u727I9U0RZR8mHsxTnwDBKXNUHw52NmA1iZHHitlXZ"
		authkey2        = "RHzXPbi4Ll0Sr4IvwEvBKzn6tEP1pIt3vtGUAdlJU0kxOzYbKXkJpTVWd2Z2ZdPDTu2koXalAXRc8o3HVp8K8S5rjAAtFE"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&notifyConfigObj,
		getNotificationConfigFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationConfig_basic(domainName, authKey1),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "domain_name", "huaweicloud_live_domain.test", "name"),
					resource.TestCheckResourceAttr(rName, "url", "http://mycallback.com/notify_config"),
					resource.TestCheckResourceAttr(rName, "auth_sign_key", authKey1),
					resource.TestCheckResourceAttr(rName, "call_back_area", "mainland_china"),
				),
			},
			{
				Config: testAccNotificationConfig_update(domainName, authkey2),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "url", "https://mycallback.com.cn/notify_config"),
					resource.TestCheckResourceAttr(rName, "auth_sign_key", authkey2),
					resource.TestCheckResourceAttr(rName, "call_back_area", "outside_mainland_china"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateIdFunc: testAccNotificationConfigImportState(rName),
			},
		},
	})
}

func testAccNotificationConfig_basic(name, signKey string) string {
	return fmt.Sprintf(`
resource "huaweicloud_live_domain" "test" {
  name = "%[1]s"
  type = "push"
}

resource "huaweicloud_live_notification_configuration" "test" {
  domain_name    = huaweicloud_live_domain.test.name
  url            = "http://mycallback.com/notify_config"
  auth_sign_key  = "%[2]s"
  call_back_area = "mainland_china"
}
`, name, signKey)
}

func testAccNotificationConfig_update(name, signKey string) string {
	return fmt.Sprintf(`
resource "huaweicloud_live_domain" "test" {
  name = "%[1]s"
  type = "push"
}

resource "huaweicloud_live_notification_configuration" "test" {
  domain_name    = huaweicloud_live_domain.test.name
  url            = "https://mycallback.com.cn/notify_config"
  auth_sign_key  = "%[2]s"
  call_back_area = "outside_mainland_china"
}
`, name, signKey)
}

func testAccNotificationConfigImportState(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var domainName string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", rName)
		}

		domainName = rs.Primary.Attributes["domain_name"]
		if domainName == "" {
			return "", fmt.Errorf("invalid format specified for import ID, 'domain_name' is empty")
		}
		return domainName, nil
	}
}
