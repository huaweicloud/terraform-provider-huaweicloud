package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGaussDbDrLogCache_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDbDrLogCache_basic(name),
			},
			{
				Config: testAccGaussDbDrLogCache_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("huaweicloud_gaussdb_dr_relationship.test", "status", "dr_log_keep"),
				),
			},
		},
	})
}

func testAccGaussDbDrLogCache_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_gaussdb_dr_log_cache" "test" {
  depends_on = [huaweicloud_gaussdb_dr_relationship.test]

  instance_id     = huaweicloud_gaussdb_instance.test[1].id
  disaster_type   = "stream"
  xlog_keep_ratio = 50
}
`, testGaussDbDrRelationship_basic(name))
}
