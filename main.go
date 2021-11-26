package main

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
	"time"
)

type Data struct {
	Title   string
	Results []string
	Other   []int
}

/*
func process2(w http.ResponseWriter, r *http.Request) {
	data := &Data{"Lord of the Rings", []string{"a", "b", "c"}, []int{1, 2, 3}}
	t, _ := template.ParseFiles("tmpl2.html")
	t.Execute(w, data)
}
 */

func main() {
	broker := os.Args[1]

	// Create a new AdminClient.
	// AdminClient can also be instantiated using an existing
	// Producer or Consumer instance, see NewAdminClientFromProducer and
	// NewAdminClientFromConsumer.
	a, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		fmt.Printf("Failed to create Admin client: %s\n", err)
		os.Exit(1)
	}
	/*
		if len(os.Args) != 3 {
			fmt.Fprintf(os.Stderr, "Usage: %s <broker> <topic>\n",
				os.Args[0])
			os.Exit(1)
		}

		broker := os.Args[1]
		topic := os.Args[2]

		p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})

		if err != nil {
			fmt.Printf("Failed to create producer: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Created Producer %v\n", p)

		// Optional delivery channel, if not specified the Producer object's
		// .Events channel is used.
		deliveryChan := make(chan kafka.Event)

		value := "Hello Go!"
		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(value),
			Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
		}, deliveryChan)

		e := <-deliveryChan
		e.(*kafka.Message)
	*/

	//http.HandleFunc("/process2", process2)

	//http.Handle("/", http.FileServer(http.Dir("./src")))

	//if err := http.ListenAndServe(":8080", nil); err != nil {
	//	panic(err)
	}

}
