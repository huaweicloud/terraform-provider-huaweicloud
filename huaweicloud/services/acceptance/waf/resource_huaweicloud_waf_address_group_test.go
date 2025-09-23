package waf

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

func getAddressGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v1/{project_id}/waf/ip-group/{id}"
		product = "waf"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{id}", state.Primary.ID)
	if enterpriseProjectID := state.Primary.Attributes["enterprise_project_id"]; enterpriseProjectID != "" {
		requestPath += fmt.Sprintf("?enterprise_project_id=%s", enterpriseProjectID)
	}

	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving WAF address group: %s", err)
	}
	return utils.FlattenResponse(resp)
}

// Before running the test case, please ensure that there is at least one WAF instance in the current region.
func TestAccAddressGroup_basic(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_waf_address_group.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAddressGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAddressGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "description", "example_description"),
					resource.TestCheckResourceAttr(rName, "ip_addresses.0", "192.168.1.0/24"),
					resource.TestCheckResourceAttrSet(rName, "rules.#"),
				),
			},
			{
				Config: testAddressGroup_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "description", "example_description_update"),
					resource.TestCheckResourceAttr(rName, "ip_addresses.0", "192.168.1.0"),
					resource.TestCheckResourceAttr(rName, "ip_addresses.1", "192.168.2.0/12"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testWAFResourceImportState(rName),
			},
		},
	})
}

func testAddressGroup_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_address_group" "test" {
  name                  = "%[1]s"
  description           = "example_description"
  ip_addresses          = ["192.168.1.0/24"]
  enterprise_project_id = "%[2]s"
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAddressGroup_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_address_group" "test" {
  name                  = "%[1]s_update"
  description           = "example_description_update"
  ip_addresses          = ["192.168.1.0", "192.168.2.0/12"]
  enterprise_project_id = "%[2]s"
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

// testWAFResourceImportState use to return an id with format <id> or <id>/<enterprise_project_id>
func testWAFResourceImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		epsID := rs.Primary.Attributes["enterprise_project_id"]
		if epsID == "" {
			return rs.Primary.ID, nil
		}
		return fmt.Sprintf("%s/%s", rs.Primary.ID, epsID), nil
	}
}
