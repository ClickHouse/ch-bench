#include <clickhouse/client.h>
#include <iostream>

using namespace clickhouse;

int main() {
    Client client(ClientOptions().SetHost("localhost"));

    uint64_t count = 0;
    client.Select("SELECT number FROM system.numbers LIMIT 500000000", [&count](const Block &block)
        {
            count += block.GetRowCount();
        }
    );

    std::cout << count << "\n";
}
