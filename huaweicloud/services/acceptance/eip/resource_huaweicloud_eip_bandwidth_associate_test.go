package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eip"
)

func getEipBandwidthAssociateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region      = acceptance.HW_REGION_NAME
		product     = "vpc"
		publicipId  = state.Primary.Attributes["publicip_id"]
		bandwidthId = state.Primary.Attributes["bandwidth_id"]
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating VPC EIP client: %s", err)
	}

	return eip.GetEipBandwidthAssociate(client, publicipId, bandwidthId)
}

func TestAccEipBandwidthAssociate_basic(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_eip_bandwidth_associate.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getEipBandwidthAssociateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEipBandwidthAssociate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "publicip_id", "huaweicloud_vpc_eip.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "bandwidth_id", "huaweicloud_vpc_bandwidth.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "public_ip_address", "huaweicloud_vpc_eip.test", "address"),
					resource.TestCheckResourceAttrSet(rName, "publicip_type"),
					resource.TestCheckResourceAttrSet(rName, "ip_version"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccEipBandwidthAssociateImportState(rName),
				ImportStateVerifyIgnore: []string{
					"bandwidth_charge_mode",
					"bandwidth_size",
					"bandwidth_name",
				},
			},
		},
	})
}

func testAccEipBandwidthAssociate_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_bandwidth" "test" {
  name = "%[1]s"
  size = 5
}

resource "huaweicloud_vpc_eip" "test" {
  name = "%[1]s"

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    name        = "%[1]s-eip"
    size        = 5
    charge_mode = "traffic"
  }

  lifecycle {
    ignore_changes = [bandwidth]
  }
}
`, rName)
}

func testAccEipBandwidthAssociate_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_eip_bandwidth_associate" "test" {
  publicip_id  = huaweicloud_vpc_eip.test.id
  bandwidth_id = huaweicloud_vpc_bandwidth.test.id

  bandwidth_charge_mode = "bandwidth"
  bandwidth_size        = 5

  depends_on = [
    huaweicloud_vpc_eip.test,
    huaweicloud_vpc_bandwidth.test,
  ]
}
`, testAccEipBandwidthAssociate_base(rName))
}

func testAccEipBandwidthAssociateImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		bandwidthId := rs.Primary.Attributes["bandwidth_id"]
		publicipId := rs.Primary.Attributes["publicip_id"]
		if bandwidthId == "" || publicipId == "" {
			return "", fmt.Errorf("the bandwidth_id (%s) or publicip_id (%s) is nil", bandwidthId, publicipId)
		}

		return publicipId + "/" + bandwidthId, nil
	}
}
