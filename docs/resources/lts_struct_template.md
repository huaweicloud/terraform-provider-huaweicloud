---
subcategory: "Deprecated"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_struct_template"
description: ""
---

# huaweicloud_lts_struct_template

Manage a log structuring template resource within HuaweiCloud.

!> **WARNING:** It has been deprecated, use `huaweicloud_lts_structuring_configuration` instead.

## Example Usage

### create with system template

```hcl
variable "group_id" {}
variable "stream_id" {}

resource "huaweicloud_lts_struct_template" "test" {
  log_group_id  = var.group_id
  log_stream_id = var.stream_id
  template_type = "built_in"
  template_name = "ELB"
}
```

### create with custom template

```hcl
variable "group_id" {}
variable "stream_id" {}

resource "huaweicloud_lts_struct_template" "test" {
  log_group_id  = var.group_id
  log_stream_id = var.stream_id
  template_type = "custom"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the log structuring template resource.
  If omitted, the provider-level region will be used. Changing this creates a new log stream resource.

* `log_group_id` - (Required, String, ForceNew) Specifies the ID of a log group. Changing this parameter will create
  a new resource.

* `log_stream_id` - (Required, String, ForceNew) Specifies the ID of a log stream. Changing this parameter will create
  a new resource.

* `template_type` - (Required, String, ForceNew) Specifies the type of the template. The value can be
  **built_in** (system templates) or **custom** (custom templates).
  Changing this parameter will create a new resource.

* `template_name` - (Optional, String) Specifies the system template name. The value can be **ELB**, **VPC**, **CTS**,
  **APIG**, **DDS_AUDIT**, **CDN**, and **SMN**. This parameter is mandatory when using system templates.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The log structuring template ID.

* `demo_log` - The sample log event.

## Import

The structuring templates can be imported using the template ID, lts group ID and stream ID separated by the slashes,
e.g.

```bash
$ terraform import huaweicloud_lts_struct_template.test <id>/<log_group_id>/<log_stream_id>
```
