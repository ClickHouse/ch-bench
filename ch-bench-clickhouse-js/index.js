import { createClient } from '@clickhouse/client'

(async () => {
  const client = createClient({
    compression: {
      request: false,
      response: false,
    },
  })
  const rows = await client.query({
    query: 'SELECT number FROM system.numbers_mt LIMIT 500000000',
    format: 'CSV',
  })
  const stream = rows.stream()
  stream.on('data', (rows) => {
    rows.forEach(() => {
      //
    })
  })
  await new Promise((resolve) => {
    stream.on('end', () => {
      resolve(0)
    })
  })
  await client.close()
})()
