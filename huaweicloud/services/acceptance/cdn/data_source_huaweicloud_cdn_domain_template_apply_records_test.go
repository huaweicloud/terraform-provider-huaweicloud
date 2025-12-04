package cdn

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDomainTemplateApplyRecords_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		all = "data.huaweicloud_cdn_domain_template_apply_records.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByTemplateId   = "data.huaweicloud_cdn_domain_template_apply_records.filter_by_template_id"
		dcFilterByTemplateId = acceptance.InitDataSourceCheck(filterByTemplateId)

		filterByTemplateName   = "data.huaweicloud_cdn_domain_template_apply_records.filter_by_template_name"
		dcFilterByTemplateName = acceptance.InitDataSourceCheck(filterByTemplateName)

		filterByOperatorId   = "data.huaweicloud_cdn_domain_template_apply_records.filter_by_operator_id"
		dcFilterByOperatorId = acceptance.InitDataSourceCheck(filterByOperatorId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCdnDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataDomainTemplateApplyRecords_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "records.#", regexp.MustCompile("^[1-9]([0-9]+)?$")),
					// Filter by 'template_id' parameter.
					dcFilterByTemplateId.CheckResourceExists(),
					resource.TestCheckOutput("is_template_id_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(filterByTemplateId, "records.0.operator_id"),
					resource.TestCheckResourceAttrSet(filterByTemplateId, "records.0.status"),
					resource.TestCheckResourceAttrSet(filterByTemplateId, "records.0.template_id"),
					resource.TestCheckResourceAttrSet(filterByTemplateId, "records.0.template_name"),
					resource.TestCheckResourceAttrSet(filterByTemplateId, "records.0.description"),
					resource.TestMatchResourceAttr(filterByTemplateId, "records.0.apply_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(filterByTemplateId, "records.0.type"),
					resource.TestCheckResourceAttrSet(filterByTemplateId, "records.0.account_id"),
					resource.TestCheckResourceAttrSet(filterByTemplateId, "records.0.resources.#"),
					resource.TestCheckResourceAttrSet(filterByTemplateId, "records.0.resources.0.status"),
					resource.TestCheckResourceAttrSet(filterByTemplateId, "records.0.resources.0.domain_name"),
					resource.TestCheckResourceAttrSet(filterByTemplateId, "records.0.configs"),
					// Filter by 'template_name' parameter.
					dcFilterByTemplateName.CheckResourceExists(),
					resource.TestCheckOutput("is_template_name_filter_useful", "true"),
					// Filter by 'operator_id' parameter.
					dcFilterByOperatorId.CheckResourceExists(),
					resource.TestCheckOutput("is_operator_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataDomainTemplateApplyRecords_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_domain_template" "test" {
  name        = "%[1]s"
  description = "Created by terraform for template apply records test"
  configs     = jsonencode({
    "cache_rules": [
      {
        "force_cache": "on",
        "follow_origin": "off",
        "match_type": "all",
        "priority": 1,
        "stale_while_revalidate": "off",
        "ttl": 20,
        "ttl_unit": "d",
        "url_parameter_type": "full_url",
        "url_parameter_value": ""
      }
    ],
    "origin_follow302_status": "off",
    "compress": {
      "type": "gzip",
      "status": "on",
      "file_type": ".js,.html,.css"
    }
  })
}

resource "huaweicloud_cdn_domain_template_apply" "test" {
  template_id = huaweicloud_cdn_domain_template.test.id
  resources   = "%[2]s"
}

# Without any filter parameter.
data "huaweicloud_cdn_domain_template_apply_records" "all" {
  depends_on = [
    huaweicloud_cdn_domain_template_apply.test
  ]
}

# Filter by 'template_id' parameter.
locals {
  template_id = huaweicloud_cdn_domain_template.test.id
}

data "huaweicloud_cdn_domain_template_apply_records" "filter_by_template_id" {
  depends_on = [
    huaweicloud_cdn_domain_template_apply.test
  ]

  template_id = local.template_id
}

locals {
  template_id_filter_result = [for v in data.huaweicloud_cdn_domain_template_apply_records.filter_by_template_id.records[*].template_id :
    v == local.template_id]
}

output "is_template_id_filter_useful" {
  value = length(local.template_id_filter_result) > 0 && alltrue(local.template_id_filter_result)
}

# Filter by 'template_name' parameter.
locals {	
  template_name = huaweicloud_cdn_domain_template.test.name
}

data "huaweicloud_cdn_domain_template_apply_records" "filter_by_template_name" {
  depends_on = [
    huaweicloud_cdn_domain_template_apply.test
  ]

  template_name = local.template_name
}

locals {
  template_name_filter_result = [for v in data.huaweicloud_cdn_domain_template_apply_records.filter_by_template_name.records[*].template_name :
    v == local.template_name]
}

output "is_template_name_filter_useful" {
  value = length(local.template_name_filter_result) > 0 && alltrue(local.template_name_filter_result)
}

# Filter by 'operator_id' parameter.
locals {
  operator_id = try(data.huaweicloud_cdn_domain_template_apply_records.all.records[0].operator_id, "NOT_FOUND")
}

data "huaweicloud_cdn_domain_template_apply_records" "filter_by_operator_id" {
  operator_id = local.operator_id
}

locals {
  operator_id_filter_result = [for v in data.huaweicloud_cdn_domain_template_apply_records.filter_by_operator_id.records[*].operator_id :
    v == local.operator_id]
}

output "is_operator_id_filter_useful" {
  value = length(local.operator_id_filter_result) > 0 && alltrue(local.operator_id_filter_result)
}
`, name, acceptance.HW_CDN_DOMAIN_NAME)
}
