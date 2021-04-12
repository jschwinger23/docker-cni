# Docker CNI

This repo aims to integrate [CNI](https://github.com/containernetworking/cni) with Dockerd.

There is, [according to CNI repo](https://github.com/containernetworking/cni/blob/master/scripts/docker-run.sh), an approach to integrate by running a [pause](https://groups.google.com/g/kubernetes-users/c/jVjv0QK4b_o) equivalent container ahead of the application container, but that's too pod-like for those who resent pod models.

Let's figure out yet another solution.

# Usage

## 0. Prepare your CNI plugin

Make sure you have everything ready:

1. CNI binaries in the right place: for example, `/opt/cni/bin/calico` and `/opt/cni/bin/calico-ipam` binaries
2. CNI configures in the right place: for exmaple, `/etc/cni/net.d/10-calico.conf`
3. Other services needed: for example, `calico-node` container

Notes:

1. Provided there are multiple CNI configures in the dir, `docker-cni` will only use the first config in alphabet order.

## 1. Configure docker-cni

### 1.1 install docker-cni binary

TODO

### 1.2 docker-cni configuration

```shell
mkdir -p /etc/docker/
cat <<!
cni_conf_dir: /etc/cni/net.d/
cni_bin_dir: /opt/cni/bin/
log_driver: file:///var/run/log/docker-cni.log
log_level: debug
!
```

You may revise the aforementioned configure with YOUR `cni_conf_dir` and `cni_bin_dir`.

## 2. Configure dockerd

### 2.1 dockerd daemon configuration

Add the additional `runtime` in docker daemon configure, which is usually located at `/etc/docker/daemon.json`:

```
{
    ...
    "runtimes": {
        "cni": {
            "path": "/usr/local/bin/docker-cni",
            "runtimeArgs": [
                "--config",
                "/etc/docker/cni.yaml",
                "--runtime-path",
                "/usr/bin/runc"]
        }
    }
}
```

### 2.2 restart dockerd

```
systemctl restart dockerd
```

## 3. Create docker container with CNI

```
docker run -td --runtime cni --net none bash bash
```

That's everything.
