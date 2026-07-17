package dsc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dsc"
)

func getScanTemplateClassificationResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dsc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DSC client: %s", err)
	}

	return dsc.GetScanTemplateClassificationById(client, state.Primary.Attributes["template_id"], state.Primary.ID)
}

// Before this test, please ensure that the DSC instance has been created.
func TestAccResourceScanTemplateClassification_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_dsc_scan_template_classification.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getScanTemplateClassificationResourceFunc)

		rNameWithChild = "huaweicloud_dsc_scan_template_classification.test2"
		rcWithChild    = acceptance.InitResourceCheck(rNameWithChild, &obj, getScanTemplateClassificationResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDscInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rc.CheckResourceDestroy(),
			rcWithChild.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScanTemplateClassification_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "classification_name", name),
					resource.TestCheckResourceAttrSet(rName, "template_id"),
					resource.TestCheckResourceAttr(rName, "depth", "1"),
					resource.TestMatchResourceAttr(rName, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(rName, "update_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					rcWithChild.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rNameWithChild, "parent_id", rName, "id"),
					resource.TestCheckResourceAttr(rNameWithChild, "depth", "2"),
				),
			},
			{
				Config: testAccResourceScanTemplateClassification_basic_step2(name, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "classification_name", updateName),
					resource.TestCheckResourceAttrSet(rName, "template_id"),
					rcWithChild.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rNameWithChild, "parent_id", rName, "id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccScanTemplateClassificationImportStateIdFunc(rName),
			},
			{
				ResourceName:      rNameWithChild,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccScanTemplateClassificationImportStateIdFunc(rNameWithChild),
			},
		},
	})
}

func testAccScanTemplateClassificationImportStateIdFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", rName)
		}

		templateId := rs.Primary.Attributes["template_id"]
		classificationId := rs.Primary.ID
		if templateId == "" || classificationId == "" {
			return "", fmt.Errorf("missing some attributes, want '<template_id>/<id>', but got '%s/%s'", templateId, classificationId)
		}

		return fmt.Sprintf("%s/%s", templateId, classificationId), nil
	}
}

func testAccResourceScanTemplateClassification_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dsc_scan_template" "test" {
  action      = "ADD"
  name        = "%[1]s"
  description = "Created_by_terraform_script"
}
`, name)
}

func testAccResourceScanTemplateClassification_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dsc_scan_template_classification" "test" {
  template_id         = huaweicloud_dsc_scan_template.test.id
  classification_name = "%[2]s"
}

resource "huaweicloud_dsc_scan_template_classification" "test2" {
  template_id         = huaweicloud_dsc_scan_template.test.id
  classification_name = "%[2]s_child"
  parent_id           = huaweicloud_dsc_scan_template_classification.test.id
}
`, testAccResourceScanTemplateClassification_base(name), name)
}

func testAccResourceScanTemplateClassification_basic_step2(name, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dsc_scan_template_classification" "test" {
  template_id         = huaweicloud_dsc_scan_template.test.id
  classification_name = "%[2]s"
}

resource "huaweicloud_dsc_scan_template_classification" "test2" {
  template_id         = huaweicloud_dsc_scan_template.test.id
  classification_name = "%[2]s_child"
  parent_id           = huaweicloud_dsc_scan_template_classification.test.id
}
`, testAccResourceScanTemplateClassification_base(name), updateName)
}
