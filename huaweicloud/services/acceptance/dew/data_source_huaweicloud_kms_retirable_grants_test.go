package dew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRetirableGrants_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_kms_retirable_grants.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDewFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRetirableGrants_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "grants.#"),
					resource.TestCheckResourceAttrSet(dataSource, "grants.0.key_id"),
					resource.TestCheckResourceAttrSet(dataSource, "grants.0.grant_id"),
					resource.TestCheckResourceAttrSet(dataSource, "grants.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "grants.0.grantee_principal"),
					resource.TestCheckResourceAttrSet(dataSource, "grants.0.grantee_principal_type"),
					resource.TestCheckResourceAttrSet(dataSource, "grants.0.operations.#"),
					resource.TestCheckResourceAttrSet(dataSource, "grants.0.issuing_principal"),
					resource.TestCheckResourceAttrSet(dataSource, "grants.0.creation_date"),
				),
			},
		},
	})
}

const testAccDataSourceRetirableGrants_basic = `
data "huaweicloud_kms_retirable_grants" "test" {
  sequence = "919c82d4-8046-4722-9094-35c3c6524cff"
}
`
