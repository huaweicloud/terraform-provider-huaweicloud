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

func getRefererValidationFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return nil, fmt.Errorf("error creating Live client: %s", err)
	}

	getHttpUrl := "v1/{project_id}/guard/referer-chain"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = fmt.Sprintf("%s?domain=%v", getPath, state.Primary.Attributes["domain_name"])

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, fmt.Errorf("error retrieving referer validation: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	refererAuthList := utils.PathSearch("referer_auth_list", getRespBody, make([]interface{}, 0)).([]interface{})
	if len(refererAuthList) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func TestAccRefererValidation_basic(t *testing.T) {
	var (
		refererValidationObj interface{}
		rName                = "huaweicloud_live_referer_validation.test"
		domainName           = fmt.Sprintf("%s.huaweicloud.com", acceptance.RandomAccResourceNameWithDash())
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&refererValidationObj,
		getRefererValidationFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRefererValidation_basic(domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "domain_name", "huaweicloud_live_domain.test", "name"),
					resource.TestCheckResourceAttr(rName, "referer_config_empty", "true"),
					resource.TestCheckResourceAttr(rName, "referer_white_list", "true"),
					resource.TestCheckResourceAttr(rName, "referer_auth_list.#", "1"),
				),
			},
			{
				Config: testAccRefererValidation_update(domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "referer_config_empty", "false"),
					resource.TestCheckResourceAttr(rName, "referer_white_list", "false"),
					resource.TestCheckResourceAttr(rName, "referer_auth_list.#", "2"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateIdFunc: testAccRefererValidationImportState(rName),
			},
		},
	})
}

func testAccRefererValidation_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_live_domain" "test" {
  name = "%s"
  type = "pull"
}

resource "huaweicloud_live_referer_validation" "test" {
  domain_name          = huaweicloud_live_domain.test.name
  referer_config_empty = "true"
  referer_white_list   = "true"
  referer_auth_list    = ["www.test.com"]
}
`, name)
}

func testAccRefererValidation_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_live_domain" "test" {
  name = "%s"
  type = "pull"
}

resource "huaweicloud_live_referer_validation" "test" {
  domain_name          = huaweicloud_live_domain.test.name
  referer_config_empty = "false"
  referer_white_list   = "false"
  referer_auth_list    = ["www.test.com","www.*com"]
}
`, name)
}

func testAccRefererValidationImportState(rName string) resource.ImportStateIdFunc {
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
