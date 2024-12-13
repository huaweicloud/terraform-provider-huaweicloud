package live

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	blocking "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/live"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getResourceGeoBlockingFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region     = acceptance.HW_REGION_NAME
		product    = "live"
		domainName = state.Primary.Attributes["domain_name"]
		app        = state.Primary.Attributes["app_name"]
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Live client: %s", err)
	}

	respBody, err := blocking.ReadGeoBlocking(client, domainName)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Live geo blocking: %s", err)
	}

	expression := fmt.Sprintf("apps[?app == '%s']|[0].area_whitelist", app)
	areaWhitelist := utils.PathSearch(expression, respBody, nil)
	if areaWhitelist == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return respBody, nil
}

func TestAccResourceGeoBlocking_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_live_geo_blocking.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getResourceGeoBlockingFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLiveStreamingDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceGeoBlocking_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "app_name", "live"),
					resource.TestCheckResourceAttr(rName, "area_whitelist.#", "5"),
				),
			},
			{
				Config: testResourceGeoBlocking_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "app_name", "live"),
					resource.TestCheckResourceAttr(rName, "area_whitelist.#", "3"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateIdFunc: testAccGeoBlockingImportState(rName),
			},
		},
	})
}

func testResourceGeoBlocking_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_live_geo_blocking" "test" {
  domain_name    = "%s"
  app_name       = "live"
  area_whitelist = ["AE", "AF", "CN-IN", "CN-HK", "CN-MO"]
}
`, acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
}

func testResourceGeoBlocking_basic_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_live_geo_blocking" "test" {
  domain_name    = "%s"
  app_name       = "live"
  area_whitelist = ["AE", "AF", "CN-IN"]
}
`, acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
}

func testAccGeoBlockingImportState(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", rName)
		}

		domainName := rs.Primary.Attributes["domain_name"]
		appName := rs.Primary.Attributes["app_name"]
		if domainName == "" || appName == "" {
			return "", fmt.Errorf("the imported ID format is invalid, want '<domain_name>/<app_name>',"+
				" but got '%s/%s'", domainName, appName)
		}
		return fmt.Sprintf("%s/%s", domainName, appName), nil
	}
}
