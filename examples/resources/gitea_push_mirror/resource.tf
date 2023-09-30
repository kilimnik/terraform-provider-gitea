resource "gitea_repository" "test" {
  username     = "lerentis"
  name         = "test"
  private      = true
  issue_labels = "Default"
  license      = "MIT"
  gitignores   = "Go"
}

resource "gitea_push_mirrir" "mirror" {
  owner           = "lerentis"
  repo            = "test"
  remote_address  = "https://git.uploadfilter24.eu/lerentis/terraform-provider-gitea.git"
  remote_password = var.remote_password
  remote_username = "lerentis"
}
