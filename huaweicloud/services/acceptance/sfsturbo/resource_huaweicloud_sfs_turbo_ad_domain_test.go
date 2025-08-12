package sfsturbo

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

func getSFSTurboAdDomainResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "sfs-turbo"
		httpUrl = "v1/{project_id}/sfs-turbo/shares/{share_id}/fs/active-directory-domain"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating SFS Turbo client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{share_id}", state.Primary.ID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error getting SFS Turbo AD domain: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("error parsing response: %s", err)
	}
	return respBody, nil
}

func TestAccSFSTurboAdDomain_basic(t *testing.T) {
	var obj interface{}
	organizationUnit := "cn=Computers,dc=" + acceptance.HW_SFS_TURBO_AD_DOMAIN_NAME + ",dc=com"
	resourceName := "huaweicloud_sfs_turbo_ad_domain.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getSFSTurboAdDomainResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckSFSTurboShareId(t)
			acceptance.TestAccPrecheckSFSTurboADDomin(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurboAdDomain_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "domain_name", fmt.Sprintf("%s.com", acceptance.HW_SFS_TURBO_AD_DOMAIN_NAME)),
					resource.TestCheckResourceAttr(resourceName, "system_name", "sfs"),
					resource.TestCheckResourceAttr(resourceName, "service_account", "administrator"),
					resource.TestCheckResourceAttr(resourceName, "dns_server.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dns_server.0", acceptance.HW_SFS_TURBO_AD_DOMAIN_DNS_SERVER_IP),
					resource.TestCheckResourceAttr(resourceName, "organization_unit", organizationUnit),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				Config: testAccSFSTurboAdDomain_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "domain_name", fmt.Sprintf("%s.com", acceptance.HW_SFS_TURBO_AD_DOMAIN_NAME)),
					resource.TestCheckResourceAttr(resourceName, "system_name", "sfs-update"),
					resource.TestCheckResourceAttr(resourceName, "service_account", "administrator"),
					resource.TestCheckResourceAttr(resourceName, "dns_server.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dns_server.0", acceptance.HW_SFS_TURBO_AD_DOMAIN_DNS_SERVER_IP),
					resource.TestCheckResourceAttr(resourceName, "organization_unit", organizationUnit),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"share_id", "service_account", "password", "overwrite_same_account",
				},
			},
		},
	})
}

func testAccSFSTurboAdDomain_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_sfs_turbo_ad_domain" "test" {
  share_id               = "%[1]s"
  service_account        = "administrator"
  password               = "%[2]s"
  domain_name            = "%[3]s.com"
  system_name            = "sfs"
  dns_server             = ["%[4]s"]
  overwrite_same_account = false
  organization_unit      = "cn=Computers,dc=%[5]s,dc=com"
}
`, acceptance.HW_SFS_TURBO_SHARE_ID, acceptance.HW_SFS_TURBO_AD_DOMAIN_PW, acceptance.HW_SFS_TURBO_AD_DOMAIN_NAME,
		acceptance.HW_SFS_TURBO_AD_DOMAIN_DNS_SERVER_IP, acceptance.HW_SFS_TURBO_AD_DOMAIN_NAME)
}

func testAccSFSTurboAdDomain_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_sfs_turbo_ad_domain" "test" {
  share_id          = "%[1]s"
  service_account   = "administrator"
  password          = "%[2]s"
  domain_name       = "%[3]s.com"
  system_name       = "sfs-update"
  dns_server        = ["%[4]s"]
  organization_unit = "cn=Computers,dc=%[5]s,dc=com"
}
`, acceptance.HW_SFS_TURBO_SHARE_ID, acceptance.HW_SFS_TURBO_AD_DOMAIN_PW, acceptance.HW_SFS_TURBO_AD_DOMAIN_NAME,
		acceptance.HW_SFS_TURBO_AD_DOMAIN_DNS_SERVER_IP, acceptance.HW_SFS_TURBO_AD_DOMAIN_NAME)
}
