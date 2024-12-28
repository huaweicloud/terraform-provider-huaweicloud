package dbss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceInstances_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dbss_instances.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDbssInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceInstances_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.availability_zone"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.charge_model"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.config_num"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.connect_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.cpu"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.database_limit"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.effect"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.expired_at"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.keep_days"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.port_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.ram"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.region"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.remain_days"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.resource_spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.security_group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.specification"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.subnet_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.task"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.vpc_id"),
				),
			},
		},
	})
}

const testDataSourceDataSourceInstances_basic = `data "huaweicloud_dbss_instances" "test" {}`
