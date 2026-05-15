package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/geminidb"
)

func getResourceMemoryRuleFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("geminidb", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating GeminiDB client: %s", err)
	}

	return geminidb.GetMemoryRuleInfo(client, state.Primary.Attributes["dbcache_mapping_id"], state.Primary.ID)
}

func TestAccResourceMemoryRule_basic(t *testing.T) {
	var (
		rName  = "huaweicloud_geminidb_memory_rule.test"
		name   = acceptance.RandomAccResourceName()
		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceMemoryRuleFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccMemoryRule_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "dbcache_mapping_id",
						"huaweicloud_geminidb_memory_mapping.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "source_db_schema", "test"),
					resource.TestCheckResourceAttr(rName, "source_db_table", "user"),
					resource.TestCheckResourceAttr(rName, "storage_type", "hash"),
					resource.TestCheckResourceAttr(rName, "target_database", "0"),
					resource.TestCheckResourceAttr(rName, "key_prefix", "prefix"),
					resource.TestCheckResourceAttr(rName, "key_columns.#", "2"),
					resource.TestCheckResourceAttr(rName, "value_columns.#", "2"),
					resource.TestCheckResourceAttr(rName, "key_separator", "."),
					resource.TestCheckResourceAttr(rName, "value_separator", ","),
					resource.TestCheckResourceAttr(rName, "ttl", "360000000"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				Config: testAccMemoryRule_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "value_columns.#", "4"),
					resource.TestCheckResourceAttr(rName, "value_separator", ";"),
					resource.TestCheckResourceAttr(rName, "ttl", "180000000"),
				),
			},

			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccMemoryRuleImportStateFunc(rName),
			},
		},
	})
}

func testAccMemoryRule_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_geminidb_memory_rule" "test" {
  dbcache_mapping_id = huaweicloud_geminidb_memory_mapping.test.id
  name               = "%[2]s"
  source_db_schema   = "test"
  source_db_table    = "user"
  storage_type       = "hash"
  target_database    = "0"
  key_prefix         = "prefix"
  key_columns        = ["uno","uname"]
  value_columns      = ["uno","uname"]
  key_separator      = "."
  value_separator    = ","
  ttl                = "360000000"
}
`, testAccMemoryMapping_basic(name), name)
}

func testAccMemoryRule_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_geminidb_memory_rule" "test" {
  dbcache_mapping_id = huaweicloud_geminidb_memory_mapping.test.id
  name               = "%[2]s"
  source_db_schema   = "test"
  source_db_table    = "user"
  storage_type       = "hash"
  target_database    = "0"
  key_prefix         = "prefix"
  key_columns        = ["uno","uname"]
  value_columns      = ["uno","uname","age","grader"]
  key_separator      = "."
  value_separator    = ";"
  ttl                = "180000000"
}
`, testAccMemoryMapping_basic(name), name)
}

func testAccMemoryRuleImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var mappingId, ruleId string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		mappingId = rs.Primary.Attributes["dbcache_mapping_id"]
		ruleId = rs.Primary.ID

		if mappingId == "" || ruleId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<dbcache_mapping_id>/<id>', but got '%s/%s'",
				mappingId, ruleId)
		}

		return fmt.Sprintf("%s/%s", mappingId, ruleId), nil
	}
}
