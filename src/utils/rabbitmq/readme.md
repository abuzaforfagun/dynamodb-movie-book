### RabbitMQ

We are using Standard Classic Queue. But there are bunch of different types of quees are available. Quorum queues could be another better choice for us.

Standard Classic Queue: It store the queue into a single node.
Quorum Queue: it store the queue in multiple node.

Trade offs of quorum queue:

- Pros: It does not lose the queue when the RabbitMQ server restarts.
- Pros: When we have multi node RabbitMQ, and one nodes goes down, it still works as expected.
- Cons: It require high resources.
- Cons: As quorum queues works with distributed nodes, it add latency by the charactaristic of distrubed systems.
- Cons: This is overkil to use quorum queue with single node deployment.

Decisions:
As we are using single node, at least for now we should stay with standard class queue.
