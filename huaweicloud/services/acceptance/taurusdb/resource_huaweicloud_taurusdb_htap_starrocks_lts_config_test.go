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

func getTaurusDBHtapStarrocksLtsConfigResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		instanceId = state.Primary.Attributes["instance_id"]
		logType    = state.Primary.Attributes["log_type"]
	)

	client, err := cfg.NewServiceClient("gaussdb", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating TaurusDB client: %s", err)
	}

	ltsConfig, err := taurusdb.GetTaurusDBHtapStarrocksLtsConfig(client, instanceId, logType)
	if err != nil {
		return nil, &golangsdk.ErrDefault404{}
	}
	return ltsConfig, nil
}

func TestAccTaurusDBHtapStarrocksLtsConfig_basic(t *testing.T) {
	var obj interface{}
	logType := "error_log"
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_taurusdb_htap_starrocks_lts_config.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getTaurusDBHtapStarrocksLtsConfigResourceFunc,
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
				Config: testAccTaurusDBHtapStarrocksLtsConfig_basic(rName, logType),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "instance_id", acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID),
					resource.TestCheckResourceAttrPair(resourceName, "lts_group_id",
						"huaweicloud_lts_group.test.0", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "lts_stream_id",
						"huaweicloud_lts_stream.test.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "log_type", logType),
				),
			},
			{
				Config: testAccTaurusDBHtapStarrocksLtsConfig_update(rName, logType),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "instance_id", acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID),
					resource.TestCheckResourceAttrPair(resourceName, "lts_group_id",
						"huaweicloud_lts_group.test.1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "lts_stream_id",
						"huaweicloud_lts_stream.test.1", "id"),
					resource.TestCheckResourceAttr(resourceName, "log_type", logType),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccTaurusDBHtapStarrocksLtsConfigImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccTaurusDBHtapStarrocksLtsConfig_slowLog(t *testing.T) {
	var obj interface{}
	logType := "slow_log"
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_taurusdb_htap_starrocks_lts_config.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getTaurusDBHtapStarrocksLtsConfigResourceFunc,
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
				Config: testAccTaurusDBHtapStarrocksLtsConfig_basic(rName, logType),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "instance_id", acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID),
					resource.TestCheckResourceAttrPair(resourceName, "lts_group_id",
						"huaweicloud_lts_group.test.0", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "lts_stream_id",
						"huaweicloud_lts_stream.test.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "log_type", logType),
				),
			},
			{
				Config: testAccTaurusDBHtapStarrocksLtsConfig_update(rName, logType),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "instance_id", acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID),
					resource.TestCheckResourceAttrPair(resourceName, "lts_group_id",
						"huaweicloud_lts_group.test.1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "lts_stream_id",
						"huaweicloud_lts_stream.test.1", "id"),
					resource.TestCheckResourceAttr(resourceName, "log_type", logType),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccTaurusDBHtapStarrocksLtsConfigImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccTaurusDBHtapStarrocksLtsConfig_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  count = 2

  group_name  = "%[1]s_${count.index}"
  ttl_in_days = 1
}

resource "huaweicloud_lts_stream" "test" {
  count = 2

  group_id    = huaweicloud_lts_group.test[count.index].id
  stream_name = "%[1]s_${count.index}"
}
`, rName)
}

func testAccTaurusDBHtapStarrocksLtsConfig_basic(rName string, logType string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_taurusdb_htap_starrocks_lts_config" "test" {
  instance_id   = "%[2]s"
  log_type      = "%[3]s"
  lts_group_id  = huaweicloud_lts_group.test[0].id
  lts_stream_id = huaweicloud_lts_stream.test[0].id
}`, testAccTaurusDBHtapStarrocksLtsConfig_base(rName), acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID, logType)
}

func testAccTaurusDBHtapStarrocksLtsConfig_update(rName string, logType string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_taurusdb_htap_starrocks_lts_config" "test" {
  instance_id   = "%[2]s"
  log_type      = "%[3]s"
  lts_group_id  = huaweicloud_lts_group.test[1].id
  lts_stream_id = huaweicloud_lts_stream.test[1].id
}`, testAccTaurusDBHtapStarrocksLtsConfig_base(rName), acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID, logType)
}

func testAccTaurusDBHtapStarrocksLtsConfigImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("the resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["instance_id"] == "" {
			return "", fmt.Errorf("attribute (instance_id) of resource(%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["log_type"] == "" {
			return "", fmt.Errorf("attribute (log_type) of Resource (%s) not found: %s", name, rs)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["log_type"]), nil
	}
}
