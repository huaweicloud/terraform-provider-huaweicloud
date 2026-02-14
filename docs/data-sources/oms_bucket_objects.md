---
subcategory: "Object Storage Migration Service (OMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_oms_bucket_objects"
description: |-
  Use this data source to query the list of buckets.
---

# huaweicloud_oms_bucket_objects

Use this data source to query the list of buckets.

## Example Usage

```hcl
variable "access_key" {}
variable "secret_key" {}
variable "bucket_name" {}
variable "file_path" {}

data "huaweicloud_oms_bucket_objects" "test" {
  cloud_type  = "HuaweiCloud"
  ak          = var.access_key
  sk          = var.secret_key
  bucket_name = var.bucket_name
  file_path   = var.file_path
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cloud_type` - (Required, String) Specifies the cloud service provider.
  The valid values are as follows:
  + **AWS**
  + **Aliyun**
  + **Qiniu**
  + **QingCloud**
  + **Tencent**
  + **Baidu**
  + **KingsoftCloud**
  + **Azure**
  + **UCloud**
  + **HuaweiCloud**
  + **Google**
  + **URLSource**
  + **HEC**

* `ak` - (Required, String) Specifies the AK for accessing the source bucket.

* `sk` - (Required, String) Specifies the SK for accessing the source bucket.

* `bucket_name` - (Required, String) Specifies the bucket name.

* `file_path` - (Required, String) Specifies the path of object files to be queried in the destination bucket.
  The value must end with a slash (/). e.g. **test/**.

* `json_auth_file` - (Optional, String) Specifies the JSON auth file.
  This parameter is mandatory when `cloud_type` is set to **Google**.
  Used for Google Cloud Storage authentication.

* `connection_string` - (Optional, String) Specifies the connection string.
  This parameter is mandatory when `cloud_type` is set to **Azure**.
  Used for Azure Blob authentication.

* `app_id` - (Optional, String) Specifies the app ID.
  This parameter is mandatory when `cloud_type` is set to **Tencent**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The queried object information.
  The [records](#oms_bucket_objects_records_struct) structure is documented below.

<a name="oms_bucket_objects_records_struct"></a>
The `records` block supports:

* `name` - The object name.

* `size` - The object size.
