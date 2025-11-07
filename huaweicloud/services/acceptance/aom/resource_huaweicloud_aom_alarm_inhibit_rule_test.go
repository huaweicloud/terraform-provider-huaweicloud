package aom

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/aom"
)

func getAlarmInhibitRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("aom", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating AOM client: %s", err)
	}

	return aom.GetAlarmInhibitRuleByName(client, state.Primary.Attributes["enterprise_project_id"], state.Primary.ID)
}

func TestAccAlarmInhibitRule_basic(t *testing.T) {
	var (
		rName = "huaweicloud_aom_alarm_inhibit_rule.test"
		obj   interface{}
		rc    = acceptance.InitResourceCheck(rName, &obj, getAlarmInhibitRuleResourceFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAlarmInhibitRule_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(rName, "source_matches.#", "3"),
					resource.TestCheckResourceAttr(rName, "source_matches.0.conditions.#", "2"),
					resource.TestCheckResourceAttr(rName, "source_matches.0.conditions.0.key", "event_severity"),
					resource.TestCheckResourceAttr(rName, "source_matches.0.conditions.0.operate", "EXIST"),
					resource.TestCheckResourceAttr(rName, "target_matches.#", "2"),
					resource.TestCheckResourceAttr(rName, "target_matches.0.conditions.#", "2"),
					resource.TestCheckResourceAttr(rName, "target_matches.0.conditions.1.key", "owner"),
					resource.TestCheckResourceAttr(rName, "target_matches.0.conditions.1.operate", "REGEX"),
					resource.TestCheckResourceAttr(rName, "target_matches.0.conditions.1.values.0", "terraform"),
				),
			},
			{
				Config: testAccAlarmInhibitRule_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "source_matches.#", "2"),
					resource.TestCheckResourceAttr(rName, "source_matches.0.conditions.#", "1"),
					resource.TestCheckResourceAttr(rName, "target_matches.#", "3"),
					resource.TestCheckResourceAttr(rName, "target_matches.0.conditions.#", "1"),
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

func testAccAlarmInhibitRule_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_alarm_inhibit_rule" "test" {
  name        = "%[1]s"
  description = "Created by terraform script"

  source_matches {
    conditions {
      key     = "event_severity"
      operate = "EXIST"
    }
    conditions {
      key     = "event_severity"
      operate = "EQUALS"
      values  = ["Critical"]
    }
  }
  source_matches {
    conditions {
      key     = "event_severity"
      operate = "EXIST"
    }
  }
  source_matches {
    conditions {
      key     = "foo"
      operate = "REGEX"
      values  = ["bar"]
    }
  }

  target_matches {
    conditions {
      key     = "event_severity"
      operate = "EQUALS"
      values  = ["Info"]
    }
    conditions {
      key     = "owner"
      operate = "REGEX"
      values  = ["terraform"]
    }
  }
  target_matches {
    conditions {
      key     = "event_severity"
      operate = "EXIST"
    }
  }
}
`, name)
}

func testAccAlarmInhibitRule_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_alarm_inhibit_rule" "test" {
  name = "%[1]s"

  source_matches {
    conditions {
      key     = "event_severity"
      operate = "EQUALS"
      values  = ["Critical", "Major"]
    }
  }
  source_matches {
    conditions {
      key     = "foo"
      operate = "REGEX"
      values  = ["bar"]
    }
    conditions {
      key     = "event_severity"
      operate = "EXIST"
    }
  }

  target_matches {
    conditions {
      key     = "event_severity"
      operate = "EQUALS"
      values  = ["Info"]
    }
  }
  target_matches {
    conditions {
      key     = "event_severity"
      operate = "EXIST"
    }
  }
  target_matches {
    conditions {
      key     = "owner"
      operate = "REGEX"
      values  = ["terraform"]
    }
  }
}
`, name)
}

// Currently, the `match_v3` field is only supported in some regions, such as `cn-north-9`.
func TestAccAlarmInhibitRule_with_matchv3AndEpsId(t *testing.T) {
	var (
		obj                interface{}
		rNameWithGroupRule = "huaweicloud_aom_alarm_inhibit_rule.test.0"
		rcWithGroupRule    = acceptance.InitResourceCheck(rNameWithGroupRule, &obj, getAlarmInhibitRuleResourceFunc)

		rNameWithoutGroupRule = "huaweicloud_aom_alarm_inhibit_rule.test.1"
		rcWithoutGroupRule    = acceptance.InitResourceCheck(rNameWithoutGroupRule, &obj, getAlarmInhibitRuleResourceFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcWithGroupRule.CheckResourceDestroy(),
			rcWithoutGroupRule.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccAlarmInhibitRule_with_matchv3AndEpsId_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rcWithGroupRule.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameWithGroupRule, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttrPair(rNameWithGroupRule, "binding_group_rule",
						"huaweicloud_aom_alarm_group_rule.test", "name"),
					resource.TestCheckResourceAttr(rNameWithGroupRule, "source_matches.#", "0"),
					resource.TestCheckResourceAttr(rNameWithGroupRule, "target_matches.#", "0"),
					rcWithoutGroupRule.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameWithoutGroupRule, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rNameWithoutGroupRule, "binding_group_rule", ""),
					resource.TestCheckResourceAttr(rNameWithoutGroupRule, "source_matches.#", "0"),
					resource.TestCheckResourceAttr(rNameWithoutGroupRule, "target_matches.#", "0"),
				),
			},
			{
				Config: testAccAlarmInhibitRule_with_matchv3AndEpsId_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rcWithGroupRule.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameWithGroupRule, "binding_group_rule", ""),
					rcWithoutGroupRule.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rNameWithoutGroupRule, "binding_group_rule",
						"huaweicloud_aom_alarm_group_rule.test", "name"),
				),
			},
			{
				ResourceName:      "huaweicloud_aom_alarm_inhibit_rule.test[0]",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAlarmInhibitRuleWithEpsIdImportStateFunc(rNameWithGroupRule),
			},
			{
				ResourceName:      "huaweicloud_aom_alarm_inhibit_rule.test[1]",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAlarmInhibitRuleWithEpsIdImportStateFunc(rNameWithoutGroupRule),
			},
		},
	})
}

// As long as the name of the alarm inhibit rule is unique, duplicate node IDs under it will not affect the creation of other rules.
func testAccAlarmInhibitRule_with_matchv3AndEpsId_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_aom_alarm_inhibit_rule" "test" {
  count = 2

  name                  = "%[2]s${count.index}"
  enterprise_project_id = "%[3]s"
  binding_group_rule    = count.index == 0 ? huaweicloud_aom_alarm_group_rule.test.name : null

  match_v3 = jsonencode({
    nodes = [
      # The node IDs naming rule: "flow" followed by 8 random characters, separated by "-", for example, "flow-VHSvKyYY".
      # Random strings can only consist of letters, numbers, underscores, and hyphens.
      {
        id   = "flow-VHSvKyYY"
        type = "Start"
      },
      {
        id           = "flow-ufs5Zsn3"
        type         = "Card"
        businessType = "InhibitSourceCondition"

        value = {
          # The card IDs and their node IDs naming rule: "{businessType}" followed by 8 random characters, separated by "-",
          # for example, "InhibitSourceCondition-wtJBNJD7".
          id    = "InhibitSourceCondition-wtJBNJD7"
          type  = "bool"
          value = "and"
          nodes = [
            {
              id    = "InhibitSourceCondition-awPpHgHv"
              type  = "text"
              match = {
                key     = "resource_type"
                operate = "EQUALS"
                value   = ["service"]
              }
            },
            {
              id    = "InhibitSourceCondition-DmfXe2vN"
              type  = "text"
              match = {
                key     = "resource_type"
                operate = "EQUALS"
                value   = ["service"]
              }
            },
            {
              id    = "InhibitSourceCondition-0fnYAVWg"
              type  = "bool"
              value = "and"
              nodes = [
                {
                  id    = "InhibitSourceCondition-Xn7CUof5"
                  type  = "text"
                  match = {
                    key     = "event_severity"
                    operate = "EQUALS"
                    value   = ["Critical"]
                  }
                }
              ]
            },
            {
              id    = "InhibitSourceCondition-juUOsA2L"
              type  = "text"
              match = {
                key     = "resource_provider"
                operate = "EQUALS"
                value   = ["1"]
              }
            },
            {
              id    = "InhibitSourceCondition-JvtIwdL6"
              type  = "text"
              match = {
                key     = "1"
                operate = "EQUALS"
                value   = ["1"]
              }
            }
          ]
        }
      },
      {
        id           = "flow-eAheOiH5"
        type         = "Card"
        businessType = "InhibitTargetCondition"

        value = {
          id    = "InhibitSourceCondition--foPmJou"
          type  = "bool"
          value = "and"
          nodes = [
            {
              id    = "InhibitSourceCondition--m9jJEAN"
              type  = "text"
              match = {
                key     = "resource_type"
                operate = "EQUALS"
                value   = ["service"]
              }
            },
            {
              id    = "InhibitSourceCondition-MmTv1BZ5"
              type  = "text"
              match = {
                key     = "resource_provider"
                operate = "EQUALS"
                value   = ["1"]
              }
            },
            {
              id    = "InhibitSourceCondition-7oz9ddLB"
              type  = "bool"
              value = "and"
              nodes = [
                {
                  id    = "InhibitSourceCondition-oM73YTRG"
                  type  = "text"
                  match = {
                    key     = "2"
                    operate = "EQUALS"
                    value   = ["2"]
                  }
                }
              ]
            }
          ]
        }
      },
      {
        id   = "flow-f6P36v4a"
        type = "End"
      }
    ]
    edges = [
      # The edge IDs naming rule: "edge" followed by 8 random characters, separated by "-", for example, "edge-M3pfK202".
      # Random strings can only consist of letters, numbers, underscores, and hyphens.
      # "source": The ID of the start node of the edge. The value is the ID of the node in the nodes list.
      # "target": The ID of the end node of the edge. The value is the ID of the node in the nodes list.
      {
        id     = "edge-M3pfK202"
        source = "flow-VHSvKyYY"
        target = "flow-ufs5Zsn3"
      },
      {
        id     = "edge-OrwQoSVg"
        source = "flow-ufs5Zsn3"
        target = "flow-eAheOiH5"
        value  = "true"
      },
      {
        id     = "edge-ukREdGnr"
        source = "flow-ufs5Zsn3"
        target = "flow-f6P36v4a"
        value  = "false"
      }
    ]
  })
}
`, testAlarmGroupRule_basic(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccAlarmInhibitRule_with_matchv3AndEpsId_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_aom_alarm_inhibit_rule" "test" {
  count = 2

  name                  = "%[2]s${count.index}"
  enterprise_project_id = "%[3]s"
  binding_group_rule    = count.index == 0 ? null : huaweicloud_aom_alarm_group_rule.test.name

  match_v3 = jsonencode({
    nodes = [
      {
        id   = "flow-VHSvKyYY"
        type = "Start"
      },
      {
        id           = "flow-ufs5Zsn3"
        type         = "Card"
        businessType = "InhibitSourceCondition"

        value = {
          id    = "InhibitSourceCondition-wtJBNJD7"
          type  = "bool"
          value = "and"
          nodes = [
            {
              id    = "InhibitSourceCondition-awPpHgHv"
              type  = "text"
              match = {
                key     = "resource_type"
                operate = "EQUALS"
                value   = ["service"]
              }
            },
            {
              id    = "InhibitSourceCondition-0fnYAVWg"
              type  = "bool"
              value = "and"
              nodes = [
                {
                  id    = "InhibitSourceCondition-Xn7CUof5"
                  type  = "text"
                  match = {
                    key     = "event_severity"
                    operate = "EQUALS"
                    value   = ["Critical"]
                  }
                }
              ]
            },
          ]
        }
      },
      {
        id           = "flow-eAheOiH5"
        type         = "Card"
        businessType = "InhibitTargetCondition"

        value = {
          id    = "InhibitSourceCondition--foPmJou"
          type  = "bool"
          value = "and"
          nodes = [
            {
              id    = "InhibitSourceCondition--m9jJEAN"
              type  = "text"
              match = {
                key     = "resource_type"
                operate = "EQUALS"
                value   = ["service"]
              }
            },
            {
              id    = "InhibitSourceCondition-MmTv1BZ5"
              type  = "text"
              match = {
                key     = "resource_provider"
                operate = "EQUALS"
                value   = ["1"]
              }
            },
            {
              id    = "InhibitSourceCondition-7oz9ddLB"
              type  = "bool"
              value = "and"
              nodes = [
                {
                  id    = "InhibitSourceCondition-oM73YTRG"
                  type  = "text"
                  match = {
                    key     = "2"
                    operate = "EQUALS"
                    value   = ["2"]
                  }
                }
              ]
            }
          ]
        }
      },
      {
        id   = "flow-f6P36v4a"
        type = "End"
      }
    ]
    edges = [
      {
        id     = "edge-M3pfK202"
        source = "flow-VHSvKyYY"
        target = "flow-ufs5Zsn3"
      },
      {
        id     = "edge-OrwQoSVg"
        source = "flow-ufs5Zsn3"
        target = "flow-eAheOiH5"
        value  = "true"
      },
      {
        id     = "edge-ukREdGnr"
        source = "flow-ufs5Zsn3"
        target = "flow-f6P36v4a"
        value  = "false"
      }
    ]
  })
}
`, testAlarmGroupRule_basic(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccAlarmInhibitRuleWithEpsIdImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		epsId := rs.Primary.Attributes["enterprise_project_id"]
		if epsId == "" {
			return "", fmt.Errorf("The imported ID specifies an invalid format, want '<id>/<enterprise_project_id>', "+
				"but got '%s/%s'", rs.Primary.ID, epsId)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.ID, epsId), nil
	}
}
