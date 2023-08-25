package lts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/lts/huawei/loggroups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getLtsGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.LtsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS client: %s", err)
	}

	resourceID := state.Primary.ID
	groups, err := loggroups.List(client).Extract()
	if err != nil {
		return nil, err
	}

	for _, item := range groups.LogGroups {
		if item.ID == resourceID {
			return &item, nil
		}
	}

	return nil, fmt.Errorf("the log group %s does not exist", resourceID)
}

func TestAccLtsGroup_basic(t *testing.T) {
	var group loggroups.LogGroup
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_lts_group.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&group,
		getLtsGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLtsGroup_basic(rName, 30),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "group_name", rName),
					resource.TestCheckResourceAttr(resourceName, "ttl_in_days", "30"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLtsGroup_basic(rName, 7),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_name", rName),
					resource.TestCheckResourceAttr(resourceName, "ttl_in_days", "7"),
				),
			},
			{
				Config: testAccLtsGroup_tags(rName, 60),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_name", rName),
					resource.TestCheckResourceAttr(resourceName, "ttl_in_days", "60"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
				),
			},
		},
	})
}

func testAccLtsGroup_basic(name string, ttl int) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  group_name  = "%s"
  ttl_in_days = %d
}
`, name, ttl)
}

func testAccLtsGroup_tags(name string, ttl int) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  group_name  = "%s"
  ttl_in_days = %d

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, ttl)
}
