import sys
from clickhouse_driver import Client

if __name__ == '__main__':
    query = "SELECT number FROM system.numbers LIMIT 500000000"
    client = Client.from_url('clickhouse://localhost')

    for row in client.execute_iter(query):
        pass
