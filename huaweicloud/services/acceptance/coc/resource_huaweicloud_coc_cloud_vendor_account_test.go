package coc

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getCloudVendorAccountResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v1/vendor-account?limit=100"
		product = "coc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating COC client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getResp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error retrieving COC cloud vendor account: %s", err)
	}

	getRespJson, err := json.Marshal(getResp)
	if err != nil {
		return nil, err
	}
	var getRespBody interface{}
	err = json.Unmarshal(getRespJson, &getRespBody)
	if err != nil {
		return nil, err
	}

	cloudVendorAccount := utils.PathSearch(fmt.Sprintf("data[?id=='%s']|[0]", state.Primary.ID), getRespBody, nil)
	if cloudVendorAccount == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func TestAccResourceCloudVendorAccount_basic(t *testing.T) {
	var obj interface{}
	name := acceptance.RandomAccResourceName()
	nameUpdate := acceptance.RandomAccResourceName()
	rName := "huaweicloud_coc_cloud_vendor_account.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCloudVendorAccountResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudVendorAccount_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "vendor", "HCS"),
					resource.TestCheckResourceAttr(rName, "account_id", "test_account_id"),
					resource.TestCheckResourceAttr(rName, "account_name", name),
					resource.TestCheckResourceAttr(rName, "ak", "test_ak"),
					resource.TestCheckResourceAttrSet(rName, "domain_id"),
					resource.TestCheckResourceAttrSet(rName, "sync_status"),
					resource.TestCheckResourceAttrSet(rName, "failure_msg"),
					resource.TestCheckResourceAttrSet(rName, "sync_date"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
				),
			},
			{
				Config: testAccCloudVendorAccount_update(nameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "vendor", "HCS"),
					resource.TestCheckResourceAttr(rName, "account_id", "test_account_id"),
					resource.TestCheckResourceAttr(rName, "account_name", nameUpdate),
					resource.TestCheckResourceAttr(rName, "ak", "test_ak_update"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"sk"},
			},
		},
	})
}

func testAccCloudVendorAccount_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_coc_cloud_vendor_account" "test" {
  vendor       = "HCS"
  account_id   = "test_account_id"
  account_name = "%s"
  ak           = "test_ak"
  sk           = "test_sk"
}
`, name)
}

func testAccCloudVendorAccount_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_coc_cloud_vendor_account" "test" {
  vendor       = "HCS"
  account_id   = "test_account_id"
  account_name = "%s"
  ak           = "test_ak_update"
  sk           = "test_sk_update"
}
`, name)
}
