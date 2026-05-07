package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDrsNameValidation_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_drs_name_validation.test"
		rName        = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDrsNameValidation_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "trans"),
					resource.TestCheckResourceAttrSet(resourceName, "is_valid"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func testAccDrsNameValidation_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_name_validation" "test" {
  name = "%s"
  type = "trans"
}
`, rName)
}
