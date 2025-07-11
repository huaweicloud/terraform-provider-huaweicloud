---
subcategory: "CodeArts Build"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_build_log_download"
description: |-
  Manages a CodeArts Build log download resource within HuaweiCloud.
---

# huaweicloud_codearts_build_log_download

Manages a CodeArts Build log download resource within HuaweiCloud.

## Example Usage

```hcl
variable "record_id" {}

resource "huaweicloud_codearts_build_log_download" "test" {
  record_id = var.record_id
  log_file  = "./buildLog.txt"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `record_id` - (Required, String, NonUpdatable) Specifies the record ID.

* `log_level` - (Optional, String, NonUpdatable) Specifies the log level. Value can be **INFO** and **DEBUG**.

* `log_file` - (Optional, String, NonUpdatable) Specifies the log file path. Defaults to *./{{record_id}}.txt*

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
