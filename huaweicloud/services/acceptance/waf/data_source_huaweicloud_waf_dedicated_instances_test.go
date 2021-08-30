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
	resourceName1 := "data.huaweicloud_waf_dedicated_instances.instance_1"
	resourceName2 := "data.huaweicloud_waf_dedicated_instances.instance_2"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedInstancesV1_conf(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceWafDedicatedInstanceV1Exists(resourceName1),

					resource.TestCheckResourceAttr(resourceName1, "name", name),
					resource.TestCheckResourceAttr(resourceName1, "instances.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.available_zone"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.cpu_flavor"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.cpu_architecture"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.security_group.#"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.server_id"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.service_ip"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.run_status"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.access_status"),
					resource.TestCheckResourceAttrSet(resourceName1, "instances.0.upgradable"),

					resource.TestCheckResourceAttr(resourceName2, "instances.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName2, "name"),
					resource.TestCheckResourceAttrSet(resourceName2, "instances.0.available_zone"),
				),
			},
		},
	})
}

func testAccCheckDataSourceWafDedicatedInstanceV1Exists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmtp.Errorf("Can't find WAF dedicated instance data source: %s.", r)
		}
		if rs.Primary.ID == "" {
			return fmtp.Errorf("The WAF dedicated instance data source ID not set.")
		}
		return nil
	}
}

func testAccWafDedicatedInstancesV1_conf(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_waf_dedicated_instances" "instance_1" {
  name = huaweicloud_waf_dedicated_instance.instance_1.name

  depends_on = [
    huaweicloud_waf_dedicated_instance.instance_1
  ]
}

data "huaweicloud_waf_dedicated_instances" "instance_2" {
  id   = huaweicloud_waf_dedicated_instance.instance_1.id
  name = huaweicloud_waf_dedicated_instance.instance_1.name

  depends_on = [
    huaweicloud_waf_dedicated_instance.instance_1
  ]
}
`, testAccWafDedicatedInstanceV1_conf(name))
}
