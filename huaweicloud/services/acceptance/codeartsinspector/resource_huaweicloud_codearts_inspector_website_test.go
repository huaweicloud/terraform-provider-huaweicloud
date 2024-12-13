package codeartsinspector

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

func getInspectorWebsiteResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v3/{project_id}/webscan/domains"
		product = "vss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CodeArts inspector client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += fmt.Sprintf("?domain_id=%s", state.Primary.ID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CodeArts inspector website: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}
	domainInfo := utils.PathSearch("domains|[0]", getRespBody, nil)
	if domainInfo == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return domainInfo, nil
}

func TestAccInspectorWebsite_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_inspector_website.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getInspectorWebsiteResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testInspectorWebsite_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "website_name", name),
					resource.TestCheckResourceAttr(rName, "auth_type", "free"),
					resource.TestCheckResourceAttr(rName, "website_address", "https://demo.test.com"),
					resource.TestCheckResourceAttr(rName, "login_url", "https://demo.test.com/login"),
					resource.TestCheckResourceAttr(rName, "login_username", "test-name"),
					resource.TestCheckResourceAttr(rName, "login_password", "test-password"),
					resource.TestCheckResourceAttr(rName, "login_cookie", "test-cookie"),
					resource.TestCheckResourceAttr(rName, "verify_url", "https://verify.cn"),
					resource.TestCheckResourceAttr(rName, "http_headers.test-key1", "test-value1"),
					resource.TestCheckResourceAttr(rName, "http_headers.test-key2", "test-value2"),
					resource.TestCheckResourceAttrSet(rName, "high"),
					resource.TestCheckResourceAttrSet(rName, "middle"),
					resource.TestCheckResourceAttrSet(rName, "low"),
					resource.TestCheckResourceAttrSet(rName, "hint"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "auth_status"),
				),
			},
			{
				Config: testInspectorWebsite_basic_update1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "website_name", name),
					resource.TestCheckResourceAttr(rName, "auth_type", "free"),
					resource.TestCheckResourceAttr(rName, "website_address", "https://demo.test.com"),
					resource.TestCheckResourceAttr(rName, "login_url", "https://demo.test.com/register"),
					resource.TestCheckResourceAttr(rName, "login_username", "name-update"),
					resource.TestCheckResourceAttr(rName, "login_password", "password-update"),
					resource.TestCheckResourceAttr(rName, "login_cookie", "cookie-update"),
					resource.TestCheckResourceAttr(rName, "verify_url", "https://verify-update.cn"),
					resource.TestCheckResourceAttr(rName, "http_headers.test-key3", "test-value3"),
					resource.TestCheckResourceAttrSet(rName, "high"),
					resource.TestCheckResourceAttrSet(rName, "middle"),
					resource.TestCheckResourceAttrSet(rName, "low"),
					resource.TestCheckResourceAttrSet(rName, "hint"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "auth_status"),
				),
			},
			{
				Config: testInspectorWebsite_basic_update2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "website_name", name),
					resource.TestCheckResourceAttr(rName, "auth_type", "free"),
					resource.TestCheckResourceAttr(rName, "website_address", "https://demo.test.com"),
					resource.TestCheckResourceAttr(rName, "login_url", ""),
					resource.TestCheckResourceAttr(rName, "login_username", ""),
					resource.TestCheckResourceAttr(rName, "login_password", ""),
					resource.TestCheckResourceAttr(rName, "login_cookie", ""),
					resource.TestCheckResourceAttr(rName, "verify_url", ""),
					resource.TestCheckResourceAttr(rName, "http_headers.#", "0"),
					resource.TestCheckResourceAttrSet(rName, "high"),
					resource.TestCheckResourceAttrSet(rName, "middle"),
					resource.TestCheckResourceAttrSet(rName, "low"),
					resource.TestCheckResourceAttrSet(rName, "hint"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "auth_status"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"auth_type",
					"login_password",
					"login_cookie",
					"http_headers",
				},
			},
		},
	})
}

func TestAccInspectorWebsite_publicIPAddress(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_inspector_website.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getInspectorWebsiteResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// make sure the IP address is a public IPv4 address
			acceptance.TestAccPreCheckCodeArtsPublicIPAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testInspectorWebsite_publicIPAddress(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "website_name", name),
					resource.TestCheckResourceAttr(rName, "auth_type", "free"),
					resource.TestCheckResourceAttr(rName, "website_address", fmt.Sprintf("http://%s",
						acceptance.HW_CODEARTS_PUBLIC_IP_ADDRESS)),
					resource.TestCheckResourceAttr(rName, "login_url", fmt.Sprintf("http://%s/login",
						acceptance.HW_CODEARTS_PUBLIC_IP_ADDRESS)),
					resource.TestCheckResourceAttr(rName, "login_username", "test-name"),
					resource.TestCheckResourceAttr(rName, "login_password", "test-password"),
					resource.TestCheckResourceAttr(rName, "login_cookie", "test-cookie"),
					resource.TestCheckResourceAttr(rName, "verify_url", "https://verify.cn"),
					resource.TestCheckResourceAttr(rName, "http_headers.test-key1", "test-value1"),
					resource.TestCheckResourceAttr(rName, "http_headers.test-key2", "test-value2"),
					resource.TestCheckResourceAttrSet(rName, "high"),
					resource.TestCheckResourceAttrSet(rName, "middle"),
					resource.TestCheckResourceAttrSet(rName, "low"),
					resource.TestCheckResourceAttrSet(rName, "hint"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "auth_status"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"auth_type",
					"login_password",
					"login_cookie",
					"http_headers",
				},
			},
		},
	})
}

func testInspectorWebsite_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_codearts_inspector_website" "test" {
  website_name    = "%s"
  auth_type       = "free"
  website_address = "https://demo.test.com"
  login_url       = "https://demo.test.com/login"
  login_username  = "test-name"
  login_password  = "test-password"
  login_cookie    = "test-cookie"
  verify_url      = "https://verify.cn"

  http_headers = {
    "test-key1" = "test-value1"
    "test-key2" = "test-value2"
  }
}
`, name)
}

func testInspectorWebsite_basic_update1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_codearts_inspector_website" "test" {
  website_name    = "%s"
  auth_type       = "free"
  website_address = "https://demo.test.com"
  login_url       = "https://demo.test.com/register"
  login_username  = "name-update"
  login_password  = "password-update"
  login_cookie    = "cookie-update"
  verify_url      = "https://verify-update.cn"

  http_headers = {
    "test-key3" = "test-value3"
  }
}
`, name)
}

func testInspectorWebsite_basic_update2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_codearts_inspector_website" "test" {
  website_name    = "%s"
  auth_type       = "free"
  website_address = "https://demo.test.com"
}
`, name)
}

func testInspectorWebsite_publicIPAddress(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_codearts_inspector_website" "test" {
  website_name    = "%[1]s"
  auth_type       = "free"
  website_address = "http://%[2]s"
  login_url       = "http://%[2]s/login"
  login_username  = "test-name"
  login_password  = "test-password"
  login_cookie    = "test-cookie"
  verify_url      = "https://verify.cn"

  http_headers = {
    "test-key1" = "test-value1"
    "test-key2" = "test-value2"
  }
}
`, name, acceptance.HW_CODEARTS_PUBLIC_IP_ADDRESS)
}
