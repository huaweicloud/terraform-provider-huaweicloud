package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceApplicationAcl_basic(t *testing.T) {
	var (
		rName = "data.huaweicloud_apig_application_acl.test"
		dc    = acceptance.InitDataSourceCheck(rName)

		rNameNotFound = "data.huaweicloud_apig_application_acl.not_found"
		dcNotFound    = acceptance.InitDataSourceCheck(rNameNotFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceApplicationAcl_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "type", "PERMIT"),
					resource.TestCheckResourceAttr(rName, "values.#", "3"),
					resource.TestCheckResourceAttr(rName, "values.0", "172.16.0.0/20"),
					resource.TestCheckResourceAttr(rName, "values.1", "192.168.0.0/18"),
					resource.TestCheckResourceAttr(rName, "values.2", "127.0.0.1"),
					dcNotFound.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameNotFound, "type", ""),
					resource.TestCheckResourceAttr(rNameNotFound, "values.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceApplicationAcl_basic_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
variable "cidrs" {
  type    = list(string)
  default = ["172.16.0.0/20", "192.168.0.0/18", "127.0.0.1"]
}

data "huaweicloud_apig_instances" "test" {
  instance_id = "%[1]s"
}

locals {
  instance_id = data.huaweicloud_apig_instances.test.instances[0].id
}

resource "huaweicloud_apig_application" "test" {
  count = 2

  instance_id = local.instance_id
  name        = format("%[2]s_%%d", count.index)
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccDataSourceApplicationAcl_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_application_acl" "test" {
  instance_id    = local.instance_id
  application_id = huaweicloud_apig_application.test[0].id
  type           = "PERMIT"
  values         = var.cidrs
}

data "huaweicloud_apig_application_acl" "test" {
  depends_on = [
    huaweicloud_apig_application_acl.test
  ]

  instance_id    = local.instance_id
  application_id = huaweicloud_apig_application.test[0].id
}

data "huaweicloud_apig_application_acl" "not_found" {
  instance_id    = local.instance_id
  application_id = huaweicloud_apig_application.test[1].id
}
`, testAccDataSourceApplicationAcl_basic_base())
}

func TestAccDataSourceApplicationAcl_expectError(t *testing.T) {
	randUUID, _ := uuid.GenerateUUID()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceApplicationAcl_expectError(randUUID),
				ExpectError: regexp.MustCompile(fmt.Sprintf("App %s does not exist", randUUID)),
			},
		},
	})
}

func testAccDataSourceApplicationAcl_expectError(uuid string) string {
	return fmt.Sprintf(`
data "huaweicloud_apig_application_acl" "test" {
  instance_id    = "%[1]s"
  application_id = "%[2]s"
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, uuid)
}
