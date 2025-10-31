package cdn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cdn"
)

func getBillingOptionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	productType := state.Primary.Attributes["product_type"]
	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return nil, fmt.Errorf("error creating CDN client: %s", err)
	}

	return cdn.GetBillingOptionDetail(client, productType)
}

func TestAccBillingOption_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_cdn_billing_option.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getBillingOptionResourceFunc,
	)

	// Avoid CheckDestroy, because there is nothing in the resource destroy method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBillingOption_basic_step1,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "charge_mode", "flux"),
					resource.TestCheckResourceAttr(rName, "product_type", "base"),
					resource.TestCheckResourceAttr(rName, "service_area", "mainland_china"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "effective_time"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "current_charge_mode"),
				),
			},
			{
				Config: testAccBillingOption_basic_step2,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "charge_mode", "bw"),
					resource.TestCheckResourceAttr(rName, "product_type", "base"),
					resource.TestCheckResourceAttr(rName, "service_area", "mainland_china"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "effective_time"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "current_charge_mode"),
				),
			},
		},
	})
}

const testAccBillingOption_basic_step1 = `
resource "huaweicloud_cdn_billing_option" "test" {
  charge_mode  = "flux"
  product_type = "base"
  service_area = "mainland_china"
}
`

const testAccBillingOption_basic_step2 = `
resource "huaweicloud_cdn_billing_option" "test" {
  charge_mode  = "bw"
  product_type = "base"
  service_area = "mainland_china"
}
`
