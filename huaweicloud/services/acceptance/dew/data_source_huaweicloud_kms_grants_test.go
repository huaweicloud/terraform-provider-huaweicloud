package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceKmsGrants_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_kms_grants.test"
		rName      = acceptance.RandomAccResourceName()
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceKmsGrants_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "grants.#"),
					resource.TestCheckResourceAttrSet(dataSource, "grants.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "grants.0.key_id"),
					resource.TestCheckResourceAttrSet(dataSource, "grants.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "grants.0.creator"),
					resource.TestCheckResourceAttrSet(dataSource, "grants.0.grantee_principal"),
					resource.TestCheckResourceAttrSet(dataSource, "grants.0.operations.#"),
					resource.TestCheckResourceAttrSet(dataSource, "grants.0.created_at"),
				),
			},
		},
	})
}

func testDataSourceKmsGrants_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_kms_grants" "test" {
  depends_on = [
    huaweicloud_kms_grant.test
  ]

  key_id = huaweicloud_kms_key.test.id
}
`, testKmsGrant_basic(name))
}
