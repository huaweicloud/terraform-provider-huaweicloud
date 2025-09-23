package cbr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVaultsByTags_basic(t *testing.T) {
	var (
		filter_by_tags_any_all       = "data.huaweicloud_cbr_vaults_by_tags.filter_by_tags_any"
		dc_tags_any                  = acceptance.InitDataSourceCheck(filter_by_tags_any_all)
		filter_by_tags_test_generate = "data.huaweicloud_cbr_vaults_by_tags.filter_by_tags_test_generate"
		dc_tags_test_generate        = acceptance.InitDataSourceCheck(filter_by_tags_test_generate)
		testName                     = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataTagsVaultsByTags_basic1(testName),
				Check: resource.ComposeTestCheckFunc(
					dc_tags_any.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(filter_by_tags_any_all, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(filter_by_tags_any_all, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(filter_by_tags_any_all, "resources.0.resource_detail.0.vault.0.id"),
					resource.TestCheckResourceAttrSet(filter_by_tags_any_all, "resources.0.resource_detail.0.vault.0.id"),
					resource.TestCheckResourceAttrSet(filter_by_tags_any_all, "resources.0.resource_detail.0.vault.0.name"),
					resource.TestCheckResourceAttrSet(filter_by_tags_any_all, "resources.0.resource_detail.0.vault.0.provider_id"),
					resource.TestCheckResourceAttrSet(filter_by_tags_any_all, "resources.0.resource_detail.0.vault.0.tags.0.key"),
					resource.TestCheckResourceAttrSet(filter_by_tags_any_all, "resources.0.resource_detail.0.vault.0.tags.0.value"),
					resource.TestCheckResourceAttrSet(filter_by_tags_any_all, "resources.0.resource_detail.0.vault.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(filter_by_tags_any_all, "resources.0.resource_detail.0.vault.0.auto_bind"),
					resource.TestCheckResourceAttrSet(filter_by_tags_any_all, "resources.0.resource_detail.0.vault.0.created_at"),
					resource.TestCheckResourceAttrSet(filter_by_tags_any_all, "resources.0.resource_detail.0.vault.0.auto_expand"),
					resource.TestCheckResourceAttrSet(filter_by_tags_any_all, "resources.0.resource_detail.0.vault.0.smn_notify"),
					resource.TestCheckResourceAttrSet(filter_by_tags_any_all, "resources.0.resource_detail.0.vault.0.threshold"),
					resource.TestCheckResourceAttrSet(filter_by_tags_any_all, "resources.0.resource_detail.0.vault.0.locked"),
					resource.TestCheckResourceAttrSet(filter_by_tags_any_all, "resources.0.resource_detail.0.vault.0.billing.0.allocated"),
					resource.TestCheckResourceAttrSet(filter_by_tags_any_all, "resources.0.resource_detail.0.vault.0.billing.0.charging_mode"),
					resource.TestCheckResourceAttrSet(filter_by_tags_any_all, "resources.0.resource_detail.0.vault.0.billing.0.cloud_type"),
					resource.TestCheckResourceAttrSet(filter_by_tags_any_all, "resources.0.resource_detail.0.vault.0.billing.0.consistent_level"),
					resource.TestCheckResourceAttrSet(filter_by_tags_any_all, "resources.0.resource_detail.0.vault.0.billing.0.object_type"),
					resource.TestCheckResourceAttrSet(filter_by_tags_any_all, "resources.0.resource_detail.0.vault.0.billing.0.protect_type"),
					resource.TestCheckResourceAttrSet(filter_by_tags_any_all, "resources.0.resource_detail.0.vault.0.billing.0.size"),
					resource.TestCheckResourceAttrSet(filter_by_tags_any_all, "resources.0.resource_detail.0.vault.0.billing.0.spec_code"),
					resource.TestCheckResourceAttrSet(filter_by_tags_any_all, "resources.0.resource_detail.0.vault.0.billing.0.status"),
					resource.TestCheckResourceAttrSet(filter_by_tags_any_all, "resources.0.resource_detail.0.vault.0.billing.0.used"),
					resource.TestCheckResourceAttrSet(filter_by_tags_any_all, "resources.0.resource_detail.0.vault.0.billing.0.is_multi_az"),
				),
			},
			{
				Config: testAccDataTagsVaultsByTags_basic2(testName),
				Check: resource.ComposeTestCheckFunc(
					dc_tags_test_generate.CheckResourceExists(),
					resource.TestCheckResourceAttr(filter_by_tags_test_generate, "total_count", "2"),
				),
			},
		},
	})
}

func testAccDataVaultsByTags_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_vault" "test" {
  count = 2
  
  name                  = format("%[1]s_%%d", count.index)
  type                  = "server"
  consistent_level      = "crash_consistent"
  protection_type       = "backup"
  size                  = 200
  enterprise_project_id = "0"
  
  tags = {
	"foo${count.index}" = "bar${count.index}"
	foo_all = "bar_all"
  }
}
`, name)
}

func testAccDataTagsVaultsByTags_basic1(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cbr_vaults_by_tags" "filter_by_tags" {
  depends_on = [huaweicloud_cbr_vault.test]
  action     = "filter"
  tags {
    key    = "foo0"
    values = ["bar0"]
  }
  tags {
    key    = "foo1"
    values = ["bar1"]
  }
}

data "huaweicloud_cbr_vaults_by_tags" "filter_by_tags_any" {
  depends_on = [huaweicloud_cbr_vault.test]
  action     = "filter"
  tags_any {
    key    = "foo0"
    values = ["bar0"]
  }
  tags_any {
    key    = "foo1"
    values = ["bar1"]
  }
}

output "generated_names" {
  value = [for vault in huaweicloud_cbr_vault.test[*] : vault.name]
}

output "retrieved_names" {
  value = [for vault in data.huaweicloud_cbr_vaults_by_tags.filter_by_tags.resources[*] : vault.resource_name]
}`, testAccDataVaultsByTags_base(name))
}

func testAccDataTagsVaultsByTags_basic2(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cbr_vaults_by_tags" "filter_by_tags_test_generate" {
  depends_on = [huaweicloud_cbr_vault.test]
  action     = "count"
  tags_any {
    key    = "foo0"
    values = ["bar0"]
  }
  tags_any {
    key    = "foo1"
    values = ["bar1"]
  }
}

`, testAccDataVaultsByTags_base(name))
}
