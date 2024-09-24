package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceClusterParameterGroups_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dws_cluster_parameter_groups.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testDataSourceClusterParameterGroups_notExist(),
				ExpectError: regexp.MustCompile("Cluster does not exist or has been deleted"),
			},
			{
				Config: testDataSourceClusterParameterGroups_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "groups.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.status"),
				),
			},
		},
	})
}

func testDataSourceClusterParameterGroups_notExist() string {
	randUUID, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dws_cluster_parameter_groups" "test" {
  cluster_id = "%s"
}
`, randUUID)
}

func testDataSourceClusterParameterGroups_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dws_cluster_parameter_groups" "test" {
  cluster_id = "%s"
}
`, acceptance.HW_DWS_CLUSTER_ID)
}
