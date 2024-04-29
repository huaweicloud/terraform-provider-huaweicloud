---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_networking_secgroups"
description: ""
---

# huaweicloud_networking_secgroups

Use this data source to get the list of the available HuaweiCloud security groups.

## Example Usage

### Filter the list of security groups by a description keyword

```hcl
variable "key_word" {}

data "huaweicloud_networking_secgroups" "test" {
  description = var.key_word
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to obtain the security group list.
  If omitted, the provider-level region will be used.

* `id` - (Optional, String) Specifies the id of the desired security group.

* `name` - (Optional, String) Specifies the name of the security group.

* `description` - (Optional, String) Specifies the description of the security group. The security groups can be
  filtered by keywords in the description.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of the security group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `security_groups` - The list of security groups. The [object](#security_groups) is documented below.

<a name="security_groups"></a>
The `security_groups` block supports:

* `id` - The security group ID.

* `description`- The description of the security group.

* `name` - The name of the security group.

* `enterprise_project_id` - The enterprise project ID of the security group.

* `created_at` - The creation time, in UTC format.

* `updated_at` - The last update time, in UTC format.
