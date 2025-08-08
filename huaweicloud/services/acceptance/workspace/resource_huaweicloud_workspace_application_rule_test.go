package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/workspace"
)

func getApplicationRuleFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("workspace", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace client: %s", err)
	}

	return workspace.GetApplicationRuleById(client, state.Primary.ID)
}

func TestAccResourceApplicationRule_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		withProductRule = "huaweicloud_workspace_application_rule.with_product_rule"
		withPathRule    = "huaweicloud_workspace_application_rule.with_path_rule"
		applicationRule interface{}
		rcWithProduct   = acceptance.InitResourceCheck(withProductRule, &applicationRule, getApplicationRuleFunc)
		rcWithPath      = acceptance.InitResourceCheck(withPathRule, &applicationRule, getApplicationRuleFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcWithProduct.CheckResourceDestroy(),
			rcWithPath.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceApplicationRule_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rcWithProduct.CheckResourceExists(),
					resource.TestCheckResourceAttr(withProductRule, "name", name),
					resource.TestCheckResourceAttr(withProductRule, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(withProductRule, "detail.0.scope", "PRODUCT"),
					resource.TestCheckResourceAttr(withProductRule, "detail.0.product_rule.0.identify_condition", "process"),
					resource.TestCheckResourceAttr(withProductRule, "detail.0.product_rule.0.publisher", "Microsoft Corporation"),
					resource.TestCheckResourceAttr(withProductRule, "detail.0.product_rule.0.product_name", "Microsoft Office"),
					resource.TestCheckResourceAttr(withProductRule, "detail.0.product_rule.0.process_name", "WINWORD.EXE"),
					resource.TestCheckResourceAttr(withProductRule, "detail.0.product_rule.0.support_os", "Windows"),
					resource.TestCheckResourceAttr(withProductRule, "detail.0.product_rule.0.version", "1.0"),
					resource.TestCheckResourceAttr(withProductRule, "detail.0.product_rule.0.product_version", "2019"),
					resource.TestCheckResourceAttr(withProductRule, "detail.0.path_rule.#", "0"),
					rcWithPath.CheckResourceExists(),
					resource.TestCheckResourceAttr(withPathRule, "name", name+"_path"),
					resource.TestCheckResourceAttr(withPathRule, "description", "Created by terraform script for path rule"),
					resource.TestCheckResourceAttr(withPathRule, "detail.0.scope", "PATH"),
					resource.TestCheckResourceAttr(withPathRule, "detail.0.path_rule.0.path",
						"C:\\Program Files\\Microsoft Office\\root\\Office16\\WINWORD.EXE"),
					resource.TestCheckResourceAttr(withPathRule, "detail.0.product_rule.#", "0"),
				),
			},
			{
				Config: testAccResourceApplicationRule_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rcWithProduct.CheckResourceExists(),
					resource.TestCheckResourceAttr(withProductRule, "name", name),
					resource.TestCheckResourceAttr(withProductRule, "description", ""),
					resource.TestCheckResourceAttr(withProductRule, "detail.0.scope", "PRODUCT"),
					resource.TestCheckResourceAttr(withProductRule, "detail.0.product_rule.0.identify_condition", "product"),
					resource.TestCheckResourceAttr(withProductRule, "detail.0.product_rule.0.publisher", "Adobe Inc."),
					resource.TestCheckResourceAttr(withProductRule, "detail.0.product_rule.0.product_name", "Adobe Photoshop"),
					resource.TestCheckResourceAttr(withProductRule, "detail.0.product_rule.0.process_name", "Photoshop.exe"),
					resource.TestCheckResourceAttr(withProductRule, "detail.0.product_rule.0.support_os", "Windows,Mac"),
					resource.TestCheckResourceAttr(withProductRule, "detail.0.product_rule.0.version", "2.0"),
					resource.TestCheckResourceAttr(withProductRule, "detail.0.product_rule.0.product_version", "2022"),
					resource.TestCheckResourceAttr(withProductRule, "detail.0.path_rule.#", "0"),
					rcWithPath.CheckResourceExists(),
					resource.TestCheckResourceAttr(withPathRule, "name", name+"_path"),
					resource.TestCheckResourceAttr(withPathRule, "description", "Updated by terraform script for path rule"),
					resource.TestCheckResourceAttr(withPathRule, "detail.0.scope", "PATH"),
					resource.TestCheckResourceAttr(withPathRule, "detail.0.path_rule.0.path",
						"C:\\Program Files\\Adobe\\Adobe Photoshop 2022\\Photoshop.exe"),
					resource.TestCheckResourceAttr(withPathRule, "detail.0.product_rule.#", "0"),
				),
			},
			{
				ResourceName:      withProductRule,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      withPathRule,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceApplicationRule_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_application_rule" "with_product_rule" {
  name        = "%[1]s"
  description = "Created by terraform script"

  detail {
    scope = "PRODUCT"

    product_rule {
      identify_condition = "process"
      publisher          = "Microsoft Corporation"
      product_name       = "Microsoft Office"
      process_name       = "WINWORD.EXE"
      support_os         = "Windows"
      version            = "1.0"
      product_version    = "2019"
    }
  }
}

resource "huaweicloud_workspace_application_rule" "with_path_rule" {
  name        = "%[1]s_path"
  description = "Created by terraform script for path rule"

  detail {
    scope = "PATH"

    path_rule {
      path = "C:\\Program Files\\Microsoft Office\\root\\Office16\\WINWORD.EXE"
    }
  }
}
`, name)
}

func testAccResourceApplicationRule_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_application_rule" "with_product_rule" {
  name = "%[1]s"

  detail {
    scope = "PRODUCT"

    product_rule {
      identify_condition = "product"
      publisher          = "Adobe Inc."
      product_name       = "Adobe Photoshop"
      process_name       = "Photoshop.exe"
      support_os         = "Windows,Mac"
      version            = "2.0"
      product_version    = "2022"
    }
  }
}

resource "huaweicloud_workspace_application_rule" "with_path_rule" {
  name        = "%[1]s_path"
  description = "Updated by terraform script for path rule"

  detail {
    scope = "PATH"

    path_rule {
      path = "C:\\Program Files\\Adobe\\Adobe Photoshop 2022\\Photoshop.exe"
    }
  }
}
`, name)
}
