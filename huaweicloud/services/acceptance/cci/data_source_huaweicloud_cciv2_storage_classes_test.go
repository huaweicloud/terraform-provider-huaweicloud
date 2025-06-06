package cci

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceV2StorageClasses_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_cciv2_storage_classes.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceV2StorageClasses_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "storage_classes.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_classes.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_classes.0.creation_timestamp"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_classes.0.parameters.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_classes.0.provisioner"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_classes.0.resource_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_classes.0.uid"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_classes.0.reclaim_policy"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_classes.0.volume_binding_mode"),
				),
			},
		},
	})
}

const testAccDataSourceV2StorageClasses_basic = `data "huaweicloud_cciv2_storage_classes" "test" {}`
