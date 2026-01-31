---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_service_linked_agency"
description: |-
  Use this resource to manage a service-linked agency within HuaweiCloud.
---

# huaweicloud_identityv5_service_linked_agency

Use this resource to manage a service-linked agency within HuaweiCloud.

-> 1. Service-linked agency can only be deleted by services.
   <br/>2. IAM administrators only have permission to view them in IAM, this prevents accidental deletion
    and service failure.
   <br/>3. This resource is only a one-time action resource for creating the service-linked agency. Deleting this
   resource will not clear the corresponding request record, but will only remove the resource information from the
   tfstate file.

## Example Usage

```hcl
variable "service_principal" {}

resource "huaweicloud_identityv5_service_linked_agency" "test" {
  service_principal = var.service_principal
}
```

## Argument Reference

The following arguments are supported:

* `service_principal` - (Required, String, ForceNew) Specifies the service principal of the service-linked agency.  
  The service principal must start with `service.` and follow by a string of `1` to `56` characters that contains
  only letters, digits, and hyphens (-).

* `description` - (Optional, String, ForceNew) Specifies the description of a service-linked agency.  
  The description cannot contain special characters: `@#%&<>\$^*`.  
  The maximum length is **1000** characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `urn` - The uniform resource name.

* `trust_policy` - The policy document of the service-linked agency, in JSON format.  
  The following characters `_=<>()|` are special characters in the syntax and are not included in the trust policy.

* `agency_id` - The ID of the service-linked agency.

* `agency_name` - The name of the service-linked agency.

* `path` - The resource path of the service-linked agency, in `service-linked-agency/<service_principal>/` format.

* `created_at` - The creation time of the service-linked agency.

* `max_session_duration` - The maximum session duration of the service-linked agency, in seconds.
