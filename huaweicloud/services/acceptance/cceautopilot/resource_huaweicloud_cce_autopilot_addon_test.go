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
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAddon_basic(rName),
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
func testAccAddon_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

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
      "swr_user" : "autopilot-official",
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
      "caCert" : "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURHVENDQWdHZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREF0TVJjd0ZRWURWUVFLRXc1RFEwVWcKVkdWamFHNXZiRzluZVRFU01CQUdBMVVFQXhNSmJHOW5MV0ZuWlc1ME1DQVhEVEkwTVRFeU9EQTNNekkwTWxvWQpEekl3TlRReE1USXhNRGN6TWpReVdqQXRNUmN3RlFZRFZRUUtFdzVEUTBVZ1ZHVmphRzV2Ykc5bmVURVNNQkFHCkExVUVBeE1KYkc5bkxXRm5aVzUwTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUEKMEpyRkwzc3ZVbnlGRVFCRGpqS1NkNFRmLzZCamliQytGWmdEa25IZDZ4U3VJS29EN1BNM3dhamF5M3ZXR2pzdgpKRzJPcTVxZWZnc0JRZ3RJVjJDVVc5amdFTm5vdjVNakJ6YVpPSExNdEFQc21tc29EYTJmdCtMeVUzUkRBTHlBCmhWc05rQmdNcWozeFZLUjk2T0ZxVWFyWXQ5b3N5dXRhNFc3MGE5MVRJdVZBS3NJd0V4V3pSZ3N0Q2wyaFVQN3kKQlUvb0M1SVJhSTloNUlmTnl5YXdpMHpaZ1pLbjBuTXo0Ukw3disvM3lmczZYc1pNektxbkFXRWtsWjZ0S0oxSgpEQk9VN3NJTkd6MS9JOWtMZlZzZzFzT01aSUYrR1ZnT2ZjZFdhWi9VUEZ5MFJGWFNaZEwraDQ3V0pwWnRZZzQyCkxEVkZ0SDAzMFozR2R5N205ZW5ZblFJREFRQUJvMEl3UURBT0JnTlZIUThCQWY4RUJBTUNBcVF3RHdZRFZSMFQKQVFIL0JBVXdBd0VCL3pBZEJnTlZIUTRFRmdRVWg4K3dyazhuRjNRc1RtZ0g0Z3NXUEp0TkNHSXdEUVlKS29aSQpodmNOQVFFTEJRQURnZ0VCQUF5ZDZRb2dQVU5MM2xKemh6UWFXRW5kamZ1NW4rYURPSGh4YXBudHVQaEN0RlNMCnpJRkhiYWtYOE1OSFdLaDN4Y1JBd3hzSWw2TDIvT1FCczVCZ01Rc2pGbnpZL2pydEpYTXVqYmpWUHZVeG5lMHkKQm5qYjBreFNPYTBISVUzbXJ3Q1Jhc1JYMW9XS1NkeWM3NTRsVnI0NG05NzE1bEVsR2pCOHhicWZQOW5HSEdQbwpmN0pvaDZqWTNTOWJqSWx6UkRBaGZBRWJ6TzBxcDdOalBOMHFqU1NJNm9DbnFOWnVxQWdTN3VjWkNhTUVjRnVGCkhNalh1VWZEeFVDcS9adDFNdnRGYmEyeVBrd3V2QVZCYTFtKzBjYVAxeW5MWDRKbC9yUnNDV25XQW9IYmZFZGQKWGNtbnR4ZEJ4UmlUTTlyL1NuUEJHcTUwZVBzUE9vMkh3RmFUN1Q0PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==",
      "clusterID" : "a6e67a56-ad2f-11ef-bc98-0255ac1005c8",
      "clusterName" : "zhangjishu-test",
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
      "projectID" : "0970dd7a1300f5672ff2c003c60ae115",
      "secretKey" : "",
      "securityToken" : "",
      "serverCert" : "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURUekNDQWplZ0F3SUJBZ0lJWlBDU25MdnAvZ0F3RFFZSktvWklodmNOQVFFTEJRQXdMVEVYTUJVR0ExVUUKQ2hNT1EwTkZJRlJsWTJodWIyeHZaM2t4RWpBUUJnTlZCQU1UQ1d4dlp5MWhaMlZ1ZERBZ0Z3MHlOREV4TWpndwpOek15TkRKYUdBOHlNRFUwTVRFeU1UQTNNekkwTWxvd0x6RVpNQmNHQTFVRUNoTVFRME5GSUZSbFkyaHViMnh2CloybGxjekVTTUJBR0ExVUVBeE1KYkc5bkxXRm5aVzUwTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEEKTUlJQkNnS0NBUUVBMEpyRkwzc3ZVbnlGRVFCRGpqS1NkNFRmLzZCamliQytGWmdEa25IZDZ4U3VJS29EN1BNMwp3YWpheTN2V0dqc3ZKRzJPcTVxZWZnc0JRZ3RJVjJDVVc5amdFTm5vdjVNakJ6YVpPSExNdEFQc21tc29EYTJmCnQrTHlVM1JEQUx5QWhWc05rQmdNcWozeFZLUjk2T0ZxVWFyWXQ5b3N5dXRhNFc3MGE5MVRJdVZBS3NJd0V4V3oKUmdzdENsMmhVUDd5QlUvb0M1SVJhSTloNUlmTnl5YXdpMHpaZ1pLbjBuTXo0Ukw3disvM3lmczZYc1pNektxbgpBV0VrbFo2dEtKMUpEQk9VN3NJTkd6MS9JOWtMZlZzZzFzT01aSUYrR1ZnT2ZjZFdhWi9VUEZ5MFJGWFNaZEwrCmg0N1dKcFp0WWc0MkxEVkZ0SDAzMFozR2R5N205ZW5ZblFJREFRQUJvMjh3YlRBT0JnTlZIUThCQWY4RUJBTUMKQmFBd0hRWURWUjBsQkJZd0ZBWUlLd1lCQlFVSEF3SUdDQ3NHQVFVRkJ3TUJNQjhHQTFVZEl3UVlNQmFBRklmUApzSzVQSnhkMExFNW9CK0lMRmp5YlRRaGlNQnNHQTFVZEVRUVVNQktDRUNvdWJXOXVhWFJ2Y21sdVp5NXpkbU13CkRRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFKcDlSbUNpbjVUdmYvZEZVaU5PcUc3ajU3TUJpY2xuYndkaE1ISk8KRmdXUjFCYzRXNmlDWTMwU0tEd01sUVN1TzN0OExrNWFuN0x5WElhZWR2K1ZwdlNPaS9PdkcxNksvL0FxMFVxdgpPTnp1TDJHRHVuLzV6ekpXVUV4RmNCbEdicnhJWW5YYmV4MXMwek93YnVyNElIc0c1WC9pUkVKVkZPaFRsZjNUCjdWT2ZoOHU2NXF0MURMUkdwalUwSGNwMHRqRzd3cDRTUHBTNXVCS0t0R2NjdUNDTkFpVnBOUTJPUHhSZmVBbWMKL1B0cWFwQUFORE9ma2xlOGRxZk1FYkVpS0ZzcC8yOWFQR3VUYmlPaVRDNS82WkFFWjZSN3JkaS8raVM2dlRRNApIQVlyc2w0MGJjVWEraWNjV1paTG9lVGRRUHRUcHZTdkRWejdhSTEzL21Oa3Z5bz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=",
      "serverKey" : "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcFFJQkFBS0NBUUVBMEpyRkwzc3ZVbnlGRVFCRGpqS1NkNFRmLzZCamliQytGWmdEa25IZDZ4U3VJS29ECjdQTTN3YWpheTN2V0dqc3ZKRzJPcTVxZWZnc0JRZ3RJVjJDVVc5amdFTm5vdjVNakJ6YVpPSExNdEFQc21tc28KRGEyZnQrTHlVM1JEQUx5QWhWc05rQmdNcWozeFZLUjk2T0ZxVWFyWXQ5b3N5dXRhNFc3MGE5MVRJdVZBS3NJdwpFeFd6UmdzdENsMmhVUDd5QlUvb0M1SVJhSTloNUlmTnl5YXdpMHpaZ1pLbjBuTXo0Ukw3disvM3lmczZYc1pNCnpLcW5BV0VrbFo2dEtKMUpEQk9VN3NJTkd6MS9JOWtMZlZzZzFzT01aSUYrR1ZnT2ZjZFdhWi9VUEZ5MFJGWFMKWmRMK2g0N1dKcFp0WWc0MkxEVkZ0SDAzMFozR2R5N205ZW5ZblFJREFRQUJBb0lCQUYzSEZob3dVS2ZPWHF1ego2S3JHUlY0Qm1BbDgrd0p0T0NiUS9kb1o0bC9LSGpXRStOck94Q1FGV3NiYlZ2Ylg0R3VKN1Bkc1BSQUF0b0lRClBHYzdmYmFFbXNZNGtBOS9mK0hBUThWQ3BvL09xOUVIbHl2Ky82eFZGQWM4WHRxMzR6Y1FKZHEvVlFJN2NvQlEKcW1IRTVGenVaeHJQdEE5TkdyLzVkMXYrVlY5N3pBeVg3cU51WG0xNVZoVVFRaEdSYjhMcG1CTU85SGFXU3IvVQo0VGtEbXVYelpPYWRYd2pqWmRCbHM1MlJvSXBuL3J0VitMVEY3L2V6NDRjakV1aTNkdDFkOGErUUJ4S2xPcXBBCjJQQjNwZlZPbElzdXBZTGlRdnZRcXNKQVoxN0NBM3Q0Wlp5UnYybCtHVGpHV3hpUFZWaXpkSWpKUWxjdjNhRGcKazJLWFhha0NnWUVBOGV4SDNWamFuUjB1OGppU044Y29OOXVoYnJDUTNzdWRrUUhBNE5hc3ZVK0dYbCtiTHF0NAp0cktEU3lHb2MyUDZIRG9LUXM2Mjd2TTloRStWVDFETFkyZjNrWmxaRXRBSlpYUmR4Vm90Z2tld2NJYWcybWZmCk8rcjV3eUtJOEovYUlKTlgzUEF3eEJSR3FhZ2s2dU1aZWJCeUtTbXVDLzhvWVJSWTN4U3krazhDZ1lFQTNMNHMKZ24wVnl6NDRaQTRHWGZkS3JtQ1YwckVXZ2ZkbzJJZis5aVluWG44UHhud29aWEZuOHV2NHhzZ3ZZTFArWExHaAp2c29mMFpaTXJpclZsQ3RqbFViSDVJK1lYT1I2UjZ4SnhTbFhwMkQzSFJwR01UbnlKRWtmT3lUdVVHTzVHNlY5CkZ0VzVWM1laK2QwQW80RzdtTi8rU3BOR09HSDcrK2tkMWZFUy8xTUNnWUVBNUVMWW05VVdrRi9VeDk3d3Q0aEcKUGs0UXgydjVoUDRCc2F4QjNPTXhJWDVEZmhBZlQ2Mml2Rjg2Mmt6cnI5U0pUTkRHbGJxTmlIQWhmeEhJQTRwcwpIV01maUZWMFlmZkFwZVZpQksvTmVMdERreWl6NU45VkZpZmplV2JBWnFtdEdrZHNBNTd0cEZTdFI2N0xCb1U0CnFFVC9zaThOZFd4UElTb2RvSDdiVUtrQ2dZRUFsWDYxNWltUWVQVEtlL2lEbDEvQzFCWFZZYnRNNHZnTHFabHcKc29Oa1pqcm5GQ1ZCdG5IM1ZDMDdibVJrc2JrMHF0SWlHSFFLMklaUnFDS2FRcDZmOHBqZEI0MjRRakQ2SDFBdgpKYmU2QlVGR0dnK1JPZ1ZrVis2dG1BQ0s1U2FrVm5UZElubmI2Nyt3RitmMFpzZVZwUk1OeExPNCtyWmhVVm12Ck94VHBLTUVDZ1lFQXk0SzBUcjlDR1lPV1ZUU24wVHRZOTAvU2RxSVNmQ2JNVkxGU1Zna1FiSHN3T3hHT1JML2YKNmdZM1M3cDBaYkxnRXozalgxdkF1aEVSalBKdFl6U1dJWjRKcXk3VlNZNGIwWk56OStublI1aHowUVB3a1AwZApWbHRhYkQyY2F2MnNnSWNNK0R4YVVkbTV2ZE5uOWtXYWVKWVhBS3h3dDRHYVpxdXNuVmZHMGxnPQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo="
    })
  }
}
`, testAccCluster_basic(rName), rName)
}
