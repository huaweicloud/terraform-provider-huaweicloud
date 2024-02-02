# Authorize users to manage OBS service-related resources

This example shows how to obtain the information required to authorize OBS operation permissions through IAM resources
and how to authorize users.

## Introduction

Currently, there are three ways to grant OBS permissions to users:

+ **Agency**: Permissions used to authorize specific services or specific domain accounts to operate cloud services.
+ **Policy and group assignment**: Permissions used to authorize specific users to operate cloud services.

This example will introduce how to configure script resources from the perspective of an authorized user.

## Policy and group assignment

This method creates OBS permissions and authorizes relevant permissions to users through user group authorization
(excluding high-risk operation permissions: deletion related operations).
This method involves the following script resources:

+ **huaweicloud_identity_role.obs_role**: A strategy that includes all OBS operations except high-risk operations.

+ **huaweicloud_identity_group.test**: A group that manages authorization policies.

+ **huaweicloud_identity_group_role_assignment.test**: Used to associate OBS policy to user group and grant global
  permission.

+ **data.huaweicloud_identity_projects.test**: Query information about MOS, a special project used to manage OBS
  service billing.

+ **huaweicloud_identity_group_role_assignment.mos**: Used to associate OBS policy to user group and grant specified
  billing permission MOS.

+ **data.huaweicloud_identity_users.test**: Query user information using its name.

+ **huaweicloud_identity_group_membership.test**: Used to associate user to user group.

## Usage

```bash
$ terraform init
$ terraform plan
$ terraform apply
$ terraform destroy
```

The creation of the ECS instance takes about few minutes. After the creation is successful, the provisioner starts to run.

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.2.0 |
| huaweicloud | >= 1.48.0 |
