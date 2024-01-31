package organizations

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getTrustedServiceResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	// getTrustedService: Query Organizations trusted service
	var (
		region                   = acceptance.HW_REGION_NAME
		getTrustedServiceHttpUrl = "v1/organizations/trusted-services"
		getTrustedServiceProduct = "organizations"
	)
	getTrustedServiceClient, err := cfg.NewServiceClient(getTrustedServiceProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Organizations Client: %s", err)
	}

	getTrustedServiceBasePath := getTrustedServiceClient.Endpoint + getTrustedServiceHttpUrl

	getTrustedServicePath := getTrustedServiceBasePath + buildGetTrustedServiceQueryParams("")

	getTrustedServiceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	var serviceName string
	var getTrustedServiceRespBody interface{}
getTrustedServicesLoop:
	for {
		getTrustedServiceResp, err := getTrustedServiceClient.Request("GET", getTrustedServicePath,
			&getTrustedServiceOpt)

		if err != nil {
			return nil, err
		}
		getTrustedServiceRespBody, err = utils.FlattenResponse(getTrustedServiceResp)
		if err != nil {
			return nil, err
		}

		trustedServices := utils.PathSearch("trusted_services", getTrustedServiceRespBody, nil)
		if trustedServices == nil {
			return nil, fmt.Errorf("error retrieving Organizations trusted service")
		}

		for _, trustedService := range trustedServices.([]interface{}) {
			servicePrincipal := utils.PathSearch("service_principal", trustedService, "").(string)
			if servicePrincipal == state.Primary.ID {
				serviceName = servicePrincipal
				break getTrustedServicesLoop
			}
		}
		marker := utils.PathSearch("page_info.next_marker", getTrustedServiceRespBody, nil)
		if marker == nil {
			break
		}
		getTrustedServicePath = getTrustedServiceBasePath + buildGetTrustedServiceQueryParams("")
	}

	if serviceName == "" {
		return nil, fmt.Errorf("error retrieving Organizations trusted service")
	}
	return getTrustedServiceRespBody, nil
}

func buildGetTrustedServiceQueryParams(marker string) string {
	res := "?limit=100"
	if marker != "" {
		res = fmt.Sprintf("%s&marker=%s", res, marker)
	}
	return res
}

func TestAccTrustedService_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_organizations_trusted_service.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getTrustedServiceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testTrustedService_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "service", "service.SecMaster"),
					resource.TestCheckResourceAttrSet(rName, "enabled_at"),
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

func testTrustedService_basic() string {
	return `
resource "huaweicloud_organizations_trusted_service" "test" {
  service = "service.SecMaster"
}
`
}
