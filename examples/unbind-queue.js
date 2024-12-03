import {Connection} from "k6/x/amqp"; // import Amqp extension

export default function () {
  const url = "amqp://guest:guest@localhost:5672/"
  const conn = new Connection({
    connection_url: url,
  });
  
  const queueName = 'K6 queue'
  const exchangeName = 'K6 exchange'

  conn.unbindQueue({
    queue_name: queueName,
    routing_key: 'rkey123',
    exchange_name: exchangeName,
    args: null
  })

  console.log(queueName + " queue unbinded from " + exchangeName)
}
