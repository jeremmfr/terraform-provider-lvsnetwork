# lvsnetwork Provider

Terraform's provider for setup network interface and keepalived vrrp_instance on two server (master/slave) with [lvsnetwork-api](https://github.com/jeremmfr/lvsnetwork-api)

~> Apply or change can be long because sleep between ifup or keepalived reload

## Example Usage

```hcl
provider "lvsnetwork" {
  firewall_ip          = "192.168.0.1"
  port                 = 8443
  https                = true
  insecure             = true
  vault_enable         = true
  default_id_vrrp      = 10
  default_track_script = ["check_custom"]
}
```

## Argument Reference

* **firewall_ip** : (Required) IP for firewall API (lvsnetwork-api)
* **default_id_vrrp** : (Required) Default id for parameter id_vrrp in resource lvsnetwork_ifacevrrp
* **port** : (Optional) [Def: 8080] Port for firewall API (lvsnetwork-api)
* **https** : (Optional) [Def: false] Use HTTPS for firewall API
* **insecure** : (Optional) [Def: false] Don't check certificate for HTTPS
* **login** : (Optional) [Def: ""] User for http basic authentication
* **password** : (Optional) [Def: ""] Password for http basic authentication
* **vault_enable** : (Optional) [Def: false] Read login/password in secret/$vault_path/$firewall_ip or secret/$vault_path/$vault_key  
(For server and token, read environnement variables "VAULT_ADDR", "VAULT_TOKEN") ConflictWith **login**/**password**
* **vault_path** : (Optional) [Def: "lvs"] Path where the key are
* **vault_key** : (Optional) [Def: ""] Name of key in vault path
* **default_advert_int** : (Optional) [Def: 1 ] Default interval for parameter advert_int in resource lvsnetwork_ifacevrrp
* **default_auth_pass** : (Optional) [Def: "word"] Default password to auth vrrp
* **default_track_script** : (Optional) Default track_script in resource lvsnetwork_ifacevrrp
* **default_vrrp_group** : (Optional) [Def: "VG_1"] Default VG for parameter vrrp_group in resource lvsnetwork_ifacevrrp
