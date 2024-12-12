package cse

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataNacosNamespaces_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_cse_nacos_namespaces.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCSENacosMicroserviceEngineID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataNacosNamespaces_basic_invalidEngine(),
				ExpectError: regexp.MustCompile(`the Nacos engine \([0-9a-f-]+\) does not exist`),
			},
			{
				Config: testAccDataNacosNamespaces_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "namespaces.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckOutput("all_custom_namespace_ids_set", "true"),
				),
			},
		},
	})
}

func testAccDataNacosNamespaces_basic_invalidEngine() string {
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
data "huaweicloud_cse_nacos_namespaces" "test" {
  engine_id = "%[1]s"
}
`, randUUID)
}

func testAccDataNacosNamespaces_basic() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
resource "huaweicloud_cse_nacos_namespace" "test" {
  engine_id = "%[1]s"
  name      = "%[2]s"
}

data "huaweicloud_cse_nacos_namespaces" "test" {
  depends_on = [huaweicloud_cse_nacos_namespace.test]

  engine_id = "%[1]s"
}

# Check whether custom namespace ID is set
locals {
  namespace_id_validate_result = [
    for o in data.huaweicloud_cse_nacos_namespaces.test.namespaces : o.id != "" if o.name != "public" 
  ]
}

output "all_custom_namespace_ids_set" {
  value = length(local.namespace_id_validate_result) > 0 && alltrue(local.namespace_id_validate_result)
}
`, acceptance.HW_CSE_NACOS_MICROSERVICE_ENGINE_ID, name)
}
