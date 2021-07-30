# lvsnetwork_vrrp_script

Create vrrp_script configuration on two servers (MASTER/SLAVE)

## Example Usage

```hcl
resource lvsnetwork_vrrp_script check_custom {
  name     = "check_custom"
  script   = "/usr/local/bin/scripts/check_custom bond0"
  user     = "guest"
  interval = 5
  fall     = 3
  rise     = 2
}
```

## Argument Reference

* **name** (Required) name of vrrp_script
* **script** (Required) program and arguments for this vrrp_script
* **interval** (Optional) seconds between script invocations
* **timeout** (Optional) seconds after which script is considered to have failed
* **weight** (Optional) [Default: 0] adjust priority by this weight
* **weight_reverse** (Optional) reverse causes the direction of the adjustment of the priority to be reversed
* **rise** (Optional) [Default: 3] number of successes for OK transition
* **fall** (Optional) [Default: 3] number of successes for KO transition
* **user** (Optional) user to run script under
* **init_fail** (Optional) assume script initially is in failed state
