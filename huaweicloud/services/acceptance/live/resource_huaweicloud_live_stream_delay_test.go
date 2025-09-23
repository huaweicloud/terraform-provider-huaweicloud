package live

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/live"
)

func getResourceStreamDelayFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region     = acceptance.HW_REGION_NAME
		domainName = state.Primary.Attributes["domain_name"]
	)
	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return nil, fmt.Errorf("error creating Live client: %s", err)
	}

	respBody, err := live.ReadStreamDelay(client, domainName)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Live stream delay time: %s", err)
	}

	return respBody, nil
}

func TestAccResourceStreamDelay_basic(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_live_stream_delay.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getResourceStreamDelayFunc,
	)

	// Avoid CheckDestroy, because there is nothing in the resource destroy method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLiveStreamingDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceStreamDelay_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "domain_name", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "delay", "2000"),
					resource.TestCheckResourceAttr(resourceName, "app_name", "live"),
				),
			},
			{
				Config: testResourceStreamDelay_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "domain_name", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "delay", "4000"),
					resource.TestCheckResourceAttr(resourceName, "app_name", "live"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testStreamDelayImportState(resourceName),
			},
		},
	})
}

func testResourceStreamDelay_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_live_stream_delay" "test" {
  domain_name = "%s"
  delay       = 2000
  app_name    = "live"
}
`, acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
}

func testResourceStreamDelay_basic_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_live_stream_delay" "test" {
  domain_name = "%s"
  delay       = 4000
  app_name    = "live"
}
`, acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
}

func testStreamDelayImportState(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var domainName string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", rName)
		}

		domainName = rs.Primary.Attributes["domain_name"]
		if domainName == "" {
			return "", fmt.Errorf("invalid format specified for import ID, `domain_name` is empty")
		}
		return domainName, nil
	}
}
