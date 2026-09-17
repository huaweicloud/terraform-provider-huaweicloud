---
subcategory: "DRS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_files"
description: |-
  Use this data source to query the list of downloadable DRS files within HuaweiCloud.
---

# huaweicloud_drs_files

Use this data source to query the list of downloadable DRS files within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_drs_files" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the downloadable files.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `files` - The list of downloadable files.  
  The [files](#drs_files_files) structure is documented below.

<a name="drs_files_files"></a>
The `files` block supports:

* `file_name` - The name of the file.

* `last_modified` - The last modified time of the file.
