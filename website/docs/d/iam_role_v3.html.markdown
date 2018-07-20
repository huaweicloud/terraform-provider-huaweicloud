---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iam_role_v3"
sidebar_current: "docs-huaweicloud-datasource-iam-role-v3"
description: |-
  Get all the IAM roles including global and a domain's custom ones.
---

# huaweicloud\_iam\_role_v3

Use this data source to get all the IAM roles a domain can use.

## Example Usage

```hcl

data "huaweicloud_iam_role_v3" "roles" {
}

```

**Note**: It can not set `tenant_name` in `provider "huaweicloud"` when
   using this data source.

## Argument Reference


## Attributes Reference

* `projects` - The list of roles which can be granted only to a project. Each
    role will include its name and description.

* `domains` - The list of roles which can be granted only to a domain. Each
    role will include its name and description.

* `project_domains` - The list of roles which can be granted to a project or
    domain. Each role will include its name and description.

* `others` - The list of roles which can be granted to other service, like
    object storage service. Each role will include its name and description.
