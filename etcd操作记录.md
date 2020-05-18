*  进入etcd docker， 获取设置键值对
```
bogon:logmanager admin1$ docker exec -it d616fef8310a /bin/sh
/ # ls
bin         etcd0.etcd  media       root        srv         usr
dev         home        mnt         run         sys         var
etc         lib         proc        sbin        tmp
/ # etcd
2020-03-14 00:26:50.057229 I | etcdmain: etcd Version: 3.3.4
2020-03-14 00:26:50.057324 I | etcdmain: Git SHA: fdde8705f
2020-03-14 00:26:50.057333 I | etcdmain: Go Version: go1.9.5
2020-03-14 00:26:50.057338 I | etcdmain: Go OS/Arch: linux/amd64
2020-03-14 00:26:50.057343 I | etcdmain: setting maximum number of CPUs to 4, total number of available CPUs is 4
2020-03-14 00:26:50.057362 W | etcdmain: no data-dir provided, using default data-dir ./default.etcd
2020-03-14 00:26:50.060252 C | etcdmain: listen tcp 127.0.0.1:2380: bind: address already in use
/ # export ETCDCTL_API=3
/ # etcdctl get /logagent/192.168.0.142/logconfig
/logagent/192.168.0.142/logconfig
aaaaa
/ # etcdctl
NAME:
	etcdctl - A simple command line client for etcd3.
USAGE:
	etcdctl
VERSION:
	3.3.4
API VERSION:
	3.3
COMMANDS:
	get			Gets the key or a range of keys
	put			Puts the given key into the store
	del			Removes the specified key or range of keys [key, range_end)
	txn			Txn processes all the requests in one transaction
	compaction		Compacts the event history in etcd
	alarm disarm		Disarms all alarms
	alarm list		Lists all alarms
	defrag			Defragments the storage of the etcd members with given endpoints
	endpoint health		Checks the healthiness of endpoints specified in `--endpoints` flag
	endpoint status		Prints out the status of endpoints specified in `--endpoints` flag
	endpoint hashkv		Prints the KV history hash for each endpoint in --endpoints
	move-leader		Transfers leadership to another etcd cluster member.
	watch			Watches events stream on keys or prefixes
	version			Prints the version of etcdctl
	lease grant		Creates leases
	lease revoke		Revokes leases
	lease timetolive	Get lease information
	lease list		List all active leases
	lease keep-alive	Keeps leases alive (renew)
	member add		Adds a member into the cluster
	member remove		Removes a member from the cluster
	member update		Updates a member in the cluster
	member list		Lists all members in the cluster
	snapshot save		Stores an etcd node backend snapshot to a given file
	snapshot restore	Restores an etcd member snapshot to an etcd directory
	snapshot status		Gets backend snapshot status of a given file
	make-mirror		Makes a mirror at the destination etcd cluster
	migrate			Migrates keys in a v2 store to a mvcc store
	lock			Acquires a named lock
	elect			Observes and participates in leader election
	auth enable		Enables authentication
	auth disable		Disables authentication
	user add		Adds a new user
	user delete		Deletes a user
	user get		Gets detailed information of a user
	user list		Lists all users
	user passwd		Changes password of user
	user grant-role		Grants a role to a user
	user revoke-role	Revokes a role from a user
	role add		Adds a new role
	role delete		Deletes a role
	role get		Gets detailed information of a role
	role list		Lists all roles
	role grant-permission	Grants a key to a role
	role revoke-permission	Revokes a key from a role
	check perf		Check the performance of the etcd cluster
	help			Help about any command

OPTIONS:
      --cacert=""				verify certificates of TLS-enabled secure servers using this CA bundle
      --cert=""					identify secure client using this TLS certificate file
      --command-timeout=5s			timeout for short running command (excluding dial timeout)
      --debug[=false]				enable client-side debug logging
      --dial-timeout=2s				dial timeout for client connections
  -d, --discovery-srv=""			domain name to query for SRV records describing cluster endpoints
      --endpoints=[127.0.0.1:2379]		gRPC endpoints
  -h, --help[=false]				help for etcdctl
      --hex[=false]				print byte strings as hex encoded strings
      --insecure-discovery[=true]		accept insecure SRV records describing cluster endpoints
      --insecure-skip-tls-verify[=false]	skip server certificate verification
      --insecure-transport[=true]		disable transport security for client connections
      --keepalive-time=2s			keepalive time for client connections
      --keepalive-timeout=6s			keepalive timeout for client connections
      --key=""					identify secure client using this TLS key file
      --user=""					username[:password] for authentication (prompt if password is not supplied)
  -w, --write-out="simple"			set the output format (fields, json, protobuf, simple, table)
```

*  向项目应用发送http post请求 
用postman发送http请求
```
post: httt://localhost:8082/etcd/
body raw内容： 
{"key":"/logagent/192.168.0.142/logconfig", "value":`
[
{
    "topic":"nginx_log",
    "log_path":"D:\\log1",
    "service":"test_service",
    "send_rate":1000
},

{
    "topic":"nginx_log1",
    "log_path":"D:\\log2",
    "service":"test_service1",
    "send_rate":1000
}
]`
}
```
