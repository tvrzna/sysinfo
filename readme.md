# sysinfo
Simple sysinfo website for Linux

```
Usage: sysinfo [options]
Options:
	-h, --help			print this help
	-v, --version			print version
	-n, --name [NAME]		name of application to be displayed
	-p, --port [PORT]		sets port for listening
	-a, --app-url [APP_URL]		application url (if behind proxy)
```

![](screenshot.png)

## Roadmap
- [X] CPU usage and frequency widget
- [X] RAM and SWAP widget
- [X] System temperature widget
- [X] Diskusage widget
- [X] Load Average and Uptime widget
- [X] Custom interval with possible pause option
- [X] Network usage widget
- [X] Process watcher widget
- [X] Shared results, if there is more data receivers, it should be handled by only one sysinfo load.
- [X] Diskstats widget
- [ ] Custom widget layout with possible configuration of widgets