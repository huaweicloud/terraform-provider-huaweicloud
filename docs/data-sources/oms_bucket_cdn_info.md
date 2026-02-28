---
subcategory: "Object Storage Migration Service (OMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_oms_bucket_cdn_info"
description: |-
  Use this data source to query CDN information of a bucket.
---

# huaweicloud_oms_bucket_cdn_info

Use this data source to query CDN information of a bucket.

## Example Usage

```hcl
variable "access_key" {}
variable "secret_key" {}
variable "bucket_name" {}
variable "domain_name" {}

data "huaweicloud_oms_bucket_cdn_info" "test" {
  cloud_type  = "HuaweiCloud"
  ak          = var.access_key
  sk          = var.secret_key
  bucket_name = var.bucket_name

  source_cdn {
    authentication_type = "NONE"
    protocol            = "https"
    domain              = var.domain_name
  }
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
  + **URLSource**
  + **HEC**

* `ak` - (Required, String) Specifies the AK for accessing the source bucket.

* `sk` - (Required, String) Specifies the SK for accessing the source bucket.

* `bucket` - (Required, String) Specifies the bucket name.

* `source_cdn` - (Required, List) Specifies whether migration from CDN is enabled.
  The [source_cdn](#oms_bucket_cdn_info_source_cdn_struct) structure is documented below.

* `connection_string` - (Optional, String) Specifies the connection string.
  This parameter is mandatory when `cloud_type` is set to **Azure**.
  Used for Azure Blob authentication.

* `app_id` - (Optional, String) Specifies the app ID.
  This parameter is mandatory when `cloud_type` is set to **Tencent**.

* `prefix` - (Optional, List) Specifies the prefix.
  The [prefix](#oms_bucket_cdn_info_prefix_struct) structure is documented below.

<a name="oms_bucket_cdn_info_source_cdn_struct"></a>
The `source_cdn` block supports:

* `authentication_type` - (Required, String) Specifies the authentication type.
  + **NONE**: Public access without authentication type.
  + **QINIU_PRIVATE_AUTHENTICATION**: Qiniu private URL signature.
  + **ALIYUN_OSS_A**: Alibaba Cloud URL signature, simple and universal.
  + **ALIYUN_OSS_B**: Alibaba Cloud header signature, used for API calling.
  + **ALIYUN_OSS_C**: Alibaba Cloud STS temporary security token, the most secure.
  + **KSYUN_PRIVATE_AUTHENTICATION**: Kingsoft Cloud private URL signature.
  + **AZURE_SAS_TOKEN**: Microsoft Azure shared access signature, flexible and secure.
  + **TENCENT_COS_A**: Tencent Cloud multi-signature scenarios (not recommended).
  + **TENCENT_COS_B**: Tencent Cloud single-signature scenarios, the most secure.
  + **TENCENT_COS_C**:  Tencent Cloud header signature, used for API calling.
  + **TENCENT_COS_D**:  Tencent Cloud header signature, used for API calling.

* `domain` - (Required, String) Specifies the domain name used to obtain objects to be migrated.

* `protocol` - (Required, String) Specifies the protocol type.
  + **http**
  + **https**

* `authentication_key` - (Optional, String) Specifies the CDN authentication key.
  This parameter is mandatory if CDN authentication is required.

<a name="oms_bucket_cdn_info_prefix_struct"></a>
The `prefix` block supports:

* `keys` - (Required, List) Specifies the keys.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `is_same_cloud_type` - Whether the CDN is provided by the source cloud vendor.

* `is_download_available` - Whether source data can be downloaded from the CDN.

* `checked_keys` - The list of checked objects returned.
  The [checked_keys](#oms_bucket_cdn_info_checked_keys_struct) structure is documented below.

<a name="oms_bucket_cdn_info_checked_keys_struct"></a>
The `checked_keys` block supports:

* `key` - The key.

* `is_etag_matching` - Whether the etag is matched.

* `is_object_existing` - Whether the object is found.
