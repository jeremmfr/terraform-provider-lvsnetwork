# terraform-provider-lvsnetwork
![GitHub release (latest by date)](https://img.shields.io/github/v/release/jeremmfr/terraform-provider-lvsnetwork)
[![Go Status](https://github.com/jeremmfr/terraform-provider-lvsnetwork/workflows/Go%20Tests/badge.svg)](https://github.com/jeremmfr/terraform-provider-lvsnetwork/actions)
[![Lint Status](https://github.com/jeremmfr/terraform-provider-lvsnetwork/workflows/GolangCI-Lint/badge.svg)](https://github.com/jeremmfr/terraform-provider-lvsnetwork/actions)
[![GoDoc](https://godoc.org/github.com/jeremmfr/terraform-provider-lvsnetwork?status.svg)](https://godoc.org/github.com/jeremmfr/terraform-provider-lvsnetwork)
[![Go Report Card](https://goreportcard.com/badge/github.com/jeremmfr/terraform-provider-lvsnetwork)](https://goreportcard.com/report/github.com/jeremmfr/terraform-provider-lvsnetwork)

Terraform's provider for setup network interface and keepalived vrrp_instance on two server (master/slave) with lvsnetwork-api (https://github.com/jeremmfr/lvsnetwork-api)

Compile:
========

go build -o terraform-provider-lvsnetwork && mv terraform-provider-lvsnetwork /usr/bin/

Config:
=======

Setup information for contact server :
```
provider "lvsnetwork" {
    firewall_ip = "192.168.0.1"
    port		= 8443
	https		= true
	insecure	= true
	vault_enable = true
	default_id_vrrp = 10
	default_track_script = [ "check_custom" ]
}
```

* **firewall_ip** : (Required) IP for firewall API (lvsnetwork-api)
* **port** : (Optional) [Def: 8080] Port for firewall API (lvsnetwork-api)
* **https** : (Optional) [Def: false] Use HTTPS for firewall API
* **insecure** : (Optional) [Def: false] Don't check certificate for HTTPS
* **login** : (Optional) [Def: ""] User for http basic authentication
* **password** : (Optional) [Def: ""] Password for http basic authentication
* **vault_enable** : (Optional) [Def: false] Read login/password in secret/$vault_path/$firewall_ip or secret/$vault_path/$vault_key (For server and token, read environnement variables "VAULT_ADDR", "VAULT_TOKEN") ConflictWith **login**/**password**
* **vault_path** : (Optional) [Def: "lvs"] Path where the key are
* **vault_key** : (Optional) [Def: ""] Name of key in vault path
* **default_id_vrrp** : (Required) Default id for parameter id_vrrp in resource lvsnetwork_ifacevrrp
* **default_vrrp_group** : (Optional) [Def: "VG_1"] Default VG for parameter vrrp_group in resource lvsnetwork_ifacevrrp
* **default_advert_int** : (Optional) [Def: 1 ] Default interval for parameter advert_int in resource lvsnetwork_ifacevrrp
* **default_track_script** : (Optional) Default track_script in resource lvsnetwork_ifacevrrp

Resource:
=========

** ifacevrrp **
---------------

Create iface and/or vrrp configuration on two servers (MASTER/SLAVE)

```
resource lvsnetwork_ifacevrrp "vlan471" {
	iface = "vlan471"
	ip_master = "10.0.71.253"
	ip_slave = "10.0.71.252"
	mask = "24"
	ip_vip = [ "10.0.71.150", "10.0.71.100", "10.0.71.254" ]
}
```
* **iface** : (Required) name of interface for iface configuration and vrrp configuration
* **ip_master** : (Optional) [ Def : "" ] IPv4 for iface configuration on master server
* **ip_slave** : (Optional) [ Def : "" ] IPv4 for iface configuration on slave server
* **mask** : (Optional) [ Def : "" ] short netmask for iface configuration on master/slave server
* **vlan_device** : (Optional) [ Computed : bond1 (if iface ~= vlan) ] vlan-raw-device for iface configuration
* **post_up** : (Optional) [ Def : [""] ] list of post-up line for iface configuration
* **default_gw** : (Optional) [ Def : "" ] default gateway, gateway parameter in iface configuration
* **lacp_slaves** : (Optional) [ Def : "" ] 802.3ad configuration with slaves ifaces
* **lacp_slaves_slave** : (Optional) [ Def : "" ] 802.3ad configuration with slaves ifaces for backup router only if different on master router
* **ip_vip_only** : (Optional) [ Def : false ] only configure vrrp
* **ip_vip** : (Optional) [ Def : [""]] list of IPv4 in vrrp configuration
* **prio_master** : (Optional) [ Computed : 150] priority for vrrp configuration on master server
* **prio_slave** : (Optional) [ Computed : 100] priority for vrrp configuration on slave server
* **vrrp_group** : (Optional) [ Computed : `default_vrrp_group` ] vrrp_sync_group for vrrp configuration
* **iface_vrrp** : (Optional) [ Def : "" ] interface for vrrp configuration (default same as iface)
* **id_vrrp** : (Optional) [ Computed : $default_id_vrrp from provider ] id for vrrp configuration (must be unique for iface)
* **auth_type** : (Optional) [ Computed : "PASS" ] authentication auth_type
* **auth_pass** : (Optional) [ Computed : "word" ] authentication auth_pass
* **sync_iface** : (Optional) [ Def : "" ] lvs_sync_daemon_interface parameter for vrrp configuration (must be unique on server)
* **garp_m_delay** : (Optional) [ Computed : 5 ] garp_master_delay parameter for vrrp configuration
* **advert_int** : (Optional) [ Computed : `default_advert_int` ] advert_int parameter for vrrp configuration
* **garp_master_refresh** : (Optional) [ Computed : 60 ] garp_master_refresh parameter for vrrp configuration
* **use_vmac** : (Optional) [ Def : true ] Use vmac for vrrp if possible
* **track_script** : (Optional) [ Computed : `default_track_script` ] List of vrrp_script to track

** vrrp_script **
---------------

Create vrrp_script configuration on two servers (MASTER/SLAVE)

```
resource lvsnetwork_vrrp_script check_custom {
        name = "check_custom"
        script = "/usr/local/bin/scripts/check_custom bond0"
        user = "guest"
        interval = 5
        fall = 3
        rise = 2
}
```

* **name** (Required) name of vrrp_script
* **script** (Required) programm and arguments for this vrrp_script
* **interval** (Optional) seconds between script invocations
* **timeout** (Optional) seconds after which script is considered to have failed
* **weight** (Optional) [Default: 0] adjust priority by this weight
* **weight_reverse** (Optional) reverse causes the direction of the adjustment of the priority to be reversed
* **rise** (Optional) [Default: 3] number of successes for OK transition
* **fall** (Optional) [Default: 3] number of successes for KO transition
* **user** (Optional) user to run script under
* **init_fail** (Optional) assume script initially is in failed state

Apply or change can be long because sleep between ifup or keepalived reload
