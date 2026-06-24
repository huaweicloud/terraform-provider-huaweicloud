package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/taurusdb"
)

func getHtapStarrocksUserResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	instanceID := state.Primary.Attributes["instance_id"]
	userName := state.Primary.Attributes["user_name"]

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return nil, fmt.Errorf("error creating TaurusDB client: %s", err)
	}

	users, err := taurusdb.QueryHtapStarrocksUsers(client, instanceID, userName)
	if err != nil {
		return nil, err
	}
	// If the length of users is 0, the resource does not exist.
	if len(users) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return users[0], nil
}

func TestAccTaurusDBHtapStarrocksUser_basic(t *testing.T) {
	var obj interface{}
	rName := "huaweicloud_taurusdb_htap_starrocks_user.test"
	userName := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getHtapStarrocksUserResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBHtapInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBHtapStarrocksUser_basic(userName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "user_name", userName),
					resource.TestCheckResourceAttr(rName, "dml", "3"),
					resource.TestCheckResourceAttr(rName, "ddl", "1"),
					resource.TestCheckResourceAttr(rName, "databases.#", "2"),
				),
			},
			{
				Config: testAccTaurusDBHtapStarrocksUser_update(userName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "user_name", userName),
					resource.TestCheckResourceAttr(rName, "dml", "0"),
					resource.TestCheckResourceAttr(rName, "ddl", "0"),
					resource.TestCheckResourceAttr(rName, "databases.#", "1"),
					resource.TestCheckResourceAttr(rName, "databases.0", "*"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccTaurusDBHtapStarrocksUser_basic(userName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_taurusdb_htap_starrocks_user" "test" {
  instance_id = "%[1]s"
  user_name   = "%[2]s"
  password    = "Test@12345678"
  databases   = ["sys", "information_schema"]
  dml         = 3
  ddl         = 1
}
`, acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID, userName)
}

func testAccTaurusDBHtapStarrocksUser_update(userName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_taurusdb_htap_starrocks_user" "test" {
  instance_id = "%[1]s"
  user_name   = "%[2]s"
  password    = "Test@87654321"
  dml         = 0
  ddl         = 0
  databases   = ["*"]
}
`, acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID, userName)
}
