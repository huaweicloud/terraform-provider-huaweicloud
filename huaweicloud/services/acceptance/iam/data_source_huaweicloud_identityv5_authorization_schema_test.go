package iam_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataV5AuthorizationSchema_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_identityv5_authorization_schema.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataV5AuthorizationSchema_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "actions.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "actions.0.name"),
					resource.TestCheckResourceAttrSet(all, "actions.0.access_level"),
					resource.TestMatchResourceAttr(all, "actions.0.description.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestMatchResourceAttr(all, "conditions.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "conditions.0.key"),
					resource.TestCheckResourceAttrSet(all, "conditions.0.value_type"),
					resource.TestCheckResourceAttrSet(all, "conditions.0.multi_valued"),
					resource.TestMatchResourceAttr(all, "conditions.0.description.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestMatchResourceAttr(all, "operations.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "operations.0.operation_action"),
					resource.TestCheckResourceAttrSet(all, "operations.0.operation_id"),
					resource.TestMatchResourceAttr(all, "resources.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "resources.0.urn_template"),
					resource.TestCheckResourceAttrSet(all, "resources.0.type_name"),
					resource.TestCheckResourceAttrSet(all, "version"),
				),
			},
		},
	})
}

const testAccDataV5AuthorizationSchema_basic = `
data "huaweicloud_identityv5_authorization_schema" "test" {
  service_code = "iam"
}
`
