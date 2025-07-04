package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGroupNameCheck_basic(t *testing.T) {
	var (
		name      = acceptance.RandomAccResourceName()
		checkName = acceptance.RandomAccResourceName()
	)

	// Avoid CheckDestroy because this resource is a one-time resource.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			// Check APIG instance does not exist.
			{
				Config:      testAccGroupNameCheck_basic_step1(name),
				ExpectError: regexp.MustCompile("The instance does not exist"),
			},
			// Check the API group name already exists.
			{
				Config:      testAccGroupNameCheck_basic_step2(name),
				ExpectError: regexp.MustCompile("The API group name already exists"),
			},
			// Check the API group name does not exist.
			{
				Config: testAccGroupNameCheck_basic_step3(name, checkName),
			},
		},
	})
}

func testAccGroupNameCheck_basic_step1(name string) string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
resource "huaweicloud_apig_group_name_check" "test" {
  instance_id = "%[1]s"
  group_name  = "%[2]s"
}
`, randomId, name)
}

func testAccGroupNameCheck_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_group" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccGroupNameCheck_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_group_name_check" "test" {
  instance_id = "%[2]s"
  group_name  = huaweicloud_apig_group.test.name

  depends_on = [huaweicloud_apig_group.test]
}
`, testAccGroupNameCheck_base(name), acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccGroupNameCheck_basic_step3(name, checkName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_group_name_check" "test" {
  instance_id = "%[2]s"
  group_name  = "%[3]s"
}
`, testAccGroupNameCheck_base(name), acceptance.HW_APIG_DEDICATED_INSTANCE_ID, checkName)
}
