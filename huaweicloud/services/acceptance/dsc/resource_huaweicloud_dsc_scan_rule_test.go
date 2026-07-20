package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dsc"
)

func getScanRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dsc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DSC client: %s", err)
	}

	return dsc.GetScanRuleById(client, state.Primary.ID)
}

// Before running the test, please ensure that you have created a DSC instance.
func TestAccScanRule_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_dsc_scan_rule.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getScanRuleResourceFunc)

		rNameImportWithTemplate = "huaweicloud_dsc_scan_rule.import_with_template"
		rcImportWithTemplate    = acceptance.InitResourceCheck(rNameImportWithTemplate, &obj, getScanRuleResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDscSecurityLevelId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rc.CheckResourceDestroy(),
			rcImportWithTemplate.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccScanRule_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "rule_name", name),
					resource.TestCheckResourceAttr(rName, "rule_type", "REGEX"),
					resource.TestCheckResourceAttr(rName, "category", "BUILT_SELF"),
					resource.TestCheckResourceAttr(rName, "logic_operator", "AND"),
					resource.TestCheckResourceAttr(rName, "match_rate", "1"),
					resource.TestCheckResourceAttr(rName, "min_match", "1"),
					resource.TestCheckResourceAttr(rName, "content.#", "2"),
					resource.TestCheckResourceAttr(rName, "content.0.effective_mode", "NOT_IN"),
					resource.TestCheckResourceAttr(rName, "content.0.location", "NAME"),
					resource.TestCheckResourceAttr(rName, "content.0.rule_content", "bphone"),
					resource.TestCheckResourceAttr(rName, "content.1.effective_mode", "IN"),
					resource.TestCheckResourceAttr(rName, "content.1.location", "REMARK"),
					resource.TestCheckResourceAttr(rName, "content.1.rule_content", "telephone number"),
					resource.TestCheckResourceAttr(rName, "rule_desc", "Created_by_acceptance_test"),
					resource.TestCheckResourceAttr(rName, "templates.#", "2"),
					resource.TestCheckResourceAttrPair(rName, "templates.0.template_id",
						"huaweicloud_dsc_scan_template_classification.test.0", "template_id"),
					resource.TestCheckResourceAttrPair(rName, "templates.0.classification_id",
						"huaweicloud_dsc_scan_template_classification.test.0", "id"),
					resource.TestCheckResourceAttr(rName, "templates.0.security_level_id", acceptance.HW_DSC_SECURITY_LEVEL_ID),
					resource.TestCheckResourceAttr(rName, "templates.0.is_used", "true"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.id"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.template_name"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.classification_name"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.security_level_color"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.security_level_name"),
					resource.TestCheckResourceAttrPair(rName, "templates.1.template_id",
						"huaweicloud_dsc_scan_template_classification.test.1", "template_id"),
					resource.TestCheckResourceAttrPair(rName, "templates.1.classification_id",
						"huaweicloud_dsc_scan_template_classification.test.1", "id"),
					resource.TestCheckResourceAttr(rName, "templates.1.security_level_id", acceptance.HW_DSC_SECURITY_LEVEL_ID),
				),
			},
			{
				Config: testAccScanRule_step2(name, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "rule_name", updateName),
					resource.TestCheckResourceAttr(rName, "rule_type", "KEYWORD"),
					resource.TestCheckResourceAttr(rName, "category", "BUILT_SELF"),
					resource.TestCheckResourceAttr(rName, "logic_operator", "OR"),
					resource.TestCheckResourceAttr(rName, "match_rate", "10"),
					resource.TestCheckResourceAttr(rName, "min_match", "10"),
					resource.TestCheckResourceAttr(rName, "content.#", "1"),
					resource.TestCheckResourceAttr(rName, "content.0.effective_mode", "KEYWORD"),
					resource.TestCheckResourceAttr(rName, "content.0.location", "CONTENT"),
					resource.TestCheckResourceAttr(rName, "content.0.rule_content", "content1,content2,content3"),
					resource.TestCheckResourceAttr(rName, "rule_desc", "Updated_by_acceptance_test"),
					resource.TestCheckResourceAttr(rName, "templates.#", "2"),
					resource.TestCheckResourceAttrPair(rName, "templates.0.template_id",
						"huaweicloud_dsc_scan_template_classification.test.1", "template_id"),
					resource.TestCheckResourceAttrPair(rName, "templates.0.classification_id",
						"huaweicloud_dsc_scan_template_classification.test.1", "id"),
					resource.TestCheckResourceAttr(rName, "templates.0.security_level_id", acceptance.HW_DSC_SECURITY_LEVEL_ID),
					resource.TestCheckResourceAttrPair(rName, "templates.1.template_id",
						"huaweicloud_dsc_scan_template_classification.test.2", "template_id"),
					resource.TestCheckResourceAttrPair(rName, "templates.1.classification_id",
						"huaweicloud_dsc_scan_template_classification.test.2", "id"),
					resource.TestCheckResourceAttr(rName, "templates.1.security_level_id", acceptance.HW_DSC_SECURITY_LEVEL_ID),
					resource.TestCheckResourceAttr(rName, "templates.1.is_used", "false"),
				),
			},
			{
				Config: testAccScanRule_step3(name, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "rule_desc", ""),
					resource.TestCheckResourceAttr(rName, "templates.#", "0"),
					rcImportWithTemplate.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameImportWithTemplate, "templates.#", "1"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// The API may return templates in a different order than the configuration, so importing
			// a resource with multiple templates can cause a plan diff. Cover the import case with
			// templates configured by using a separate resource that has only one template.
			{
				ResourceName:      rNameImportWithTemplate,
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore `templates_origin` only for this test, real terraform import is unaffected and will not
				// cause a configuration drift.
				ImportStateVerifyIgnore: []string{"templates_origin"},
			},
		},
	})
}

func testAccScanRule_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dsc_scan_template" "test" {
  count = 3

  action      = "ADD"
  name        = "%[1]s_${count.index}"
  description = "Created_by_terraform_script"
}

resource "huaweicloud_dsc_scan_template_classification" "test" {
  count = 3

  template_id         = huaweicloud_dsc_scan_template.test[count.index].id
  classification_name = "%[1]s"
}
`, name)
}

func testAccScanRule_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dsc_scan_rule" "test" {
  rule_name      = "%[2]s"
  rule_type      = "REGEX"
  category       = "BUILT_SELF"
  logic_operator = "AND"
  match_rate     = 1
  min_match      = 1
  rule_desc      = "Created_by_acceptance_test"

  content {
    effective_mode = "NOT_IN"
    location       = "NAME"
    rule_content   = "bphone"
  }
  content {
    effective_mode = "IN"
    location       = "REMARK"
    rule_content   = "telephone number"
  }

  dynamic "templates" {
    for_each = slice(huaweicloud_dsc_scan_template_classification.test[*], 0, 2)

    content {
      template_id       = templates.value.template_id
      classification_id = templates.value.id
      security_level_id = "%[3]s"
    }
  }
}
`, testAccScanRule_base(name), name, acceptance.HW_DSC_SECURITY_LEVEL_ID)
}

func testAccScanRule_step2(name, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dsc_scan_rule" "test" {
  rule_name      = "%[2]s"
  rule_type      = "KEYWORD"
  category       = "BUILT_SELF"
  logic_operator = "OR"
  match_rate     = 10
  min_match      = 10
  rule_desc      = "Updated_by_acceptance_test"

  content {
    effective_mode = "KEYWORD"
    location       = "CONTENT"
    rule_content   = "content1,content2,content3"
  }

  dynamic "templates" {
    for_each = slice(huaweicloud_dsc_scan_template_classification.test[*], 1, 3)

    content {
      template_id       = templates.value.template_id
      classification_id = templates.value.id
      security_level_id = "%[3]s"
      is_used           = "false"
    }
  }
}
`, testAccScanRule_base(name), updateName, acceptance.HW_DSC_SECURITY_LEVEL_ID)
}

func testAccScanRule_step3(name, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dsc_scan_rule" "test" {
  rule_name      = "%[2]s"
  rule_type      = "KEYWORD"
  category       = "BUILT_SELF"
  logic_operator = "OR"
  match_rate     = 10
  min_match      = 10

  content {
    effective_mode = "KEYWORD"
    location       = "CONTENT"
    rule_content   = "content1,content2,content3"
  }
}

resource "huaweicloud_dsc_scan_rule" "import_with_template" {
  rule_name      = "%[2]s_import_with_template"
  rule_type      = "KEYWORD"
  category       = "BUILT_SELF"
  logic_operator = "OR"
  match_rate     = 10
  min_match      = 10

  content {
    effective_mode = "KEYWORD"
    location       = "CONTENT"
    rule_content   = "content1,content2,content3"
  }

  templates {
    template_id       = try(huaweicloud_dsc_scan_template_classification.test[0].template_id, null)
    classification_id = try(huaweicloud_dsc_scan_template_classification.test[0].id, null)
    security_level_id = "%[3]s"
  }
}
`, testAccScanRule_base(name), updateName, acceptance.HW_DSC_SECURITY_LEVEL_ID)
}
