import {Publisher, Consumer} from "k6/x/amqp"; // import Amqp extension

const url = "amqp://guest:guest@localhost:5672/"
const exchangeName = 'K6 exchange'
const queueName = 'K6 general'
const routingKey = 'rkey123'

const exchange = {
  name: exchangeName,
  kind: "direct",
  durable: true,
};

const queue = {
  name: queueName,
  routing_key: routingKey,
  durable: true,
};

const publisher = new Publisher({
  connection_url: url,
  exchange: exchange,
  queue: queue,
});

const consumer = new Consumer({
  connection_url: url,
  exchange: exchange,
  queue: queue,
});

export default function () {
  console.log("Publisher is ready")

  let body = {
    metadata: {
      header1: "Performance Test Message"
    },
    body: {
      field1: "some value"
    }
  }

  publisher.publish({
    exchange  : exchangeName,
    routing_key: routingKey,
    body: JSON.stringify(body),
    content_type: "application/x-msgpack"
  })

  const data = consumer.consume({
    read_timeout: '3s',
    consume_limit: 1,
  })
  console.log('received data: ' + data[0].body) 
}
