package live

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	acl "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/live"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getResourceIpAclFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region     = acceptance.HW_REGION_NAME
		product    = "live"
		domainName = state.Primary.Attributes["domain_name"]
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Live client: %s", err)
	}

	respBody, err := acl.ReadIPAddressAcl(client, domainName)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Live IP address acl: %s", err)
	}

	authType := utils.PathSearch("auth_type", respBody, "").(string)
	if authType == "NONE" {
		return nil, golangsdk.ErrDefault404{}
	}

	return respBody, nil
}

func TestAccResourceIpAcl_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_live_ip_acl.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getResourceIpAclFunc,
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
				Config: testResourceIpAcl_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "auth_type", "WHITE"),
					resource.TestCheckResourceAttr(rName, "ip_auth_list", "192.168.0.0;192.168.0.8;127.0.0.1/24"),
				),
			},
			{
				Config: testResourceIpAcl_basic_update1(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "auth_type", "WHITE"),
					resource.TestCheckResourceAttr(rName, "ip_auth_list", "192.168.0.0;192.168.0.8"),
				),
			},
			{
				Config: testResourceIpAcl_basic_update2(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "auth_type", "BLACK"),
					resource.TestCheckResourceAttr(rName, "ip_auth_list", "192.168.0.0"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateIdFunc: testAccIPAclImportState(rName),
			},
		},
	})
}

func testResourceIpAcl_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_live_ip_acl" "test" {
  domain_name  = "%s"
  auth_type    = "WHITE"
  ip_auth_list = "192.168.0.0;192.168.0.8;127.0.0.1/24"
}
`, acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
}

func testResourceIpAcl_basic_update1() string {
	return fmt.Sprintf(`
resource "huaweicloud_live_ip_acl" "test" {
  domain_name  = "%s"
  auth_type    = "WHITE"
  ip_auth_list = "192.168.0.0;192.168.0.8"
}
`, acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
}

func testResourceIpAcl_basic_update2() string {
	return fmt.Sprintf(`
resource "huaweicloud_live_ip_acl" "test" {
  domain_name  = "%s"
  auth_type    = "BLACK"
  ip_auth_list = "192.168.0.0"
}
`, acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
}

func testAccIPAclImportState(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", rName)
		}

		domainName := rs.Primary.Attributes["domain_name"]
		if domainName == "" {
			return "", fmt.Errorf("the imported ID format is invalid, 'domain_name' is empty")
		}
		return domainName, nil
	}
}
