package dds

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatabaseUpgrade_basic(t *testing.T) {
	// Lack of testing conditions
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDSInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testDatabaseUpgrade_basic(),
				ExpectError: regexp.MustCompile("error upgrading the database patch"),
			},
		},
	})
}

func testDatabaseUpgrade_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dds_database_upgrade" "test" {
  instance_id  = "%[1]s"
  upgrade_mode = "minimized_interrupt_time"
  is_delayed   = false
}
`, acceptance.HW_DDS_INSTANCE_ID)
}
