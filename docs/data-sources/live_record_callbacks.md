---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_record_callbacks"
description: |-
  Use this data source to get a list of the Live record callback configurations.
---

# huaweicloud_live_record_callbacks

Use this data source to get a list of the Live record callback configurations.

## Example Usage

```hcl
data "huaweicloud_live_record_callbacks" "test" {
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `domain_name` - (Optional, String) Specifies the ingest domain name.

* `app_name` - (Optional, String) Specifies the application name.
  To match all applications, set this parameter to a wildcard character *****.
  Exact application matching is preferred. If no application is matched, all applications will be matched.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `callbacks` - The callback configurations.

  The [callbacks](#callbacks_struct) structure is documented below.

<a name="callbacks_struct"></a>
The `callbacks` block supports:

* `sign_type` - The encryption type. Contains the following values:
  + **HMACSHA256**
  + **MD5**

* `created_at` - The creation time in the format of **yyyy-mm-ddThh:mm:ssZ** (UTC time).

* `updated_at` - The latest modification time in the format of **yyyy-mm-ddThh:mm:ssZ** (UTC time).

* `id` - The recording callback ID.

* `domain_name` - The ingest domain name.

* `app_name` - The application name.

* `url` - The callback URL for sending recording notifications.

* `types` - The types of recording notifications. Contains the following values:
  + **RECORD_NEW_FILE_START**: Recording started.
  + **RECORD_FILE_COMPLETE**: Recording file generated.
  + **RECORD_OVER**: Recording completed.
  + **RECORD_FAILED**: Recording failed.
