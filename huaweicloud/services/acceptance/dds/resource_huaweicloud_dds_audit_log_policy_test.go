package dds

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

func getDdsAuditLogPolicyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getAuditLog: Query DDS audit log
	var (
		getAuditLogPolicyHttpUrl = "v3/{project_id}/instances/{instance_id}/auditlog-policy"
		getAuditLogPolicyProduct = "dds"
	)
	getAuditLogPolicyClient, err := cfg.NewServiceClient(getAuditLogPolicyProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DDS Client: %s", err)
	}

	instanceID := state.Primary.ID
	getAuditLogPolicyPath := getAuditLogPolicyClient.Endpoint + getAuditLogPolicyHttpUrl
	getAuditLogPolicyPath = strings.ReplaceAll(getAuditLogPolicyPath, "{project_id}",
		getAuditLogPolicyClient.ProjectID)
	getAuditLogPolicyPath = strings.ReplaceAll(getAuditLogPolicyPath, "{instance_id}", instanceID)

	getAuditLogPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	getAuditLogPolicyResp, err := getAuditLogPolicyClient.Request("GET", getAuditLogPolicyPath, &getAuditLogPolicyOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DDS audit log policy: %s", err)
	}

	getAuditLogPolicyRespBody, err := utils.FlattenResponse(getAuditLogPolicyResp)
	if err != nil {
		return nil, err
	}

	keepDays := utils.PathSearch("keep_days", getAuditLogPolicyRespBody, 0)
	if keepDays.(float64) == 0 {
		return nil, fmt.Errorf("the instance %s has no audit log policy", instanceID)
	}

	return getAuditLogPolicyRespBody, nil
}

func TestAccDdsAuditLogPolicy_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dds_audit_log_policy.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDdsAuditLogPolicyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDdsAuditLogPolicy_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "keep_days", "7"),
					resource.TestCheckResourceAttr(rName, "audit_types.#", "1"),
					resource.TestCheckResourceAttr(rName, "audit_types.0", "auth"),
				),
			},
			{
				Config: testDdsAuditLogPolicy_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "keep_days", "15"),
					resource.TestCheckResourceAttr(rName, "audit_types.#", "1"),
					resource.TestCheckResourceAttr(rName, "audit_types.0", "insert"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id"},
			},
		},
	})
}

func testDdsAuditLogPolicy_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dds_audit_log_policy" "test" {
  instance_id = huaweicloud_dds_instance.instance.id
  keep_days   = 7

  audit_types = [
    "auth"
  ]
}
`, testAccDDSInstanceV3Config_basic(name, 8800))
}

func testDdsAuditLogPolicy_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dds_audit_log_policy" "test" {
  instance_id = huaweicloud_dds_instance.instance.id
  keep_days   = 15
  
  audit_types = [
    "insert"
  ]
}
`, testAccDDSInstanceV3Config_basic(name, 8800))
}
