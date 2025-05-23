package antiddos

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDomainResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		product           = "aad"
		getDomainsHttpUrl = "v1/aad/protected-domains"
	)
	client, err := cfg.NewServiceClient(product, acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating AAD client: %s", err)
	}

	getDomainsPath := client.Endpoint + getDomainsHttpUrl

	getDomainsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getDomainsResp, err := client.Request("GET", getDomainsPath, &getDomainsOpt)
	if err != nil {
		return nil, fmt.Errorf("error getting AAD domains: %s", err)
	}

	respBody, err := utils.FlattenResponse(getDomainsResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening AAD domains response: %s", err)
	}

	ids := fmt.Sprintf("items[?domain_id=='%s']|[0]", state.Primary.ID)
	domains := utils.PathSearch(ids, respBody, nil)
	if domains == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return domains, nil
}

func TestAccResourceDomain_basic(t *testing.T) {
	var (
		object       interface{}
		resourceName = "huaweicloud_aad_domain.test0"
		rc           = acceptance.InitResourceCheck(resourceName, &object, getDomainResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare the AAD instance created and the domain which is record-keeping before running this test.
			acceptance.TestAccPreCheckAADIncetance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDomain_base(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "cname"),
					resource.TestCheckResourceAttrSet(resourceName, "waf_status"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_ids.#"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", acceptance.HW_AAD_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "vips.0", acceptance.HW_AAD_IP_ADDRESS),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "real_server_type", "0"),
					resource.TestCheckResourceAttr(resourceName, "real_server", "124.1.5.9"),
					resource.TestCheckResourceAttr(resourceName, "port_http.0", "140"),
					resource.TestCheckResourceAttr(resourceName, "port_https.0", "9999"),
					resource.TestCheckResourceAttr(resourceName, "protocol.#", "2"),
				),
			},
			{
				Config: testAccResourceDomain_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "real_server_type", "1"),
					resource.TestCheckResourceAttr(resourceName, "real_server", "change.test.com"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"vips",
					"port_https",
					"port_http",
				},
			},
		},
	})
}

func testAccResourceDomain_base() string {
	return fmt.Sprintf(`
resource "huaweicloud_aad_domain" "test0" {
  domain_name           = "%[1]s"
  enterprise_project_id = 0
  vips                  = ["%[2]s"]
  real_server_type      = 0
  real_server           = "124.1.5.9"

  port_http = [
    140
  ]

  port_https = [
    9999
  ]
}
`, acceptance.HW_AAD_DOMAIN_NAME, acceptance.HW_AAD_IP_ADDRESS)
}

func testAccResourceDomain_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_aad_domain" "test0" {
  domain_name           = "%[1]s"
  enterprise_project_id = 0
  vips                  = ["%[2]s"]
  real_server_type      = 1
  real_server           = "change.test.com"

  port_http  = [
    140
  ]

  port_https = [
    9999
  ]
}
`, acceptance.HW_AAD_DOMAIN_NAME, acceptance.HW_AAD_IP_ADDRESS)
}
