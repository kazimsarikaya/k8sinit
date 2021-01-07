# k8sinit A Kubernetes Init System

This project's aim is create automous Kubernetes manager with web panel.

## Building

Building requires alpine 3.12 with zfs support. Additionaly genisoimage binary is required.

For building go 1.15 and make is essential.

For local builds

```
make build
```

For remote builds copy rbuild.env.sample to rbuild.env and modify it as required. Then

```
make remote
```

Build creates minimal initramfs with build host only support. For addinional hosts please modprobe required kernel modules. Initramfs will be builded with modules from lsmod output.
