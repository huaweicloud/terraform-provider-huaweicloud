---
subcategory: Dedicated Load Balance (Dedicated ELB)
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_ipgroups"
description: |-
  Use this data source to get the list of ELB IP groups.
---

# huaweicloud_elb_ipgroups

Use this data source to get the list of ELB IP groups.

## Example Usage

```hcl
variable "ipgroup_name" {}

data "huaweicloud_elb_ipgroups" "test" {
  name = var.ipgroup_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the IP address group.

* `ipgroup_id` - (Optional, String) Specifies the ID of the IP address group.

* `description` - (Optional, String) Specifies the description of the IP address group.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `ip_address` - (Optional, String) Specifies the IP address of the IP address group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `ipgroups` - Lists the IP groups.
  The [ipgroups](#Elb_ipgroups) structure is documented below.

<a name="Elb_ipgroups"></a>
The `ipgroups` block supports:

* `id` - The ID of the IP address group.

* `name` - The name of the IP address group.

* `description` - The description of the IP address group.

* `enterprise_project_id` - The enterprise project ID.

* `project_id` - The project ID of the IP address group.

* `listeners` - The IDs of listeners with which the IP address group is associated. The [listeners](#Elb_ipgroups_listeners)
  structure is documented below.

* `ip_list` - The IP addresses or CIDR blocks in the IP address group. The [ip_list](#Elb_ipgroups_ip_list) structure is
  documented below.

* `created_at` - The time when the IP address group was created.

* `updated_at` - The time when the IP address group was updated.

<a name="Elb_ipgroups_listeners"></a>
The `listeners` block supports:

* `id` - The listener ID.

<a name="Elb_ipgroups_ip_list"></a>
The `ip_list` block supports:

* `ip` - The IP addresses.

* `description` - The description of the IP address group.
