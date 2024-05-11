package cdn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getBillingOptionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	hcCdnClient, err := cfg.HcCdnV2Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating CDN v2 client: %s", err)
	}

	request := model.ShowChargeModesRequest{
		ProductType: state.Primary.Attributes["product_type"],
	}

	resp, err := hcCdnClient.ShowChargeModes(&request)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CDN billing option: %s", err)
	}

	if resp == nil || resp.Result == nil || len(*resp.Result) == 0 {
		return nil, fmt.Errorf("error retrieving CDN billing option: Result is not found in API response")
	}

	resultArray := *resp.Result
	return resultArray[0], nil
}

func TestAccBillingOption_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_cdn_billing_option.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getBillingOptionResourceFunc,
	)

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testBillingOption_basic,
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
				Config: testBillingOption_basic_update,
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

const testBillingOption_basic = `
resource "huaweicloud_cdn_billing_option" "test" {
  charge_mode  = "flux"
  product_type = "base"
  service_area = "mainland_china"
}
`

const testBillingOption_basic_update = `
resource "huaweicloud_cdn_billing_option" "test" {
  charge_mode  = "bw"
  product_type = "base"
  service_area = "mainland_china"
}
`
