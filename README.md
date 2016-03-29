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
In this file, put your list of nodes.  The syntax of this file is:
```
<nodename>,<ipaddress>[:port],<remote user name>,<priv key>,<group>,<group>,..
```

As an example, this could be: 

```
node01,10.93.234.5,ubuntu,~/.ssh/t4,aws,prod,us-east
```
The first values are fixed: A short name, IP address, username, and private
key.  These are the basic things you use to SSH into machines.  The subsequent 
fields are group names for your servers.  

## Syntax and Noderange support

This is a list of supported noderanges that you can use psh to run 
the ```date``` command on a slew of nodes.  Really you can run any
remote command but the ```date``` command is used here for an example.


#### Comma separated list of nodes

```
psh node01,node99,node13 date
```

#### Range of nodes syntax: (node01-node99) or (node01:node99)

```
psh node01-node55 date
psh node01:node55 date # the ':' and '-' are both proper. 
```
or you can use: 
```
psh node[01-55] date
```

#### Exclusions

You can use the '-' in front of a node to not include a node or range of 
nodes from a group:

```
psh node01-node99,-node05-node10 date
```
This is equivalent to doing node01-node04,node11-node99
