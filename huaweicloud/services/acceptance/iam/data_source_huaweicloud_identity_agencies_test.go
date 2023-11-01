package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityAgenciesDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identity_agencies.all"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityAgenciesDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "agencies.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "agencies.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "agencies.0.name"),
				),
			},
		},
	})
}

const testAccIdentityAgenciesDataSourceBasic string = `
data "huaweicloud_identity_agencies" "all" {
}
`

func TestAccIdentityAgenciesDataSource_byName(t *testing.T) {
	dataSourceName := "data.huaweicloud_identity_agencies.query_by_name"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityAgenciesDataSource_byName(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "agencies.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "agencies.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "agencies.0.name"),
				),
			},
		},
	})
}

func testAccIdentityAgenciesDataSource_byName(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_agency" "test" {
  name                  = "%s"
  description           = "This is a test agency"
  delegated_domain_name = "%s"
  duration              = "ONEDAY"

  project_role {
    project = "%s"
    roles   = ["CCE Administrator"]
  }
}

data "huaweicloud_identity_agencies" "query_by_name" {
  depends_on = [
    huaweicloud_identity_agency.test
  ]
  name = huaweicloud_identity_agency.test.name
}
`, rName, acceptance.HW_DOMAIN_NAME, acceptance.HW_REGION_NAME)
}

func TestAccIdentityAgenciesDataSource_byTrustDomainId(t *testing.T) {
	dataSourceName := "data.huaweicloud_identity_agencies.query_by_trust_domain_id"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	randUUID, _ := uuid.GenerateUUID()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
			acceptance.TestAccPrecheckDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityAgenciesDataSource_byTrustDomainId(randUUID),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "agencies.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "agencies.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "agencies.0.name"),
				),
			},
		},
	})
}

func testAccIdentityAgenciesDataSource_byTrustDomainId(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_agency" "test" {
  name                  = "%s"
  description           = "This is a test agency"
  delegated_domain_name = "%s"
  duration              = "ONEDAY"

  project_role {
    project = "%s"
    roles   = ["CCE Administrator"]
  }
}

data "huaweicloud_identity_agencies" "query_by_trust_domain_id" {
  depends_on = [
    huaweicloud_identity_agency.test
  ]
  trust_domain_id = "%s"
}
`, rName, acceptance.HW_DOMAIN_NAME, acceptance.HW_REGION_NAME, acceptance.HW_DOMAIN_ID)
}
