package fgs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/fgs"
)

func getApplicationFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("fgs", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating FunctionGraph client: %s", err)
	}
	return fgs.GetApplicationById(client, state.Primary.ID)
}

func TestAccApplication_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_fgs_application.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getApplicationFunc)

		name        = acceptance.RandomAccResourceName()
		randUUID, _ = uuid.GenerateUUID()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please read the instructions carefully before use to ensure sufficient permissions.
			acceptance.TestAccPreCheckFgsAgency(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testAccApplication_invalidTemplateId(randUUID),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`initAppModel failed.*templateId\[%s\]`, randUUID)),
			},
			{
				Config: testAccApplication_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestMatchResourceAttr(resourceName, "stack_resources.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckResourceAttr(resourceName, "repository.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"template_id",
					"agency_name",
					"params",
				},
			},
		},
	})
}

func testAccApplication_invalidTemplateId(randUUID string) string {
	return fmt.Sprintf(`
resource "huaweicloud_fgs_application" "invalid_template_id" {
  name        = "tf_test_invalid_template_id"
  agency_name = "%[1]s"
  template_id = "%[2]s"
  description = "The template ID is invalid"
}
`, acceptance.HW_FGS_AGENCY_NAME, randUUID)
}

func testAccApplication_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_fgs_application_templates" "test" {}

resource "huaweicloud_fgs_application" "test" {
  name        = "%[1]s"
  agency_name = "%[2]s"
  template_id = data.huaweicloud_fgs_application_templates.test.templates[0].id
  description = "Created by terraform script"
}
`, name, acceptance.HW_FGS_AGENCY_NAME)
}
