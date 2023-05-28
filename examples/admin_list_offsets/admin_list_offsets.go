/**
 * Copyright 2018 Confluent Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Delete topics
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {

	bootstrapServers := "localhost:9092"

	// Create a new AdminClient.
	// AdminClient can also be instantiated using an existing
	// Producer or Consumer instance, see NewAdminClientFromProducer and
	// NewAdminClientFromConsumer.
	a, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": bootstrapServers})
	if err != nil {
		fmt.Printf("Failed to create Admin client: %s\n", err)
		os.Exit(1)
	}

	// Contexts are used to abort or limit the amount of time
	// the Admin call blocks waiting for a result.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var requests map[kafka.TopicPartition]int
	requests[kafka.TopicPartitions({Topic:"topicname",Partition : 0})] = kafka.EarliestOffsetSpec;


	results, err := a.ListOffsets(ctx, requests,kafka.SetAdminIsolationLevel(kafka.ReadCommitted))
	if err != nil {
		fmt.Printf("Failed to List offsets: %v\n", err)
		os.Exit(1)
	}
	// map[TopicPartition]ListOffsetResultInfo
	// Print results
	for tp,info := range results {
		fmt.Printf("Topic: %s Partition_Index : %d\n",tp.Topic,tp.Partition)
		if info.err.code {
			fmt.Printf("	ErrorCode : %d ErrorMessage : %s\n\n",info.err.code,info.err.str)
		} else {
			fmt.Prinitf("	Offset : %d Timestamp : %d\n\n",info.offset,info.timestamp)
		}
	}

	a.Close()
}
