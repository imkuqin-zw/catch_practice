#!/bin/bash

./etcdctl put zw.com/shop/inventory-srv/app '{"name":"zw.com.shop.inventory","version":"v1","address":""}'

./etcdctl put zw.com/shop/inventory-srv/zap '{"level":"debug","development":false}'

./etcdctl put zw.com/shop/etcd '{"addrs":["127.0.0.1:2379"]}'