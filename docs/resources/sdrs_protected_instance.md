---
subcategory: "Storage Disaster Recovery Service (SDRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sdrs_protected_instance"
description: ""
---

# huaweicloud_sdrs_protected_instance

Manages an SDRS protected instance resource within HuaweiCloud.

## Example Usage

```hcl
variable "protection_group_id" {}
variable "server_id" {}

resource "huaweicloud_sdrs_protected_instance" "test" {
  name                 = "test-name"
  group_id             = var.protection_group_id
  server_id            = var.server_id
  delete_target_server = true
  delete_target_eip    = true
  description          = "test description"

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of a protected instance. The name can contain a maximum of 64
  characters. The value can contain only letters (a to z and A to Z), digits (0 to 9), dots (.),
  underscores (_), and hyphens (-).

* `group_id` - (Required, String, ForceNew) Specifies the ID of the protection group where a protected instance is
  added.

  Changing this parameter will create a new resource.

* `server_id` - (Required, String, ForceNew) Specifies the ID of the production site server. The server and protection
  group must belong to the same VPC. And each server can only belong to one protection group.
  
  The protection instance has limited requirements for the operating system version,
  you can refer to this [document](https://support.huaweicloud.com/intl/en-us/productdesc-sdrs/sdrs_pro_0007.html).

  Changing this parameter will create a new resource.

* `cluster_id` - (Optional, String, ForceNew) Specifies the DSS storage pool ID.
  This parameter needs to be specified if the disaster recovery site disk uses distributed storage.

  Changing this parameter will create a new resource.

* `primary_subnet_id` - (Optional, String, ForceNew) Specifies the network ID of the subnet for the primary NIC on the
  DR site server.

  Changing this parameter will create a new resource.

* `primary_ip_address` - (Optional, String, ForceNew) Specifies the IP address of the primary NIC on the DR site server.
  This parameter is valid only when primary_subnet_id is specified. If this parameter is not specified when
  primary_subnet_id is specified, the system automatically assigns an IP address to the primary NIC on the DR site
  server.

  Changing this parameter will create a new resource.

* `delete_target_server` - (Optional, Bool) Specifies whether to delete the DR site server.
  The default value is **false**.

* `delete_target_eip` - (Optional, Bool) Specifies whether to delete the EIP of the DR site server.
  The default value is **false**.

* `description` - (Optional, String, ForceNew) Specifies the description of a protected instance. The description can
  contain a maximum of 64 characters. The value angle brackets (<) and (>) are not allowed.

  Changing this parameter will create a new resource.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the SDRS protected instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `target_server` - ID of the target server.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The SDRS protected instance can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_sdrs_protected_instance.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include:
`cluster_id`, `primary_subnet_id`, `primary_ip_address`, `delete_target_server` and `delete_target_eip`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_sdrs_protected_instance" "test" {
  ...
  
  lifecycle {
    ignore_changes = [
      cluster_id,
      primary_subnet_id,
      primary_ip_address,
      delete_target_server,
      delete_target_eip,
    ]
  }
}
```
