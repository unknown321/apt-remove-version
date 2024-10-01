apt-remove-version
==================

Removes package records from local apt cache.

```
Usage of ./apt-remove-version:
  -keep_cache
        keep apt cache file?
  -lists string
        path to apt lists directory with *_Packages (default "/var/lib/apt/lists/")
  -out string
        path to output directory (default "/var/lib/apt/lists/")
  -pkgcache string
        full path to apt cache file (default "/var/cache/apt/pkgcache.bin")

Package name must be provided in https://golang.org/s/re2syntax format
Package version must be an exact match

Example:
        ./apt-remove-version -out /tmp/no-new-nvidia/ nvidia-driver=550.90.12-1 ".*=545.23.08-1"
```

### Why?

Let's say you want to downgrade `nvidia-driver` from 560* to 550*:

```shell
$ apt-cache madison nvidia-driver
nvidia-driver | 560.35.03-1 | https://developer.download.nvidia.com/compute/cuda/repos/debian12/x86_64  Packages
nvidia-driver | 560.28.03-1 | https://developer.download.nvidia.com/compute/cuda/repos/debian12/x86_64  Packages
nvidia-driver | 555.42.06-1 | https://developer.download.nvidia.com/compute/cuda/repos/debian12/x86_64  Packages
nvidia-driver | 555.42.02-1 | https://developer.download.nvidia.com/compute/cuda/repos/debian12/x86_64  Packages
nvidia-driver | 550.90.12-1 | https://developer.download.nvidia.com/compute/cuda/repos/debian12/x86_64  Packages
nvidia-driver | 550.90.07-1 | https://developer.download.nvidia.com/compute/cuda/repos/debian12/x86_64  Packages

$ sudo apt-get install nvidia-driver=550.90.12-1
Reading package lists... Done
Building dependency tree... Done
Reading state information... Done
Some packages could not be installed. This may mean that you have
requested an impossible situation or if you are using the unstable
distribution that some required packages have not yet been created
or been moved out of Incoming.
The following information may help to resolve the situation:

The following packages have unmet dependencies:
 nvidia-driver : Depends: nvidia-driver-libs (= 550.90.12-1) but it is not going to be installed
                 Depends: nvidia-driver-bin (= 550.90.12-1) but 560.35.03-1 is to be installed
                 Depends: xserver-xorg-video-nvidia (= 550.90.12-1) but it is not going to be installed
                 Depends: nvidia-vdpau-driver (= 550.90.12-1) but 560.35.03-1 is to be installed
                 Depends: nvidia-alternative (= 550.90.12-1)
                 Depends: nvidia-kernel-dkms (= 550.90.12-1) but it is not going to be installed or
                          nvidia-kernel-545.23.08 or
                          nvidia-kernel-open-dkms (= 550.90.12-1) but it is not going to be installed
                 Recommends: nvidia-settings (>= 545) but it is not going to be installed
                 Recommends: libnvidia-cfg1 (= 550.90.12-1) but 560.35.03-1 is to be installed
E: Unable to correct problems, you have held broken packages.
```

You'll have to specify version not only for main package, but also for dependencies (and dependencies' dependencies!): 

```shell
$ sudo apt-get install nvidia-driver-libs=550.90.12-1 \
  nvidia-driver-bin=550.90.12-1 xserver-xorg-video-nvidia=550.90.12-1 \
  nvidia-vdpau-driver=550.90.12-1 nvidia-alternative=550.90.12-1 \
  nvidia-kernel-dkms=550.90.12-1 ...
```

Aptitude provides weird solutions with mixed versions.

At this point comes this program, which removes unneeded versions from apt files:

```shell
$ sudo ./apt-remove-version ".*=560.35.03-1" ".*=560.28.03-1" ".*=555.42.06-1" ".*=555.42.02-1" 
2024/10/01 18:25:55 INFO removing package name=cuda-compat-12-6 version=560.35.03-1 filename=/var/lib/apt/lists/developer.download.nvidia.com_compute_cuda_repos_debian12_x86%5f64_Packages
2024/10/01 18:25:55 INFO removing package name=cuda-drivers version=560.35.03-1 filename=/var/lib/apt/lists/developer.download.nvidia.com_compute_cuda_repos_debian12_x86%5f64_Packages
...
2024/10/01 18:25:55 INFO saving to=/var/lib/apt/lists/developer.download.nvidia.com_compute_cuda_repos_debian12_x86%5f64_Packages from=/var/lib/apt/lists/developer.download.nvidia.com_compute_cuda_repos_debian12_x86%5f64_Packages
2024/10/01 18:25:55 INFO removed apt cache file path=/var/cache/apt/pkgcache.bin

$ apt-cache madison nvidia-driver
nvidia-driver | 550.90.12-1 | https://developer.download.nvidia.com/compute/cuda/repos/debian12/x86_64  Packages
nvidia-driver | 550.90.07-1 | https://developer.download.nvidia.com/compute/cuda/repos/debian12/x86_64  Packages
nvidia-driver | 550.54.15-1 | https://developer.download.nvidia.com/compute/cuda/repos/debian12/x86_64  Packages
```

Now you can install `nvidia-driver` without specifying version, `apt-mark hold`ing them, etc.