package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
)

func getOrchestrationRuleFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("apig", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG client: %s", err)
	}

	return apig.GetOrchestrationRuleById(client, state.Primary.Attributes["instance_id"], state.Primary.ID)
}

func TestAccOrchestrationRule_basic(t *testing.T) {
	var (
		obj interface{}

		typeList   = "huaweicloud_apig_orchestration_rule.type_list"
		rcTypeList = acceptance.InitResourceCheck(typeList, &obj, getOrchestrationRuleFunc)

		typeRange   = "huaweicloud_apig_orchestration_rule.type_range"
		rcTypeRange = acceptance.InitResourceCheck(typeRange, &obj, getOrchestrationRuleFunc)

		typeHash   = "huaweicloud_apig_orchestration_rule.type_hash"
		rcTypeHash = acceptance.InitResourceCheck(typeHash, &obj, getOrchestrationRuleFunc)

		typeHashRange   = "huaweicloud_apig_orchestration_rule.type_hash_range"
		rcTypeHashRange = acceptance.InitResourceCheck(typeHashRange, &obj, getOrchestrationRuleFunc)

		typeNoneValue   = "huaweicloud_apig_orchestration_rule.type_none_value"
		rcTypeNoneValue = acceptance.InitResourceCheck(typeNoneValue, &obj, getOrchestrationRuleFunc)

		typeDefault   = "huaweicloud_apig_orchestration_rule.type_default"
		rcTypeDefault = acceptance.InitResourceCheck(typeDefault, &obj, getOrchestrationRuleFunc)

		typeHeadN   = "huaweicloud_apig_orchestration_rule.type_head_n"
		rcTypeHeadN = acceptance.InitResourceCheck(typeHeadN, &obj, getOrchestrationRuleFunc)

		typeTailN   = "huaweicloud_apig_orchestration_rule.type_tail_n"
		rcTypeTailN = acceptance.InitResourceCheck(typeTailN, &obj, getOrchestrationRuleFunc)

		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcTypeList.CheckResourceDestroy(),
			rcTypeRange.CheckResourceDestroy(),
			rcTypeHash.CheckResourceDestroy(),
			rcTypeHashRange.CheckResourceDestroy(),
			rcTypeNoneValue.CheckResourceDestroy(),
			rcTypeDefault.CheckResourceDestroy(),
			rcTypeHeadN.CheckResourceDestroy(),
			rcTypeTailN.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				// Check whether illegal type ​​can be intercepted normally (create phase).
				Config:      testAccPlugin_basic_invalidType(),
				ExpectError: regexp.MustCompile(`Invalid parameter value.*orchestration_strategy`),
			},
			{
				Config: testAccOrchestrationRule_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Orchestration rule: type list.
					rcTypeList.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeList, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeList, "name", name+"_type_list"),
					resource.TestCheckResourceAttr(typeList, "strategy", "list"),
					resource.TestCheckResourceAttr(typeList, "is_preprocessing", "false"),
					resource.TestCheckResourceAttr(typeList, "mapped_param",
						"{\"mapped_param_location\":\"header\",\"mapped_param_name\":\"listParam\",\"mapped_param_type\":\"string\"}"),
					resource.TestCheckResourceAttr(typeList, "map.#", "3"),
					resource.TestCheckResourceAttr(typeList, "map.0",
						"{\"map_param_list\":[\"ValueA\"],\"mapped_param_value\":\"ValueAA\"}"),
					resource.TestCheckResourceAttr(typeList, "map.1",
						"{\"map_param_list\":[\"ValueC\"],\"mapped_param_value\":\"ValueCC\"}"),
					resource.TestCheckResourceAttr(typeList, "map.2",
						"{\"map_param_list\":[\"ValueB\"],\"mapped_param_value\":\"ValueBB\"}"),
					resource.TestMatchResourceAttr(typeList, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Orchestration rule: type range.
					rcTypeRange.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeRange, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeRange, "name", name+"_type_range"),
					resource.TestCheckResourceAttr(typeRange, "strategy", "range"),
					resource.TestCheckResourceAttr(typeRange, "is_preprocessing", "false"),
					resource.TestCheckResourceAttr(typeRange, "mapped_param",
						"{\"mapped_param_location\":\"query\",\"mapped_param_name\":\"rangeParam\",\"mapped_param_type\":\"number\"}"),
					resource.TestCheckResourceAttr(typeRange, "map.#", "3"),
					resource.TestCheckResourceAttr(typeRange, "map.0",
						"{\"map_param_range\":{\"range_end\":\"1999\",\"range_start\":\"1000\"},\"mapped_param_value\":\"10001\"}"),
					resource.TestCheckResourceAttr(typeRange, "map.1",
						"{\"map_param_range\":{\"range_end\":\"3999\",\"range_start\":\"3000\"},\"mapped_param_value\":\"10003\"}"),
					resource.TestCheckResourceAttr(typeRange, "map.2",
						"{\"map_param_range\":{\"range_end\":\"2999\",\"range_start\":\"2000\"},\"mapped_param_value\":\"10002\"}"),
					resource.TestMatchResourceAttr(typeRange, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Orchestration rule: type hash.
					rcTypeHash.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeHash, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeHash, "name", name+"_type_hash"),
					resource.TestCheckResourceAttr(typeHash, "strategy", "hash"),
					resource.TestCheckResourceAttr(typeHash, "is_preprocessing", "false"),
					resource.TestCheckResourceAttr(typeHash, "mapped_param",
						"{\"mapped_param_location\":\"header\",\"mapped_param_name\":\"hashParam\",\"mapped_param_type\":\"string\"}"),
					resource.TestMatchResourceAttr(typeHash, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Orchestration rule: type hash range.
					rcTypeHashRange.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeHashRange, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeHashRange, "name", name+"_type_hash_range"),
					resource.TestCheckResourceAttr(typeHashRange, "strategy", "hash_range"),
					resource.TestCheckResourceAttr(typeHashRange, "is_preprocessing", "false"),
					resource.TestCheckResourceAttr(typeHashRange, "mapped_param",
						"{\"mapped_param_location\":\"header\",\"mapped_param_name\":\"hashRangeParam\",\"mapped_param_type\":\"number\"}"),
					resource.TestCheckResourceAttr(typeHashRange, "map.#", "3"),
					resource.TestCheckResourceAttr(typeHashRange, "map.0",
						"{\"map_param_range\":{\"range_end\":\"1999\",\"range_start\":\"1000\"},\"mapped_param_value\":\"10001\"}"),
					resource.TestCheckResourceAttr(typeRange, "map.1",
						"{\"map_param_range\":{\"range_end\":\"3999\",\"range_start\":\"3000\"},\"mapped_param_value\":\"10003\"}"),
					resource.TestCheckResourceAttr(typeRange, "map.2",
						"{\"map_param_range\":{\"range_end\":\"2999\",\"range_start\":\"2000\"},\"mapped_param_value\":\"10002\"}"),
					resource.TestMatchResourceAttr(typeRange, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Orchestration rule: type none value.
					rcTypeNoneValue.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeNoneValue, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeNoneValue, "name", name+"_type_none_value"),
					resource.TestCheckResourceAttr(typeNoneValue, "strategy", "none_value"),
					resource.TestCheckResourceAttr(typeNoneValue, "is_preprocessing", "false"),
					resource.TestCheckResourceAttr(typeNoneValue, "mapped_param",
						"{\"mapped_param_location\":\"query\",\"mapped_param_name\":\"noneValueParam\",\"mapped_param_type\":\"string\"}"),
					resource.TestCheckResourceAttr(typeNoneValue, "map.#", "1"),
					resource.TestCheckResourceAttr(typeNoneValue, "map.0", "{\"mapped_param_value\":\"NoneValueReturned\"}"),
					resource.TestMatchResourceAttr(typeNoneValue, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Orchestration rule: type default.
					rcTypeDefault.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeDefault, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeDefault, "name", name+"_type_default"),
					resource.TestCheckResourceAttr(typeDefault, "strategy", "default"),
					resource.TestCheckResourceAttr(typeDefault, "is_preprocessing", "false"),
					resource.TestCheckResourceAttr(typeDefault, "mapped_param",
						"{\"mapped_param_location\":\"query\",\"mapped_param_name\":\"defaultParam\",\"mapped_param_type\":\"string\"}"),
					resource.TestCheckResourceAttr(typeDefault, "map.#", "1"),
					resource.TestCheckResourceAttr(typeDefault, "map.0", "{\"mapped_param_value\":\"DefaultValueReturned\"}"),
					resource.TestMatchResourceAttr(typeDefault, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Orchestration rule: type head N.
					rcTypeHeadN.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeHeadN, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeHeadN, "name", name+"_type_head_n"),
					resource.TestCheckResourceAttr(typeHeadN, "strategy", "head_n"),
					resource.TestCheckResourceAttr(typeHeadN, "is_preprocessing", "false"),
					resource.TestCheckResourceAttr(typeHeadN, "mapped_param",
						"{\"mapped_param_location\":\"query\",\"mapped_param_name\":\"headNParam\",\"mapped_param_type\":\"string\"}"),
					resource.TestCheckResourceAttr(typeHeadN, "map.#", "1"),
					resource.TestCheckResourceAttr(typeHeadN, "map.0", "{\"intercept_length\":5,\"mapped_param_value\":\"\"}"),
					resource.TestMatchResourceAttr(typeHeadN, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Orchestration rule: type tail N.
					rcTypeTailN.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeTailN, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeTailN, "name", name+"_type_tail_n"),
					resource.TestCheckResourceAttr(typeTailN, "strategy", "tail_n"),
					resource.TestCheckResourceAttr(typeTailN, "is_preprocessing", "false"),
					resource.TestCheckResourceAttr(typeTailN, "mapped_param",
						"{\"mapped_param_location\":\"query\",\"mapped_param_name\":\"tailNParam\",\"mapped_param_type\":\"string\"}"),
					resource.TestCheckResourceAttr(typeTailN, "map.#", "1"),
					resource.TestCheckResourceAttr(typeTailN, "map.0", "{\"intercept_length\":5,\"mapped_param_value\":\"\"}"),
					resource.TestMatchResourceAttr(typeTailN, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccOrchestrationRule_basic(updateName),
				Check: resource.ComposeTestCheckFunc(
					// Orchestration rule: type list.
					rcTypeList.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeList, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeList, "name", updateName+"_type_list"),
					resource.TestMatchResourceAttr(typeList, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Orchestration rule: type range.
					rcTypeRange.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeRange, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeRange, "name", updateName+"_type_range"),
					resource.TestMatchResourceAttr(typeRange, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Orchestration rule: type hash.
					rcTypeHash.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeHash, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeHash, "name", updateName+"_type_hash"),
					resource.TestMatchResourceAttr(typeHash, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Orchestration rule: type hash range.
					rcTypeHashRange.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeHashRange, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeHashRange, "name", updateName+"_type_hash_range"),
					resource.TestMatchResourceAttr(typeHashRange, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Orchestration rule: type none value.
					rcTypeNoneValue.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeNoneValue, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeNoneValue, "name", updateName+"_type_none_value"),
					resource.TestMatchResourceAttr(typeNoneValue, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Orchestration rule: type default.
					rcTypeDefault.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeDefault, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeDefault, "name", updateName+"_type_default"),
					resource.TestMatchResourceAttr(typeDefault, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Orchestration rule: type head N.
					rcTypeHeadN.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeHeadN, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeHeadN, "name", updateName+"_type_head_n"),
					resource.TestMatchResourceAttr(typeHeadN, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Orchestration rule: type tail N.
					rcTypeTailN.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeTailN, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeTailN, "name", updateName+"_type_tail_n"),
					resource.TestMatchResourceAttr(typeTailN, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      typeList,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccOrchestrationRuleImportStateFunc(typeList),
			},
			{
				ResourceName:      typeRange,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccOrchestrationRuleImportStateFunc(typeRange),
			},
			{
				ResourceName:      typeHash,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccOrchestrationRuleImportStateFunc(typeHash),
			},
			{
				ResourceName:      typeHashRange,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccOrchestrationRuleImportStateFunc(typeHashRange),
			},
			{
				ResourceName:      typeNoneValue,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccOrchestrationRuleImportStateFunc(typeNoneValue),
			},
			{
				ResourceName:      typeDefault,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccOrchestrationRuleImportStateFunc(typeDefault),
			},
			{
				ResourceName:      typeHeadN,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccOrchestrationRuleImportStateFunc(typeHeadN),
			},
			{
				ResourceName:      typeTailN,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccOrchestrationRuleImportStateFunc(typeTailN),
			},
		},
	})
}

func testAccOrchestrationRuleImportStateFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", rsName)
		}
		if rs.Primary.Attributes["instance_id"] == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<id>', but got '%s/%s'",
				rs.Primary.Attributes["instance_id"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.ID), nil
	}
}

func testAccPlugin_basic_invalidType() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
resource "huaweicloud_apig_orchestration_rule" "type_invalid" {
  instance_id = "%[1]s"
  name        = "%[2]s_type_invalid"
  strategy    = "invalid"

  mapped_param = jsonencode({
    "mapped_param_name": "invalidParam",
    "mapped_param_type": "string",
    "mapped_param_location": "header"
  })
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccOrchestrationRule_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_orchestration_rule" "type_list" {
  instance_id = "%[1]s"
  name        = "%[2]s_type_list"
  strategy    = "list"

  mapped_param = jsonencode({
    "mapped_param_name": "listParam",
    "mapped_param_type": "string",
    "mapped_param_location": "header"
  })

  map = [
    jsonencode({
      "mapped_param_value": "ValueAA",
      "map_param_list": ["ValueA"]
    }),
    jsonencode({
      "mapped_param_value": "ValueCC",
      "map_param_list": ["ValueC"]
    }),
    jsonencode({
      "mapped_param_value": "ValueBB",
      "map_param_list": ["ValueB"]
    })
  ]
}

resource "huaweicloud_apig_orchestration_rule" "type_range" {
  instance_id = "%[1]s"
  name        = "%[2]s_type_range"
  strategy    = "range"

  mapped_param = jsonencode({
    "mapped_param_name": "rangeParam",
    "mapped_param_type": "number",
    "mapped_param_location": "query"
  })

  map = [
    jsonencode({
      "mapped_param_value": "10001",
      "map_param_range": {
        "range_start": "1000",
        "range_end": "1999"
      }
    }),
    jsonencode({
      "mapped_param_value": "10003",
      "map_param_range": {
        "range_start": "3000",
        "range_end": "3999"
      }
    }),
    jsonencode({
      "mapped_param_value": "10002",
      "map_param_range": {
        "range_start": "2000",
        "range_end": "2999"
      }
    })
  ]
}

resource "huaweicloud_apig_orchestration_rule" "type_hash" {
  instance_id = "%[1]s"
  name        = "%[2]s_type_hash"
  strategy    = "hash"

  mapped_param = jsonencode({
    "mapped_param_name": "hashParam",
    "mapped_param_type": "string",
    "mapped_param_location": "header"
  })
}

resource "huaweicloud_apig_orchestration_rule" "type_hash_range" {
  instance_id = "%[1]s"
  name        = "%[2]s_type_hash_range"
  strategy    = "hash_range"

  mapped_param = jsonencode({
    "mapped_param_name": "hashRangeParam",
    "mapped_param_type": "number",
    "mapped_param_location": "header"
  })

  map = [
    jsonencode({
      "mapped_param_value": "10001",
      "map_param_range": {
        "range_start": "1000",
        "range_end": "1999"
      }
    }),
    jsonencode({
      "mapped_param_value": "10003",
      "map_param_range": {
        "range_start": "3000",
        "range_end": "3999"
      }
    }),
    jsonencode({
      "mapped_param_value": "10002",
      "map_param_range": {
        "range_start": "2000",
        "range_end": "2999"
      }
    })
  ]
}

resource "huaweicloud_apig_orchestration_rule" "type_none_value" {
  instance_id = "%[1]s"
  name        = "%[2]s_type_none_value"
  strategy    = "none_value"

  mapped_param = jsonencode({
    "mapped_param_name": "noneValueParam",
    "mapped_param_type": "string",
    "mapped_param_location": "query"
  })

  map = [
    jsonencode({
      "mapped_param_value": "NoneValueReturned"
    })
  ]
}

resource "huaweicloud_apig_orchestration_rule" "type_default" {
  instance_id = "%[1]s"
  name        = "%[2]s_type_default"
  strategy    = "default"

  mapped_param = jsonencode({
    "mapped_param_name": "defaultParam",
    "mapped_param_type": "string",
    "mapped_param_location": "query"
  })

  map = [
    jsonencode({
      "mapped_param_value": "DefaultValueReturned"
    })
  ]
}

resource "huaweicloud_apig_orchestration_rule" "type_head_n" {
  instance_id = "%[1]s"
  name        = "%[2]s_type_head_n"
  strategy    = "head_n"

  mapped_param = jsonencode({
    "mapped_param_name": "headNParam",
    "mapped_param_type": "string",
    "mapped_param_location": "query"
  })

  map = [
    jsonencode({
      "intercept_length": 5,
      "mapped_param_value": ""
    })
  ]
}

resource "huaweicloud_apig_orchestration_rule" "type_tail_n" {
  instance_id = "%[1]s"
  name        = "%[2]s_type_tail_n"
  strategy    = "tail_n"

  mapped_param = jsonencode({
    "mapped_param_name": "tailNParam",
    "mapped_param_type": "string",
    "mapped_param_location": "query"
  })

  map = [
    jsonencode({
      "intercept_length": 5,
      "mapped_param_value": ""
    })
  ]
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}

func TestAccOrchestrationRule_preprocessing(t *testing.T) {
	var (
		obj interface{}

		typeList   = "huaweicloud_apig_orchestration_rule.type_list"
		rcTypeList = acceptance.InitResourceCheck(typeList, &obj, getOrchestrationRuleFunc)

		typeRange   = "huaweicloud_apig_orchestration_rule.type_range"
		rcTypeRange = acceptance.InitResourceCheck(typeRange, &obj, getOrchestrationRuleFunc)

		typeHash   = "huaweicloud_apig_orchestration_rule.type_hash"
		rcTypeHash = acceptance.InitResourceCheck(typeHash, &obj, getOrchestrationRuleFunc)

		typeHashRange   = "huaweicloud_apig_orchestration_rule.type_hash_range"
		rcTypeHashRange = acceptance.InitResourceCheck(typeHashRange, &obj, getOrchestrationRuleFunc)

		typeHeadN   = "huaweicloud_apig_orchestration_rule.type_head_n"
		rcTypeHeadN = acceptance.InitResourceCheck(typeHeadN, &obj, getOrchestrationRuleFunc)

		typeTailN   = "huaweicloud_apig_orchestration_rule.type_tail_n"
		rcTypeTailN = acceptance.InitResourceCheck(typeTailN, &obj, getOrchestrationRuleFunc)

		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcTypeList.CheckResourceDestroy(),
			rcTypeRange.CheckResourceDestroy(),
			rcTypeHash.CheckResourceDestroy(),
			rcTypeHashRange.CheckResourceDestroy(),
			rcTypeHeadN.CheckResourceDestroy(),
			rcTypeTailN.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccOrchestrationRule_preprocessing(name),
				Check: resource.ComposeTestCheckFunc(
					// Orchestration rule: type list.
					rcTypeList.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeList, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeList, "name", name+"_type_list"),
					resource.TestCheckResourceAttr(typeList, "strategy", "list"),
					resource.TestCheckResourceAttr(typeList, "is_preprocessing", "true"),
					resource.TestCheckResourceAttr(typeList, "mapped_param", ""),
					resource.TestCheckResourceAttr(typeList, "map.#", "3"),
					resource.TestCheckResourceAttr(typeList, "map.0",
						"{\"map_param_list\":[\"ValueA\"],\"mapped_param_value\":\"ValueAA\"}"),
					resource.TestCheckResourceAttr(typeList, "map.1",
						"{\"map_param_list\":[\"ValueC\"],\"mapped_param_value\":\"ValueCC\"}"),
					resource.TestCheckResourceAttr(typeList, "map.2",
						"{\"map_param_list\":[\"ValueB\"],\"mapped_param_value\":\"ValueBB\"}"),
					resource.TestMatchResourceAttr(typeList, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Orchestration rule: type range.
					rcTypeRange.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeRange, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeRange, "name", name+"_type_range"),
					resource.TestCheckResourceAttr(typeRange, "strategy", "range"),
					resource.TestCheckResourceAttr(typeRange, "is_preprocessing", "true"),
					resource.TestCheckResourceAttr(typeRange, "mapped_param", ""),
					resource.TestCheckResourceAttr(typeRange, "map.#", "3"),
					resource.TestCheckResourceAttr(typeRange, "map.0",
						"{\"map_param_range\":{\"range_end\":\"1999\",\"range_start\":\"1000\"},\"mapped_param_value\":\"10001\"}"),
					resource.TestCheckResourceAttr(typeRange, "map.1",
						"{\"map_param_range\":{\"range_end\":\"3999\",\"range_start\":\"3000\"},\"mapped_param_value\":\"10003\"}"),
					resource.TestCheckResourceAttr(typeRange, "map.2",
						"{\"map_param_range\":{\"range_end\":\"2999\",\"range_start\":\"2000\"},\"mapped_param_value\":\"10002\"}"),
					resource.TestMatchResourceAttr(typeRange, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Orchestration rule: type hash.
					rcTypeHash.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeHash, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeHash, "name", name+"_type_hash"),
					resource.TestCheckResourceAttr(typeHash, "strategy", "hash"),
					resource.TestCheckResourceAttr(typeHash, "is_preprocessing", "true"),
					resource.TestCheckResourceAttr(typeHash, "mapped_param", ""),
					resource.TestMatchResourceAttr(typeHash, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Orchestration rule: type hash range.
					rcTypeHashRange.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeHashRange, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeHashRange, "name", name+"_type_hash_range"),
					resource.TestCheckResourceAttr(typeHashRange, "strategy", "hash_range"),
					resource.TestCheckResourceAttr(typeHashRange, "is_preprocessing", "true"),
					resource.TestCheckResourceAttr(typeHashRange, "mapped_param", ""),
					resource.TestCheckResourceAttr(typeHashRange, "map.#", "3"),
					resource.TestCheckResourceAttr(typeHashRange, "map.0",
						"{\"map_param_range\":{\"range_end\":\"1999\",\"range_start\":\"1000\"},\"mapped_param_value\":\"10001\"}"),
					resource.TestCheckResourceAttr(typeHashRange, "map.1",
						"{\"map_param_range\":{\"range_end\":\"3999\",\"range_start\":\"3000\"},\"mapped_param_value\":\"10003\"}"),
					resource.TestCheckResourceAttr(typeHashRange, "map.2",
						"{\"map_param_range\":{\"range_end\":\"2999\",\"range_start\":\"2000\"},\"mapped_param_value\":\"10002\"}"),
					resource.TestMatchResourceAttr(typeHashRange, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Orchestration rule: type head N.
					rcTypeHeadN.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeHeadN, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeHeadN, "name", name+"_type_head_n"),
					resource.TestCheckResourceAttr(typeHeadN, "strategy", "head_n"),
					resource.TestCheckResourceAttr(typeHeadN, "is_preprocessing", "true"),
					resource.TestCheckResourceAttr(typeHeadN, "mapped_param", ""),
					resource.TestCheckResourceAttr(typeHeadN, "map.#", "1"),
					resource.TestCheckResourceAttr(typeHeadN, "map.0", "{\"intercept_length\":5,\"mapped_param_value\":\"\"}"),
					resource.TestMatchResourceAttr(typeHeadN, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Orchestration rule: type tail N.
					rcTypeTailN.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeTailN, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeTailN, "name", name+"_type_tail_n"),
					resource.TestCheckResourceAttr(typeTailN, "strategy", "tail_n"),
					resource.TestCheckResourceAttr(typeTailN, "is_preprocessing", "true"),
					resource.TestCheckResourceAttr(typeTailN, "mapped_param", ""),
					resource.TestCheckResourceAttr(typeTailN, "map.#", "1"),
					resource.TestCheckResourceAttr(typeTailN, "map.0", "{\"intercept_length\":5,\"mapped_param_value\":\"\"}"),
					resource.TestMatchResourceAttr(typeTailN, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccOrchestrationRule_basic(updateName),
				Check: resource.ComposeTestCheckFunc(
					// Orchestration rule: type list.
					rcTypeList.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeList, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeList, "name", updateName+"_type_list"),
					resource.TestMatchResourceAttr(typeList, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Orchestration rule: type range.
					rcTypeRange.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeRange, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeRange, "name", updateName+"_type_range"),
					resource.TestMatchResourceAttr(typeRange, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Orchestration rule: type hash.
					rcTypeHash.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeHash, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeHash, "name", updateName+"_type_hash"),
					resource.TestMatchResourceAttr(typeHash, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Orchestration rule: type hash range.
					rcTypeHashRange.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeHashRange, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeHashRange, "name", updateName+"_type_hash_range"),
					resource.TestMatchResourceAttr(typeHashRange, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Orchestration rule: type head N.
					rcTypeHeadN.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeHeadN, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeHeadN, "name", updateName+"_type_head_n"),
					resource.TestMatchResourceAttr(typeHeadN, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Orchestration rule: type tail N.
					rcTypeTailN.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeTailN, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeTailN, "name", updateName+"_type_tail_n"),
					resource.TestMatchResourceAttr(typeTailN, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      typeList,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccOrchestrationRuleImportStateFunc(typeList),
			},
			{
				ResourceName:      typeRange,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccOrchestrationRuleImportStateFunc(typeRange),
			},
			{
				ResourceName:      typeHash,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccOrchestrationRuleImportStateFunc(typeHash),
			},
			{
				ResourceName:      typeHashRange,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccOrchestrationRuleImportStateFunc(typeHashRange),
			},
			{
				ResourceName:      typeHeadN,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccOrchestrationRuleImportStateFunc(typeHeadN),
			},
			{
				ResourceName:      typeTailN,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccOrchestrationRuleImportStateFunc(typeTailN),
			},
		},
	})
}

func testAccOrchestrationRule_preprocessing(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_orchestration_rule" "type_list" {
  instance_id      = "%[1]s"
  name             = "%[2]s_type_list"
  strategy         = "list"
  is_preprocessing = true

  map = [
    jsonencode({
      "mapped_param_value": "ValueAA",
      "map_param_list": ["ValueA"]
    }),
    jsonencode({
      "mapped_param_value": "ValueCC",
      "map_param_list": ["ValueC"]
    }),
    jsonencode({
      "mapped_param_value": "ValueBB",
      "map_param_list": ["ValueB"]
    })
  ]
}

resource "huaweicloud_apig_orchestration_rule" "type_range" {
  instance_id      = "%[1]s"
  name             = "%[2]s_type_range"
  strategy         = "range"
  is_preprocessing = true

  map = [
    jsonencode({
      "mapped_param_value": "10001",
      "map_param_range": {
        "range_start": "1000",
        "range_end": "1999"
      }
    }),
    jsonencode({
      "mapped_param_value": "10003",
      "map_param_range": {
        "range_start": "3000",
        "range_end": "3999"
      }
    }),
    jsonencode({
      "mapped_param_value": "10002",
      "map_param_range": {
        "range_start": "2000",
        "range_end": "2999"
      }
    })
  ]
}

resource "huaweicloud_apig_orchestration_rule" "type_hash" {
  instance_id      = "%[1]s"
  name             = "%[2]s_type_hash"
  strategy         = "hash"
  is_preprocessing = true
}

resource "huaweicloud_apig_orchestration_rule" "type_hash_range" {
  instance_id      = "%[1]s"
  name             = "%[2]s_type_hash_range"
  strategy         = "hash_range"
  is_preprocessing = true

  map = [
    jsonencode({
      "mapped_param_value": "10001",
      "map_param_range": {
        "range_start": "1000",
        "range_end": "1999"
      }
    }),
    jsonencode({
      "mapped_param_value": "10003",
      "map_param_range": {
        "range_start": "3000",
        "range_end": "3999"
      }
    }),
    jsonencode({
      "mapped_param_value": "10002",
      "map_param_range": {
        "range_start": "2000",
        "range_end": "2999"
      }
    })
  ]
}

resource "huaweicloud_apig_orchestration_rule" "type_head_n" {
  instance_id      = "%[1]s"
  name             = "%[2]s_type_head_n"
  strategy         = "head_n"
  is_preprocessing = true

  map = [
    jsonencode({
      "intercept_length": 5,
      "mapped_param_value": ""
    })
  ]
}

resource "huaweicloud_apig_orchestration_rule" "type_tail_n" {
  instance_id      = "%[1]s"
  name             = "%[2]s_type_tail_n"
  strategy         = "tail_n"
  is_preprocessing = true

  map = [
    jsonencode({
      "intercept_length": 5,
      "mapped_param_value": ""
    })
  ]
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}

func TestAccOrchestrationRule_strategyConvert(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_apig_orchestration_rule.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getOrchestrationRuleFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccOrchestrationRule_strategyConvert_step1(name),
				Check: resource.ComposeTestCheckFunc(
					// Orchestration rule: type hash.
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(resourceName, "name", name+"_type_hash"),
					resource.TestCheckResourceAttr(resourceName, "strategy", "hash"),
					resource.TestCheckResourceAttr(resourceName, "is_preprocessing", "false"),
					resource.TestCheckResourceAttr(resourceName, "mapped_param",
						"{\"mapped_param_location\":\"header\",\"mapped_param_name\":\"strategyConvertParam\",\"mapped_param_type\":\"string\"}"),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccOrchestrationRule_strategyConvert_step2(name),
				Check: resource.ComposeTestCheckFunc(
					// Orchestration rule: type hash range.
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(resourceName, "name", name+"_type_hash_range"),
					resource.TestCheckResourceAttr(resourceName, "strategy", "hash_range"),
					resource.TestCheckResourceAttr(resourceName, "is_preprocessing", "false"),
					resource.TestCheckResourceAttr(resourceName, "mapped_param",
						"{\"mapped_param_location\":\"header\",\"mapped_param_name\":\"strategyConvertParam\",\"mapped_param_type\":\"string\"}"),
					resource.TestCheckResourceAttr(resourceName, "map.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "map.0",
						"{\"map_param_range\":{\"range_end\":\"1999\",\"range_start\":\"1000\"},\"mapped_param_value\":\"10001\"}"),
					resource.TestCheckResourceAttr(resourceName, "map.1",
						"{\"map_param_range\":{\"range_end\":\"3999\",\"range_start\":\"3000\"},\"mapped_param_value\":\"10003\"}"),
					resource.TestCheckResourceAttr(resourceName, "map.2",
						"{\"map_param_range\":{\"range_end\":\"2999\",\"range_start\":\"2000\"},\"mapped_param_value\":\"10002\"}"),
					resource.TestMatchResourceAttr(resourceName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccOrchestrationRuleImportStateFunc(resourceName),
			},
		},
	})
}

func testAccOrchestrationRule_strategyConvert_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_orchestration_rule" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s_type_hash"
  strategy    = "hash"

  mapped_param = jsonencode({
    "mapped_param_name": "strategyConvertParam",
    "mapped_param_type": "string",
    "mapped_param_location": "header"
  })
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccOrchestrationRule_strategyConvert_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_orchestration_rule" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s_type_hash_range"
  strategy    = "hash_range"

  mapped_param = jsonencode({
    "mapped_param_name": "strategyConvertParam",
    "mapped_param_type": "string",
    "mapped_param_location": "header"
  })

  map = [
    jsonencode({
      "mapped_param_value": "10001",
      "map_param_range": {
        "range_start": "1000",
        "range_end": "1999"
      }
    }),
    jsonencode({
      "mapped_param_value": "10003",
      "map_param_range": {
        "range_start": "3000",
        "range_end": "3999"
      }
    }),
    jsonencode({
      "mapped_param_value": "10002",
      "map_param_range": {
        "range_start": "2000",
        "range_end": "2999"
      }
    })
  ]
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}
