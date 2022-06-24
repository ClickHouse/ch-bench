import java.io.EOFException;
import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.sql.Statement;
import java.util.Properties;
import java.util.concurrent.ExecutionException;

import com.clickhouse.client.ClickHouseChecker;
import com.clickhouse.client.ClickHouseClient;
import com.clickhouse.client.ClickHouseFormat;
import com.clickhouse.client.ClickHouseInputStream;
import com.clickhouse.client.ClickHouseNode;
import com.clickhouse.client.ClickHouseNodeSelector;
import com.clickhouse.client.ClickHouseNodes;
import com.clickhouse.client.ClickHouseProtocol;
import com.clickhouse.client.ClickHouseRecord;
import com.clickhouse.client.ClickHouseResponse;
import com.clickhouse.client.config.ClickHouseClientOption;
import com.clickhouse.jdbc.ClickHouseConnection;
import com.clickhouse.jdbc.ClickHouseDriver;

public class Main {
    static int useJavaClient(String url, String sql, String[] args) throws Exception {
        System.out.println("Using Java client");
        int count = 0;

        String deser = args[0];
        ClickHouseNodes servers = ClickHouseNodes.of(url);
        ClickHouseNode server = servers.apply(ClickHouseNodeSelector.EMPTY);
        try (ClickHouseClient client = ClickHouseClient.newInstance(server.getProtocol());
                ClickHouseResponse response = client.connect(server).query(sql).executeAndWait();
                ClickHouseInputStream input = response.getInputStream()) {
            if ("byte".equalsIgnoreCase(deser)) {
                while (true) {
                    try {
                        input.readByte();
                    } catch (EOFException e) {
                        break;
                    }
                    count++;
                }
            } else if ("bytes".equalsIgnoreCase(deser)) {
                int batchSize = 1000;
                if (args.length > 1 && !ClickHouseChecker.isNullOrBlank(args[1])) {
                    batchSize = Integer.parseInt(args[1]);
                }
                while (true) {
                    try {
                        count += input.readBuffer(batchSize).length();
                    } catch (EOFException e) {
                        break;
                    }
                    count++;
                }
            } else if ("long".equalsIgnoreCase(deser)) {
                while (true) {
                    try {
                        input.readBuffer(8).asLong();
                    } catch (EOFException e) {
                        break;
                    }
                    count++;
                }
            } else if ("longs".equalsIgnoreCase(deser)) {
                int batchSize = 1000;
                if (args.length > 1 && !ClickHouseChecker.isNullOrBlank(args[1])) {
                    batchSize = Integer.parseInt(args[1]);
                }
                while (true) {
                    try {
                        count += input.readBuffer(batchSize).asLongArray().length;
                    } catch (EOFException e) {
                        break;
                    }
                }
            } else if ("string".equalsIgnoreCase(deser)) {
                while (true) {
                    try {
                        input.readUnicodeString();
                    } catch (EOFException e) {
                        break;
                    }
                    count++;
                }
            } else {
                for (ClickHouseRecord r : response.records()) {
                    count++;
                }
            }
        }

        return count;
    }

    static int useJdbc(String url, String sql) throws SQLException {
        System.out.println("Using jdbc");
        int count = 0;

        try (Connection conn = DriverManager.getConnection("jdbc:ch:" + url);
                Statement stmt = conn.createStatement();
                ResultSet rs = stmt.executeQuery(sql)) {
            while (rs.next()) {
                count++;
            }
        }

        return count;
    }

    public static void main(String[] args) throws Exception {
        long time = System.nanoTime();
        String url = System.getProperty("url");
        if (ClickHouseChecker.isNullOrBlank(url)) {
            url = "http://localhost?!compress&format=RowBinaryWithNamesAndTypes";
        }
        String sql = System.getProperty("sql");
        if (ClickHouseChecker.isNullOrBlank(sql)) {
            sql = "SELECT number FROM numbers(500000000)";
        }

        int count = args != null && args.length > 0 ? useJavaClient(url, sql, args) : useJdbc(url, sql);
        System.out.println(String.format("%fs\t%d", (System.nanoTime() - time) / 1000000000.0, count));
    }
}
