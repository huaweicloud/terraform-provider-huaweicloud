---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_template_versions"
description: |-
  Use this data source to get the list of template versions.
---

# huaweicloud_compute_template_versions

Use this data source to get the list of template versions.

## Example Usage

```hcl
data "huaweicloud_compute_template_versions" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `launch_template_id` - (Optional, String) Specifies the template ID.

* `flavor_id` - (Optional, String) Specifies the flavor ID of the template.

* `version` - (Optional, List) Specifies the template versions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `launch_template_versions` - Indicates the list of template versions.

  The [launch_template_versions](#launch_template_versions_struct) structure is documented below.

<a name="launch_template_versions_struct"></a>
The `launch_template_versions` block supports:

* `launch_template_id` - Indicates the template ID.

* `version_id` - Indicates the template version ID.

* `version_number` - Indicates the template version.

* `version_description` - Indicates the template version description.

* `template_data` - Indicates the data info of the template.

  The [template_data](#launch_template_versions_template_data_struct) structure is documented below.

* `created_at` - Indicates the time when the template version was created.

<a name="launch_template_versions_template_data_struct"></a>
The `template_data` block supports:

* `flavor_id` - Indicates the flavor ID of the ECS created based on the template.

* `name` - Indicates the template name.

* `description` - Indicates the template description.

* `availability_zone_id` - Indicates the AZ of the template.

* `enterprise_project_id` - Indicates the enterprise project ID of the template.

* `auto_recovery` - Indicates whether enable auto-recovery.

* `os_profile` - Indicates the OS attributes.

  The [os_profile](#template_data_os_profile_struct) structure is documented below.

* `security_group_ids` - Indicates the security group ID list.

* `network_interfaces` - Indicates the network interfaces.

  The [network_interfaces](#template_data_network_interfaces_struct) structure is documented below.

* `block_device_mappings` - Indicates the BDM mounting information.

  The [block_device_mappings](#template_data_block_device_mappings_struct) structure is documented below.

* `market_options` - Indicates the billing information.

  The [market_options](#template_data_market_options_struct) structure is documented below.

* `internet_access` - Indicates the public network access.

  The [internet_access](#template_data_internet_access_struct) structure is documented below.

* `metadata` - Indicates the metadata.

* `tag_options` - Indicates the VM tags.
  Currently, only VMs can be tagged. In the future, associated resources such as volumes can be tagged, too.

  The [tag_options](#template_data_tag_options_struct) structure is documented below.

<a name="template_data_os_profile_struct"></a>
The `os_profile` block supports:

* `key_name` - Indicates the key name.

* `user_data` - Indicates the custom user data to be injected into the instance during instance creation.
  Text and text files can be injected.

* `iam_agency_name` - Indicates the agency name.

* `enable_monitoring_service` - Indicates whether enable HSS.

<a name="template_data_network_interfaces_struct"></a>
The `network_interfaces` block supports:

* `virsubnet_id` - Indicates the subnet ID.

* `attachment` - Indicates the network interface details.

  The [attachment](#network_interfaces_attachment_struct) structure is documented below.

<a name="network_interfaces_attachment_struct"></a>
The `attachment` block supports:

* `device_index` - Indicates the loading sequence. The value 0 indicates the primary network interface.

<a name="template_data_block_device_mappings_struct"></a>
The `block_device_mappings` block supports:

* `volume_type` - Indicates the volume type.

* `volume_size` - Indicates the volume size.

* `attachment` - Indicates the disk interface

  The [attachment](#block_device_mappings_attachment_struct) structure is documented below.

* `source_id` - Indicates the data source type of the ECS volume.

* `source_type` - Indicates the source type of the volume device.

* `encrypted` - Indicates the encrypted or not.

* `cmk_id` - Indicates the key ID.

<a name="block_device_mappings_attachment_struct"></a>
The `attachment` block supports:

* `boot_index` - Indicates the loading sequence.

* `delete_on_termination` - Indicates whether the disk is released along with the instance.

<a name="template_data_market_options_struct"></a>
The `market_options` block supports:

* `market_type` - Indicates the billing mode.

* `spot_options` - Indicates spot instance parameters.

  The [spot_options](#market_options_spot_options_struct) structure is documented below.

<a name="market_options_spot_options_struct"></a>
The `spot_options` block supports:

* `spot_price` - Indicates  the highest price per hour you are willing to pay for a spot ECS.

* `block_duration_minutes` - Indicates the predefined duration of the spot ECS.

* `instance_interruption_behavior` - Indicates the spot ECS interruption policy, which can only be set to **immediate**
  currently.

<a name="template_data_internet_access_struct"></a>
The `internet_access` block supports:

* `publicip` - Indicates the public network access.

  The [publicip](#internet_access_publicip_struct) structure is documented below.

<a name="internet_access_publicip_struct"></a>
The `publicip` block supports:

* `publicip_type` - Indicates the EIP type.

* `charging_mode` - Indicates the EIP billing mode.

* `bandwidth` - Indicates the EIP bandwidth.

  The [bandwidth](#publicip_bandwidth_struct) structure is documented below.

* `delete_on_termination` - Indicates whether the EIP is released along with the instance.

<a name="publicip_bandwidth_struct"></a>
The `bandwidth` block supports:

* `share_type` - Indicates the bandwidth type.

* `size` - Indicates the bandwidth size.

* `charge_mode` - Indicates the billing mode.

* `id` - Indicates the bandwidth ID.
  You can use an existing shared bandwidth when applying for an EIP for the bandwidth of type **WHOLE**.

<a name="template_data_tag_options_struct"></a>
The `tag_options` block supports:

* `tags` - Indicates the tags.

  The [tags](#tag_options_tags_struct) structure is documented below.

<a name="tag_options_tags_struct"></a>
The `tags` block supports:

* `value` - Indicates the tag key.

* `key` - Indicates the tag value.
