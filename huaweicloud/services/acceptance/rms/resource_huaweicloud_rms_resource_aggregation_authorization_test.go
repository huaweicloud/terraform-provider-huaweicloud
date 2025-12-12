package rms

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getAggregationAuthResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region                    = acceptance.HW_REGION_NAME
		getAggregationAuthHttpUrl = "v1/resource-manager/domains/{domain_id}/aggregators/aggregation-authorization"
		getAggregationAuthProduct = "rms"
	)

	getAggregationAuthClient, err := cfg.NewServiceClient(getAggregationAuthProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RMS Client: %s", err)
	}

	getAggregationAuthPath := getAggregationAuthClient.Endpoint + getAggregationAuthHttpUrl
	getAggregationAuthPath = strings.ReplaceAll(getAggregationAuthPath, "{domain_id}", cfg.DomainID)
	getAggregationAuthPath += fmt.Sprintf("?account_id=%s", state.Primary.ID)

	getAggregationAuthOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getAggregationAuthResp, err := getAggregationAuthClient.Request("GET", getAggregationAuthPath, &getAggregationAuthOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving aggregation authorization: %s", err)
	}

	respBody, err := utils.FlattenResponse(getAggregationAuthResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving aggregation authorization: %s", err)
	}

	item := utils.PathSearch("aggregation_authorizations[0]", respBody, nil)
	if item == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return item, nil
}

func TestAccAggregationAuthorization_basic(t *testing.T) {
	var obj interface{}

	// the RMS service does not validate the account ID, so we can randomly generate it.
	accountID := acctest.RandStringFromCharSet(32, randomCharSet)
	rName := "huaweicloud_rms_resource_aggregation_authorization.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAggregationAuthResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAggregationAuthorizationorization_basic(accountID),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "account_id", accountID),
					resource.TestCheckResourceAttrSet(rName, "urn"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttr(rName, "tags.k1", "v1"),
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

func testAggregationAuthorizationorization_basic(account string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rms_resource_aggregation_authorization" "test" {
  account_id = "%s"
  tags       = {
    k1 = "v1"
  }
}
`, account)
}
