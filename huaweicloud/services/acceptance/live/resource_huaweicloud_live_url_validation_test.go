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

func getUrlValidationFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return nil, fmt.Errorf("error creating Live client: %s", err)
	}

	getHttpUrl := "v1/{project_id}/guard/key-chain"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = fmt.Sprintf("%s?domain=%v", getPath, state.Primary.Attributes["domain_name"])

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, fmt.Errorf("error retrieving URL validation: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	key := utils.PathSearch("key", getRespBody, "").(string)
	if key == "" {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func TestAccURLValidation_basic(t *testing.T) {
	var (
		urlValidationObj interface{}
		rName            = "huaweicloud_live_url_validation.test"
		domainName       = fmt.Sprintf("%s.huaweicloud.com", acceptance.RandomAccResourceNameWithDash())
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&urlValidationObj,
		getUrlValidationFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccUrlValidation_basic(domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "domain_name", "huaweicloud_live_domain.test", "name"),
					resource.TestCheckResourceAttr(rName, "key", "IbBIzklRGCyMEd18oPV9sxAuuwNIzT81"),
					resource.TestCheckResourceAttr(rName, "auth_type", "d_sha256"),
					resource.TestCheckResourceAttr(rName, "timeout", "800"),
				),
			},
			{
				Config: testAccUrlValidation_update(domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "key", "IbBIzklRGCyMEd74oMS9sxAuuwNIzT59"),
					resource.TestCheckResourceAttr(rName, "auth_type", "c_aes"),
					resource.TestCheckResourceAttr(rName, "timeout", "1200"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateIdFunc: testAccUrlValidationImportState(rName),
			},
		},
	})
}

func testAccUrlValidation_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_live_domain" "test" {
  name = "%s"
  type = "push"
}

resource "huaweicloud_live_url_validation" "test" {
  domain_name = huaweicloud_live_domain.test.name
  key         = "IbBIzklRGCyMEd18oPV9sxAuuwNIzT81"
  auth_type   = "d_sha256"
  timeout     = 800
}
`, name)
}

func testAccUrlValidation_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_live_domain" "test" {
  name = "%s"
  type = "push"
}

resource "huaweicloud_live_url_validation" "test" {
  domain_name = huaweicloud_live_domain.test.name
  key         = "IbBIzklRGCyMEd74oMS9sxAuuwNIzT59"
  auth_type   = "c_aes"
  timeout     = 1200
}
`, name)
}

func testAccUrlValidationImportState(rName string) resource.ImportStateIdFunc {
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
