package cceautopilot

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

func getAutopilotAddonFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	var (
		getAddonHttpUrl = "autopilot/v3/addons/{id}"
		getAddonProduct = "cce"
	)
	getAddonClient, err := cfg.NewServiceClient(getAddonProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CCE Client: %s", err)
	}

	getAddonPath := getAddonClient.Endpoint + getAddonHttpUrl
	getAddonPath = strings.ReplaceAll(getAddonPath, "{id}", state.Primary.ID)

	getAddonOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getAddonResp, err := getAddonClient.Request("GET", getAddonPath, &getAddonOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CCE autopolit addon: %s", err)
	}

	return utils.FlattenResponse(getAddonResp)
}

func TestAccAutopilotAddon_basic(t *testing.T) {
	var (
		cluster      interface{}
		resourceName = "huaweicloud_cce_autopilot_addon.test"

		rc = acceptance.InitResourceCheck(
			resourceName,
			&cluster,
			getAutopilotAddonFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectID(t)
			acceptance.TestAccPreCheckHSSCCEProtection(t)
			acceptance.TestAccPreCheckCertificateBase(t)
			acceptance.TestAccPreCheckCertificateRootCA(t)
			acceptance.TestAccPreCheckSWRUser(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAddon_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "cluster_id",
						"huaweicloud_cce_autopilot_cluster.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "addon_template_name", "log-agent"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.4.3"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"values"},
			},
		},
	})
}

// nolint:revive
func testAccAddon_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cce_autopilot_addon" "test" {
  cluster_id          = huaweicloud_cce_autopilot_cluster.test.id
  version             = "1.4.3"
  addon_template_name = "log-agent"

  values = {
    "basic" = jsonencode({
      "aomEndpoint" : "https://aom.cn-north-4.myhuaweicloud.com",
      "iam_url" : "iam.cn-north-4.myhuaweicloud.com",
      "ltsAccessEndpoint" : "https://lts-access.cn-north-4.myhuaweicloud.com:8102",
      "ltsEndpoint" : "https://lts.cn-north-4.myhuaweicloud.com",
      "region" : "cn-north-4",
      "swr_addr" : "swr.cn-north-4.myhuaweicloud.com",
      "swr_user" : "%s",
      "rbac_enabled" : true,
      "cluster_version" : "v1.28"
    })

    "flavor" = jsonencode({
      "category" : [
        "Autopilot"
      ],
      "description" : "Recommanded when the number of logs per second does not exceed 5000.",
      "is_default" : true,
      "name" : "Autopilot-Low",
      "replicas" : 2,
      "resources" : [
        {
          "name" : "log-operator",
          "limitsCpu" : "1000m",
          "requestsCpu" : "1000m",
          "limitsMem" : "2048Mi",
          "requestsMem" : "2048Mi"
        },
        {
          "name" : "otel-collector",
          "limitsCpu" : "2000m",
          "requestsCpu" : "2000m",
          "limitsMem" : "4096Mi",
          "requestsMem" : "4096Mi"
        },
        {
          "name" : "fluent-bit",
          "limitsMem" : "400Mi",
          "requestsMem" : "50Mi"
        }
      ]
    })

    "custom" = jsonencode({
      "accessKey" : "",
      "aomEndpoint" : "https://aom.cn-north-4.myhuaweicloud.com",
      "aomPrivateEndpointIP" : "",
      "caCert" : "%s",
      "clusterID" : "%s",
      "clusterName" : "%s",
      "cluster_category" : "CCE",
      "createAudit" : false,
      "createDefaultEvent" : false,
      "createDefaultEventToAOM" : true,
      "createDefaultStdout" : false,
      "createKubeApiserver" : false,
      "createKubeControllerManager" : false,
      "createKubeScheduler" : false,
      "enable_autopilot" : true,
      "ltsAccessEndpoint" : "https://lts-access.cn-north-4.myhuaweicloud.com:8102",
      "ltsAuditStreamID" : "",
      "ltsEndpoint" : "https://lts.cn-north-4.myhuaweicloud.com",
      "ltsEventStreamID" : "",
      "ltsGroupID" : "",
      "ltsKubeApiserverStreamID" : "",
      "ltsKubeControllerManagerStreamID" : "",
      "ltsKubeSchedulerStreamID" : "",
      "ltsLogReportDomain" : "",
      "ltsPrivateEndpointIP" : "",
      "ltsStdoutStreamID" : "",
      "multiAZEnabled" : false,
      "paasakskEnable" : true,
      "projectID" : "%s",
      "secretKey" : "",
      "securityToken" : "",
      "serverCert" : "%s",
      "serverKey" : "%s"
    })
  }
}
`, acceptance.HW_SWR_USER, acceptance.HW_CERTIFICATE_ROOT_CA, acceptance.HW_CCE_CLUSTER_ID, acceptance.HW_CCE_CLUSTER_NAME, acceptance.HW_PROJECT_ID, acceptance.HW_CERTIFICATE_CONTENT, acceptance.HW_CERTIFICATE_PRIVATE_KEY)
}
