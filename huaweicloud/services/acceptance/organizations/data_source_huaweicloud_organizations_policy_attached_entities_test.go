package organizations

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOrganizationsPolicyAttachedEntities_basic(t *testing.T) {
	dataSource := "data.huaweicloud_organizations_policy_attached_entities.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceOrganizationsPolicyAttachedEntities_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "attached_entities.#"),
					resource.TestCheckResourceAttrSet(dataSource, "attached_entities.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "attached_entities.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "attached_entities.0.name"),
				),
			},
		},
	})
}

func testDataSourceOrganizationsPolicyAttachedEntities_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_organizations_organization" "test" {}

resource "huaweicloud_organizations_organizational_unit" "test" {
  name      = "%[1]s"
  parent_id = data.huaweicloud_organizations_organization.test.root_id

  tags = {
    "key1" = "value1"
    "key2" = "value2"
  }
}

resource "huaweicloud_organizations_policy" "test" {
  name        = "%[1]s"
  type        = "service_control_policy"
  description = "test service control policy description"
  content     = jsonencode(
{
	"Version":"5.0",
	"Statement":[
		{
			"Effect":"Deny",
			"Action":[]
		}
	]
}
)

  tags = {
    "key1" = "value1"
    "key2" = "value2"
  }
}

resource "huaweicloud_organizations_policy_attach" "test" {
  policy_id = huaweicloud_organizations_policy.test.id
  entity_id = huaweicloud_organizations_organizational_unit.test.id
}
`, name)
}

func testDataSourceDataSourceOrganizationsPolicyAttachedEntities_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_organizations_policy_attached_entities" "test" {
  depends_on = [huaweicloud_organizations_policy_attach.test]

  policy_id = huaweicloud_organizations_policy.test.id
}
`, testDataSourceOrganizationsPolicyAttachedEntities_base(name))
}
