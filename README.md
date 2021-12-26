# Overview of Milestone 4 Implementation 

# What Was Created?
In this milestone:

created a database called `liveproject` that contains the following table definition:
```sql
-- DROP TABLE processed_events;

CREATE TABLE processed_events (
	id varchar(255) NOT NULL,
	processed_timestamp timestamp NOT NULL,
	event_name varchar(256) NOT NULL
);
```
