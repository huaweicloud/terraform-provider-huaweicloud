package accessanalyzer

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

func getAnalyzerResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("accessanalyzer", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Access Analyzer client: %s", err)
	}
	getAnalyzerHttpUrl := "v5/analyzers/{analyzer_id}"
	getAnalyzerPath := client.Endpoint + getAnalyzerHttpUrl
	getAnalyzerPath = strings.ReplaceAll(getAnalyzerPath, "{analyzer_id}", state.Primary.ID)
	getAnalyzerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getAnalyzerResp, err := client.Request("GET", getAnalyzerPath, &getAnalyzerOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving access annlyzer: %s", err)
	}
	return utils.FlattenResponse(getAnalyzerResp)
}

func TestAccAnalyzer_basic(t *testing.T) {
	var object interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_access_analyzer.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&object,
		getAnalyzerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAnalyzer_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "account"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				Config: testAccAnalyzer_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "account"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key_update", "value_update"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAnalyzer_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_access_analyzer" "test" {
  name = "%s"
  type = "account"
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName)
}

func testAccAnalyzer_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_access_analyzer" "test" {
  name = "%s"
  type = "account"
  tags = {
    foo        = "bar_update"
    key_update = "value_update"
  }
}
`, rName)
}

func TestAccAnalyzer_unused(t *testing.T) {
	var object interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_access_analyzer.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&object,
		getAnalyzerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAnalyzer_unused(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "account_unused_access"),
					resource.TestCheckResourceAttr(resourceName,
						"configuration.0.unused_access.0.unused_access_age", "30"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_analyzed_resource", "last_resource_analyzed_at"},
			},
		},
	})
}

func testAccAnalyzer_unused(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_access_analyzer" "test" {
  name = "%s"
  type = "account_unused_access"

  configuration {
    unused_access {
      unused_access_age = 30
    }
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName)
}
