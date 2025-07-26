---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_device_authorization"
description: ""
---

# huaweicloud_identitycenter_device_authorization

Manages an Identity Center device authorization resource within HuaweiCloud.

## Example Usage

```hcl
variable "client_id" {}
variable "client_secret" {}
variable "start_url" {}
resource "huaweicloud_identitycenter_device_authorization" "test"{
	client_id       = var.client_id
	client_secret   = var.client_secret
    start_url		= var.start_url
}
```

## Argument Reference

The following arguments are supported:

* `client_id` - (Required, String, ForceNew) Unique ID of the client registered in the IAM Identity Center.

  Changing this parameter will create a new resource.

* `client_secret` - (Required, String, ForceNew) Secret string generated for the client to obtain authorization from services in subsequent calls.

  Changing this parameter will create a new resource.

* `start_url` - (Required, String, ForceNew) User Portal URL.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `device_code` - Device code used by the device to poll session tokens. 

* `verification_uri_complete` - Alternate URL that the client can use to automatically start the browser. This procedure skips the manual steps of the user accessing the validation page and entering the code.
