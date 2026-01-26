package iam

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataV5ServicePrincipals_basic(t *testing.T) {
	var (
		byLanguageWithZhCn   = "data.huaweicloud_identityv5_service_principals.filter_by_language_with_zh_cn"
		dcByLanguageWithZhCn = acceptance.InitDataSourceCheck(byLanguageWithZhCn)

		byLanguageWithEnUs   = "data.huaweicloud_identityv5_service_principals.filter_by_language_with_en_us"
		dcByLanguageWithEnUs = acceptance.InitDataSourceCheck(byLanguageWithEnUs)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataV5ServicePrincipals_basic,
				Check: resource.ComposeTestCheckFunc(
					dcByLanguageWithZhCn.CheckResourceExists(),
					resource.TestMatchResourceAttr(byLanguageWithZhCn, "service_principals.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(byLanguageWithZhCn, "service_principals.0.service_principal"),
					resource.TestCheckResourceAttrSet(byLanguageWithZhCn, "service_principals.0.description"),
					resource.TestCheckResourceAttrSet(byLanguageWithZhCn, "service_principals.0.display_name"),
					resource.TestCheckResourceAttrSet(byLanguageWithZhCn, "service_principals.0.service_catalog"),
					dcByLanguageWithEnUs.CheckResourceExists(),
					resource.TestMatchResourceAttr(byLanguageWithEnUs, "service_principals.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestMatchResourceAttr(byLanguageWithEnUs, "service_principals.0.display_name", regexp.MustCompile(`^[a-zA-Z]`)),
					resource.TestMatchResourceAttr(byLanguageWithEnUs, "service_principals.0.description", regexp.MustCompile(`^[a-zA-Z]`)),
				),
			},
		},
	})
}

const testAccDataV5ServicePrincipals_basic = `
# Filter by 'language' parameter with 'zh-cn', Default value is 'zh-cn'.
data "huaweicloud_identityv5_service_principals" "filter_by_language_with_zh_cn" {}

# Filter by 'language' parameter with 'en-us'.
data "huaweicloud_identityv5_service_principals" "filter_by_language_with_en_us" {
  language = "en-us"
}
`
