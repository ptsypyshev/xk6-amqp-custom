import {Connection} from "k6/x/amqp"; // import Amqp extension

export default function () {
  const url = "amqp://guest:guest@localhost:5672/"
  const conn = new Connection({
    connection_url: url,
  });
  
  const queueName = 'K6 queue'
  console.log('Inspecting ' + queueName)

  console.log(JSON.stringify(conn.inspectQueue({name: queueName}), null, 2))
}
