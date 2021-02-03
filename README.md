# Usage of auto deployment
## Scenario

* You have several nodes (less than 1000) need to manipulate.
* The tasks are almost the same. (Transfer/update files, excute bash)
* Your node has limited performance (In this scenario the nodes considered as iot devices) and no internet access.
* You can connect your node via ssh and scp (scp should be installed in the node)
* You can use this little tool to reach your nodes, well, but with limited functionalities

## Usage

You can use this auto-deployment-tool by writing 3 kinds of files:

* Bash files (*.sh) that will be implemented by the nodes in node-domain.
* Domain file (nodeDomain.txt), describe the node IP address, indices, ssh-username, ssh-password.
* Operation file (*.csv), describe the task pipe line including command and bash files.

Domain file has the following syntax:

[index], [IP], [ssh-username], [ssh-password]

Bash files are trivial and as usual

Operation file has three keywords (for now)
* command, the command that directly implemented via ssh, with syntax(e.g.):


**command, touch swh.txt**



* copy, send a local file to remote node, with syntax(e.g. followings). (Specify the destination will be included later)


**copy, /temp/swh.txt**



* copyN, send a local file with a pattern name. For now only one pattern is supported, a string with a symbol "%" in, and "%" will be substituted by node index. (e.g. our domain file is as below: )


2,  192.168.233.104,       pi,     raspberry
3,  192.168.233.105,    pi, raspberry
4,  192.168.233.106,     pi,     raspberry
5,  192.168.233.107,    pi,     raspberry 

task copyN is used:


**copyN, node%.config**


Then, node2.config, node3.config, node4.config, node5.config will be found and transfered to corresponding nodes.

## Execution

Very simple, only one line command


**autoDeployment**(built target excutable file) **operation.csv**(operation)