package tms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/tms/v1/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTmsTag_basic(t *testing.T) {
	resourceName := "huaweicloud_tms_tags.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTmsTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTmsTag_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTmsTagExists("foo", "bar"),
					testAccCheckTmsTagExists("k", "v"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
				),
			},
		},
	})
}

func testAccCheckTmsTagDestroy(s *terraform.State) error {
	conf := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := conf.HcTmsV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating TMS client: %s", err)
	}

	tags := map[string]string{"foo": "bar", "k": "v"}
	for key, value := range tags {
		request := &model.ListPredefineTagsRequest{
			Key:   &key,
			Value: &value,
		}

		response, err := client.ListPredefineTags(request)
		if err != nil {
			return err
		}
		tagsFromResponse := *response.Tags
		if len(tagsFromResponse) != 0 {
			return fmt.Errorf("huaweicloud_tms_tags %s/%s still exists", key, value)
		}
	}

	return nil
}

func testAccCheckTmsTagExists(key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conf := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := conf.HcTmsV1Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating TMS client: %s", err)
		}

		request := &model.ListPredefineTagsRequest{
			Key:   &key,
			Value: &value,
		}

		response, err := client.ListPredefineTags(request)
		if err != nil {
			return err
		}
		tags := *response.Tags
		if len(tags) != 0 {
			return nil
		}

		return fmt.Errorf("huaweicloud_tms_tags %s/%s does not exist", key, value)
	}
}

const testTmsTag_basic = `
resource "huaweicloud_tms_tags" "test" {
  tags {
    key   = "foo"
    value = "bar"
  }
  tags {
    key   = "k"
    value = "v"
  }
}
`
