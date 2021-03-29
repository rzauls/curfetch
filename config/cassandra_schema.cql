CREATE KEYSPACE IF NOT EXISTS "curfetch" WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': 1};
CREATE TABLE curfetch.currencies (code VARCHAR, value VARCHAR, pubdate TIMESTAMP, PRIMARY KEY(code, pubdate));
CREATE INDEX IF NOT EXISTS date ON curfetch.currencies (pubDate);