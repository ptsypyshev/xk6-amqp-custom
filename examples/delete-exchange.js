import {Connection} from "k6/x/amqp"; // import Amqp extension

export default function () {
  const url = "amqp://guest:guest@localhost:5672/"
  const conn = new Connection({
    connection_url: url,
  });
  
  const exchangeName = 'K6 exchange'

  conn.deleteExchange({
    name: exchangeName
  })

  console.log(exchangeName + " exchange deleted")
}
