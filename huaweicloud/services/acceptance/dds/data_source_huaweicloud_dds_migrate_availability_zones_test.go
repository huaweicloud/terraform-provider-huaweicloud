package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAZs_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dds_migrate_availability_zones.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAZs_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "names.#"),
				),
			},
		},
	})
}

func testAccDataSourceAZs_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dds_migrate_availability_zones" "test" {
  depends_on = [huaweicloud_dds_instance.instance]

  instance_id = huaweicloud_dds_instance.instance.id
}
`, testAccDDSInstanceReplicaSetBasic(rName))
}
