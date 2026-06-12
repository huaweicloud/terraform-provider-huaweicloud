package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGaussDbDrDrill_basic(t *testing.T) {
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
				Config: testAccGaussDbDrDrill_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.huaweicloud_gaussdb_dr_relationships.test",
						"relations.0.status", "simulation"),
				),
			},
		},
	})
}

func testAccGaussDbDrDrill_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_gaussdb_dr_drill" "test" {
  depends_on = [huaweicloud_gaussdb_dr_relationship.test]

  instance_id     = huaweicloud_gaussdb_instance.test[0].id
  disaster_type   = "stream"
  xlog_keep_ratio = 50
}

data "huaweicloud_gaussdb_dr_relationships" "test" {
  depends_on = [huaweicloud_gaussdb_dr_drill.test]

  instance_id = huaweicloud_gaussdb_instance.test[0].id
}
`, testGaussDbDrRelationship_basic(name))
}
