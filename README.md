A powerful linux command control ByPanel for golang


# Install
```bash
wget https://mirrors.linuxeye.com/bypanel/bypctl-linux-amd64 -O /usr/local/bin/bypctl
chmod +x /usr/local/bin/bypctl
```
Run dependency: `/opt/bypanel`

Install Recommend: `curl https://raw.githubusercontent.com/linuxeye/bypanel/main/quick_install.sh | bash`

# Help
```bash
bypctl help
```
```
bypanel deployment command management

Usage:
  bypctl [flags]
  bypctl [command]

Available Commands:
  config      Configuration of deployment parameters for bypanel
  down        Stop and remove containers, networks
  exec        Execute a command in a running container
  help        Help about any command
  logs        View output from containers
  mkcfg       Make web config
  ps          List containers
  pull        Pull service images
  reload      Reload Web service
  restart     Restart service containers
  start       Start services
  status      List containers
  stop        Stop services
  up          Create and start containers
  upgrade     Upgrade bypanel
  version     Show the bypanel version information

Flags:
  -c, --config string   config file (default is /opt/bypanel/.env)
  -h, --help            help for bypctl
  -l, --lang string     set language

Use "bypctl [command] --help" for more information about a command.
```

# Download
* AMD64: https://mirrors.linuxeye.com/bypanel/bypctl-linux-amd64
* AArch64: https://mirrors.linuxeye.com/bypanel/bypctl-linux-arm64

