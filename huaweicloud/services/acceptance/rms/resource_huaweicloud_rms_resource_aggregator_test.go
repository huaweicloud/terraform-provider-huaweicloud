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

const randomCharSet string = "abcdef012346789"

func getAggregatorResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region               = acceptance.HW_REGION_NAME
		getAggregatorHttpUrl = "v1/resource-manager/domains/{domain_id}/aggregators/{id}"
		getAggregatorProduct = "rms"
	)

	getAggregatorClient, err := cfg.NewServiceClient(getAggregatorProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RMS Client: %s", err)
	}

	getAggregatorPath := getAggregatorClient.Endpoint + getAggregatorHttpUrl
	getAggregatorPath = strings.ReplaceAll(getAggregatorPath, "{domain_id}", cfg.DomainID)
	getAggregatorPath = strings.ReplaceAll(getAggregatorPath, "{id}", state.Primary.ID)

	getAggregatorOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getAggregatorResp, err := getAggregatorClient.Request("GET", getAggregatorPath, &getAggregatorOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving aggregator: %s", err)
	}
	return utils.FlattenResponse(getAggregatorResp)
}

func TestAccAggregator_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_rms_resource_aggregator.test"

	// the RMS service does not validate the account IDs, so we can randomly generate them.
	oneAccount := acctest.RandStringFromCharSet(32, randomCharSet)
	twoAccount := acctest.RandStringFromCharSet(32, randomCharSet)
	basicAccountIDs := []string{oneAccount}
	updateAccountIDs := []string{twoAccount, oneAccount}

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAggregatorResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAggregator_config(name, basicAccountIDs),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "ACCOUNT"),
					resource.TestCheckResourceAttrSet(rName, "urn"),
					resource.TestCheckResourceAttr(rName, "tags.k1", "v1"),
				),
			},
			{
				Config: testAggregator_config_update(name, updateAccountIDs),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "ACCOUNT"),
					resource.TestCheckResourceAttrSet(rName, "urn"),
					resource.TestCheckResourceAttr(rName, "tags.k2", "v2"),
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

func testAggregator_config(name string, accounts []string) string {
	var hclAccounts string
	for _, id := range accounts {
		hclAccounts += fmt.Sprintf("\"%s\", ", id)
	}

	return fmt.Sprintf(`
resource "huaweicloud_rms_resource_aggregator" "test" {
  name        = "%s"
  type        = "ACCOUNT"
  account_ids = [%s]
  
  tags = {
    k1 = "v1"
  }
}
`, name, hclAccounts)
}

func testAggregator_config_update(name string, accounts []string) string {
	var hclAccounts string
	for _, id := range accounts {
		hclAccounts += fmt.Sprintf("\"%s\", ", id)
	}

	return fmt.Sprintf(`
resource "huaweicloud_rms_resource_aggregator" "test" {
  name        = "%s"
  type        = "ACCOUNT"
  account_ids = [%s]
  
  tags = {
    k2 = "v2"
  }
}
`, name, hclAccounts)
}
