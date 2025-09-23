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

func getHlsConfigurationFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region     = acceptance.HW_REGION_NAME
		domainName = state.Primary.Attributes["domain_name"]
	)
	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return nil, fmt.Errorf("error creating Live client: %s", err)
	}

	getRespBody, err := live.ReadHlsConfiguration(client, domainName)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Live HLS configuration: %s", err)
	}

	return getRespBody, nil
}

func TestAccHlsConfiguration_basic(t *testing.T) {
	var (
		hlsConfigurationObj interface{}
		rName               = "huaweicloud_live_hls_configuration.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&hlsConfigurationObj,
		getHlsConfigurationFunc,
	)

	// Avoid CheckDestroy, because there is nothing in the resource destroy method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLiveIngestDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccHlsConfiguration_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_LIVE_INGEST_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "application.#", "1"),
					resource.TestCheckResourceAttr(rName, "application.0.name", "live"),
					resource.TestCheckResourceAttr(rName, "application.0.hls_fragment", "5"),
					resource.TestCheckResourceAttr(rName, "application.0.hls_ts_count", "5"),
					resource.TestCheckResourceAttr(rName, "application.0.hls_min_frags", "5"),
				),
			},
			{
				Config: testAccHlsConfiguration_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "application.0.name", "live"),
					resource.TestCheckResourceAttr(rName, "application.0.hls_fragment", "4"),
					resource.TestCheckResourceAttr(rName, "application.0.hls_ts_count", "4"),
					resource.TestCheckResourceAttr(rName, "application.0.hls_min_frags", "4"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateIdFunc: testAccHlsConfigurationImportState(rName),
			},
		},
	})
}

func testAccHlsConfiguration_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_live_hls_configuration" "test" {
  domain_name = "%s"

  application {
    name          = "live"
    hls_fragment  = 5
    hls_ts_count  = 5
    hls_min_frags = 5
  }
}
`, acceptance.HW_LIVE_INGEST_DOMAIN_NAME)
}

func testAccHlsConfiguration_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_live_hls_configuration" "test" {
  domain_name = "%s"

  application {
    name          = "live"
    hls_fragment  = 4
    hls_ts_count  = 4
    hls_min_frags = 4
  }
}
`, acceptance.HW_LIVE_INGEST_DOMAIN_NAME)
}

func testAccHlsConfigurationImportState(rName string) resource.ImportStateIdFunc {
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
