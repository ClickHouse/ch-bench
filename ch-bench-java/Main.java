import java.sql.Connection;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.sql.Statement;
import java.util.Properties;
import java.util.concurrent.ExecutionException;

import com.clickhouse.client.ClickHouseClient;
import com.clickhouse.client.ClickHouseFormat;
import com.clickhouse.client.ClickHouseNode;
import com.clickhouse.client.ClickHouseProtocol;
import com.clickhouse.client.ClickHouseRecord;
import com.clickhouse.client.ClickHouseResponse;
import com.clickhouse.client.config.ClickHouseClientOption;
import com.clickhouse.jdbc.ClickHouseConnection;
import com.clickhouse.jdbc.ClickHouseDriver;

public class Main {
    static int useJavaClient(String host, int port, String sql) throws InterruptedException, ExecutionException {
        int count = 0;

        ClickHouseNode server = ClickHouseNode.of(host, ClickHouseProtocol.HTTP, port, "system");
        try (ClickHouseClient client = ClickHouseClient.newInstance(server.getProtocol());
                ClickHouseResponse response = client.connect(server).option(ClickHouseClientOption.ASYNC, true)
                        .option(ClickHouseClientOption.FORMAT, ClickHouseFormat.RowBinaryWithNamesAndTypes).query(sql)
                        .execute().get()) {
            for (ClickHouseRecord r : response.records()) {
                count++;
            }
        }

        return count;
    }

    static int useJdbc(String host, int port, String sql) throws SQLException {
        int count = 0;

        String url = new StringBuilder().append("jdbc:ch://").append(host).append(':').append(port).append("/system")
                .toString();
        try (Connection conn = new ClickHouseDriver().connect(url, new Properties());
                Statement stmt = conn.createStatement();
                ResultSet rs = stmt.executeQuery(sql)) {
            while (rs.next()) {
                count++;
            }
        }

        return count;
    }

    public static void main(String[] args) throws Exception {
        String host = System.getProperty("dbHost", "localhost");
        int port = Integer.parseInt(System.getProperty("dbPort", "8123"));
        String sql = "SELECT number FROM system.numbers LIMIT 500000000";

        int count = args != null && args.length > 0 && "client".equals(args[0]) ? useJavaClient(host, port, sql)
                : useJdbc(host, port, sql);
        // System.out.println(count);
    }
}
