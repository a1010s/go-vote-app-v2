
// Datacenter
resource "ionoscloud_datacenter" "golang-app" {
  name = "Quick-Server (rename)"
  location = var.location
}

// Public Lan
resource "ionoscloud_lan" "internet" {
  name = "Public Lan Connection"
  public = true
  datacenter_id = ionoscloud_datacenter.golang-app.id
}

locals {
  ssh_key_content = <<-EOT
    ${secrets.SSH_PUB_KEY_FOR_TF_SERVER}
  EOT
}

data "template_file" "ssh_key" {
  template = local.ssh_key_content
}

// Server
resource "ionoscloud_server" "server" {
  #count             = 1
  name              = "server${var.id_name} ${var.location}"
  datacenter_id     = ionoscloud_datacenter.golang-app.id
  availability_zone = "AUTO"
  cores             = 1
  ram               = 2048
  cpu_family        = "INTEL_SKYLAKE"
  image_name        = "ubuntu:latest"
  ssh_key_path      = data.template_file.ssh_key.rendered

  volume {
    name      = "server-volume${var.location} boot"
    size      = 10
    disk_type = "SSD"
  }
  nic {
    lan             = ionoscloud_lan.internet.id
    dhcp            = true
    firewall_active = false

  }
  lifecycle {
    ignore_changes = [nic, ssh_key_path]
  }
}
