package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataAppcodes_basic(t *testing.T) {
	var (
		rName = acceptance.RandomAccResourceName()

		dataSource = "data.huaweicloud_apig_appcodes.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		notFound   = "data.huaweicloud_apig_appcodes.not_found"
		dcNotFound = acceptance.InitDataSourceCheck(notFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataAppcodes_nonExistentInstance(),
				ExpectError: regexp.MustCompile(`error querying APPCODEs`),
			},
			{
				Config:      testAccDataAppcodes_nonExistentApplication(),
				ExpectError: regexp.MustCompile(`App [0-9a-fA-F-]{36} does not exist`),
			},
			{
				Config: testAccDataAppcodes_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttrPair(dataSource, "application_id", "huaweicloud_apig_application.test.0", "id"),
					resource.TestCheckResourceAttr(dataSource, "appcodes.#", "1"),
					resource.TestCheckResourceAttrSet(dataSource, "appcodes.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "appcodes.0.value"),
					resource.TestCheckResourceAttrPair(dataSource, "appcodes.0.application_id", "huaweicloud_apig_application.test.0", "id"),
					resource.TestMatchResourceAttr(dataSource, "appcodes.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					dcNotFound.CheckResourceExists(),
					resource.TestCheckResourceAttr(notFound, "appcodes.#", "0"),
				),
			},
		},
	})
}

func testAccDataAppcodes_nonExistentInstance() string {
	randomUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
data "huaweicloud_apig_appcodes" "test" {
  instance_id    = "%[1]s"
  application_id = "%[1]s"
}
`, randomUUID.String())
}

func testAccDataAppcodes_nonExistentApplication() string {
	randomUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
data "huaweicloud_apig_appcodes" "test" {
  instance_id    = "%[1]s"
  application_id = "%[2]s"
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, randomUUID.String())
}

func testAccDataAppcodes_base(name string) string {
	return fmt.Sprintf(`
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

resource "huaweicloud_apig_appcode" "test" {
  instance_id    = local.instance_id
  application_id = huaweicloud_apig_application.test[0].id
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccDataAppcodes_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_apig_appcodes" "test" {
  depends_on = [
    huaweicloud_apig_appcode.test,
  ]

  instance_id    = local.instance_id
  application_id = huaweicloud_apig_application.test[0].id
}

data "huaweicloud_apig_appcodes" "not_found" {
  instance_id    = local.instance_id
  application_id = huaweicloud_apig_application.test[1].id
}
`, testAccDataAppcodes_base(name))
}
