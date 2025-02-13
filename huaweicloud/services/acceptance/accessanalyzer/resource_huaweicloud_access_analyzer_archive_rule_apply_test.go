package accessanalyzer

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccArchiveRuleApply_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccArchiveRuleApply_basic(rName),
			},
		},
	})
}

func testAccArchiveRuleApply_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_access_analyzer_archive_rule_apply" "test" {
  analyzer_id     = huaweicloud_access_analyzer.test.id
  archive_rule_id = huaweicloud_access_analyzer_archive_rule.test.id
}
`, testAccArchiveRule_basic(rName))
}
