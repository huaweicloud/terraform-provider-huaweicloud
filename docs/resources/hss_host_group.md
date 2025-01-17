---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_host_group"
description: |-
  Manages an HSS host group resource within HuaweiCloud.
---

# huaweicloud_hss_host_group

Manages an HSS host group resource within HuaweiCloud.

## Example Usage

### Create an HSS host group and bind some ECS instances

```hcl
variable "host_group_name" {}
variable "host_ids" {
  type = list(string)
}

resource "huaweicloud_hss_host_group" "test" {
  name     = var.host_group_name
  host_ids = var.host_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region where the host group is located.
  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the host group.  
  The valid length is limited from `1` to `64`, only Chinese characters, English letters, digits, hyphens (-),
  underscores (_), dots (.), pluses (+) and asterisks (*) are allowed.  
  The Chinese characters must be in **UTF-8** or **Unicode** format.

* `host_ids` - (Required, List) Specifies the list of host IDs.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the ID of the enterprise project to which the host
  group belongs.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `host_num` - The total host number.

* `risk_host_num` - The number of hosts at risk.

* `unprotect_host_num` - The number of unprotect hosts.

* `unprotect_host_ids` - The ID list of the unprotect hosts.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.

## Import

The host group resource can be imported using `enterprise_project_id` and `id`, separated by a slash, e.g.

### Import resource under the default enterprise project

```bash
$ terraform import huaweicloud_hss_host_group.test 0/<id>
```

### Import resource from non default enterprise project

```bash
$ terraform import huaweicloud_hss_host_group.test <enterprise_project_id>/<id>
```
