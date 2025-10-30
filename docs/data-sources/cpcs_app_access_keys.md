---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cpcs_app_access_keys"
description: |-
  Use this data source to get the list of CPCS application access keys.
---

# huaweicloud_cpcs_app_access_keys

Use this data source to get the list of CPCS application access keys.

-> Currently, this data source is valid only in cn-north-9 region.

## Example Usage

```hcl
variable "app_id" {}

data "huaweicloud_cpcs_app_access_keys" "test" {
  app_id = var.app_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `app_id` - (Required, String) Specifies the application ID.

* `key_name` - (Optional, String) Specifies the access key name used for filtering.

* `sort_key` - (Optional, String) Specifies the sort attribute.
  The default value is **create_time**.

* `sort_dir` - (Optional, String) Specifies the sort direction.
  The default value is **DESC**. Valid values are **ASC** and **DESC**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `access_keys` - Indicates the access keys list.
  The [access_keys](#CPCS_app_access_keys) structure is documented below.

<a name="CPCS_app_access_keys"></a>
The `access_keys` block supports:

* `access_key_id` - The access key ID.

* `key_name` - The access key name.

* `access_key` - The access key AK.

* `app_name` - The name of the application to which the access key belongs.

* `status` - The access key status.

* `create_time` - The creation time of the access key.

* `download_time` - The download time of the access key.

* `is_downloaded` - Whether the access key has been downloaded.

* `is_imported` - Whether the access key is imported.
