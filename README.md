# Parallel Shell

Execute Parallel SSH Commands to Multiple Hosts

## NOT READY FOR USE YET

## History

This is a port of a very good parallel shell we used in xCAT
for Linux clusters.  I had the need with multiple machines to 
use it on a workstation node for various system administration
tasks. 

The original xCAT code was written by Jarrod Johnson and Egan Ford.
The source can be seen [here](https://sourceforge.net/p/xcat/xcat-core/ci/master/tree/xCAT-client/bin/psh)
and [here](https://sourceforge.net/p/xcat/xcat-core/ci/master/tree/perl-xCAT/xCAT/NodeRange.pm)

## Usage

Create a file ```~/.psh```
In this file, put your list of nodes.  Right now its: 
```
<nodename>,<ipaddress>[:port],<remote user name>,<priv key>,<group>,<group>,..
```

## Syntax and Noderange support

This is a list of supported noderange ways to call psh: 

#### Comma separated list: (node01,node02,node03)

```
psh node01,node02,node03 date
```

#### Range syntax: (node01-node99) or (node01:node99)

```
psh node01-node55 date
```

