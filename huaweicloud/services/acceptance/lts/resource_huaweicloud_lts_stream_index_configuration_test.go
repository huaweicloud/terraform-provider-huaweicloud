package lts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/lts"
)

func getStreamIndexConfiguration(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("lts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS client: %s", err)
	}

	return lts.GetStreamIndexConfiguration(client, state.Primary.Attributes["group_id"], state.Primary.Attributes["stream_id"])
}

func TestAccStreamIndexConfiguration_basic(t *testing.T) {
	var (
		name               = acceptance.RandomAccResourceName()
		indexConfiguration interface{}
		deafultConfig      = "huaweicloud_lts_stream_index_configuration.default"
		rName              = "huaweicloud_lts_stream_index_configuration.test"
		rcDeafultConfig    = acceptance.InitResourceCheck(deafultConfig, &indexConfiguration, getStreamIndexConfiguration)
		rc                 = acceptance.InitResourceCheck(rName, &indexConfiguration, getStreamIndexConfiguration)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcDeafultConfig.CheckResourceDestroy(),
			rc.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccStreamIndexConfiguration_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rcDeafultConfig.CheckResourceExists(),
					resource.TestCheckResourceAttr(deafultConfig, "full_text_index.0.tokenizer", ""),
					resource.TestCheckResourceAttr(deafultConfig, "full_text_index.0.enable", "true"),
					resource.TestCheckResourceAttr(deafultConfig, "full_text_index.0.case_sensitive", "false"),
					resource.TestCheckResourceAttr(deafultConfig, "full_text_index.0.include_chinese", "true"),
					resource.TestCheckResourceAttr(deafultConfig, "full_text_index.0.ascii.#", "0"),
					resource.TestCheckResourceAttr(deafultConfig, "fields.#", "0"),
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "stream_id", "huaweicloud_lts_stream.test.1", "id"),
					resource.TestCheckResourceAttr(rName, "full_text_index.0.tokenizer", ""),
					resource.TestCheckResourceAttr(rName, "full_text_index.0.enable", "true"),
					resource.TestCheckResourceAttr(rName, "full_text_index.0.case_sensitive", "false"),
					resource.TestCheckResourceAttr(rName, "full_text_index.0.include_chinese", "true"),
					resource.TestCheckResourceAttr(rName, "fields.#", "2"),
					resource.TestCheckResourceAttr(rName, "fields.0.field_name", "terraform"),
					resource.TestCheckResourceAttr(rName, "fields.0.field_type", "json"),
					resource.TestCheckResourceAttr(rName, "fields.0.quick_analysis", "false"),
					resource.TestCheckResourceAttr(rName, "fields.0.case_sensitive", "false"),
					resource.TestCheckResourceAttr(rName, "fields.0.include_chinese", "false"),
					resource.TestCheckResourceAttr(rName, "fields.0.tokenizer", ""),
					resource.TestCheckResourceAttr(rName, "fields.0.ascii.#", "1"),
					resource.TestCheckResourceAttr(rName, "fields.0.ascii.0", "14"),
					resource.TestCheckResourceAttr(rName, "fields.0.lts_sub_fields_info_list.0.field_name", "resource_id"),
					resource.TestCheckResourceAttr(rName, "fields.0.lts_sub_fields_info_list.0.field_type", "string"),
					resource.TestCheckResourceAttr(rName, "fields.0.lts_sub_fields_info_list.0.quick_analysis", "false"),
					resource.TestCheckResourceAttr(rName, "fields.1.field_type", "string"),
					resource.TestCheckResourceAttr(rName, "fields.1.case_sensitive", "true"),
					resource.TestCheckResourceAttr(rName, "fields.1.include_chinese", "true"),
					resource.TestCheckResourceAttr(rName, "fields.1.tokenizer", "?|"),
					resource.TestCheckResourceAttr(rName, "fields.1.quick_analysis", "true"),
					resource.TestCheckResourceAttr(rName, "fields.1.ascii.#", "2"),
				),
			},
			{
				Config: testAccStreamIndexConfiguration_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rcDeafultConfig.CheckResourceExists(),
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "full_text_index.0.tokenizer", ",|"),
					resource.TestCheckResourceAttr(rName, "full_text_index.0.enable", "true"),
					resource.TestCheckResourceAttr(rName, "full_text_index.0.case_sensitive", "true"),
					resource.TestCheckResourceAttr(rName, "full_text_index.0.include_chinese", "true"),
					resource.TestCheckResourceAttr(rName, "full_text_index.0.ascii.#", "2"),
					resource.TestCheckResourceAttr(rName, "fields.0.field_name", "terraform_update"),
					resource.TestCheckResourceAttr(rName, "fields.0.quick_analysis", "false"),
					resource.TestCheckResourceAttr(rName, "fields.0.lts_sub_fields_info_list.0.field_name", "count"),
					resource.TestCheckResourceAttr(rName, "fields.0.lts_sub_fields_info_list.0.field_type", "long"),
					resource.TestCheckResourceAttr(rName, "fields.0.lts_sub_fields_info_list.0.quick_analysis", "true"),
					resource.TestCheckResourceAttr(rName, "fields.1.tokenizer", ""),
					resource.TestCheckResourceAttr(rName, "fields.1.field_type", "long"),
					resource.TestCheckResourceAttr(rName, "fields.1.case_sensitive", "false"),
					resource.TestCheckResourceAttr(rName, "fields.1.include_chinese", "false"),
					resource.TestCheckResourceAttr(rName, "fields.1.tokenizer", ""),
					resource.TestCheckResourceAttr(rName, "fields.1.quick_analysis", "false"),
					resource.TestCheckResourceAttr(rName, "fields.1.ascii.#", "0")),
			},
			{
				Config: testAccStreamIndexConfiguration_basic_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "full_text_index.0.tokenizer", ""),
					resource.TestCheckResourceAttr(rName, "full_text_index.0.enable", "false"),
					resource.TestCheckResourceAttr(rName, "full_text_index.0.case_sensitive", "false"),
					resource.TestCheckResourceAttr(rName, "full_text_index.0.include_chinese", "false"),
					resource.TestCheckResourceAttr(rName, "full_text_index.0.ascii.#", "0"),
					resource.TestCheckResourceAttr(rName, "fields.#", "0"),
				),
			},
			{
				ResourceName:      deafultConfig,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccStreamIndexConfiguration_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test" {
  count = 2

  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[1]s_${count.index}"
}
`, name)
}

func testAccStreamIndexConfiguration_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_stream_index_configuration" "default" {
  group_id  = huaweicloud_lts_group.test.id
  stream_id = huaweicloud_lts_stream.test[0].id

  full_text_index {}
}

resource "huaweicloud_lts_stream_index_configuration" "test" {
  group_id  = huaweicloud_lts_group.test.id
  stream_id = huaweicloud_lts_stream.test[1].id

  fields {
    field_name = "terraform"
    field_type = "json"
    ascii      = ["14"]

    lts_sub_fields_info_list {
      field_name = "resource_id"
      field_type = "string"
    }
  }
  fields {
    field_name      = "field"
    field_type      = "string"
    case_sensitive  = true
    include_chinese = true
    tokenizer       = "?|"
    quick_analysis  = true
    ascii           = ["4", "3"]
  }
}
`, testAccStreamIndexConfiguration_base(name))
}

func testAccStreamIndexConfiguration_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_stream_index_configuration" "default" {
  group_id  = huaweicloud_lts_group.test.id
  stream_id = huaweicloud_lts_stream.test[0].id

  full_text_index {}
}

resource "huaweicloud_lts_stream_index_configuration" "test" {
  group_id  = huaweicloud_lts_group.test.id
  stream_id = huaweicloud_lts_stream.test[1].id

  full_text_index {
    tokenizer      = ",|"
    case_sensitive = true
    ascii          = ["2", "1"]
  }

  fields {
    field_name = "terraform_update"
    field_type = "json"

    lts_sub_fields_info_list {
      field_name     = "count"
      field_type     = "long"
      quick_analysis = true
    }
  }
  fields {
    field_name = "field"
    field_type = "long"
  }
}
`, testAccStreamIndexConfiguration_base(name))
}

func testAccStreamIndexConfiguration_basic_step3(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_stream_index_configuration" "default" {
  group_id  = huaweicloud_lts_group.test.id
  stream_id = huaweicloud_lts_stream.test[0].id

  full_text_index {}
}

resource "huaweicloud_lts_stream_index_configuration" "test" {
  group_id  = huaweicloud_lts_group.test.id
  stream_id = huaweicloud_lts_stream.test[1].id

  full_text_index {
    enable          = false
    include_chinese = false
  }
}
`, testAccStreamIndexConfiguration_base(name))
}
