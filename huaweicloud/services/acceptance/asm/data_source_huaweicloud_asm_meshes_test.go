package asm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAsmMeshes_basic(t *testing.T) {
	dataSource := "data.huaweicloud_asm_meshes.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceAsmMeshes_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceAsmMeshes_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_asm_meshes" "test" {
  depends_on = [huaweicloud_asm_mesh.test]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_asm_meshes.test.items) > 0
}
`, testResourceResourceAsmMesh_basic(name))
}
