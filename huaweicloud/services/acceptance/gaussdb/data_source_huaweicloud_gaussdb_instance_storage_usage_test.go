package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbInstanceStorageUsage_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_instance_storage_usage.test"
	dc := acceptance.InitDataSourceCheck(dataSource)
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbInstanceStorageUsage_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "used"),
					resource.TestCheckResourceAttrSet(dataSource, "total"),
				),
			},
		},
	})
}

func testDataSourceGaussdbInstanceStorageUsage_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_instance_storage_usage" "test" {
  instance_id = huaweicloud_gaussdb_instance.test.id
}
`, testDataSourceGaussdbInstanceMetrics_base(name))
}
