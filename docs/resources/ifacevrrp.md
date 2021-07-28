# lvsnetwork_ifacevrrp

Create iface and/or vrrp configuration on two servers (MASTER/SLAVE)

## Example Usage

```hcl
resource lvsnetwork_ifacevrrp "vlan471" {
  iface     = "vlan471"
  ip_master = "10.0.71.253"
  ip_slave  = "10.0.71.252"
  mask      = "24"
  ip_vip    = ["10.0.71.150", "10.0.71.100", "10.0.71.254"]
}
```

## Argument Reference

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
