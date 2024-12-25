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

func getOriginPullConfigurationFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region     = acceptance.HW_REGION_NAME
		domainName = state.Primary.Attributes["domain_name"]
	)
	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return nil, fmt.Errorf("error creating Live client: %s", err)
	}

	getRespBody, err := live.ReadOriginPullConfiguration(client, domainName)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Live origin pull configuration: %s", err)
	}

	return getRespBody, nil
}

func TestAccOriginPullConfiguration_basic(t *testing.T) {
	var (
		originPullConfigurationObj interface{}
		rName                      = "huaweicloud_live_origin_pull_configuration.test"
		domainName                 = fmt.Sprintf("%s.huaweicloud.com", acceptance.RandomAccResourceNameWithDash())
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&originPullConfigurationObj,
		getOriginPullConfigurationFunc,
	)

	// Avoid CheckDestroy, because there is nothing in the resource destroy method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOriginPullConfiguration_basic(domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "domain_name", domainName),
					resource.TestCheckResourceAttr(rName, "source_type", "domain"),
					resource.TestCheckResourceAttr(rName, "sources.#", "2"),
					resource.TestCheckResourceAttr(rName, "source_port", "8888"),
					resource.TestCheckResourceAttr(rName, "scheme", "rtmp"),
					resource.TestCheckResourceAttr(rName, "additional_args.param1", "value1"),
					resource.TestCheckResourceAttr(rName, "additional_args.param2", "value2"),
				),
			},
			{
				Config: testAccOriginPullConfiguration_update1(domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "domain_name", domainName),
					resource.TestCheckResourceAttr(rName, "source_type", "ipaddr"),
					resource.TestCheckResourceAttr(rName, "sources_ip.#", "2"),
					resource.TestCheckResourceAttr(rName, "source_port", "9999"),
					resource.TestCheckResourceAttr(rName, "scheme", "http"),
					resource.TestCheckResourceAttr(rName, "additional_args.param1", "value1"),
					resource.TestCheckResourceAttr(rName, "additional_args.param2", "value2"),
					resource.TestCheckResourceAttr(rName, "additional_args.param3", "value3"),
				),
			},
			{
				Config: testAccOriginPullConfiguration_update2(domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "domain_name", domainName),
					resource.TestCheckResourceAttr(rName, "source_type", "huawei"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"additional_args"},
				ImportStateIdFunc:       testAccOriginPullConfigurationImportState(rName),
			},
		},
	})
}

func testAccOriginPullConfiguration_base(domainName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_live_domain" "test" {
  name = "%s"
  type = "pull"
}
`, domainName)
}

func testAccOriginPullConfiguration_basic(domainName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_live_origin_pull_configuration" "test" {
  domain_name = huaweicloud_live_domain.test.name
  source_type = "domain"
  sources     = ["play.tftest.huaweicloud1.com", "play.tftest.huaweicloud2.com"]
  source_port = 8888
  scheme      = "rtmp"

  additional_args = {
    param1 = "value1"
    param2 = "value2"
  }
}
`, testAccOriginPullConfiguration_base(domainName))
}

func testAccOriginPullConfiguration_update1(domainName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_live_origin_pull_configuration" "test" {
  domain_name = huaweicloud_live_domain.test.name
  source_type = "ipaddr"
  sources     = ["play.tftest.huaweicloud1.com"]
  sources_ip  = ["192.127.0.124", "192.168.0.123"]
  source_port = 9999
  scheme      = "http"

  additional_args = {
    param1 = "value1"
    param3 = "value3"
    param2 = "value2"
  }
}
`, testAccOriginPullConfiguration_base(domainName))
}

func testAccOriginPullConfiguration_update2(domainName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_live_origin_pull_configuration" "test" {
  domain_name = huaweicloud_live_domain.test.name
  source_type = "huawei"
}
`, testAccOriginPullConfiguration_base(domainName))
}

func testAccOriginPullConfigurationImportState(rName string) resource.ImportStateIdFunc {
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
