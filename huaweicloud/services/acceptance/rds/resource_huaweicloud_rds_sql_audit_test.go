package rds

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getSQLAuditResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getSQLAudit: Query the RDS SQL audit
	var (
		getSQLAuditHttpUrl = "v3/{project_id}/instances/{instance_id}/auditlog-policy"
		getSQLAuditProduct = "rds"
	)
	getSQLAuditClient, err := cfg.NewServiceClient(getSQLAuditProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	getSQLAuditPath := getSQLAuditClient.Endpoint + getSQLAuditHttpUrl
	getSQLAuditPath = strings.ReplaceAll(getSQLAuditPath, "{project_id}", getSQLAuditClient.ProjectID)
	getSQLAuditPath = strings.ReplaceAll(getSQLAuditPath, "{instance_id}", state.Primary.ID)

	getSQLAuditOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getSQLAuditResp, err := getSQLAuditClient.Request("GET", getSQLAuditPath, &getSQLAuditOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS SQL audit: %s", err)
	}

	getSQLAuditRespBody, err := utils.FlattenResponse(getSQLAuditResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS SQL audit: %s", err)
	}

	keepDays := utils.PathSearch("keep_days", getSQLAuditRespBody, float64(0)).(float64)
	if keepDays == 0 {
		return nil, fmt.Errorf("error retrieving RDS SQL audit: %s", err)
	}

	return getSQLAuditRespBody, nil
}

func TestAccSQLAudit_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_rds_sql_audit.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getSQLAuditResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSQLAudit_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "keep_days", "5"),
					resource.TestCheckResourceAttr(rName, "audit_types.#", "4"),
				),
			},
			{
				Config: testSQLAudit_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "keep_days", "9"),
					resource.TestCheckResourceAttr(rName, "audit_types.#", "5"),
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

func testSQLAudit_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_sql_audit" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  keep_days   = "5"
  audit_types = [
    "CREATE_USER",
    "RENAME_USER",
    "CREATE",
    "PREPARED_STATEMENT"
  ]
}
`, testAccRdsInstance_mysql_step1(name))
}

func testSQLAudit_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_sql_audit" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  keep_days   = "9"
  audit_types = [
    "CREATE_USER",
    "DROP_USER",
    "DROP",
    "INSERT",
    "BEGIN/COMMIT/ROLLBACK"
  ]
}
`, testAccRdsInstance_mysql_step1(name))
}
