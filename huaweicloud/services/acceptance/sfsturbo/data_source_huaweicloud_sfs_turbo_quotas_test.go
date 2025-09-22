package sfsturbo

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSfsTurboQuotas_basic(t *testing.T) {
	resourceName := "data.huaweicloud_sfs_turbo_quotas.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSfsTurboQuotas_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestCheckResourceAttrSet(resourceName, "quotas.#"),
					resource.TestCheckResourceAttrSet(resourceName, "quotas.0.resources.#"),
					resource.TestCheckResourceAttrSet(resourceName, "quotas.0.resources.0.type"),
					resource.TestCheckResourceAttrSet(resourceName, "quotas.0.resources.0.max"),
					resource.TestCheckResourceAttrSet(resourceName, "quotas.0.resources.0.min"),
					resource.TestCheckResourceAttrSet(resourceName, "quotas.0.resources.0.quota"),
					resource.TestCheckResourceAttrSet(resourceName, "quotas.0.resources.0.used"),
				),
			},
		},
	})
}

const testAccDataSourceSfsTurboQuotas_basic = `data "huaweicloud_sfs_turbo_quotas" "test" {}`
