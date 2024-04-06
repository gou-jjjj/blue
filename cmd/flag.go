package main

import "flag"

var defCluster = "127.0.0.1:13141"
var defConfPath = "./blue-server.json"

var confPath = flag.String("c", defConfPath, "config file path")
var clusterPath = flag.String("p", defCluster, "cluster path")
