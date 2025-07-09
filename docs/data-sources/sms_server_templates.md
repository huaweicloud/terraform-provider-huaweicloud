---
subcategory: "Server Migration Service (SMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sms_server_templates"
description: |-
  Use this data source to get the list of SMS templates used for setting up target servers.
---

# huaweicloud_sms_server_templates

Use this data source to get the list of SMS templates used for setting up target servers.

## Example Usage

```hcl
data "huaweicloud_sms_server_templates" "test" {}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Specifies the template name.

* `availability_zone` - (Optional, String) Specifies the availability zone.

* `region` - (Optional, String) Specifies the region ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `templates` - Indicates the template information.

  The [templates](#templates_struct) structure is documented below.

<a name="templates_struct"></a>
The `templates` block supports:

* `id` - Indicates the template ID.

* `name` - Indicates the template name.

* `is_template` - Indicates whether the template is general.

* `region` - Indicates the region.

* `project_id` - Indicates the project ID.

* `target_server_name` - Indicates the name of the target server.

* `availability_zone` - Indicates the availability zone.

* `volume_type` - Indicates the disk type.

* `flavor` - Indicates the server flavor.

* `vpc` - Indicates the VPC information.

  The [vpc](#templates_vpc_struct) structure is documented below.

* `nics` - Indicates the NIC information.

  The [nics](#templates_nics_struct) structure is documented below.

* `security_groups` - Indicates the security group information.

  The [security_groups](#templates_security_groups_struct) structure is documented below.

* `publicip` - Indicates the public IP address information.

  The [publicip](#templates_publicip_struct) structure is documented below.

* `disk` - Indicates the disk information.

  The [disk](#templates_disk_struct) structure is documented below.

* `data_volume_type` - Indicates the data disk type.

* `target_password` - Indicates the server login password.

* `image_id` - Indicates the ID of the selected image.

<a name="templates_vpc_struct"></a>
The `vpc` block supports:

* `id` - Indicates the VPC ID.

* `name` - Indicates the VPC name.

* `cidr` - Indicates the VPC CIDR block.

<a name="templates_nics_struct"></a>
The `nics` block supports:

* `id` - Indicates the subnet ID.

* `name` - Indicates the subnet name.

* `cidr` - Indicates the subnet gateway/mask.

* `ip` - Indicates the server IP address.

<a name="templates_security_groups_struct"></a>
The `security_groups` block supports:

* `id` - Indicates the security group ID.

* `name` - Indicates the security group name.

<a name="templates_publicip_struct"></a>
The `publicip` block supports:

* `type` - Indicates the EIP type.

* `bandwidth_size` - Indicates the bandwidth size, the unit is Mbit/s.

* `bandwidth_share_type` - Indicates the bandwidth type.

<a name="templates_disk_struct"></a>
The `disk` block supports:

* `id` - Indicates the disk ID.

* `index` - Indicates the disk serial number.

* `name` - Indicates the disk name.

* `disktype` - Indicates the disk type.

* `size` - Indicates the disk size, the unit is GB.

* `device_use` - Indicates the used disk space.
