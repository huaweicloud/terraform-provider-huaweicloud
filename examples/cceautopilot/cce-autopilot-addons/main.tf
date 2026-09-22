# Create VPC
resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

# Create subnet
resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr == "" ? cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0) : var.subnet_cidr
  gateway_ip = var.gateway_ip == "" ? cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1) : var.gateway_ip
}

# Create CCE Autopilot cluster
resource "huaweicloud_cce_autopilot_cluster" "test" {
  name        = var.cluster_name
  flavor      = "cce.autopilot.cluster"
  version     = var.cluster_version
  alias       = var.cluster_alias
  description = var.cluster_description
  category    = var.cluster_category
  type        = var.cluster_type
  custom_san  = var.cluster_custom_san
  eip_id      = var.cluster_eip_id

  enable_snat             = var.cluster_enable_snat
  enable_swr_image_access = var.cluster_enable_swr_image_access

  delete_efs   = var.cluster_delete_efs
  delete_eni   = var.cluster_delete_eni
  delete_net   = var.cluster_delete_net
  delete_obs   = var.cluster_delete_obs
  delete_sfs30 = var.cluster_delete_sfs_turbo

  lts_reclaim_policy = var.cluster_lts_reclaim_policy

  host_network {
    vpc    = huaweicloud_vpc.test.id
    subnet = huaweicloud_vpc_subnet.test.id
  }

  container_network {
    mode = "eni"
  }

  eni_network {
    subnets {
      subnet_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id
    }
  }

  extend_param {
    enterprise_project_id = var.cluster_enterprise_project_id
  }

  tags = var.cluster_tags

  lifecycle {
    ignore_changes = [
      enable_snat,
      delete_efs,
      delete_obs,
    ]
  }
}

# Create SWR organization for log-agent addon
resource "huaweicloud_swr_organization" "test" {
  name = var.swr_organization_name
}

# Create CCE Autopilot addon
resource "huaweicloud_cce_autopilot_addon" "test" {
  cluster_id          = huaweicloud_cce_autopilot_cluster.test.id
  addon_template_name = var.addon_template_name
  version             = var.addon_version
  name                = var.addon_name
  alias               = var.addon_alias

  values = {
    "basic" = jsonencode(merge({
      "aomEndpoint" : "https://aom.${var.region_name}.myhuaweicloud.com",
      "iam_url" : "iam.${var.region_name}.myhuaweicloud.com",
      "ltsAccessEndpoint" : "https://lts-access.${var.region_name}.myhuaweicloud.com:8102",
      "ltsEndpoint" : "https://lts.${var.region_name}.myhuaweicloud.com",
      "region" : var.region_name,
      "swr_addr" : "swr.${var.region_name}.myhuaweicloud.com",
      "swr_user" : huaweicloud_swr_organization.test.name,
      "rbac_enabled" : true,
      "cluster_version" : huaweicloud_cce_autopilot_cluster.test.version
    }, var.addon_values_basic))

    "flavor" = jsonencode(merge({
      "category" : ["Autopilot"],
      "description" : "Applicable to clusters where the logs of a single pod are less than 10000/s or 5 MB/s.",
      "is_default" : true,
      "name" : "custom-resources",
      "resources" : [
        {
          "name" : "log-operator",
          "limitsCpu" : "1000m",
          "requestsCpu" : "1000m",
          "replicas" : 3,
          "limitsMem" : "2048Mi",
          "requestsMem" : "2048Mi"
        },
        {
          "name" : "otel-collector-event",
          "limitsCpu" : "1000m",
          "requestsCpu" : "1000m",
          "replicas" : 3,
          "limitsMem" : "2048Mi",
          "requestsMem" : "2048Mi"
        }
      ],
      "size" : "custom"
    }, var.addon_values_flavor))

    "custom" = jsonencode(merge({
      "accessKey" : "",
      "agency_name" : "",
      "aomEndpoint" : "https://aom.${var.region_name}.myhuaweicloud.com",
      "aomPrivateEndpointIP" : "",
      "bufferChunkSize" : "128k",
      "bufferMaxSize" : "512k",
      "caCert" : "",
      "clusterID" : huaweicloud_cce_autopilot_cluster.test.id,
      "clusterName" : huaweicloud_cce_autopilot_cluster.test.name,
      "cluster_category" : "CCE",
      "createAudit" : true,
      "createDefaultEvent" : true,
      "createDefaultEventToAOM" : true,
      "createDefaultStdout" : true,
      "createKubeApiserver" : false,
      "createKubeControllerManager" : false,
      "createKubeScheduler" : false,
      "enableEventReport" : true,
      "enableFullPathCollection" : false,
      "enableGcrypto" : true,
      "enableLogOperatorHA" : true,
      "enableLogReport" : true,
      "featureGates" : ["sendAOMAllEvent", "outputKafka", "podLabelExclude", "fullPathCollection", "independentEvents"],
      "host_network" : false,
      "ltsAccessEndpoint" : "https://lts-access.${var.region_name}.myhuaweicloud.com:8102",
      "ltsAuditStreamID" : "",
      "ltsEndpoint" : "https://lts.${var.region_name}.myhuaweicloud.com",
      "ltsEnterpriseProjectID" : var.cluster_enterprise_project_id,
      "ltsEventStreamID" : "",
      "ltsGroupID" : "",
      "ltsKubeApiserverStreamID" : "",
      "ltsKubeControllerManagerStreamID" : "",
      "ltsKubeSchedulerStreamID" : "",
      "ltsLogReportDomain" : "",
      "ltsPrivateEndpointIP" : "",
      "ltsStdoutStreamID" : "",
      "maxEventAgeSeconds" : "",
      "memBufLimit" : "40mb",
      "multiAZEnabled" : false,
      "otelReportLogs" : true,
      "paasakskEnable" : true,
      "podDisruptionBudget" : { "create" : true, "maxUnavailable" : 1 },
      "projectID" : "",
      "secretKey" : "",
      "securityToken" : "",
      "serverCert" : "",
      "serverKey" : ""
    }, var.addon_values_custom))
  }

  depends_on = [
    huaweicloud_cce_autopilot_cluster.test,
    huaweicloud_swr_organization.test
  ]

  lifecycle {
    ignore_changes = [
      values
    ]
  }
}
