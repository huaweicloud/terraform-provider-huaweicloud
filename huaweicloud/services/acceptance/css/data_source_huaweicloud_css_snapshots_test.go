package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCssSnapshots_basic(t *testing.T) {
	dataSource := "data.huaweicloud_css_snapshots.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCssSnapshots_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "backups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.indices"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.end_time"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCssSnapshots_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_css_snapshots" "test" {
  depends_on = [huaweicloud_css_snapshot.snapshot]

  cluster_id = huaweicloud_css_cluster.test.id
}
`, testAccCssSnapshot_basic(name))
}
