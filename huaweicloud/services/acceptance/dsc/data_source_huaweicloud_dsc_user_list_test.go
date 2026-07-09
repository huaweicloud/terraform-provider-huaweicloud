package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscUserList_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dsc_user_list.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscUserList_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "user_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "user_list.0.user_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "user_list.0.user_name"),
				),
			},
		},
	})
}

const testAccDataSourceDscUserList_basic = `
data "huaweicloud_dsc_user_list" "test" {}
`
