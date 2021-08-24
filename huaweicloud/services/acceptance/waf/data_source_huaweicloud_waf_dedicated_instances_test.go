package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccDataSourceWafDedicatedInstancesV1_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	resourceName := "data.huaweicloud_waf_dedicated_instances.instance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedInstancesV1_conf(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceWafDedicatedInstanceV1Exists(resourceName),

					resource.TestCheckResourceAttr(resourceName, "name", name+"_updated"),
					resource.TestCheckResourceAttr(resourceName, "instances.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "instances.0.available_zone"),
					resource.TestCheckResourceAttrSet(resourceName, "instances.0.cpu_flavor"),
					resource.TestCheckResourceAttrSet(resourceName, "instances.0.cpu_architecture"),
					resource.TestCheckResourceAttrSet(resourceName, "instances.0.security_group.#"),
					resource.TestCheckResourceAttrSet(resourceName, "instances.0.server_id"),
					resource.TestCheckResourceAttrSet(resourceName, "instances.0.service_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "instances.0.subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName, "instances.0.vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName, "instances.0.run_status"),
					resource.TestCheckResourceAttrSet(resourceName, "instances.0.access_status"),
					resource.TestCheckResourceAttrSet(resourceName, "instances.0.upgradable"),
				),
			},
		},
	})
}

func testAccCheckDataSourceWafDedicatedInstanceV1Exists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmtp.Errorf("Can't find WAF policies data source: %s.", r)
		}
		if rs.Primary.ID == "" {
			return fmtp.Errorf("The WAF policies data source ID does not set.")
		}
		return nil
	}
}

func testAccWafDedicatedInstancesV1_conf(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_waf_dedicated_instances" "instance" {
  name = huaweicloud_waf_dedicated_instance.instance_1.name
 
  depends_on = [huaweicloud_waf_dedicated_instance.instance_1]
}
`, testAccWafDedicatedInstanceV1_conf(name))
}
