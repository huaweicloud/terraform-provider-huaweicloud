---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_service_linked_agency"
description: |-
  Create a service-linked agency within HuaweiCloud.
---

# huaweicloud_identityv5_service_linked_agency

Create a service-linked agency within HuaweiCloud.

->**Note** The service-linked agency resource can not be destroyed. Service-linked agencies can only be deleted by services.
  IAM administrators only have permission to view them in IAM. This prevents accidental deletion and service failure.

## Example Usage

```hcl
variable "service_principal" {}
variable "description" {}

resource "huaweicloud_identityv5_service_linked_agency" "test" {
  service_principal = var.service_principal
  description       = var.description
}
```

## Argument Reference

The following arguments are supported:

* `service_principal` - (Required, String, ForceNew) Specifies the service principal, which starts with `service.` and
  is followed by a string of 1 to 56 characters containing only letters, digits, and hyphens `-`.

* `description` - (Optional, String, ForceNew) Specifies the description of a service-linked agency.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `urn` - Indicates the uniform resource name.

* `trust_policy` - Indicates the JSON format of the policy document of a trust agency's trust policy.

* `agency_id` - Indicates the service-linked agency ID.

* `agency_name` - Indicates the service-linked agency name .

* `path` - Indicates the resource path in the format of `service-linked-agency/<service_principal>/`.

* `created_at` - Indicates the time when the service-linked agency was created.

* `max_session_duration` - Indicates the maximum session duration of the service-linked agency.
