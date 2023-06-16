package sdrs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/sdrs/v1/replications"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getReplicationPairResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.SdrsV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating SDRS Client: %s", err)
	}
	return replications.Get(client, state.Primary.ID).Extract()
}

func TestAccReplicationPair_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_sdrs_replication_pair.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getReplicationPairResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testReplicationPair_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "delete_target_volume", "true"),
					resource.TestCheckResourceAttrPair(rName, "group_id", "huaweicloud_sdrs_protection_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "volume_id", "huaweicloud_evs_volume.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "replication_model"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "target_volume_id"),
				),
			},
			{
				Config: testReplicationPair_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"delete_target_volume",
				},
			},
		},
	})
}

func testReplicationPair_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}
data "huaweicloud_sdrs_domain" "test" {}

resource "huaweicloud_sdrs_protection_group" "test" {
  name                     = "%[2]s"
  source_availability_zone = data.huaweicloud_availability_zones.test.names[0]
  target_availability_zone = data.huaweicloud_availability_zones.test.names[1]
  domain_id                = data.huaweicloud_sdrs_domain.test.id
  source_vpc_id            = huaweicloud_vpc.test.id
  description              = "test description"
}

resource "huaweicloud_evs_volume" "test" {
  name              = "%[2]s"
  description       = "test volume for sdrs replication pair"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  volume_type       = "SSD"
  size              = 100
}
`, common.TestVpc(name), name)
}

func testReplicationPair_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_sdrs_replication_pair" "test" {
  name                 = "%[2]s"
  group_id             = huaweicloud_sdrs_protection_group.test.id
  volume_id            = huaweicloud_evs_volume.test.id
  description          = "test description"
  delete_target_volume = true
}
`, testReplicationPair_base(name), name)
}

func testReplicationPair_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_sdrs_replication_pair" "test" {
  name                 = "%[2]s_update"
  group_id             = huaweicloud_sdrs_protection_group.test.id
  volume_id            = huaweicloud_evs_volume.test.id
  description          = "test description"
  delete_target_volume = true
}
`, testReplicationPair_base(name), name)
}
