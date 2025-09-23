---
subcategory: "ServiceStage"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_servicestagev3_environment"
description: |-
  Manages an environment resource within HuaweiCloud.
---

# huaweicloud_servicestagev3_environment

Manages an environment resource within HuaweiCloud.

## Example Usage

```hcl
variable "vpc_id" {}
variable "env_name" {}
variable "enterprise_project_id" {}

resource "huaweicloud_servicestagev3_environment" "test" {
  vpc_id                = var.vpc_id
  name                  = var.env_name
  deploy_mode           = "mixed"
  description           = "Created by terraform script"
  enterprise_project_id = var.enterprise_project_id

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the environment is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID to which the environment belongs.  
  Changing this will create a new resource.

* `name` - (Required, String) Specifies the name of the environment.  
  The valid length is limited from `2` to `64`, only letters, digits, hyphens (-) and underscores (_) are allowed.
  The name must start with a letter and end with a letter or a digit.

* `deploy_mode` - (Optional, String, ForceNew) Specifies the deploy mode of the environment.  
  + **virtualmachine**: Virtual machine.
  + **container**: Kubernetes.
  + **mixed**: Both virtual machine and Kubernetes.

  Defaults to **mixed**.
  Changing this will create a new resource.

* `description` - (Optional, String) Specifies the description of the environment.  
  The value can contain a maximum of `128` characters.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the environment belongs.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the environment that used to filter resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, in UUID format.

* `creator` - The creator name of the environment.

* `created_at` - The creation time of the environment, in RFC3339 format.

* `updated_at` - The latest update time of the environment, in RFC3339 format.

## Import

Environments can be imported using their `id`, e.g.

```bash
$ terraform import huaweicloud_servicestagev3_environment.test <id>
```
