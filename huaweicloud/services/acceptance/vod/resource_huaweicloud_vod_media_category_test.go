package vod

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	vod "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vod/v1/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getResourceCategory(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcVodV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating VOD client: %s", err)
	}

	id, err := strconv.ParseInt(state.Primary.ID, 10, 32)
	if err != nil {
		return nil, err
	}

	resp, err := client.ListAssetCategory(&vod.ListAssetCategoryRequest{Id: int32(id)})
	if err != nil {
		return nil, fmt.Errorf("error retrieving VOD media category: %d", id)
	}

	categoryList := *resp.Body
	if len(categoryList) == 0 {
		return nil, fmt.Errorf("unable to retrieve VOD media category: %d", id)
	}
	category := categoryList[0]

	return category, nil
}

func TestAccMediaCategory_basic(t *testing.T) {
	var category vod.QueryCategoryRsp
	rName := acceptance.RandomAccResourceName()
	updateName := rName + "-update"
	resourceName := "huaweicloud_vod_media_category.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&category,
		getResourceCategory,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccMediaCategory_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				Config: testAccMediaCategory_basic(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"parent_id",
				},
			},
		},
	})
}

func testAccMediaCategory_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vod_media_category" "test" {
  name = "%s"
}
`, rName)
}
