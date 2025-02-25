package rabbitmq

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDmsRabbitmqVhosts_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dms_rabbitmq_vhosts.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDmsRabbitmqVhosts_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "vhosts.#"),
					resource.TestCheckResourceAttrSet(dataSource, "vhosts.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "vhosts.0.tracing"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDmsRabbitmqVhosts_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_rabbitmq_vhosts" "test" {
  depends_on = [huaweicloud_dms_rabbitmq_vhost.test]

  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
}
`, testRabbitmqVhost_basic(name))
}
