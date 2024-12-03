import {Connection} from "k6/x/amqp"; // import Amqp extension

export default function () {
  const url = "amqp://guest:guest@localhost:5672/"
  const conn = new Connection({
    connection_url: url,
  });
  
  const queueName = 'K6 queue'
  
  conn.deleteQueue({
    name: queueName
  })

  console.log(queueName + " queue deleted")
}
