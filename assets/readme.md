# min_max_versions.yaml #

This is a definition of minOsVersion and maxOsVersion that will be
used to produce the entries for the new builds from the full.json
(server's version manifest).

Each entry in the `flavors` list is a Linux "flavor", and for each one
of the versions, a `min` and `max` is defined.

Let's assume a new `debian` version 10 was released, then a new entry
in the `osVersions` list of the `name: debian` flavor should be
entered.

This file helps with the **transformation** from the source to
destination file.

# ops_manager_host_support.yaml #

This file describes the different Ops Manager versions that we
support, currently 4.0 and 3.6. For each version we have different
platforms (Linux, Windows and OSX).

For now we assume all the OSX and Windows builds will be included.

On the other hand, for each Linux build, its version needs to be
supported (included in the `supportedVersions` list) and its
architecture needs to be supported as well (included in the `arch`
list).

A good example of how this file works is by looking at the archs that
are only supported in 4.0, like `s390x` or the `ubuntu1804` builds
only available in Ops Manager 4.0
