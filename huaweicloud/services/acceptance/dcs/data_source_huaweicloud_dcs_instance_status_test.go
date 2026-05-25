package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDcsInstanceStatus_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dcs_instance_status.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsInstanceStatusConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "running_count"),
					resource.TestCheckResourceAttrSet(dataSourceName, "creating_count"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis.0.running_count"),
					resource.TestCheckResourceAttrSet(dataSourceName, "memcached.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "memcached.0.running_count"),
				),
			},
		},
	})
}
func testAccDcsInstanceStatusConfig_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dcs_instance_status" "test" {
  depends_on = [huaweicloud_dcs_instance.instance_1]

  include_failure = "true"
}
`, testAccDcsInstance_base(name))
}
