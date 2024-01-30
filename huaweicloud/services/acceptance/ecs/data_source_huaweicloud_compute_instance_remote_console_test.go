package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEcsInstanceRemoteConsoleDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_compute_instance_remote_console.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEcsInstanceRemoteConsoleDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "protocol"),
					resource.TestCheckResourceAttrSet(dataSourceName, "type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "url"),
				),
			},
		},
	})
}

func testAccEcsInstanceRemoteConsoleDataSource_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_compute_instance_remote_console" "test" {
  instance_id = huaweicloud_compute_instance.test.id
}
`, testAccComputeInstanceDataSource_basic(acceptance.RandomAccResourceNameWithDash()))
}
