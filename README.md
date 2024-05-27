# <center>Disman Display Manager :desktop_computer:</center>
<center><b>Disman</b> is a CLI display manager (login manager) for Xorg written in golang.</center>

#
Project Status: **Early Alpha stage**

# Installation

1. `sudo mv -v res/pam /etc/pam.d/disman`
2. Create disman service configuration for your *Init System* (Systemd, OpenRC, SysVinit, runit, dinit, etc.). *They are not ready yet.*
3. `make && sudo make install` inside `src` directory

To test installation, just reboot the system.
