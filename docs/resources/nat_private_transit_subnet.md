---
subcategory: "NAT Gateway (NAT)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_nat_private_transit_subnet"
description: |-
  Manages a transit subnet resource of the **private** NAT within HuaweiCloud.
---

# huaweicloud_nat_private_transit_subnet

Manages a transit subnet resource of the **private** NAT within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "description" {}
variable "virsubnet_id" {}
variable "virsubnet_project_id" {}

resource "huaweicloud_nat_private_transit_subnet" "test" {
  name                 = var.name
  description          = var.description
  virsubnet_id         = var.virsubnet_id
  virsubnet_project_id = var.virsubnet_project_id

  tags = {
    foo = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the transit subnet is located.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `name` - (Required, String) Specifies the transit subnet name.
  The name can contain only digits, letters, underscores (_), and hyphens (-).

* `description` - (Optional, String) Provides supplementary information about the transit subnet.
  The description can contain up to 255 characters and cannot contain angle brackets (<>).

* `virsubnet_id` - (Required, String, ForceNew) Specifies the transit subnet ID.
  Changing this will create a new resource.

* `virsubnet_project_id` - (Required, String, ForceNew) Specifies the ID of the project to which
  the transit subnet belongs.
  Only digits and lowercase letters are supported.
  Changing this will create a new resource.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the transit subnet.
  The [tags](#transit_subnet_tags) structure is documented below.

<a name="transit_subnet_tags"></a>
The `tags` block supports:

* `key` - The key of the resource tag.

  -> The key of tags has limits as follows:
  <br/>1. It can contain a maximum of `36` characters.
  <br/>2. It cannot be an empty string.
  <br/>3. Spaces before and after a key will be discarded.
  <br/>4. It cannot contain non-printable ASCII characters (`0`–`31`) and the following characters: =*<>,|/
  <br/>5. It can contain only letters, digits, hyphens (-), and underscores (_).

* `value` - The value of the resource tag.

  -> The value of tags has limits as follows:
  <br/>1. It is mandatory when a tag is added and optional when a tag is deleted.
  <br/>2. It can contain a maximum of `43` characters.
  <br/>3. It can be an empty string.
  <br/>4. Spaces before and after a value will be discarded.
  <br/>5. It cannot contain non-printable ASCII characters (`0`–`31`) and the following characters: =*<>,|/
  <br/>6. It can contain only letters, digits, hyphens (-), underscores (_), and periods (.).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `project_id` - The project ID.

* `vpc_id` - ID of the VPC to which the transit subnet belongs.

* `cidr` - CIDR block of the transit subnet.

* `type` - transit subnet type. The value can only be VPC.

* `status` - transit subnet status. The value can be ACTIVE, indicating the transit subnet is normal.

* `ip_count` - The number of IP addresses that has been assigned from the transit subnet.

* `status` - The status of the transit subnet.

* `created_at` - The creation time of the transit subnet for private NAT.

* `updated_at` - The latest update time of the transit subnet for private NAT.

## Import

The resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_nat_private_transit_subnet.test <id>
```
