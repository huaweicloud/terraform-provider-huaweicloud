---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_waf_access"
description: |-
  Manages an LTS access WAF logs configuration resource within HuaweiCloud.
---

# huaweicloud_lts_waf_access

Manages an LTS access WAF logs configuration resource within HuaweiCloud.

-> **NOTE:** This resource depends on WAF instances, the instance can be Cloud Mode or Dedicated Mode.

## Example Usage

```hcl
variable "lts_group_id" {}
variable "lts_attack_stream_id" {}
variable "lts_access_stream_id" {}
variable "enterprise_project_id" {}

resource "huaweicloud_lts_waf_access" "test" {
  lts_group_id          = var.lts_group_id
  lts_attack_stream_id  = var.lts_attack_stream_id
  lts_access_stream_id  = var.lts_access_stream_id
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `lts_group_id` - (Required, String) Specifies the log group ID.

* `lts_attack_stream_id` - (Optional, String) Specifies the log stream ID for attack logs.

* `lts_access_stream_id` - (Optional, String) Specifies the log stream ID for access logs.

-> The fields `lts_attack_stream_id` and `lts_access_stream_id` must be specified as different log streams.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID.  
  This parameter is only valid for enterprise users. If not specified, the default enterprise project will be used.
  The default enterprise project ID is **0**.  

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

For enterprise users, The resource can be imported using the `enterprise_project_id`, e.g.

```bash
$ terraform import huaweicloud_lts_waf_access.test <enterprise_project_id>
```

For non-enterprise users, The resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_lts_waf_access.test <id>
```
