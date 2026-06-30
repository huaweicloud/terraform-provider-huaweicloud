package cse

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
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
					resource.TestCheckOutput("all_custom_namespaces_set", "true"),
				),
			},
		},
	})
}

func testAccDataNacosNamespaces_basic_invalidEngine() string {
	randUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
variable "enterprise_project_id" {
  default = "%[1]s"
}

data "huaweicloud_cse_nacos_namespaces" "test" {
  engine_id             = "%[1]s"
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, randUUID.String())
}

func testAccDataNacosNamespaces_basic() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
variable "enterprise_project_id" {
  default = "%[1]s"
}

resource "huaweicloud_cse_nacos_namespace" "test" {
  engine_id             = "%[2]s"
  name                  = "%[3]s"
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

data "huaweicloud_cse_nacos_namespaces" "test" {
  depends_on = [huaweicloud_cse_nacos_namespace.test]

  engine_id             = "%[2]s"
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

# Check whether custom namespace ID is set
locals {
  namespaces_result = [
    for o in data.huaweicloud_cse_nacos_namespaces.test.namespaces : o if o.id == huaweicloud_cse_nacos_namespace.test.id &&
	  o.name == huaweicloud_cse_nacos_namespace.test.name && o.name != "public"
  ]
}

output "all_custom_namespaces_set" {
  value = length(local.namespaces_result) > 0
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST,
		acceptance.HW_CSE_NACOS_MICROSERVICE_ENGINE_ID, name)
}
