# k8sinit A Kubernetes Init System

This project's aim is create automous Kubernetes manager with web panel.

## Building

Building requires alpine 3.13. Minimal setup with make command is enough. Also open community repo. Required packets will be installed by scripts

For local builds

```
make build
```

For remote builds copy rbuild.env.sample to rbuild.env and modify it as required. Allow root login with ssh key. Then

```
make remote
```

Build creates minimal initramfs with build host only support. For addinional hosts please modprobe required kernel modules. Initramfs will be builded with modules from lsmod output.
