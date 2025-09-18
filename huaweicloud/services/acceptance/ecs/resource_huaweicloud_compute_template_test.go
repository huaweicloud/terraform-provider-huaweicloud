package ecs

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

func getTemplateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl        = "v3/{project_id}/launch-templates?launch_template_id={launch_template_id}"
		versionHttpUrl = "v3/{project_id}/launch-template-versions?launch_template_id={launch_template_id}"
		product        = "ecs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating ECS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{launch_template_id}", state.Primary.ID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving ECS template: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}
	template := utils.PathSearch("launch_templates[0]", getRespBody, nil)
	if template == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	getTemplateVersionPath := client.Endpoint + versionHttpUrl
	getTemplateVersionPath = strings.ReplaceAll(getTemplateVersionPath, "{project_id}", client.ProjectID)
	getTemplateVersionPath = strings.ReplaceAll(getTemplateVersionPath, "{launch_template_id}", state.Primary.ID)

	getTemplateVersionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getTemplateVersionResp, err := client.Request("GET", getTemplateVersionPath, &getTemplateVersionOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving ECS template version: %s", err)
	}

	getTemplateVersionRespBody, err := utils.FlattenResponse(getTemplateVersionResp)
	if err != nil {
		return nil, err
	}
	templateVersion := utils.PathSearch("launch_template_versions[0]", getTemplateVersionRespBody, nil)
	if templateVersion == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func TestAccComputeTemplate_Basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_compute_template.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeTemplate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "version_description", "test version description"),
					resource.TestCheckResourceAttr(resourceName, "template_data.0.flavor_id", "test_flavor_id"),
					resource.TestCheckResourceAttr(resourceName, "template_data.0.name", "test_name"),
					resource.TestCheckResourceAttr(resourceName, "template_data.0.description", "test_description"),
					resource.TestCheckResourceAttr(resourceName, "template_data.0.availability_zone_id",
						"test_availability_zone_id"),
					resource.TestCheckResourceAttr(resourceName, "template_data.0.enterprise_project_id",
						"test_enterprise_project_id"),
					resource.TestCheckResourceAttr(resourceName, "template_data.0.auto_recovery", "true"),
					resource.TestCheckResourceAttr(resourceName, "template_data.0.os_profile.0.key_name",
						"test_os_profile_key_name"),
					resource.TestCheckResourceAttr(resourceName, "template_data.0.os_profile.0.user_data",
						"test_os_profile_user_data"),
					resource.TestCheckResourceAttr(resourceName, "template_data.0.os_profile.0.iam_agency_name",
						"test_os_profile_iam_agency_name"),
					resource.TestCheckResourceAttr(resourceName, "template_data.0.os_profile.0.enable_monitoring_service",
						"true"),
					resource.TestCheckResourceAttr(resourceName, "template_data.0.security_group_ids.0",
						"test_security_group_ids"),
					resource.TestCheckResourceAttr(resourceName, "template_data.0.network_interfaces.0.virsubnet_id",
						"test_network_interfaces_virsubnet_id"),
					resource.TestCheckResourceAttr(resourceName,
						"template_data.0.network_interfaces.0.attachment.0.device_index", "1"),
					resource.TestCheckResourceAttr(resourceName, "template_data.0.block_device_mappings.0.source_id",
						"test_block_device_mappings_source_id"),
					resource.TestCheckResourceAttr(resourceName, "template_data.0.block_device_mappings.0.source_type",
						"blank"),
					resource.TestCheckResourceAttr(resourceName, "template_data.0.block_device_mappings.0.encrypted",
						"true"),
					resource.TestCheckResourceAttr(resourceName, "template_data.0.block_device_mappings.0.cmk_id",
						"test_block_device_mappings_cmk_id"),
					resource.TestCheckResourceAttr(resourceName, "template_data.0.block_device_mappings.0.volume_type",
						"SATA"),
					resource.TestCheckResourceAttr(resourceName, "template_data.0.block_device_mappings.0.volume_size",
						"100"),
					resource.TestCheckResourceAttr(resourceName, "template_data.0.market_options.0.market_type",
						"postpaid"),
					resource.TestCheckResourceAttr(resourceName,
						"template_data.0.market_options.0.spot_options.0.spot_price", "3.5"),
					resource.TestCheckResourceAttr(resourceName,
						"template_data.0.market_options.0.spot_options.0.block_duration_minutes", "2"),
					resource.TestCheckResourceAttr(resourceName,
						"template_data.0.market_options.0.spot_options.0.instance_interruption_behavior", "immediate"),
					resource.TestCheckResourceAttr(resourceName,
						"template_data.0.block_device_mappings.0.attachment.0.boot_index", "1"),
					resource.TestCheckResourceAttr(resourceName,
						"template_data.0.internet_access.0.publicip.0.publicip_type", "5_bgp"),
					resource.TestCheckResourceAttr(resourceName,
						"template_data.0.internet_access.0.publicip.0.charging_mode", "postPaid"),
					resource.TestCheckResourceAttr(resourceName,
						"template_data.0.internet_access.0.publicip.0.bandwidth.0.share_type", "PER"),
					resource.TestCheckResourceAttr(resourceName,
						"template_data.0.internet_access.0.publicip.0.bandwidth.0.size", "500"),
					resource.TestCheckResourceAttr(resourceName,
						"template_data.0.internet_access.0.publicip.0.bandwidth.0.charge_mode", "bandwidth"),
					resource.TestCheckResourceAttr(resourceName,
						"template_data.0.internet_access.0.publicip.0.bandwidth.0.id", "internet_access_publicip_bandwidth_id"),
					resource.TestCheckResourceAttr(resourceName, "template_data.0.metadata.aaa", "bbb"),
					resource.TestCheckResourceAttr(resourceName, "template_data.0.metadata.ccc", "ddd"),
					resource.TestCheckResourceAttr(resourceName, "template_data.0.tag_options.0.tags.0.key", "aaa"),
					resource.TestCheckResourceAttr(resourceName, "template_data.0.tag_options.0.tags.0.value", "bbb"),
					resource.TestCheckResourceAttrSet(resourceName, "default_version"),
					resource.TestCheckResourceAttrSet(resourceName, "latest_version"),
					resource.TestCheckResourceAttrSet(resourceName, "version_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeTemplate_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_compute_template" "test" {
  name                = "%s"
  description         = "test description"
  version_description = "test version description"

  template_data {
    flavor_id             = "test_flavor_id"
    name                  = "test_name"
    description           = "test_description"
    availability_zone_id  = "test_availability_zone_id"
    enterprise_project_id = "test_enterprise_project_id"
    auto_recovery         = true

    os_profile {
      key_name                  = "test_os_profile_key_name"
      user_data                 = "test_os_profile_user_data"
      iam_agency_name           = "test_os_profile_iam_agency_name"
      enable_monitoring_service = true
    }

    security_group_ids = ["test_security_group_ids"]

    network_interfaces {
      virsubnet_id = "test_network_interfaces_virsubnet_id"

      attachment {
        device_index = 1
      }
    }


    block_device_mappings {
      source_id   = "test_block_device_mappings_source_id"
      source_type = "blank"
      encrypted   = true
      cmk_id      = "test_block_device_mappings_cmk_id"
      volume_type = "SATA"
      volume_size = 100

      attachment {
        boot_index = 1
      }
    }

    market_options {
      market_type = "postpaid"

      spot_options {
        spot_price                     = 3.5
        block_duration_minutes         = 2
        instance_interruption_behavior = "immediate"
      }
    }

    internet_access {
      publicip {
        publicip_type = "5_bgp"
        charging_mode = "postPaid"

        bandwidth {
          share_type  = "PER"
          size        = 500
          charge_mode = "bandwidth"
          id          = "internet_access_publicip_bandwidth_id"
        }
      }
    }

    metadata = {
      "aaa" = "bbb"
      "ccc" = "ddd"
    }

    tag_options {
      tags {
        key   = "aaa"
        value = "bbb"
      }
    }
  }
}
`, rName)
}
