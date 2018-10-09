res9
====
__!! under development... !!__  
__res9__ : revive container service written in golang. 


Description
----
__res9__ is the daemon that checks the health status of container and revive unhealthy container. This tool use systemd, upstart or sysvinit for daemonize. __res9__ depend on health check option on dockerfile. Because It inspects container health status for health check. If there is not health option, __res9__ can't inspect the health status of the container. 

Revive is not restart.  
Procedure of revive container is as follows.  

1. stop container
1. remove container
1. create container 
1. start container 

Require
----
- need health check option on dockerfile  
- docker client version 1.38

Example
----
Procedure of run as a service is as follows.  

1. install service
1. start service

```
$ sudo ./res9 install
$ sudo ./res9 start
```

__result__
\\\\\\\\ image \\\\\\\\\

Usage
----

__install service__  
```
$ sudo ./res9 install
```

__start service__  
```
$ sudo ./res9 start
```

__stop service__  
```
$ sudo ./res9 stop
```

__remove service__  
```
$ sudo ./res9 remove
```

__show service status__  
```
$ sudo ./res9 status
```
