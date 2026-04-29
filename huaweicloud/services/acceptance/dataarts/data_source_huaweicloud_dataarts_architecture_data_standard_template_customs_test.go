package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceArchitectureDataStandardTemplateCustoms_basic(t *testing.T) {
	rName := "data.huaweicloud_dataarts_architecture_data_standard_template_customs.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDatasourceArchitectureDataStandardTemplateCustoms_nonExistentWorkspace(),
				ExpectError: regexp.MustCompile("error querying DataArts Architecture data standard template customs"),
			},
			{
				Config: testAccDatasourceArchitectureDataStandardTemplateCustoms_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(rName, "customs.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(rName, "customs.0.id"),
					resource.TestCheckResourceAttrSet(rName, "customs.0.fd_name"),
					resource.TestCheckResourceAttrSet(rName, "customs.0.fd_name_en"),
					resource.TestCheckResourceAttrSet(rName, "customs.0.description"),
					resource.TestCheckResourceAttrSet(rName, "customs.0.actived"),
					resource.TestCheckResourceAttrSet(rName, "customs.0.required"),
					resource.TestCheckResourceAttrSet(rName, "customs.0.searchable"),
					resource.TestCheckResourceAttrSet(rName, "customs.0.optional_values"),
					resource.TestCheckResourceAttrSet(rName, "customs.0.create_time"),
					resource.TestCheckResourceAttrSet(rName, "customs.0.update_time"),
					resource.TestCheckResourceAttrSet(rName, "customs.0.create_by"),
					resource.TestCheckResourceAttrSet(rName, "customs.0.update_by"),
				),
			},
		},
	})
}

func testAccDatasourceArchitectureDataStandardTemplateCustoms_nonExistentWorkspace() string {
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
data "huaweicloud_dataarts_architecture_data_standard_template_customs" "test" {
  workspace_id = "%[1]s"
}
`, randUUID)
}

func testAccDatasourceArchitectureDataStandardTemplateCustoms_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dataarts_architecture_data_standard_template_customs" "test" {
  workspace_id = "%[1]s"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID)
}
