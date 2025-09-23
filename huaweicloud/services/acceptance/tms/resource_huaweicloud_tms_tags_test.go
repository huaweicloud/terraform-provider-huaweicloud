package tms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func TestAccTmsTag_basic(t *testing.T) {
	resourceName := "huaweicloud_tms_tags.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTmsTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTmsTags_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
					testAccCheckTmsTagsExists([]map[string]string{
						{
							"key":   "key_1",
							"value": "value_1",
						},
						{
							"key":   "key_2",
							"value": "value_2",
						},
						{
							"key":   "key_2",
							"value": "value_22",
						},
					}),
				),
			},
			{
				Config: testAccTmsTags_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tags.#", "4"),
					testAccCheckTmsTagsExists([]map[string]string{
						{
							"key":   "key_2",
							"value": "value_2",
						},
						{
							"key":   "key_2",
							"value": "value_222",
						},
						{
							"key":   "key_3",
							"value": "value_3",
						},
						{
							"key":   "key_4",
							"value": "value_4",
						},
					}),
				),
			},
		},
	})
}

func testAccCheckTmsTagsExists(rawTags []map[string]string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		region := acceptance.HW_REGION_NAME
		var (
			product = "tms"
		)
		client, err := cfg.NewServiceClient(product, region)
		if err != nil {
			return fmt.Errorf("error creating TMS client: %s", err)
		}
		tags, err := getTmsTags(client)
		if err != nil {
			return err
		}
		for _, rawTag := range rawTags {
			exists := false
			for _, tag := range tags {
				if rawTag["key"] == tag["key"] && rawTag["value"] == tag["value"] {
					exists = true
					break
				}
			}
			if !exists {
				return fmt.Errorf("the tag(key: %s, value: %s) does not exist", rawTag["key"], rawTag["value"])
			}
		}
		return nil
	}
}

func testAccCheckTmsTagDestroy(_ *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	region := acceptance.HW_REGION_NAME
	var (
		product = "tms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return fmt.Errorf("error creating TMS client: %s", err)
	}
	tags, err := getTmsTags(client)
	if err != nil {
		return err
	}

	rawTags := []map[string]string{
		{
			"key":   "key_2",
			"value": "value_2",
		},
		{
			"key":   "key_2",
			"value": "value_222",
		},
		{
			"key":   "key_3",
			"value": "value_3",
		},
		{
			"key":   "key_4",
			"value": "value_4",
		},
	}
	for _, rawTag := range rawTags {
		exists := false
		for _, tag := range tags {
			if rawTag["key"] == tag["key"] && rawTag["value"] == tag["value"] {
				exists = true
				break
			}
		}
		if exists {
			return fmt.Errorf("the tag(key: %s, value: %s) still exist", rawTag["key"], rawTag["value"])
		}
	}

	return nil
}

func getTmsTags(client *golangsdk.ServiceClient) ([]map[string]string, error) {
	var (
		httpUrl = "v1.0/predefine_tags"
	)

	getBasePath := client.Endpoint + httpUrl
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	res := make([]map[string]string, 0)
	marker := ""
	for {
		getPath := getBasePath + buildGetPath(marker)
		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving TMS tags: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return nil, err
		}
		curJson := utils.PathSearch("tags", getRespBody, make([]interface{}, 0))
		curArray := curJson.([]interface{})
		for _, v := range curArray {
			res = append(res, map[string]string{
				"key":   utils.PathSearch("key", v, "").(string),
				"value": utils.PathSearch("value", v, "").(string),
			})
		}
		marker = utils.PathSearch("marker", getRespBody, "").(string)
		if marker == "" {
			break
		}
	}

	return res, nil
}

func buildGetPath(marker string) string {
	res := "?limit=100"
	if marker != "" {
		res = fmt.Sprintf("%s&marker=%s", res, marker)
	}
	return res
}

func testAccTmsTags_basic() string {
	return `
resource "huaweicloud_tms_tags" "test" {
  tags {
    key   = "key_1"
    value = "value_1"
  }
  tags {
    key   = "key_2"
    value = "value_2"
  }
  tags {
    key   = "key_2"
    value = "value_22"
  }
}
`
}

func testAccTmsTags_update() string {
	return `
resource "huaweicloud_tms_tags" "test" {
  tags {
    key   = "key_2"
    value = "value_2"
  }
  tags {
    key   = "key_2"
    value = "value_222"
  }
  tags {
    key   = "key_3"
    value = "value_3"
  }
  tags {
    key   = "key_4"
    value = "value_4"
  }
}
`
}
