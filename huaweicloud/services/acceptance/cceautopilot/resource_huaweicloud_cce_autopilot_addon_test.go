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
		rName        = acceptance.RandomAccResourceNameWithDash()

		rc = acceptance.InitResourceCheck(
			resourceName,
			&cluster,
			getAutopilotAddonFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAutopilotAddon_basic(rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "cluster_id",
						"huaweicloud_cce_autopilot_cluster.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "addon_template_name", "log-agent"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccAutopilotAddon_update(rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "cluster_id",
						"huaweicloud_cce_autopilot_cluster.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "addon_template_name", "log-agent"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
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

// testAccAutopilotAddon_base builds the shared infrastructure for addon tests:
//   - A CCE Autopilot cluster (provides cluster_id, cluster_name, cluster_version)
//   - An SWR organization (provides swr_user, which equals the organization name)
//
// The region is derived from the provider configuration.
// The enterprise project ID (epsId) is passed via environment variable HW_ENTERPRISE_PROJECT_ID_TEST.
func testAccAutopilotAddon_base(rName, epsId string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_swr_organization" "test" {
  name = "%[2]s"
}

locals {
  # The SWR user equals the organization name, as the creator of the organization is the SWR user.
  swr_user = huaweicloud_swr_organization.test.name

  # The enterprise project ID used for LTS log reporting.
  eps_id = "%[3]s"

  # The cluster ID and name are derived from the CCE Autopilot cluster resource.
  region_name     = huaweicloud_cce_autopilot_cluster.test.region
  cluster_id      = huaweicloud_cce_autopilot_cluster.test.id
  cluster_name    = huaweicloud_cce_autopilot_cluster.test.name
  cluster_version = huaweicloud_cce_autopilot_cluster.test.version
}
`, testAccCluster_basic(rName), rName, epsId)
}

func testAccAutopilotAddon_basic(rName, epsId string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_autopilot_addon" "test" {
  cluster_id          = local.cluster_id
  addon_template_name = "log-agent"

  values = {
    "basic" = jsonencode({
      "aomEndpoint" : "https://aom.${local.region_name}.myhuaweicloud.com",
      "iam_url" : "iam.${local.region_name}.myhuaweicloud.com",
      "ltsAccessEndpoint" : "https://lts-access.${local.region_name}.myhuaweicloud.com:8102",
      "ltsEndpoint" : "https://lts.${local.region_name}.myhuaweicloud.com",
      "region" : local.region_name,
      "swr_addr" : "swr.${local.region_name}.myhuaweicloud.com",
      "swr_user" : local.swr_user,
      "rbac_enabled" : true,
      "cluster_version" : local.cluster_version
    })

    "flavor" = jsonencode({
      "category" : [
        "Autopilot"
      ],
      "description" : "Applicable to clusters where the logs of a single pod are less than 10000/s or 5 MB/s.",
      "is_default" : true,
      "name" : "custom-resources",
      "resources" : [
        {
          "name" : "log-operator",
          "limitsCpu" : "500m",
          "requestsCpu" : "500m",
          "replicas" : 2,
          "limitsMem" : "2048Mi",
          "requestsMem" : "2048Mi"
        },
        {
          "name" : "otel-collector-event",
          "limitsCpu" : "500m",
          "requestsCpu" : "500m",
          "replicas" : 2,
          "limitsMem" : "2048Mi",
          "requestsMem" : "2048Mi"
        }
      ],
      "size" : "custom"
    })

    "custom" = jsonencode({
      "accessKey" : "",
      "agency_name" : "",
      "aomEndpoint" : "https://aom.${local.region_name}.myhuaweicloud.com",
      "aomPrivateEndpointIP" : "",
      "bufferChunkSize" : "128k",
      "bufferMaxSize" : "512k",
      "caCert" : "",
      "clusterID" : local.cluster_id,
      "clusterName" : local.cluster_name,
      "cluster_category" : "CCE",
      "createAudit" : false,
      "createDefaultEvent" : false,
      "createDefaultEventToAOM" : true,
      "createDefaultStdout" : false,
      "createKubeApiserver" : false,
      "createKubeControllerManager" : false,
      "createKubeScheduler" : false,
      "enableEventReport" : true,
      "enableFullPathCollection" : false,
      "enableGcrypto" : true,
      "enableLogOperatorHA" : true,
      "enableLogReport" : true,
      "featureGates" : ["sendAOMAllEvent", "outputKafka", "podLabelExclude", "fullPathCollection", "independentEvents"],
      "host_network" : false,
      "ltsAccessEndpoint" : "https://lts-access.${local.region_name}.myhuaweicloud.com:8102",
      "ltsAuditStreamID" : "",
      "ltsEndpoint" : "https://lts.${local.region_name}.myhuaweicloud.com",
      "ltsEnterpriseProjectID" : local.eps_id,
      "ltsEventStreamID" : "",
      "ltsGroupID" : "",
      "ltsKubeApiserverStreamID" : "",
      "ltsKubeControllerManagerStreamID" : "",
      "ltsKubeSchedulerStreamID" : "",
      "ltsLogReportDomain" : "",
      "ltsPrivateEndpointIP" : "",
      "ltsStdoutStreamID" : "",
      "maxEventAgeSeconds" : "",
      "memBufLimit" : "40mb",
      "multiAZEnabled" : false,
      "otelReportLogs" : true,
      "paasakskEnable" : true,
      "podDisruptionBudget" : {"create" : true, "maxUnavailable" : 1},
      "projectID" : "",
      "secretKey" : "",
      "securityToken" : "",
      "serverCert" : "",
      "serverKey" : ""
    })
  }
}
`, testAccAutopilotAddon_base(rName, epsId))
}

// testAccAutopilotAddon_update changes the log-agent configuration to test the Update flow.
// The key change is enabling the audit log creation (createAudit: true) and
// disabling the default event to AOM (createDefaultEventToAOM: false).
func testAccAutopilotAddon_update(rName, epsId string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_autopilot_addon" "test" {
  cluster_id          = local.cluster_id
  addon_template_name = "log-agent"

  values = {
    "basic" = jsonencode({
      "aomEndpoint" : "https://aom.${local.region_name}.myhuaweicloud.com",
      "iam_url" : "iam.${local.region_name}.myhuaweicloud.com",
      "ltsAccessEndpoint" : "https://lts-access.${local.region_name}.myhuaweicloud.com:8102",
      "ltsEndpoint" : "https://lts.${local.region_name}.myhuaweicloud.com",
      "region" : local.region_name,
      "swr_addr" : "swr.${local.region_name}.myhuaweicloud.com",
      "swr_user" : local.swr_user,
      "rbac_enabled" : true,
      "cluster_version" : local.cluster_version
    })

    "flavor" = jsonencode({
      "category" : [
        "Autopilot"
      ],
      "description" : "Applicable to clusters where the logs of a single pod are less than 10000/s or 5 MB/s.",
      "is_default" : true,
      "name" : "custom-resources",
      "resources" : [
        {
          "name" : "log-operator",
          "limitsCpu" : "1000m",
          "requestsCpu" : "1000m",
          "replicas" : 3,
          "limitsMem" : "2048Mi",
          "requestsMem" : "2048Mi"
        },
        {
          "name" : "otel-collector-event",
          "limitsCpu" : "1000m",
          "requestsCpu" : "1000m",
          "replicas" : 3,
          "limitsMem" : "2048Mi",
          "requestsMem" : "2048Mi"
        }
      ],
      "size" : "custom"
    })

    "custom" = jsonencode({
      "accessKey" : "",
      "agency_name" : "",
      "aomEndpoint" : "https://aom.${local.region_name}.myhuaweicloud.com",
      "aomPrivateEndpointIP" : "",
      "bufferChunkSize" : "128k",
      "bufferMaxSize" : "512k",
      "caCert" : "",
      "clusterID" : local.cluster_id,
      "clusterName" : local.cluster_name,
      "cluster_category" : "CCE",
      "createAudit" : true,
      "createDefaultEvent" : false,
      "createDefaultEventToAOM" : false,
      "createDefaultStdout" : false,
      "createKubeApiserver" : false,
      "createKubeControllerManager" : false,
      "createKubeScheduler" : false,
      "enableEventReport" : true,
      "enableFullPathCollection" : false,
      "enableGcrypto" : true,
      "enableLogOperatorHA" : true,
      "enableLogReport" : true,
      "featureGates" : ["sendAOMAllEvent", "outputKafka", "podLabelExclude", "fullPathCollection", "independentEvents"],
      "host_network" : false,
      "ltsAccessEndpoint" : "https://lts-access.${local.region_name}.myhuaweicloud.com:8102",
      "ltsAuditStreamID" : "",
      "ltsEndpoint" : "https://lts.${local.region_name}.myhuaweicloud.com",
      "ltsEnterpriseProjectID" : local.eps_id,
      "ltsEventStreamID" : "",
      "ltsGroupID" : "",
      "ltsKubeApiserverStreamID" : "",
      "ltsKubeControllerManagerStreamID" : "",
      "ltsKubeSchedulerStreamID" : "",
      "ltsLogReportDomain" : "",
      "ltsPrivateEndpointIP" : "",
      "ltsStdoutStreamID" : "",
      "maxEventAgeSeconds" : "",
      "memBufLimit" : "40mb",
      "multiAZEnabled" : false,
      "otelReportLogs" : true,
      "paasakskEnable" : true,
      "podDisruptionBudget" : {"create" : true, "maxUnavailable" : 1},
      "projectID" : "",
      "secretKey" : "",
      "securityToken" : "",
      "serverCert" : "",
      "serverKey" : ""
    })
  }
}
`, testAccAutopilotAddon_base(rName, epsId))
}
