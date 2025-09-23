package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func TestAccDataSourceSignatures_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_apig_signatures.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byId   = "data.huaweicloud_apig_signatures.filter_by_id"
		dcById = acceptance.InitDataSourceCheck(byId)

		byName   = "data.huaweicloud_apig_signatures.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byType   = "data.huaweicloud_apig_signatures.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byAlgorithm   = "data.huaweicloud_apig_signatures.filter_by_algorithm"
		dcByAlgorithm = acceptance.InitDataSourceCheck(byAlgorithm)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSignatures_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "signatures.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "signatures.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "signatures.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "signatures.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "signatures.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "signatures.0.secret"),
					resource.TestCheckResourceAttrSet(dataSource, "signatures.0.algorithm"),
					resource.TestCheckResourceAttrSet(dataSource, "signatures.0.bind_num"),
					resource.TestCheckResourceAttrSet(dataSource, "signatures.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "signatures.0.updated_at"),

					dcById.CheckResourceExists(),
					resource.TestCheckOutput("signature_id_filter_is_useful", "true"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					dcByAlgorithm.CheckResourceExists(),
					resource.TestCheckOutput("algorithm_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceSignatures_base() string {
	name := acceptance.RandomAccResourceName()
	signKey := acctest.RandString(16)
	signSecret := utils.Reverse(signKey)

	return fmt.Sprintf(`
data "huaweicloud_apig_instances" "test" {
  instance_id = "%[1]s"
}

locals {
  instance_id = data.huaweicloud_apig_instances.test.instances[0].id
}

resource "huaweicloud_apig_signature" "with_key" {
  instance_id = local.instance_id
  name        = "%[2]s_with_key"
  type        = "aes"
  algorithm   = "aes-128-cfb"
  key         = "%[3]s"
  secret      = "%[4]s"
}

resource "huaweicloud_apig_signature" "without_key" {
  instance_id = local.instance_id
  name        = "%[2]s_without_key"
  type        = "aes"
  algorithm   = "aes-256-cfb"
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name, signKey, signSecret)
}

func testAccDataSourceSignatures_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_apig_signatures" "test" {
  depends_on = [
    huaweicloud_apig_signature.with_key
  ]

  instance_id = local.instance_id
}

# Filter by ID
locals {
  signature_id = huaweicloud_apig_signature.with_key.id
}

data "huaweicloud_apig_signatures" "filter_by_id" {
  instance_id  = local.instance_id
  signature_id = local.signature_id
}

locals {
  id_filter_result = [
    for v in data.huaweicloud_apig_signatures.filter_by_id.signatures[*].id : v == local.signature_id
  ]
}

output "signature_id_filter_is_useful" {
  value = length(local.id_filter_result) > 0 && alltrue(local.id_filter_result)
}

# Filter by name
locals {
  name = huaweicloud_apig_signature.with_key.name
}

data "huaweicloud_apig_signatures" "filter_by_name" {
  depends_on = [
    huaweicloud_apig_signature.with_key
  ]

  instance_id = local.instance_id
  name        = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_apig_signatures.filter_by_name.signatures[*].name : v == local.name
  ]
}

output "name_filter_is_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by type
locals {
  type = huaweicloud_apig_signature.with_key.type
}

data "huaweicloud_apig_signatures" "filter_by_type" {
  depends_on = [
    huaweicloud_apig_signature.with_key
  ]

  instance_id = local.instance_id
  type        = local.type
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_apig_signatures.filter_by_type.signatures[*].type : v == local.type
  ]
}

output "type_filter_is_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}

# Filter by algorithm
locals {
  algorithm = huaweicloud_apig_signature.with_key.algorithm
}

data "huaweicloud_apig_signatures" "filter_by_algorithm" {
  depends_on = [
    huaweicloud_apig_signature.with_key
  ]

  instance_id = local.instance_id
  algorithm   = local.algorithm
}

locals {
  algorithm_filter_result = [
    for v in data.huaweicloud_apig_signatures.filter_by_algorithm.signatures[*].algorithm : v == local.algorithm
  ]
}

output "algorithm_filter_is_useful" {
  value = length(local.algorithm_filter_result) > 0 && alltrue(local.algorithm_filter_result)
}
`, testAccDataSourceSignatures_base())
}
