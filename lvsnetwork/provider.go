package lvsnetwork

import (
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	defaultFirewallPort = 8080
	defaultAdvertInt    = 1
)

// Provider lvsnetwork for terraform.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"firewall_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  defaultFirewallPort,
			},
			"https": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"insecure": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"login": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vault_enable": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"login", "password"},
			},
			"vault_path": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "lvs",
			},
			"vault_key": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"default_advert_int": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      defaultAdvertInt,
				ValidateFunc: validation.IntBetween(1, 10),
			},
			"default_auth_pass": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "word",
				ValidateFunc: validation.StringLenBetween(1, 7),
			},
			"default_id_vrrp": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 255),
			},
			"default_track_script": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"default_vrrp_group": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "VG_1",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"lvsnetwork_ifacevrrp":   resourceIfaceVrrp(),
			"lvsnetwork_vrrp_script": resourceVrrpScript(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	defaultTrackScript := make([]string, 0)
	for _, elem := range d.Get("default_track_script").([]interface{}) {
		defaultTrackScript = append(defaultTrackScript, elem.(string))
	}
	config := Config{
		firewallIP:         d.Get("firewall_ip").(string),
		firewallPortAPI:    d.Get("port").(int),
		https:              d.Get("https").(bool),
		insecure:           d.Get("insecure").(bool),
		logname:            os.Getenv("USER"),
		login:              d.Get("login").(string),
		password:           d.Get("password").(string),
		vaultEnable:        d.Get("vault_enable").(bool),
		vaultPath:          d.Get("vault_path").(string),
		vaultKey:           d.Get("vault_key").(string),
		defaultIDVrrp:      d.Get("default_id_vrrp").(int),
		defaultAuthPass:    d.Get("default_auth_pass").(string),
		defaultVrrpGroup:   d.Get("default_vrrp_group").(string),
		defaultAdvertInt:   d.Get("default_advert_int").(int),
		defaultTrackScript: defaultTrackScript,
	}

	return config.Client(), nil
}
