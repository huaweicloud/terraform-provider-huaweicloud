package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceApplicationAcl_basic(t *testing.T) {
	var (
		rName = "data.huaweicloud_apig_application_acl.test"
		dc    = acceptance.InitDataSourceCheck(rName)

		rNameNotFound = "data.huaweicloud_apig_application_acl.not_found"
		dcNotFound    = acceptance.InitDataSourceCheck(rNameNotFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
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

data "huaweicloud_availability_zones" "test" {}

%[1]s

resource "huaweicloud_apig_instance" "test" {
  name                  = "%[2]s"
  edition               = "BASIC"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "0"

  availability_zones = try(slice(data.huaweicloud_availability_zones.test.names, 0, 1), null)
}

resource "huaweicloud_apig_application" "test" {
  count = 2

  instance_id = huaweicloud_apig_instance.test.id
  name        = format("%[2]s_%%d", count.index)
}
`, common.TestBaseNetwork(name), name)
}

func testAccDataSourceApplicationAcl_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_application_acl" "test" {
  instance_id    = huaweicloud_apig_instance.test.id
  application_id = huaweicloud_apig_application.test[0].id
  type           = "PERMIT"
  values         = var.cidrs
}

data "huaweicloud_apig_application_acl" "test" {
  depends_on = [
    huaweicloud_apig_application_acl.test
  ]

  instance_id    = huaweicloud_apig_instance.test.id
  application_id = huaweicloud_apig_application.test[0].id
}

data "huaweicloud_apig_application_acl" "not_found" {
  instance_id    = huaweicloud_apig_instance.test.id
  application_id = huaweicloud_apig_application.test[1].id
}
`, testAccDataSourceApplicationAcl_basic_base())
}

func TestAccDataSourceApplicationAcl_expectError(t *testing.T) {
	randUUID, _ := uuid.GenerateUUID()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceApplicationAcl_expectError(randUUID),
				ExpectError: regexp.MustCompile(fmt.Sprintf("App %s does not exist", randUUID)),
			},
		},
	})
}

func testAccDataSourceApplicationAcl_expectError_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

%[1]s

resource "huaweicloud_apig_instance" "test" {
  name                  = "%[2]s"
  edition               = "BASIC"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "0"

  availability_zones = try(slice(data.huaweicloud_availability_zones.test.names, 0, 1), null)
}
`, common.TestBaseNetwork(name), name)
}

func testAccDataSourceApplicationAcl_expectError(uuid string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_apig_application_acl" "test" {
  instance_id    = huaweicloud_apig_instance.test.id
  application_id = "%[2]s"
}
`, testAccDataSourceApplicationAcl_expectError_base(), uuid)
}
