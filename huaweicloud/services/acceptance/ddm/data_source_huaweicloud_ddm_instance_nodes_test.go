package ddm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDdmInstanceNodes_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	name = strings.ReplaceAll(name, "_", "-")
	rName := "data.huaweicloud_ddm_instance_nodes.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDdmInstanceNodes_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "nodes.#", "2"),
					resource.TestCheckResourceAttr(rName, "nodes.0.status", "normal"),
				),
			},
		},
	})
}

func testAccDatasourceDdmInstanceNodes_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_ddm_instance_nodes" "test" {
  instance_id = huaweicloud_ddm_instance.test.id
}
`, testDdmInstance_basic(name))
}
