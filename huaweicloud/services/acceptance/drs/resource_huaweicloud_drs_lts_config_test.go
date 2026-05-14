package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/drs"
)

func getLtsConfigFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("drs", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DRS client: %s", err)
	}

	return drs.GetLtsConfig(client, state.Primary.ID)
}

func TestAccResourceLtsConfig_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_drs_lts_config.test"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getLtsConfigFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLtsConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "job_id", acceptance.HW_DRS_JOB_ID),
					resource.TestCheckResourceAttrPair(resourceName, "log_group_id", "huaweicloud_lts_group.test_1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "log_stream_id", "huaweicloud_lts_stream.test_1", "id"),
				),
			},
			{
				Config: testAccLtsConfig_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "job_id", acceptance.HW_DRS_JOB_ID),
					resource.TestCheckResourceAttrPair(resourceName, "log_group_id", "huaweicloud_lts_group.test_2", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "log_stream_id", "huaweicloud_lts_stream.test_2", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccLtsConfig_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test_1" {
  group_name  = "%[1]s_1"
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test_1" {
  group_id    = huaweicloud_lts_group.test_1.id
  stream_name = "%[1]s_1"
  is_favorite = true
}

resource "huaweicloud_lts_group" "test_2" {
  group_name  = "%[1]s_2"
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test_2" {
  group_id    = huaweicloud_lts_group.test_2.id
  stream_name = "%[1]s_2"
  is_favorite = true
}
`, name)
}

func testAccLtsConfig_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_drs_lts_config" "test" { 
  job_id        = "%[2]s"
  log_group_id  = huaweicloud_lts_group.test_1.id
  log_stream_id = huaweicloud_lts_stream.test_1.id
}
`, testAccLtsConfig_base(name), acceptance.HW_DRS_JOB_ID)
}

func testAccLtsConfig_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_drs_lts_config" "test" { 
  job_id        = "%[2]s"
  log_group_id  = huaweicloud_lts_group.test_2.id
  log_stream_id = huaweicloud_lts_stream.test_2.id
}
`, testAccLtsConfig_base(name), acceptance.HW_DRS_JOB_ID)
}
