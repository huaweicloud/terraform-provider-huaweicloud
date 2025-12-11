---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_template"
description: |-
  Manages an ECS template resource within HuaweiCloud.
---

# huaweicloud_compute_template

Manages an ECS template resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_compute_template" "test" {
  name = "test_template_name"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the name of the template.

* `template_data` - (Optional, List, NonUpdatable) Specifies the data info of the template.
  The [template_data](#template_data_struct) structure is documented below.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the template.

* `version_description` - (Optional, String, NonUpdatable)  Specifies the version description of the template.

<a name="template_data_struct"></a>
The `template_data` block supports:

* `flavor_id` - (Optional, String, NonUpdatable) Specifies the flavor ID of the ECS to be created based on the template.

* `name` - (Optional, String, NonUpdatable) Specifies the ECS name.

* `description` - (Optional, String, NonUpdatable) Specifies the ECS description.

* `availability_zone_id` - (Optional, String, NonUpdatable) Specifies the AZ.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.

* `auto_recovery` - (Optional, Bool, NonUpdatable) Specifies whether enable auto-recovery.

* `os_profile` - (Optional, List, NonUpdatable) Specifies the image attribute.
  The [os_profile](#os_profile_struct) structure is documented below.

* `security_group_ids` - (Optional, List, NonUpdatable) Specifies the security group ID list.

* `network_interfaces` - (Optional, List, NonUpdatable) Specifies the network interfaces.
  The [network_interfaces](#network_interfaces_struct) structure is documented below.

* `block_device_mappings` - (Optional, List, NonUpdatable) Specifies the BDM mounting information.
  The [block_device_mappings](#block_device_mappings_struct) structure is documented below.

* `market_options` - (Optional, List, NonUpdatable) Specifies the billing information.
  The [market_options](#market_options_struct) structure is documented below.

* `internet_access` - (Optional, List, NonUpdatable) Specifies the public network access.
  The [internet_access](#internet_access_struct) structure is documented below.

* `metadata` - (Optional, Map, NonUpdatable) Specifies the metadata.

* `tag_options` - (Optional, List, NonUpdatable) Specifies the VM tags. Currently, only VMs can be tagged. In the future,
  associated resources such as volumes can be tagged, too.
  The [tag_options](#tag_options_struct) structure is documented below.

<a name="os_profile_struct"></a>
The `os_profile` block supports:

* `key_name` - (Optional, String, NonUpdatable) Specifies the key name.

* `user_data` - (Optional, String, NonUpdatable) Specifies the custom user data to be injected into the instance during
  instance creation. Text and text files can be injected.

* `iam_agency_name` - (Optional, String, NonUpdatable) Specifies the agency name.

* `enable_monitoring_service` - (Optional, Bool, NonUpdatable) Specifies whether enable HSS.

<a name="network_interfaces_struct"></a>
The `network_interfaces` block supports:

* `virsubnet_id` - (Optional, String, NonUpdatable) Specifies the subnet ID.

* `attachment` - (Optional, List, NonUpdatable) Specifies the network interface details.
  The [network_interfaces_attachment](#network_interfaces_attachment_struct) structure is documented below.

<a name="network_interfaces_attachment_struct"></a>
The `network_interfaces_attachment` block supports:

* `device_index` - (Optional, Int, NonUpdatable) Specifies the loading sequence. The value **0** indicates the primary
  network interface.

<a name="block_device_mappings_struct"></a>
The `block_device_mappings` block supports:

* `source_id` - (Optional, String, NonUpdatable) Specifies the VM volume data source type.

* `source_type` - (Optional, String, NonUpdatable) Specifies the source type of the volume device.

* `encrypted` - (Optional, Bool, NonUpdatable) Specifies the encrypted or not.

* `cmk_id` - (Optional, String, NonUpdatable) Specifies the key ID.

* `volume_type` - (Optional, String, NonUpdatable) Specifies the volume type.

* `volume_size` - (Optional, Int, NonUpdatable) Specifies the volume size.

* `attachment` - (Optional, List, NonUpdatable) Specifies the disk interface.
  The [block_device_mappings_attachment](#block_device_mappings_attachment_struct) structure is documented below.

<a name="block_device_mappings_attachment_struct"></a>
The `block_device_mappings_attachment` block supports:

* `boot_index` - (Optional, Int, NonUpdatable) Specifies the loading sequence. The value **0** indicates the system disk.

* `delete_on_termination` - (Optional, Bool, NonUpdatable) Specifies whether the disk is released along with the instance.

<a name="market_options_struct"></a>
The `market_options` block supports:

* `market_type` - (Optional, String, NonUpdatable) Specifies the billing mode. Value options: **postpaid**, **spot**.

* `spot_options` - (Optional, List, NonUpdatable) Specifies the spot instance parameters.
  The [spot_options](#spot_options_struct) structure is documented below.

<a name="spot_options_struct"></a>
The `spot_options` block supports:

* `spot_price` - (Optional, Float, NonUpdatable) Specifies the highest price per hour you are willing to pay for a spot
 ECS.

* `block_duration_minutes` - (Optional, Int, NonUpdatable) Specifies the predefined duration of the spot ECS.

* `instance_interruption_behavior` - (Optional, String, NonUpdatable) Specifies the spot ECS interruption policy, which
  can only be set to **immediate** currently.

<a name="internet_access_struct"></a>
The `internet_access` block supports:

* `publicip` - (Optional, List, NonUpdatable) Specifies the public network access.
  The [publicip](#publicip_struct) structure is documented below.

<a name="publicip_struct"></a>
The `publicip` block supports:

* `publicip_type` - (Optional, String, NonUpdatable) Specifies the EIP type.

* `charging_mode` - (Optional, String, NonUpdatable) Specifies the EIP billing mode.

* `bandwidth` - (Optional, List, NonUpdatable) Specifies the bandwidth.
  The [bandwidth](#bandwidth_struct) structure is documented below.

* `delete_on_termination` - (Optional, Bool, NonUpdatable) Specifies whether the EIP is released along with the instance.

<a name="bandwidth_struct"></a>
The `bandwidth` block supports:

* `share_type` - (Optional, String, NonUpdatable) Specifies the bandwidth type.

* `size` - (Optional, Int, NonUpdatable) Specifies the bandwidth size.

* `charge_mode` - (Optional, String, NonUpdatable) Specifies the billing mode.

* `id` - (Optional, String, NonUpdatable) Specifies the bandwidth ID. You can use an existing shared bandwidth when
  applying for an EIP for the bandwidth of type **WHOLE**.

<a name="tag_options_struct"></a>
The `tag_options` block supports:

* `tags` - (Optional, List, NonUpdatable) Specifies the tags.
  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - (Optional, String, NonUpdatable) Specifies the tag key.

* `value` - (Optional, String, NonUpdatable) Specifies the tag value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `default_version` - Indicates the default version of the template.

* `latest_version` - Indicates the latest version of the template.

* `version_id` - Indicates the template version ID.

* `created_at` - Indicates the time when the template version was created.

## Import

The ECS template resource can be imported using the `id`, e.g.

```shell
$ terraform import huaweicloud_compute_template.test <id>
```
