---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_key"
description: ""
---

# huaweicloud_identity_key

Query all permanent access key data source within Huaweicloud

-> **NOTE:** You *must* have admin privileges to use this data source.

## Example Usage

```hcl
data "huaweicloud_identity_key" "test" {

}

```

## Argument Reference

* `user_id` - (Optional, String) The IAM user ID associated with the access key.

## Attribute Reference

* `user_id` -  The IAM user ID to which the access key belongs.

* `access` - The Access Key (AK) that was queried.

* `status` - The status of the access key. Valid values: **active**, The access key is enabled. **inactive**: The access key is disabled.

* `create_time` - The creation time of the access key. The value is a UTC time in the YYYY-MM-DDTHH:mm:ss.ssssssZ format.

* `description` - The description of the access key.


