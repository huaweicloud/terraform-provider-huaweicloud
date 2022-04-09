data "huaweicloud_availability_zones" "myaz" {}

data "huaweicloud_vpc" "myvpc" {
  id = var.vpc_id
}

data "huaweicloud_vpc_subnet" "mysubnet" {
  id = var.subnet_id
}

resource "huaweicloud_vpc_eip" "cce" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "cce-apiserver"
    size        = 20
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_vpc_eip" "nat" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "NAT"
    size        = 100
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_nat_gateway" "nat_1" {
  name      = var.nat_name
  spec      = "1"
  vpc_id    = data.huaweicloud_vpc.myvpc.id
  subnet_id = data.huaweicloud_vpc_subnet.mysubnet.id
}

resource "huaweicloud_nat_snat_rule" "snat_1" {
  nat_gateway_id = huaweicloud_nat_gateway.nat_1.id
  floating_ip_id = huaweicloud_vpc_eip.nat.id
  subnet_id      = data.huaweicloud_vpc_subnet.mysubnet.id
}

resource "huaweicloud_cce_cluster" "cluster" {
  name                   = var.cce_name
  cluster_type           = "VirtualMachine"
  cluster_version        = "v1.19"
  flavor_id              = "cce.s1.small"
  vpc_id                 = data.huaweicloud_vpc.myvpc.id
  subnet_id              = data.huaweicloud_vpc_subnet.mysubnet.id
  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
  eip                    = huaweicloud_vpc_eip.cce.address
  delete_all             = "true"
}

resource "huaweicloud_cce_node" "cce-node1" {
  cluster_id        = huaweicloud_cce_cluster.cluster.id
  name              = "node1"
  flavor_id         = "s6.large.2"
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
  key_pair          = var.key_pair_name

  root_volume {
    size       = 80
    volumetype = "SAS"
  }
  data_volumes {
    size       = 100
    volumetype = "SAS"
  }
}

resource "huaweicloud_cce_node" "cce-node2" {
  cluster_id        = huaweicloud_cce_cluster.cluster.id
  name              = "node2"
  flavor_id         = "s6.large.2"
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
  key_pair          = var.key_pair_name

  root_volume {
    size       = 80
    volumetype = "SAS"
  }
  data_volumes {
    size       = 100
    volumetype = "SAS"
  }
}

resource "huaweicloud_cce_node" "cce-node3" {
  cluster_id        = huaweicloud_cce_cluster.cluster.id
  name              = "node3"
  flavor_id         = "s6.large.2"
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
  key_pair          = var.key_pair_name

  root_volume {
    size       = 80
    volumetype = "SAS"
  }
  data_volumes {
    size       = 100
    volumetype = "SAS"
  }
}

resource "local_file" "kubeconfig" {
  content  = huaweicloud_cce_cluster.cluster.kube_config_raw
  filename = "$Local kubeconfig file path"
}

