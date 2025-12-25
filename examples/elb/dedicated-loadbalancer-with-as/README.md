# Integrate ELB with Auto Scaling Group

This example provides best practice code for using Terraform to integrate a dedicated Elastic Load Balance (ELB)
service with dedicated instance in HuaweiCloud.
This example demonstrates how to automatically add or remove backend servers based on scaling policies, including ELB, listener,
backend server group, scaling configuration, and scaling group.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `loadbalancer_name` - The name of the loadbalancer
* `listener_name` - The name of the listener
* `configuration_name` - The name of the scaling configuration
* `configuration_user_data` - The user data for the scaling configuration instances
* `configuration_disks` - The disk configurations for the scaling configuration instances, each disk includes:
  + `size` - The disk size in GB (required)
  + `volume_type` - The volume type, such as "SAS", "SSD", "GPSSD" (required)
  + `disk_type` - The disk type, such as "SYS" for system disk, "DATA" for data disk (required)
* `group_name` - The name of the AS group
* `security_group_name` - The name of the security group
* `instance_name` - The name of the ECS instance
* `alarm_rule_name` - The name of the CES alarm rule
* `alarm_rule_conditions` - The conditions of the CES alarm rule, each condition includes:
  + `period` - The period of the condition (required)
  + `filter` - The filter of the condition (required)
  + `comparison_operator` - The comparison operator (required)
  + `value` - The threshold value (required)
  + `unit` - The unit of the value (required)
  + `count` - The count of consecutive occurrences (required)
  + `alarm_level` - The alarm severity level, 1-4 (required)
  + `metric_name` - The metric name (required)
* `policy_name` - The name of the AS policy

#### Optional Variables

* `availability_zone` - The name of the availability zone to which the resources belong (default: "")
* `vpc_cidr` - The CIDR block of the VPC (default: "172.16.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP address of the subnet (default: "")
* `loadbalancer_cross_vpc_backend` - Whether to associate backend servers with the load balancer by using their IP
  addresses (default: false)
* `loadbalancer_description` - The description of the loadbalancer (default: null)
* `enterprise_project_id` - The enterprise project ID (default: null)
* `loadbalancer_tags` - The tags of the loadbalancer (default: {})
* `loadbalancer_force_delete` - Whether to force delete the loadbalancer (default: true)
* `eip_type` - The type of the EIP (default: "5_bgp")
* `bandwidth_name` - The name of the EIP bandwidth (default: "tf_test_eip")
* `bandwidth_size` - The bandwidth size of the EIP in Mbps (default: 5)
* `bandwidth_share_type` - The share type of the bandwidth (default: "PER")
* `bandwidth_charge_mode` - The charge mode of the bandwidth (default: "traffic")
* `listener_protocol` - The protocol of the listener (default: "HTTP")
* `listener_port` - The port of the listener (default: 8080)
* `listener_idle_timeout` - The idle timeout of the listener (default: 60)
* `listener_request_timeout` - The request timeout of the listener (default: null)
* `listener_response_timeout` - The response timeout of the listener (default: null)
* `listener_description` - The description of the listener (default: null)
* `listener_tags` - The tags of the listener (default: {})
* `pool_name` - The name of the pool (default: null)
* `pool_protocol` - The protocol of the pool (default: "HTTP")
* `pool_method` - The load balancing method of the pool (default: "ROUND_ROBIN")
* `pool_any_port_enable` - Whether to enable any port for the pool (default: false)
* `pool_description` - The description of the pool (default: null)
* `pool_persistences` - The persistence configurations for the pool (default: [])
  + `type` - The persistence type
  + `cookie_name` - The cookie name for APP_COOKIE type (optional)
  + `timeout` - The timeout for the persistence (optional)
* `instance_flavor_id` - The flavor ID of the instance (default: "")
* `instance_flavor_performance_type` - The performance type of the instance flavor (default: "normal")
* `instance_flavor_cpu_core_count` - The CPU core count of the instance flavor (default: 2)
* `instance_flavor_memory_size` - The memory size of the instance flavor (default: 4)
* `instance_image_id` - The image ID of the instance (default: "")
* `instance_image_visibility` - The visibility of the instance image (default: "public")
* `instance_image_os` - The OS of the instance image (default: "Ubuntu")
* `configuration_image_id` - The image ID of the scaling configuration (default: "")
* `configuration_flavor_id` - The flavor ID of the scaling configuration (default: "")
* `group_desire_instance_number` - The desire instance number of the AS group (default: 0)
* `group_min_instance_number` - The min instance number of the AS group (default: 0)
* `group_max_instance_number` - The max instance number of the AS group (default: 10)
* `group_delete_publicip` - Whether to delete the public IP address of the AS group (default: true)
* `group_delete_instances` - Whether to delete the instances of the AS group (default: true)
* `group_force_delete` - Whether to force delete the AS group (default: true)
* `policy_cool_down_time` - The cool down time of the AS policy (default: 900)
* `policy_operation` - The operation of the AS policy (default: "ADD")
* `policy_instance_number` - The instance number of the AS policy (default: 1)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                     = "tf_test_vpc"
  subnet_name                  = "tf_test_subnet"
  loadbalancer_name            = "tf_test_loadbalancer"
  listener_name                = "tf_test_listener"
  security_group_name          = "tf_test_security_group"
  instance_name                = "tf_test_instance"
  configuration_name           = "tf_test_configuration"
  configuration_user_data      = "Your_linux_Password"
  configuration_disks          = [{ size = 40, volume_type = "SAS", disk_type = "SYS" }]
  group_name                   = "tf_test_as_group"
  group_desire_instance_number = 1
  group_min_instance_number    = 0
  group_max_instance_number    = 10
  alarm_rule_name              = "tf_test_alarm_rule"
  alarm_rule_conditions        = [
    {
      period              = 300
      filter              = "max"
      comparison_operator = ">"
      value               = 80
      unit                = "%"
      count               = 3
      alarm_level         = 2
      metric_name         = "cpu_util"
    }
  ]

  policy_name = "tf_test_as_policy"
  ```

* Initialize Terraform:

  ```bash
  $ terraform init
  ```

* Review the Terraform plan:

  ```bash
  $ terraform plan
  ```

* Apply the configuration:

  ```bash
  $ terraform apply
  ```

* To clean up the resources:

  ```bash
  $ terraform destroy
  ```

## Note

* Make sure to keep your credentials secure and never commit them to version control
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.3.0 |
| huaweicloud | >= 1.69.0 |
