package lts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccStructCustomConfig_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_lts_structuring_custom_configuration.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getStructConfigResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testStructCustomConfig_json(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "log_group_id",
						"huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_stream_id",
						"huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttr(rName, "content", "{'code':38,'user':{'name':'testdemo'}}"),
					resource.TestCheckResourceAttr(rName, "layers", "3"),
					resource.TestCheckResourceAttr(rName, "demo_fields.#", "2"),
					resource.TestCheckResourceAttr(rName, "demo_fields.0.is_analysis", "true"),
					resource.TestCheckResourceAttr(rName, "demo_fields.0.field_name", "code"),
					resource.TestCheckResourceAttr(rName, "demo_fields.0.content", "38"),
					resource.TestCheckResourceAttr(rName, "demo_fields.0.type", "long"),
					resource.TestCheckResourceAttr(rName, "demo_fields.1.is_analysis", "true"),
					resource.TestCheckResourceAttr(rName, "demo_fields.1.field_name", "user.name"),
					resource.TestCheckResourceAttr(rName, "demo_fields.1.content", "test_demo"),
					resource.TestCheckResourceAttr(rName, "demo_fields.1.type", "string"),

					resource.TestCheckResourceAttr(rName, "tag_fields.0.is_analysis", "true"),
					resource.TestCheckResourceAttr(rName, "tag_fields.0.field_name", "hostIP"),
					resource.TestCheckResourceAttr(rName, "tag_fields.0.content", "192.168.2.134"),
					resource.TestCheckResourceAttr(rName, "tag_fields.0.type", "string"),
				),
			},
			{
				Config: testStructCustomConfig_json_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "content", "{'error_code':400,'user':{'id':'000'}}"),
					resource.TestCheckResourceAttr(rName, "layers", "3"),
					resource.TestCheckResourceAttr(rName, "demo_fields.0.is_analysis", "false"),
					resource.TestCheckResourceAttr(rName, "demo_fields.0.field_name", "error_code"),
					resource.TestCheckResourceAttr(rName, "demo_fields.0.content", "400"),
					resource.TestCheckResourceAttr(rName, "demo_fields.0.type", "long"),
					resource.TestCheckResourceAttr(rName, "demo_fields.1.is_analysis", "true"),
					resource.TestCheckResourceAttr(rName, "demo_fields.1.field_name", "user.id"),
					resource.TestCheckResourceAttr(rName, "demo_fields.1.content", "000"),
					resource.TestCheckResourceAttr(rName, "demo_fields.1.type", "string"),

					resource.TestCheckResourceAttr(rName, "tag_fields.0.is_analysis", "true"),
					resource.TestCheckResourceAttr(rName, "tag_fields.0.field_name", "hostName"),
					resource.TestCheckResourceAttr(rName, "tag_fields.0.content", "demoName"),
					resource.TestCheckResourceAttr(rName, "tag_fields.0.type", "string"),
				),
			},
			{
				Config: testStructCustomConfig_split_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "content", "2023-09-09/18:50:51 Error"),
					resource.TestCheckResourceAttr(rName, "tokenizer", " "),
					resource.TestCheckResourceAttr(rName, "demo_fields.0.is_analysis", "true"),
					resource.TestCheckResourceAttr(rName, "demo_fields.0.field_name", "b1"),
					resource.TestCheckResourceAttr(rName, "demo_fields.0.content", "2023-09-09/18:50:51"),
					resource.TestCheckResourceAttr(rName, "demo_fields.0.type", "string"),
					resource.TestCheckResourceAttr(rName, "demo_fields.1.is_analysis", "true"),
					resource.TestCheckResourceAttr(rName, "demo_fields.1.field_name", "b2"),
					resource.TestCheckResourceAttr(rName, "demo_fields.1.content", "Error"),
					resource.TestCheckResourceAttr(rName, "demo_fields.1.type", "string"),

					resource.TestCheckResourceAttr(rName, "tag_fields.#", "0"),
				),
			},
			{
				Config: testStructCustomConfig_nginx_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "content",
						"39.149.31.187 - - [12/Mar/2020:12:24:02 +0800] \"GET / HTTP/1.1\" 304 "),
					resource.TestCheckResourceAttr(rName, "log_format",
						"log_format  main   '$remote_addr - $remote_user [$time_local] \"$request\" '\n'$status ';"),
					resource.TestCheckResourceAttr(rName, "demo_fields.0.is_analysis", "true"),
					resource.TestCheckResourceAttr(rName, "demo_fields.0.field_name", "remote_addr"),
					resource.TestCheckResourceAttr(rName, "demo_fields.0.content", "39.149.31.187"),
					resource.TestCheckResourceAttr(rName, "demo_fields.0.type", "string"),
					resource.TestCheckResourceAttr(rName, "demo_fields.1.is_analysis", "true"),
					resource.TestCheckResourceAttr(rName, "demo_fields.1.field_name", "remote_user"),
					resource.TestCheckResourceAttr(rName, "demo_fields.1.content", "-"),
					resource.TestCheckResourceAttr(rName, "demo_fields.1.type", "string"),
					resource.TestCheckResourceAttr(rName, "demo_fields.2.is_analysis", "true"),
					resource.TestCheckResourceAttr(rName, "demo_fields.2.field_name", "request_method"),
					resource.TestCheckResourceAttr(rName, "demo_fields.2.content", "GET"),
					resource.TestCheckResourceAttr(rName, "demo_fields.2.type", "string"),
				),
			},
			{
				Config: testStructCustomConfig_custom_regex_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "content",
						"2023-09-09/18:15:41 this log is Error NO 6323"),
					resource.TestCheckResourceAttr(rName, "regex_rules",
						"^(?<a01>[^ ]+)(?:[^ ]* ){1}(?<a02>\\w+)(?:[^ ]* ){1}(?<a03>\\w+)(?:[^ ]* )"),
					resource.TestCheckResourceAttr(rName, "demo_fields.0.is_analysis", "true"),
					resource.TestCheckResourceAttr(rName, "demo_fields.0.field_name", "a01"),
					resource.TestCheckResourceAttr(rName, "demo_fields.0.type", "string"),
					resource.TestCheckResourceAttr(rName, "demo_fields.1.is_analysis", "true"),
					resource.TestCheckResourceAttr(rName, "demo_fields.1.field_name", "a02"),
					resource.TestCheckResourceAttr(rName, "demo_fields.1.type", "string"),
					resource.TestCheckResourceAttr(rName, "demo_fields.2.is_analysis", "true"),
					resource.TestCheckResourceAttr(rName, "demo_fields.2.field_name", "a03"),
					resource.TestCheckResourceAttr(rName, "demo_fields.2.type", "string"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testLtsStructConfigImportState(rName),
				ImportStateVerifyIgnore: []string{
					"demo_fields",
					"regex_rules",
					"layers",
					"tokenizer",
					"log_format",
					"tag_fields",
				},
			},
		},
	})
}

func testStructCustomConfig_json(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_lts_structuring_custom_configuration" "test" {
  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id
  content       = "{'code':38,'user':{'name':'testdemo'}}"
  layers        = 3
  
  demo_fields {
    is_analysis = true
    field_name  = "code"
    content     = "38"
    type        = "long"
  }

  demo_fields {
    is_analysis = true
    field_name  = "user.name"
    content     = "test_demo"
    type        = "string"
  }

  tag_fields {
    is_analysis = true
    field_name  = "hostIP"
    content     = "192.168.2.134"
    type        = "string"
  }
}
`, testAccLtsStream_basic(name))
}

func testStructCustomConfig_json_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_lts_structuring_custom_configuration" "test" {
  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id
  content       = "{'error_code':400,'user':{'id':'000'}}"
  layers        = 3
  
  demo_fields {
    field_name  = "error_code"
    content     = "400"
    type        = "long"
  }

  demo_fields {
    is_analysis = true
    field_name  = "user.id"
    content     = "000"
    type        = "string"
  }

  tag_fields {
    is_analysis = true
    field_name  = "hostName"
    content     = "demoName"
    type        = "string"
  }
}
`, testAccLtsStream_basic(name))
}

func testStructCustomConfig_split_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_lts_structuring_custom_configuration" "test" {
  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id
  content       = "2023-09-09/18:50:51 Error"
  tokenizer     = " "
  
  demo_fields {
    is_analysis = true
    field_name  = "b1"
    content     = "2023-09-09/18:50:51"
    type        = "string"
  }

  demo_fields {
    is_analysis = true
    field_name  = "b2"
    content     = "Error"
    type        = "string"
  }
}
`, testAccLtsStream_basic(name))
}

func testStructCustomConfig_nginx_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_lts_structuring_custom_configuration" "test" {
  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id
  content       = "39.149.31.187 - - [12/Mar/2020:12:24:02 +0800] \"GET / HTTP/1.1\" 304 "
  log_format    = "log_format  main   '$remote_addr - $remote_user [$time_local] \"$request\" '\n'$status ';"

  demo_fields {
    is_analysis = true
    field_name  = "remote_addr"
    content     = "39.149.31.187"
    type        = "string"
  }

  demo_fields {
    is_analysis = true
    field_name  = "remote_user"
    content     = "-"
    type        = "string"
  }

  demo_fields {
    is_analysis = true
    field_name  = "request_method"
    content     = "GET"
    type        = "string"
  }
}
`, testAccLtsStream_basic(name))
}

func testStructCustomConfig_custom_regex_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_lts_structuring_custom_configuration" "test" {
  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id
  content       = "2023-09-09/18:15:41 this log is Error NO 6323"
  regex_rules   = "^(?<a01>[^ ]+)(?:[^ ]* ){1}(?<a02>\\w+)(?:[^ ]* ){1}(?<a03>\\w+)(?:[^ ]* )"

  demo_fields {
    is_analysis = true
    field_name  = "a01"
    type        = "string"
  }

  demo_fields {
    is_analysis = true
    field_name  = "a02"
    type        = "string"
  }

  demo_fields {
    is_analysis = true
    field_name  = "a03"
    type        = "string"
  }
}
`, testAccLtsStream_basic(name))
}
