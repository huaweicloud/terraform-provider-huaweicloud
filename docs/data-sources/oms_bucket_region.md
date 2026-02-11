---
subcategory: "Object Storage Migration Service (OMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_oms_bucket_region"
description: |-
  Use this data source to query the region information corresponding to the bucket.
---

# huaweicloud_oms_bucket_region

Use this data source to query the region information corresponding to the bucket.

## Example Usage

```hcl
variable "bucket_name" {}
variable "access_key" {}
variable "secret_key" {}

data "huaweicloud_oms_bucket_region" "test" {
  cloud_type  = "HuaweiCloud"
  bucket_name = var.bucket_name
  ak          = var.access_key
  sk          = var.secret_key
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

* `bucket_name` - (Required, String) Specifies the bucket name.

* `ak` - (Required, String) Specifies the AK for accessing the source bucket.

* `sk` - (Required, String) Specifies the SK for accessing the source bucket.

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

* `id` -  The data source ID.

* `region_id` - The region ID.

* `name` - The region name.

* `support` - Whether migration is supported in the region.
