package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDdsInstance_basic(t *testing.T) {
	rName := "data.huaweicloud_dds_instances.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDdsInstance_basic(name, 8800),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "instances.0.name", name),
					resource.TestCheckResourceAttr(rName, "instances.0.mode", "Sharding"),
				),
			},
		},
	})
}

func testAccDatasourceDdsInstance_basic(name string, port int) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dds_instances" "test" {
  depends_on = [huaweicloud_dds_instance.instance]

  name = huaweicloud_dds_instance.instance.name
}
`, testAccDDSInstanceV3Config_basic(name, port))
}
