## Milestone 1

zookeeper-server-start -daemon /usr/local/etc/kafka/zookeeper.properties & kafka-server-start /usr/local/etc/kafka/server.properties

/bin/schema-registry-start etc/schema-registry/schema-registry.properties

## Define an event schema representing when an order has been received

{"namespace": "com.kakfainaction",
 "type": "record",
 "name": "Customer",
 "fields": [
     {"name":"date", "type": "date"},
     {"name":"customerId", "type": "string"},
     {"name":"firstName", type:"string"},
     {"name":"lastName", type:"string"},
     {"name":"address", type:"string"}
  ]
}

 {"namespace": "com.kakfainaction",
 "type": "record",
 "name": "Product",
 "fields": [
     {"name":"date", "type": "date"},
     {"name":"name", "type": "string"},
     {"name":"description", type:"string"},
     {"name":"price", type:"float"}
  ]
}

{"namespace": "com.kakfainaction",
 "type": "record",
 "name": "Order",
 "fields": [
     {"name":"date", "type": "date"},
     {"name":"orderId", "type": "string"},
     {"name":"customer", type:"Customer"},
     {"name":"product", type:"Product"}
     {"name":"status", "type":
        {"type":"enum",
         "name":"status",
         "symbols":["Pending","Shipped","PaidFor","Error"]},
   ]
}

## Define an event schema representing when an order has been confirmed (is not a duplicate order and can be processed).

{"namespace": "com.kakfainaction",
 "type": "record",
 "name": "OrderConfirmed",
 "fields": [
     {"name":"date", "type": "date"},
     {"name":"orderId", "type": "string"}
   ]
}

## Define an event schema representing when an order has been picked from within a warehouse and packed (ready to be shipped).
{"namespace": "com.kakfainaction",
 "type": "record",
 "name": "OrderPicked",
 "fields": [
     {"name":"date", "type": "date"},
     {"name":"orderId", "type": "string"},
     {"name":"customer", type:"Customer"}
   ]
}

## Define an event schema representing when an email notification needs to be sent out.
{"namespace": "com.kakfainaction",
 "type": "record",
 "name": "OrderReady",
 "fields": [
     {"name":"date", "type": "date"},
     {"name":"orderId", "type": "string"}
   ]
}

## Define an event schema representing when a consumer is unable to successfully process an event it has received.
# ##This “error” event should contain the event that could not be processed.

{"namespace": "com.kakfainaction",
 "type": "record",
 "name": "OrderError",
 "fields": [
     {"name":"date", "type": "date"},
     {"name":"orderId", "type": "string"}
     {"name":"reason", "type":["OrderReady", "OrderConfirmed", "OrderReady" ],
   ]
}


## Create a Kafka topic that will contain order received events, and verify it exists. The easiest way to do this is to use Step 3 in the “Apache Kafka Quickstart” guide.

bin/kafka-topics.sh --describe --topic milestone1-order --bootstrap-server localhost:9092 --replication-factor 3 --partitions 11

## Modify the topic created in Step 3 by increasing its retention time to three days. For more information on retention policies, review the last bullet point in the notes below.

## Create additional Kafka topics that use the same configuration as the topic created in Step 3.
--replication-factor <n> --partitions <n>

Create a Kafka topic that will contain order confirmed events, and verify it exists.
bin/kafka-topics.sh --describe --topic milestone1-order-confirmed --bootstrap-server localhost:9092

Create a Kafka topic that will contain order picked and packed events, and verify it exists.
bin/kafka-topics.sh --describe --topic milestone1-order-picked --bootstrap-server localhost:9092

Create a Kafka topic that will contain notification events, and verify it exists.
bin/kafka-topics.sh --describe --topic milestone1-order-ready --bootstrap-server localhost:9092

Create a Kafka topic that will contain error events, and verify it exists.
bin/kafka-topics.sh --describe --topic milestone1-order-error --bootstrap-server localhost:9092

bin/kafka-console-producer.sh --topic milestone1-order --bootstrap-server localhost:9092
{"orderId": "1234",
  "customer:" : {"firstName": "Jack", "lastName": "Jones"}.
  "product:" : {"name": "Furby", "description": "black"}, "price": "34.55"},
}
