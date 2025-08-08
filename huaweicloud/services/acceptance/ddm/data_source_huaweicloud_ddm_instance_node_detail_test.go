package ddm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDdmInstanceNodeDetail_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ddm_instance_node_detail.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDdmInstanceNodeDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "status"),
					resource.TestCheckResourceAttrSet(dataSource, "name"),
					resource.TestCheckResourceAttrSet(dataSource, "private_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "floating_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "server_id"),
					resource.TestCheckResourceAttrSet(dataSource, "subnet_name"),
					resource.TestCheckResourceAttrSet(dataSource, "datavolume_id"),
					resource.TestCheckResourceAttrSet(dataSource, "res_subnet_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "systemvolume_id"),
				),
			},
		},
	})
}

func testAccDatasourceDdmInstanceNodeDetail_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_ddm_instance_node_detail" "test" {
  instance_id = "%s"
  node_id     = "%s"
}
`, acceptance.HW_DDM_INSTANCE_ID, acceptance.HW_DDM_NODE_ID)
}
