package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestDataSourceInstanceIngressAssociatedDomains_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_apig_instance_ingress_associated_domains.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceInstanceIngressAssociatedDomains_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "domain_infos.#", regexp.MustCompile("^[1-9]([0-9]+)?")),
					resource.TestCheckResourceAttrSet(dcName, "domain_infos.0.group_id"),
					resource.TestCheckResourceAttrSet(dcName, "domain_infos.0.group_name"),
					resource.TestCheckResourceAttrSet(dcName, "domain_infos.0.domain_name"),
				),
			},
		},
	})
}

func testAccDataSourceInstanceIngressAssociatedDomains_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_apig_instances" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_apig_instance_ingress_associated_domains" "test" {
  instance_id     = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  ingress_port_id = "%[2]s"
}`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, acceptance.HW_APIG_DEDICATED_INSTANCE_INGRESS_PORT_ID)
}
