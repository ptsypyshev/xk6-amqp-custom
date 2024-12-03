import {Connection} from "k6/x/amqp"; // import Amqp extension

export default function () {
  const url = "amqp://guest:guest@localhost:5672/"
  const conn = new Connection({
    connection_url: url,
  });
  
  console.log("Connection opened: " + url)

  const exchangeName = 'K6 exchange'

  conn.declareExchange({
    name: exchangeName,
  	kind: 'direct',
    durable: false,
    auto_delete: false,
    internal: false,
    no_wait: false,
	  args: null
  })

  console.log(exchangeName + " exchange is ready")
}
