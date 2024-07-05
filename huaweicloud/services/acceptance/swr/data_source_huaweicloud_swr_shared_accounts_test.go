package swr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSharedAccounts_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_shared_accounts.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSWRDomian(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSharedAccounts_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "shared_accounts.#"),
					resource.TestCheckResourceAttrSet(dataSource, "shared_accounts.0.organization"),
					resource.TestCheckResourceAttrSet(dataSource, "shared_accounts.0.repository"),
					resource.TestCheckResourceAttrSet(dataSource, "shared_accounts.0.shared_account"),
					resource.TestCheckResourceAttrSet(dataSource, "shared_accounts.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "shared_accounts.0.permit"),
					resource.TestCheckResourceAttrSet(dataSource, "shared_accounts.0.deadline"),
					resource.TestCheckResourceAttrSet(dataSource, "shared_accounts.0.creator_id"),
					resource.TestCheckResourceAttrSet(dataSource, "shared_accounts.0.created_by"),
					resource.TestMatchResourceAttr(dataSource, "shared_accounts.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(dataSource, "shared_accounts.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccDataSourceSharedAccounts_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_swr_shared_accounts" "test" {
  depends_on = [huaweicloud_swr_repository_sharing.test]
  
  organization = huaweicloud_swr_organization.test.name
  repository   = huaweicloud_swr_repository.test.name
}
`, testAccSWRRepositorySharing_basic(rName))
}
