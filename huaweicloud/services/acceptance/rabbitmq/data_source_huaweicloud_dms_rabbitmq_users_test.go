package rabbitmq

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDmsRabbitmqUsers_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dms_rabbitmq_users.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDmsRabbitmqUsers_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "users.#"),
					resource.TestCheckResourceAttrSet(dataSource, "users.0.access_key"),
					resource.TestCheckResourceAttrSet(dataSource, "users.0.vhosts.#"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDmsRabbitmqUsers_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_rabbitmq_users" "test" {
  depends_on = [huaweicloud_dms_rabbitmq_user.test]

  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
}
`, testRabbitmqUser_basic(name))
}
