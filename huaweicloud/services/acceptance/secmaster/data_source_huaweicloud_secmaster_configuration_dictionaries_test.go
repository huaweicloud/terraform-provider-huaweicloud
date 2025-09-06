package secmaster

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceConfigurationDictionaries_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_configuration_dictionaries.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceConfigurationDictionaries_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "success_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "success_list.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "success_list.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "success_list.0.dict_id"),
					resource.TestCheckResourceAttrSet(dataSource, "success_list.0.dict_key"),
					resource.TestCheckResourceAttrSet(dataSource, "success_list.0.dict_code"),
					resource.TestCheckResourceAttrSet(dataSource, "success_list.0.dict_val"),
					resource.TestCheckResourceAttrSet(dataSource, "success_list.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "success_list.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "success_list.0.publish_time"),
					resource.TestCheckResourceAttrSet(dataSource, "success_list.0.scope"),
					resource.TestCheckResourceAttrSet(dataSource, "success_list.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "success_list.0.language"),
				),
			},
		},
	})
}

const testDataSourceConfigurationDictionaries_basic = `data "huaweicloud_secmaster_configuration_dictionaries" "test" {}`
