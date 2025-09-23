package evs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEvsV3Quotas_basic(t *testing.T) {
	dataSource := "data.huaweicloud_evsv3_quotas.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEvsV3Quotas_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "id"),
					resource.TestCheckResourceAttrSet(dataSource, "region"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.backup_gigabytes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.backups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.gigabytes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.snapshots.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.volumes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.gigabytes_sata.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.snapshots_sata.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.volumes_sata.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.gigabytes_sas.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.snapshots_sas.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.volumes_sas.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.gigabytes_ssd.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.snapshots_ssd.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.volumes_ssd.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.gigabytes_gpssd.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.snapshots_gpssd.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.volumes_gpssd.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.per_volume_gigabytes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.backup_gigabytes.0.limit"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.backup_gigabytes.0.in_use"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.backups.0.limit"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.backups.0.in_use"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.gigabytes.0.limit"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.gigabytes.0.in_use"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.snapshots.0.limit"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.snapshots.0.in_use"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.volumes.0.limit"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.volumes.0.in_use"),
				),
			},
		},
	})
}

const testDataSourceEvsV3Quotas_basic = `
data "huaweicloud_evsv3_quotas" "test" {
  usage = true
}
`
