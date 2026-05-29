package geminidb

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Due to the lack of testing conditions, only the error situations of API calls were verified.
// Before running test, prepare a GeminiDB instance that do not support upgrade patches.
func TestAccDatabaseUpgrade_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccCheckGeminidbInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// If the database does not support patchs upgrade, will return above error.
				Config:      testDatabaseUpgrade_basic(),
				ExpectError: regexp.MustCompile("error upgrading GeminiDB database patch"),
			},
		},
	})
}

func testDatabaseUpgrade_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_geminidb_database_upgrade" "test" {
  instance_id = "%[1]s"
}
`, acceptance.HW_GEMINIDB_INSATNCE_ID)
}
