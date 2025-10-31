---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_access_key"
description: |-
  Use this data source to query all permanent access key and specified access key within Huaweicloud.
---

# huaweicloud_identity_access_key

Use this data source to query all permanent access key and specified access key within Huaweicloud.

## Example Usage

### Query AccessKey By UserId

```hcl
resource "huaweicloud_identity_user" "test_user" {
  name        = "my_test_huaweicloud"
  password    = "password@123!"
  enabled     = true
  description = "tested by terraform"
}
resource "huaweicloud_identity_access_key" "key_1" {
  user_id = huaweicloud_identity_user.test_user.id
}
data "huaweicloud_identity_access_key" "test" {
  depends_on = [huaweicloud_identity_access_key.key_1]
  
  user_id = huaweicloud_identity_access_key.key_1.user_id
}
```

### Query AccessKey Info

```hcl
resource "huaweicloud_identity_user" "test_user" {
  name        = "my_test_huaweicloud"
  password    = "password@123!"
  enabled     = true
  description = "tested by terraform"
}
resource "huaweicloud_identity_access_key" "key_1" {
  user_id = huaweicloud_identity_user.test_user.id
}
data "huaweicloud_identity_access_key" "test" {
  access_key = huaweicloud_identity_access_key.key_1.id
}
```

### Query AccessKey

```hcl
data "huaweicloud_identity_access_key" "test3" {}
```

## Argument Reference

* `user_id` - (Optional, String) Specifies the IAM user ID associated with the access_key.

* `access_key` - (Optional, String) Specifies access_key associated with the identity user.

## Attribute Reference

* `credentials` - The credentials carrying user information.
  The [credentials](#IdentityAccessKey_credentials) structure is documented below.

<a name="IdentityAccessKey_credentials"></a>
The `credentials` block contains:

* `user_id` - Indicates the IAM user ID to which the access_key belongs.

* `access` - Indicates the Access Key (AK) that was queried.

* `status` - Indicates the status of the access_key. Valid values:
  + **active**: The access_key is enabled.
  + **inactive**: The access_key is disabled.

* `create_time` - Indicates the creation time of the access_key. The value is a UTC time in the
  **YYYY-MM-DDTHH:mm:ss.ssssssZ** format.

* `description` - Indicates the description of the access_key.

* `last_use_time` - Indicates the field is only present when you are certified by access_key.
