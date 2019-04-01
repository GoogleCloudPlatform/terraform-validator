resource "google_compute_disk" "my-disk-resource" {
  project = "{{.Provider.project}}"
  name    = "{{.Disk.Name}}"
  type    = "{{.Disk.Type | pastLastSlash}}"
  zone    = "{{.Disk.Zone | pastLastSlash}}"
  image   = "{{.Disk.SourceImage | pastLastSlash}}"
  labels  = {
  {{range $key, $val := .Disk.Labels}}
    {{$key}} = "{{$val}}"
  {{end}}
  }
}
