package rms

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

func getAdvancedQueryResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getAdvancedQuery: Query the RMS advanced query
	var (
		getAdvancedQueryHttpUrl = "v1/resource-manager/domains/{domain_id}/stored-queries/{query_id}"
		getAdvancedQueryProduct = "rms"
	)
	getAdvancedQueryClient, err := cfg.NewServiceClient(getAdvancedQueryProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RMS client: %s", err)
	}

	getAdvancedQueryPath := getAdvancedQueryClient.Endpoint + getAdvancedQueryHttpUrl
	getAdvancedQueryPath = strings.ReplaceAll(getAdvancedQueryPath, "{domain_id}", cfg.DomainID)
	getAdvancedQueryPath = strings.ReplaceAll(getAdvancedQueryPath, "{query_id}", state.Primary.ID)

	getAdvancedQueryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getAdvancedQueryResp, err := getAdvancedQueryClient.Request("GET", getAdvancedQueryPath, &getAdvancedQueryOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RMS advanced query: %s", err)
	}

	getAdvancedQueryRespBody, err := utils.FlattenResponse(getAdvancedQueryResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RMS advanced query: %s", err)
	}

	return getAdvancedQueryRespBody, nil
}

func TestAccAdvancedQuery_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_rms_advanced_query.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAdvancedQueryResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAdvancedQuery_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "expression", "select colume_1 from table_1"),
					resource.TestCheckResourceAttr(rName, "description", "test_description"),
				),
			},
			{
				Config: testAdvancedQuery_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "expression", "update table_1 set volume_1 = 5"),
					resource.TestCheckResourceAttr(rName, "description", "test_description_update"),
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

func TestAccAdvancedQuery_aggregator(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_rms_advanced_query.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAdvancedQueryResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAdvancedQuery_aggregator(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "aggregator"),
					resource.TestCheckResourceAttr(rName, "expression", "SELECT set_agg(domain_id) AS domain_ids FROM aggregator_resources"),
					resource.TestCheckResourceAttr(rName, "description", "test_description"),
				),
			},
			{
				Config: testAdvancedQuery_aggregator_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "type", "aggregator"),
					resource.TestCheckResourceAttr(rName, "expression", "SELECT id FROM aggregator_resources WHERE ep_id = '0'"),
					resource.TestCheckResourceAttr(rName, "description", "test_description_update"),
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

func testAdvancedQuery_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rms_advanced_query" "test" {
  name        = "%s"
  expression  = "select colume_1 from table_1"
  description = "test_description"
}
`, name)
}

func testAdvancedQuery_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rms_advanced_query" "test" {
  name        = "%s-update"
  expression  = "update table_1 set volume_1 = 5"
  description = "test_description_update"
}
`, name)
}

func testAdvancedQuery_aggregator(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rms_advanced_query" "test" {
  name        = "%s"
  type        = "aggregator"
  expression  = "SELECT set_agg(domain_id) AS domain_ids FROM aggregator_resources"
  description = "test_description"
}
`, name)
}

func testAdvancedQuery_aggregator_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rms_advanced_query" "test" {
  name        = "%s-update"
  type        = "aggregator"
  expression  = "SELECT id FROM aggregator_resources WHERE ep_id = '0'"
  description = "test_description_update"
}
`, name)
}
