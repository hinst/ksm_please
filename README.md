# Use Kernel Samepage Merging in Linux

## Initial setup
1. Install KSM-tuned
    1. `sudo apt-get install ksmtuned`
1. Edit file `/etc/ksmtuned.conf`
    1. Uncomment all lines
    1. Increase `KSM_THRES_COEF`, because KSM will stay disabled when memory usage is below this threshold, and with pagefile enabled, KSM might stay off forever

* View merged memory page statistics: run command `grep . /sys/kernel/mm/ksm/*`.
* Use zero pages: `/sys/kernel/mm/ksm/use_zero_pages:0`, it is best to leave this option DISABLED because it breaks statistics.

## Systemd setup
The initial setup is not enough to fully enable Kernel Samepage Merging, because it needs to be enabled explicitly for every process.
Next step is to edit systemd unit files:

    [Service]
    MemoryKMS=true

There are many unit files, therefore I write a program to update all unit files automatically. See `main.go`