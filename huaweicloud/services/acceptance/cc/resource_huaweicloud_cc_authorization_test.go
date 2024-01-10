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

func getAuthorizationResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		getAuthorizationHttpUrl = "v3/{domain_id}/ccaas/authorisations?id={id}"
		getAuthorizationProduct = "cc"
	)
	getAuthorizationClient, err := cfg.NewServiceClient(getAuthorizationProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CC client: %s", err)
	}

	getAuthorizationPath := getAuthorizationClient.Endpoint + getAuthorizationHttpUrl
	getAuthorizationPath = strings.ReplaceAll(getAuthorizationPath, "{domain_id}", cfg.DomainID)
	getAuthorizationPath = strings.ReplaceAll(getAuthorizationPath, "{id}", state.Primary.ID)

	getAuthorizationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getAuthorizationResp, err := getAuthorizationClient.Request("GET", getAuthorizationPath, &getAuthorizationOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving authorization: %s", err)
	}

	getAuthorizationRespBody, err := utils.FlattenResponse(getAuthorizationResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving authorization: %s", err)
	}

	jsonPath := fmt.Sprintf("authorisations[?id =='%s']|[0]", state.Primary.ID)
	getAuthorizationRespBody = utils.PathSearch(jsonPath, getAuthorizationRespBody, nil)
	if getAuthorizationRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return getAuthorizationRespBody, nil
}

func TestAccAuthorization_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cc_authorization.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAuthorizationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCCAuth(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAuthorization_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "instance_type", "vpc"),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttr(rName, "cloud_connection_domain_id", acceptance.HW_CC_PEER_DOMAIN_ID),
					resource.TestCheckResourceAttr(rName, "cloud_connection_id", acceptance.HW_CC_PEER_CONNECTION_ID),
					resource.TestCheckResourceAttr(rName, "description", "This is a test"),
				),
			},
			{
				Config: testAuthorization_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttr(rName, "instance_type", "vpc"),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttr(rName, "cloud_connection_domain_id", acceptance.HW_CC_PEER_DOMAIN_ID),
					resource.TestCheckResourceAttr(rName, "cloud_connection_id", acceptance.HW_CC_PEER_CONNECTION_ID),
					resource.TestCheckResourceAttr(rName, "description", "This is a test update"),
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

func testAuthorization_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name        = "%s"
  cidr        = "192.168.0.0/16"
}

resource "huaweicloud_cc_authorization" "test" {
   name                       = "%[1]s"
   instance_type              = "vpc"
   instance_id                = huaweicloud_vpc.test.id
   cloud_connection_domain_id = "%[2]s"
   cloud_connection_id        = "%[3]s"
   description                = "This is a test"
}
`, name, acceptance.HW_CC_PEER_DOMAIN_ID, acceptance.HW_CC_PEER_CONNECTION_ID)
}

func testAuthorization_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name        = "%s"
  cidr        = "192.168.0.0/16"
}

resource "huaweicloud_cc_authorization" "test" {
   name                       = "%[1]s_update"
   instance_type              = "vpc"
   instance_id                = huaweicloud_vpc.test.id
   cloud_connection_domain_id = "%[2]s"
   cloud_connection_id        = "%[3]s"
   description                = "This is a test update"
}
`, name, acceptance.HW_CC_PEER_DOMAIN_ID, acceptance.HW_CC_PEER_CONNECTION_ID)
}
