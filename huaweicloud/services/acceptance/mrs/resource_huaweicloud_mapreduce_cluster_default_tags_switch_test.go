package mrs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/mrs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getClusterDefaultTagsSwitchResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("mrs", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating MRS client: %s", err)
	}

	respbody, err := mrs.GetClusterDefaultTagsSwitchStatus(client, state.Primary.ID)
	if err != nil {
		return nil, err
	}

	if !utils.PathSearch("default_tags_enable", respbody, false).(bool) {
		return nil, golangsdk.ErrDefault404{}
	}

	return respbody, nil
}

func TestAccClusterDefaultTagsSwitch_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_mapreduce_cluster_default_tags_switch.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getClusterDefaultTagsSwitchResourceFunc)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMrsClusterID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterDefaultTagsSwitch_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_MRS_CLUSTER_ID),
					resource.TestCheckResourceAttr(rName, "action", "create"),
					resource.TestCheckResourceAttr(rName, "default_tags_enable", "true"),
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

func testAccClusterDefaultTagsSwitch_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_mapreduce_cluster_default_tags_switch" "test" {
  cluster_id = "%s"
  action     = "create"
}
`, acceptance.HW_MRS_CLUSTER_ID)
}

func getClusterDeleteDefaultTagsSwitchResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("mrs", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating MRS client: %s", err)
	}

	return mrs.GetClusterDefaultTagsSwitchStatus(client, state.Primary.ID)
}

func TestAccClusterDefaultTagsSwitch_delete(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_mapreduce_cluster_default_tags_switch.delete"
		rc    = acceptance.InitResourceCheck(rName, &obj, getClusterDeleteDefaultTagsSwitchResourceFunc)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMrsClusterID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// When `action` parameter is `delete`, the resource does not need to be deleted.
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterDefaultTagsSwitch_delete(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_MRS_CLUSTER_ID),
					resource.TestCheckResourceAttr(rName, "action", "delete"),
					resource.TestCheckResourceAttr(rName, "default_tags_enable", "false"),
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

func testAccClusterDefaultTagsSwitch_delete() string {
	return fmt.Sprintf(`
resource "huaweicloud_mapreduce_cluster_default_tags_switch" "create" {
  cluster_id = "%[1]s"
  action     = "create"

  lifecycle {
    ignore_changes = [action]
  }
}

resource "huaweicloud_mapreduce_cluster_default_tags_switch" "delete" {
  cluster_id = "%[1]s"
  action     = "delete"

  depends_on = [huaweicloud_mapreduce_cluster_default_tags_switch.create]
}
`, acceptance.HW_MRS_CLUSTER_ID)
}
