# Cloud Tools #

*Tools to make your Cloud life easier.*

Cloud Tools is a set of tools that allow for easier, faster and safer
operation for Cloud team. Currently it supports managing the version
manifest on a semi-automated way.

## Installation ##

You only need a Go toolchain (1.10 and up) to make use this tool. Make
sure you have set your `GOPATH` correctly and then run.

    go get gitlab.com/licorna/cloud-tools

After doing this, you should be able to reach the `cloud-tools`
program in your `PATH` (`GOBIN` needs to be part of `PATH`).

## Build ##

You need the `go-bindata` module and build static assets as part of
the go binary.

    go get -u github.com/jteeuwen/go-bindata/...
    go-bindata assets/
    go build

## Usage ##

You can use `cloud-tools` for different tasks. We are going to focus
on managing the version manifest.

### Add a new version to the version manifest ###

    cloud-tools version-manifest add-build --new 3.6.9 --merge-with https://s3.amazonaws.com/om-kubernetes-conf/4.0.json

This will print out the version manifest passed via `merge-with` with
the new version passed via the `new` parameter.

To test this, try to add the 3.6.9 version into a version manifest
that has not been updated with this version yet. As an example I've
published the
https://s3.amazonaws.com/om-kubernetes-conf/4.0_pre_adding_3.6.9.json
file which does not include this version. If you run the same command
with this file instead:

    cloud-tools version-manifest add-build --new 3.6.9 --merge-with https://s3.amazonaws.com/om-kubernetes-conf/4.0_pre_adding_3.6.9.json

`cloud-tools` will output the merging of this version manifest with
the added 3.6.9 version.

### Publishing changes to S3 ###

**Coming soon**

    cloud-tools version-manifest add-build --new 3.6.9 \
        --merge-with https://s3.amazonaws.com/om-kubernetes-conf/4.0_pre_adding_3.6.9.json \
        --into https://s3.amazonaws.com/om-kubernetes-conf/4.0_post_adding_3.6.9.json

This command will merge the new 3.6.9 version into the `--merge-with`
file and publish it with the file passed as the `--into`
parameter. This requires write access to this particular S3 bucket.
