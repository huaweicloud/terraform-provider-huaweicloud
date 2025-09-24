package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDomainRouteUpdate_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// The value of the `routes` parameter is obtained through API queries, if this value is not used, the
			// resource execution will report an error.
			// So when the value of `instance_id` changes, the execution of this test case may result in an error.
			acceptance.TestAccPreCheckWafDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainRouteUpdate_basic(),
			},
		},
	})
}

func testAccDomainRouteUpdate_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_domain_route_update" "test" {
  instance_id = "%[1]s"

  routes {
    name  = "Beijing"
    cname = "a59802a4b1c041328e1301a6c46eae7d"

    servers {
      back_protocol = "HTTP"
      address       = "167.168.0.16"
      port          = 80
    }
  }

  routes {
    name  = "Langfang"
    cname = "62371ea53d584cfea6239a4594c9c529"

    servers {
      back_protocol = "HTTP"
      address       = "167.168.0.16"
      port          = 80
    }
  }

  routes {
    name  = "Internal"
    cname = "3c360f2a90d243bb84edca4aec080722"

    servers {
      back_protocol = "HTTP"
      address       = "167.168.0.16"
      port          = 80
    }
  }
}
`, acceptance.HW_WAF_DOMAIN_ID)
}
