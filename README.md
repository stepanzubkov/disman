# <center>Disman Display Manager :desktop_computer:</center>
<center><b>Disman</b> is a CLI display manager (login manager) for Xorg written in golang.</center>

#
Project Status: **Early Alpha stage**

# Installation

1. Install disman service configuration for your *Init System* (Systemd, OpenRC, SysVinit, runit, dinit, etc.)

**SysVinit (PCLinuxOS)**

Install SysVinit configuration by running `sudo make install-pclos`. It will also set disman as preffered DM.

2. `make && sudo make install` inside `src` directory

To test installation, just reboot the system.

# Configuration

See `docs/CONFIG.md`
