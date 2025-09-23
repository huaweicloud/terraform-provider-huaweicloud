---
subcategory: "Identity and Access Management (IAM)"
---

# huaweicloud_identity_service_agency

Manages a service agency resource within HuaweiCloud.

-> **NOTE:** You *must* have admin privileges to use this resource.

## Example Usage

### Delegate another HUAWEI CLOUD service to perform operations on your resources

```hcl
variable "agency_name" {}
variable "delegated_service_name" {}
variable "policy_name" {}

resource "huaweicloud_identity_service_agency" "test" {
  name                   = var.agency_name
  delegated_service_name = var.delegated_service_name
  policy_names           = [var.policy_name]
  description            = "test demo"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Specifies the name of service agency. The name is a string of `1` to `64`
  characters. Only English letters, digits, underscores (_), plus (+), equals (=), commas (,), dots (.), ats (@) and
  hyphens (-) are allowed. Changing this will create a new service agency.

* `delegated_service_name` - (Required, String, ForceNew) Specifies the name of delegated service name.
  Prefix is `service.`. Such as **service.APIG**. Changing this will create a new service agency.

* `policy_names` - (Required, List) Specifies a string list of one or more policy names that you would like to attach to
  the service agency.

* `path` - (Optional, String, ForceNew) Specifies the resource path. It is made of several strings, each containing one
  or more English letters, digits, underscores (_), plus (+), equals (=), comma (,), dots (.), at (@) and hyphens (-),
  and must be ended with slash (/). Such as **foo/bar/**. It's a part of the uniform resource name. Default is empty.
  Changing this will create a new service agency.

* `duration` - (Optional, Int) Specifies the validity period of a service agency.
  Default value is `3,600`. The unit is seconds.

* `tags` - (Optional, Map) Specifies the tags of the service agency.

* `description` - (Optional, String) Specifies the description of the service agency.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The service agency ID.

* `trust_policy` - The trust policy of the service agency. It's a JSON string.

* `urn` - The uniform resource name of the service agency. Format is `iam::$accountID:agency:$path$agencyName` where
  `$accountID` is IAM account ID, `$path` is `path`, `$agencyName` is `name`.

* `created_at` - The time when the service agency was created.

## Import

Service agencies can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_identity_service_agency.test <id>
```
