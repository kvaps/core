// Copyright (c) 2016 Pani Networks
// All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/romana/core/agent/policy/cache"
	"github.com/romana/core/common"
	"github.com/romana/core/common/client"
)

func main() {
	endpointsStr := flag.String("etcd-endpoints", "localhost:2379", "Comma-separated list of etcd endpoints.")
	flag.Parse()

	endpoints := strings.Split(*endpointsStr, ",")

	clientConfig := common.Config{EtcdEndpoints: endpoints,
		EtcdPrefix: "romana",
	}
	client, err := client.NewClient(&clientConfig)
	if err != nil {
		panic(err)
	}
	c := cache.New(client, cache.Config{CacheTickSeconds: 10})

	stop := make(chan struct{})
	update := c.Run(stop)

	for {
		select {
		case hash := <-update:
			fmt.Printf("New policy hash %s\n", hash)
			PrintPolicy(c)
		}
	}
}

func PrintPolicy(cache cache.Interface) {
	policies := cache.List()
	for policyNum, policy := range policies {
		fmt.Printf("Policy %d ID %s\n", policyNum, policy.ID)
	}

	fmt.Printf("Detected %d romana policies\n", len(cache.List()))
}