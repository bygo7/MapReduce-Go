# Method

We chose to use RPC in order to replicate the data structures. The leader replicates in four scenarios:
1. When a new worker is registered
2. When a map task is assigned to a worker
3. When a reduce task is assigned to a worker
4. When a user request is completely finished

In all scenarios, every data structure is replicated via gRPC and entirely replaced inside the followers. We did this for simplicity to cut down on the number of different replication methods we would need to implement. We also chose to replicate the entire data structure instead of just the changes because we wanted to make sure that the followers were always in sync with the leader, i.e. out-of-order messaging, although the leader does block on replication so this is not really a concern.

Each RPC comes with a sequence number that is incremented every time a new RPC is sent. If a follower receives an RPC with a sequence number that is less than the sequence number of the last RPC it received, it will ignore the RPC as it is out of date.

# Data Structures

We chose to have three main data structures, all of which are persisted through RPC:
1. Worker Pool
2. Task Pool
3. User Request Queue

The worker pool holds all of the registered workers' statuses and is able to poll idle workers for task assignment.

The task pool holds all of the map tasks, the reduce tasks, and the active tasks that have been assigned to workers. It also holds the input, output, and status of each task.

The user request queue holds all of the user requests that have been submitted to the master.
