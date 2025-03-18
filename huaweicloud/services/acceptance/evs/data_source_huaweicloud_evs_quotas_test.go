package evs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEvsQuotas_basic(t *testing.T) {
	dataSource := "data.huaweicloud_evs_quotas.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEvsQuotas_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.backup_gigabytes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.backups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.gigabytes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.gigabytes_sas.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.gigabytes_ssd.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.snapshots.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.snapshots_sas.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.snapshots_ssd.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.volumes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.volumes_sas.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.volumes_ssd.#"),
				),
			},
		},
	})
}

const testDataSourceEvsQuotas_basic = `
data "huaweicloud_evs_quotas" "test" {
  usage = "True"
}
`
