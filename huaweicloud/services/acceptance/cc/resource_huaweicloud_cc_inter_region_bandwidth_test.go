package cc

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getInterRegionBandwidthResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		getInterRegionBandwidthHttpUrl = "v3/{domain_id}/ccaas/inter-region-bandwidths/{id}"
		getInterRegionBandwidthProduct = "cc"
	)
	getInterRegionBandwidthClient, err := cfg.NewServiceClient(getInterRegionBandwidthProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CC Client: %s", err)
	}

	getInterRegionBandwidthPath := getInterRegionBandwidthClient.Endpoint + getInterRegionBandwidthHttpUrl
	getInterRegionBandwidthPath = strings.ReplaceAll(getInterRegionBandwidthPath, "{domain_id}", cfg.DomainID)
	getInterRegionBandwidthPath = strings.ReplaceAll(getInterRegionBandwidthPath, "{id}", state.Primary.ID)

	getInterRegionBandwidthOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getInterRegionBandwidthResp, err := getInterRegionBandwidthClient.Request("GET", getInterRegionBandwidthPath, &getInterRegionBandwidthOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving inter-region bandwidth: %s", err)
	}

	getInterRegionBandwidthRespBody, err := utils.FlattenResponse(getInterRegionBandwidthResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving inter-region bandwidth: %s", err)
	}

	return getInterRegionBandwidthRespBody, nil
}

func TestAccInterRegionBandwidth_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cc_inter_region_bandwidth.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getInterRegionBandwidthResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckCustomRegion(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testInterRegionBandwidth_basic(name, 2),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "cloud_connection_id",
						"huaweicloud_cc_connection.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "bandwidth_package_id",
						"huaweicloud_cc_bandwidth_package.test", "id"),
					resource.TestCheckResourceAttr(rName, "bandwidth", "2"),
					resource.TestCheckResourceAttr(rName, "inter_region_ids.0", acceptance.HW_REGION_NAME),
					resource.TestCheckResourceAttr(rName, "inter_region_ids.1", acceptance.HW_CUSTOM_REGION_NAME),
					resource.TestCheckResourceAttr(rName, "inter_regions.#", "2"),
					resource.TestCheckResourceAttr(rName, "inter_regions.0.local_region_id", acceptance.HW_REGION_NAME),
					resource.TestCheckResourceAttr(rName, "inter_regions.0.remote_region_id", acceptance.HW_CUSTOM_REGION_NAME),
					resource.TestCheckResourceAttr(rName, "inter_regions.1.local_region_id", acceptance.HW_CUSTOM_REGION_NAME),
					resource.TestCheckResourceAttr(rName, "inter_regions.1.remote_region_id", acceptance.HW_REGION_NAME),
				),
			},
			{
				Config: testInterRegionBandwidth_basic(name, 3),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "cloud_connection_id",
						"huaweicloud_cc_connection.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "bandwidth_package_id",
						"huaweicloud_cc_bandwidth_package.test", "id"),
					resource.TestCheckResourceAttr(rName, "bandwidth", "3"),
					resource.TestCheckResourceAttr(rName, "inter_region_ids.0", acceptance.HW_REGION_NAME),
					resource.TestCheckResourceAttr(rName, "inter_region_ids.1", acceptance.HW_CUSTOM_REGION_NAME),
					resource.TestCheckResourceAttr(rName, "inter_regions.#", "2"),
					resource.TestCheckResourceAttr(rName, "inter_regions.0.local_region_id", acceptance.HW_REGION_NAME),
					resource.TestCheckResourceAttr(rName, "inter_regions.0.remote_region_id", acceptance.HW_CUSTOM_REGION_NAME),
					resource.TestCheckResourceAttr(rName, "inter_regions.1.local_region_id", acceptance.HW_CUSTOM_REGION_NAME),
					resource.TestCheckResourceAttr(rName, "inter_regions.1.remote_region_id", acceptance.HW_REGION_NAME),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testInterRegionBandwidth_basic(name string, bandwidth int) string {
	return fmt.Sprintf(`
resource "huaweicloud_cc_connection" "test" {
  name = "%[1]s"
}

resource "huaweicloud_cc_bandwidth_package" "test" {
  name           = "%[1]s"
  local_area_id  = "Chinese-Mainland"
  remote_area_id = "Chinese-Mainland"
  charge_mode    = "bandwidth"
  billing_mode   = 3
  bandwidth      = 6
  description    = "This is an accaptance test"
  resource_id    = huaweicloud_cc_connection.test.id
  resource_type  = "cloud_connection"
}

resource "huaweicloud_cc_inter_region_bandwidth" "test" {
  cloud_connection_id  = huaweicloud_cc_connection.test.id
  bandwidth_package_id = huaweicloud_cc_bandwidth_package.test.id
  bandwidth            = %[2]d
  inter_region_ids     = ["%[3]s", "%[4]s"]
}
`, name, bandwidth, acceptance.HW_REGION_NAME, acceptance.HW_CUSTOM_REGION_NAME)
}
