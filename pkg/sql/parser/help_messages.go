// Code generated by help.awk. DO NOT EDIT.
// GENERATED FILE DO NOT EDIT

package parser

var helpMessages = map[string]HelpMessageBody{
	//line sql.y: 930
	`ALTER`: {
		//line sql.y: 931
		Category: hGroup,
		//line sql.y: 932
		Text: `ALTER TABLE, ALTER INDEX, ALTER VIEW, ALTER DATABASE
`,
	},
	//line sql.y: 940
	`ALTER TABLE`: {
		ShortDescription: `change the definition of a table`,
		//line sql.y: 941
		Category: hDDL,
		//line sql.y: 942
		Text: `
ALTER TABLE [IF EXISTS] <tablename> <command> [, ...]

Commands:
  ALTER TABLE ... ADD [COLUMN] [IF NOT EXISTS] <colname> <type> [<qualifiers...>]
  ALTER TABLE ... ADD <constraint>
  ALTER TABLE ... DROP [COLUMN] [IF EXISTS] <colname> [RESTRICT | CASCADE]
  ALTER TABLE ... DROP CONSTRAINT [IF EXISTS] <constraintname> [RESTRICT | CASCADE]
  ALTER TABLE ... ALTER [COLUMN] <colname> {SET DEFAULT <expr> | DROP DEFAULT}
  ALTER TABLE ... ALTER [COLUMN] <colname> DROP NOT NULL
  ALTER TABLE ... RENAME TO <newname>
  ALTER TABLE ... RENAME [COLUMN] <colname> TO <newname>
  ALTER TABLE ... VALIDATE CONSTRAINT <constraintname>
  ALTER TABLE ... SPLIT AT <selectclause>
  ALTER TABLE ... SCATTER [ FROM ( <exprs...> ) TO ( <exprs...> ) ]

Column qualifiers:
  [CONSTRAINT <constraintname>] {NULL | NOT NULL | UNIQUE | PRIMARY KEY | CHECK (<expr>) | DEFAULT <expr>}
  FAMILY <familyname>, CREATE [IF NOT EXISTS] FAMILY [<familyname>]
  REFERENCES <tablename> [( <colnames...> )]
  COLLATE <collationname>

`,
		//line sql.y: 964
		SeeAlso: `WEBDOCS/alter-table.html
`,
	},
	//line sql.y: 975
	`ALTER VIEW`: {
		ShortDescription: `change the definition of a view`,
		//line sql.y: 976
		Category: hDDL,
		//line sql.y: 977
		Text: `
ALTER VIEW [IF EXISTS] <name> RENAME TO <newname>
`,
		//line sql.y: 979
		SeeAlso: `WEBDOCS/alter-view.html
`,
	},
	//line sql.y: 986
	`ALTER DATABASE`: {
		ShortDescription: `change the definition of a database`,
		//line sql.y: 987
		Category: hDDL,
		//line sql.y: 988
		Text: `
ALTER DATABASE <name> RENAME TO <newname>
`,
		//line sql.y: 990
		SeeAlso: `WEBDOCS/alter-database.html
`,
	},
	//line sql.y: 997
	`ALTER INDEX`: {
		ShortDescription: `change the definition of an index`,
		//line sql.y: 998
		Category: hDDL,
		//line sql.y: 999
		Text: `
ALTER INDEX [IF EXISTS] <idxname> <command>

Commands:
  ALTER INDEX ... RENAME TO <newname>
  ALTER INDEX ... SPLIT AT <selectclause>
  ALTER INDEX ... SCATTER [ FROM ( <exprs...> ) TO ( <exprs...> ) ]

`,
		//line sql.y: 1007
		SeeAlso: `WEBDOCS/alter-index.html
`,
	},
	//line sql.y: 1217
	`BACKUP`: {
		ShortDescription: `back up data to external storage`,
		//line sql.y: 1218
		Category: hCCL,
		//line sql.y: 1219
		Text: `
BACKUP <targets...> TO <location...>
       [ AS OF SYSTEM TIME <expr> ]
       [ INCREMENTAL FROM <location...> ]
       [ WITH <option> [= <value>] [, ...] ]

Targets:
   TABLE <pattern> [, ...]
   DATABASE <databasename> [, ...]

Location:
   "[scheme]://[host]/[path to backup]?[parameters]"

Options:
   INTO_DB
   SKIP_MISSING_FOREIGN_KEYS

`,
		//line sql.y: 1236
		SeeAlso: `RESTORE, WEBDOCS/backup.html
`,
	},
	//line sql.y: 1244
	`RESTORE`: {
		ShortDescription: `restore data from external storage`,
		//line sql.y: 1245
		Category: hCCL,
		//line sql.y: 1246
		Text: `
RESTORE <targets...> FROM <location...>
        [ AS OF SYSTEM TIME <expr> ]
        [ WITH <option> [= <value>] [, ...] ]

Targets:
   TABLE <pattern> [, ...]
   DATABASE <databasename> [, ...]

Locations:
   "[scheme]://[host]/[path to backup]?[parameters]"

Options:
   INTO_DB
   SKIP_MISSING_FOREIGN_KEYS

`,
		//line sql.y: 1262
		SeeAlso: `BACKUP, WEBDOCS/restore.html
`,
	},
	//line sql.y: 1276
	`IMPORT`: {
		ShortDescription: `load data from file in a distributed manner`,
		//line sql.y: 1277
		Category: hCCL,
		//line sql.y: 1278
		Text: `
IMPORT TABLE <tablename>
       { ( <elements> ) | CREATE USING <schemafile> }
       <format>
       DATA ( <datafile> [, ...] )
       [ WITH <option> [= <value>] [, ...] ]

Formats:
   CSV

Options:
   distributed = '...'
   sstsize = '...'
   temp = '...'
   comma = '...'          [CSV-specific]
   comment = '...'        [CSV-specific]
   nullif = '...'         [CSV-specific]

`,
		//line sql.y: 1296
		SeeAlso: `CREATE TABLE
`,
	},
	//line sql.y: 1391
	`CANCEL`: {
		//line sql.y: 1392
		Category: hGroup,
		//line sql.y: 1393
		Text: `CANCEL JOB, CANCEL QUERY
`,
	},
	//line sql.y: 1399
	`CANCEL JOB`: {
		ShortDescription: `cancel a background job`,
		//line sql.y: 1400
		Category: hMisc,
		//line sql.y: 1401
		Text: `CANCEL JOB <jobid>
`,
		//line sql.y: 1402
		SeeAlso: `SHOW JOBS, PAUSE JOBS, RESUME JOB
`,
	},
	//line sql.y: 1410
	`CANCEL QUERY`: {
		ShortDescription: `cancel a running query`,
		//line sql.y: 1411
		Category: hMisc,
		//line sql.y: 1412
		Text: `CANCEL QUERY <queryid>
`,
		//line sql.y: 1413
		SeeAlso: `SHOW QUERIES
`,
	},
	//line sql.y: 1421
	`CREATE`: {
		//line sql.y: 1422
		Category: hGroup,
		//line sql.y: 1423
		Text: `
CREATE DATABASE, CREATE TABLE, CREATE INDEX, CREATE TABLE AS,
CREATE USER, CREATE VIEW
`,
	},
	//line sql.y: 1437
	`DELETE`: {
		ShortDescription: `delete rows from a table`,
		//line sql.y: 1438
		Category: hDML,
		//line sql.y: 1439
		Text: `DELETE FROM <tablename> [WHERE <expr>]
              [LIMIT <expr>]
              [RETURNING <exprs...>]
`,
		//line sql.y: 1442
		SeeAlso: `WEBDOCS/delete.html
`,
	},
	//line sql.y: 1455
	`DISCARD`: {
		ShortDescription: `reset the session to its initial state`,
		//line sql.y: 1456
		Category: hCfg,
		//line sql.y: 1457
		Text: `DISCARD ALL
`,
	},
	//line sql.y: 1469
	`DROP`: {
		//line sql.y: 1470
		Category: hGroup,
		//line sql.y: 1471
		Text: `DROP DATABASE, DROP INDEX, DROP TABLE, DROP VIEW, DROP USER
`,
	},
	//line sql.y: 1480
	`DROP VIEW`: {
		ShortDescription: `remove a view`,
		//line sql.y: 1481
		Category: hDDL,
		//line sql.y: 1482
		Text: `DROP VIEW [IF EXISTS] <tablename> [, ...] [CASCADE | RESTRICT]
`,
		//line sql.y: 1483
		SeeAlso: `WEBDOCS/drop-index.html
`,
	},
	//line sql.y: 1495
	`DROP TABLE`: {
		ShortDescription: `remove a table`,
		//line sql.y: 1496
		Category: hDDL,
		//line sql.y: 1497
		Text: `DROP TABLE [IF EXISTS] <tablename> [, ...] [CASCADE | RESTRICT]
`,
		//line sql.y: 1498
		SeeAlso: `WEBDOCS/drop-table.html
`,
	},
	//line sql.y: 1510
	`DROP INDEX`: {
		ShortDescription: `remove an index`,
		//line sql.y: 1511
		Category: hDDL,
		//line sql.y: 1512
		Text: `DROP INDEX [IF EXISTS] <idxname> [, ...] [CASCADE | RESTRICT]
`,
		//line sql.y: 1513
		SeeAlso: `WEBDOCS/drop-index.html
`,
	},
	//line sql.y: 1533
	`DROP DATABASE`: {
		ShortDescription: `remove a database`,
		//line sql.y: 1534
		Category: hDDL,
		//line sql.y: 1535
		Text: `DROP DATABASE [IF EXISTS] <databasename> [CASCADE | RESTRICT]
`,
		//line sql.y: 1536
		SeeAlso: `WEBDOCS/drop-database.html
`,
	},
	//line sql.y: 1556
	`DROP USER`: {
		ShortDescription: `remove a user`,
		//line sql.y: 1557
		Category: hPriv,
		//line sql.y: 1558
		Text: `DROP USER [IF EXISTS] <user> [, ...]
`,
		//line sql.y: 1559
		SeeAlso: `CREATE USER, SHOW USERS
`,
	},
	//line sql.y: 1601
	`EXPLAIN`: {
		ShortDescription: `show the logical plan of a query`,
		//line sql.y: 1602
		Category: hMisc,
		//line sql.y: 1603
		Text: `
EXPLAIN <statement>
EXPLAIN [( [PLAN ,] <planoptions...> )] <statement>

Explainable statements:
    SELECT, CREATE, DROP, ALTER, INSERT, UPSERT, UPDATE, DELETE,
    SHOW, EXPLAIN, EXECUTE

Plan options:
    TYPES, EXPRS, METADATA, QUALIFY, INDENT, VERBOSE, DIST_SQL

`,
		//line sql.y: 1614
		SeeAlso: `WEBDOCS/explain.html
`,
	},
	//line sql.y: 1672
	`PREPARE`: {
		ShortDescription: `prepare a statement for later execution`,
		//line sql.y: 1673
		Category: hMisc,
		//line sql.y: 1674
		Text: `PREPARE <name> [ ( <types...> ) ] AS <query>
`,
		//line sql.y: 1675
		SeeAlso: `EXECUTE, DEALLOCATE, DISCARD
`,
	},
	//line sql.y: 1697
	`EXECUTE`: {
		ShortDescription: `execute a statement prepared previously`,
		//line sql.y: 1698
		Category: hMisc,
		//line sql.y: 1699
		Text: `EXECUTE <name> [ ( <exprs...> ) ]
`,
		//line sql.y: 1700
		SeeAlso: `PREPARE, DEALLOCATE, DISCARD
`,
	},
	//line sql.y: 1723
	`DEALLOCATE`: {
		ShortDescription: `remove a prepared statement`,
		//line sql.y: 1724
		Category: hMisc,
		//line sql.y: 1725
		Text: `DEALLOCATE [PREPARE] { <name> | ALL }
`,
		//line sql.y: 1726
		SeeAlso: `PREPARE, EXECUTE, DISCARD
`,
	},
	//line sql.y: 1746
	`GRANT`: {
		ShortDescription: `define access privileges`,
		//line sql.y: 1747
		Category: hPriv,
		//line sql.y: 1748
		Text: `
GRANT {ALL | <privileges...> } ON <targets...> TO <grantees...>

Privileges:
  CREATE, DROP, GRANT, SELECT, INSERT, DELETE, UPDATE

Targets:
  DATABASE <databasename> [, ...]
  [TABLE] [<databasename> .] { <tablename> | * } [, ...]

`,
		//line sql.y: 1758
		SeeAlso: `REVOKE, WEBDOCS/grant.html
`,
	},
	//line sql.y: 1766
	`REVOKE`: {
		ShortDescription: `remove access privileges`,
		//line sql.y: 1767
		Category: hPriv,
		//line sql.y: 1768
		Text: `
REVOKE {ALL | <privileges...> } ON <targets...> FROM <grantees...>

Privileges:
  CREATE, DROP, GRANT, SELECT, INSERT, DELETE, UPDATE

Targets:
  DATABASE <databasename> [, <databasename>]...
  [TABLE] [<databasename> .] { <tablename> | * } [, ...]

`,
		//line sql.y: 1778
		SeeAlso: `GRANT, WEBDOCS/revoke.html
`,
	},
	//line sql.y: 1865
	`RESET`: {
		ShortDescription: `reset a session variable to its default value`,
		//line sql.y: 1866
		Category: hCfg,
		//line sql.y: 1867
		Text: `RESET [SESSION] <var>
`,
		//line sql.y: 1868
		SeeAlso: `RESET CLUSTER SETTING, WEBDOCS/set-vars.html
`,
	},
	//line sql.y: 1880
	`RESET CLUSTER SETTING`: {
		ShortDescription: `reset a cluster setting to its default value`,
		//line sql.y: 1881
		Category: hCfg,
		//line sql.y: 1882
		Text: `RESET CLUSTER SETTING <var>
`,
		//line sql.y: 1883
		SeeAlso: `SET CLUSTER SETTING, RESET
`,
	},
	//line sql.y: 1909
	`SET CLUSTER SETTING`: {
		ShortDescription: `change a cluster setting`,
		//line sql.y: 1910
		Category: hCfg,
		//line sql.y: 1911
		Text: `SET CLUSTER SETTING <var> { TO | = } <value>
`,
		//line sql.y: 1912
		SeeAlso: `SHOW CLUSTER SETTING, RESET CLUSTER SETTING, SET SESSION,
WEBDOCS/cluster-settings.html
`,
	},
	//line sql.y: 1933
	`SET SESSION`: {
		ShortDescription: `change a session variable`,
		//line sql.y: 1934
		Category: hCfg,
		//line sql.y: 1935
		Text: `
SET [SESSION] <var> { TO | = } <values...>
SET [SESSION] TIME ZONE <tz>
SET [SESSION] CHARACTERISTICS AS TRANSACTION ISOLATION LEVEL { SNAPSHOT | SERIALIZABLE }

`,
		//line sql.y: 1940
		SeeAlso: `SHOW SESSION, RESET, DISCARD, SHOW, SET CLUSTER SETTING, SET TRANSACTION,
WEBDOCS/set-vars.html
`,
	},
	//line sql.y: 1957
	`SET TRANSACTION`: {
		ShortDescription: `configure the transaction settings`,
		//line sql.y: 1958
		Category: hTxn,
		//line sql.y: 1959
		Text: `
SET [SESSION] TRANSACTION <txnparameters...>

Transaction parameters:
   ISOLATION LEVEL { SNAPSHOT | SERIALIZABLE }
   PRIORITY { LOW | NORMAL | HIGH }

`,
		//line sql.y: 1966
		SeeAlso: `SHOW TRANSACTION, SET SESSION,
WEBDOCS/set-transaction.html
`,
	},
	//line sql.y: 2105
	`SHOW`: {
		//line sql.y: 2106
		Category: hGroup,
		//line sql.y: 2107
		Text: `
SHOW SESSION, SHOW CLUSTER SETTING, SHOW DATABASES, SHOW TABLES, SHOW COLUMNS, SHOW INDEXES,
SHOW CONSTRAINTS, SHOW CREATE TABLE, SHOW CREATE VIEW, SHOW USERS, SHOW TRANSACTION, SHOW BACKUP,
SHOW JOBS, SHOW QUERIES, SHOW SESSIONS, SHOW TRACE
`,
	},
	//line sql.y: 2132
	`SHOW SESSION`: {
		ShortDescription: `display session variables`,
		//line sql.y: 2133
		Category: hCfg,
		//line sql.y: 2134
		Text: `SHOW [SESSION] { <var> | ALL }
`,
		//line sql.y: 2135
		SeeAlso: `WEBDOCS/show-vars.html
`,
	},
	//line sql.y: 2156
	`SHOW BACKUP`: {
		ShortDescription: `list backup contents`,
		//line sql.y: 2157
		Category: hCCL,
		//line sql.y: 2158
		Text: `SHOW BACKUP <location>
`,
		//line sql.y: 2159
		SeeAlso: `WEBDOCS/show-backup.html
`,
	},
	//line sql.y: 2167
	`SHOW CLUSTER SETTING`: {
		ShortDescription: `display cluster settings`,
		//line sql.y: 2168
		Category: hCfg,
		//line sql.y: 2169
		Text: `
SHOW CLUSTER SETTING <var>
SHOW ALL CLUSTER SETTINGS
`,
		//line sql.y: 2172
		SeeAlso: `WEBDOCS/cluster-settings.html
`,
	},
	//line sql.y: 2189
	`SHOW COLUMNS`: {
		ShortDescription: `list columns in relation`,
		//line sql.y: 2190
		Category: hDDL,
		//line sql.y: 2191
		Text: `SHOW COLUMNS FROM <tablename>
`,
		//line sql.y: 2192
		SeeAlso: `WEBDOCS/show-columns.html
`,
	},
	//line sql.y: 2200
	`SHOW DATABASES`: {
		ShortDescription: `list databases`,
		//line sql.y: 2201
		Category: hDDL,
		//line sql.y: 2202
		Text: `SHOW DATABASES
`,
		//line sql.y: 2203
		SeeAlso: `WEBDOCS/show-databases.html
`,
	},
	//line sql.y: 2211
	`SHOW GRANTS`: {
		ShortDescription: `list grants`,
		//line sql.y: 2212
		Category: hPriv,
		//line sql.y: 2213
		Text: `SHOW GRANTS [ON <targets...>] [FOR <users...>]
`,
		//line sql.y: 2214
		SeeAlso: `WEBDOCS/show-grants.html
`,
	},
	//line sql.y: 2222
	`SHOW INDEXES`: {
		ShortDescription: `list indexes`,
		//line sql.y: 2223
		Category: hDDL,
		//line sql.y: 2224
		Text: `SHOW INDEXES FROM <tablename>
`,
		//line sql.y: 2225
		SeeAlso: `WEBDOCS/show-index.html
`,
	},
	//line sql.y: 2243
	`SHOW CONSTRAINTS`: {
		ShortDescription: `list constraints`,
		//line sql.y: 2244
		Category: hDDL,
		//line sql.y: 2245
		Text: `SHOW CONSTRAINTS FROM <tablename>
`,
		//line sql.y: 2246
		SeeAlso: `WEBDOCS/show-constraints.html
`,
	},
	//line sql.y: 2259
	`SHOW QUERIES`: {
		ShortDescription: `list running queries`,
		//line sql.y: 2260
		Category: hMisc,
		//line sql.y: 2261
		Text: `SHOW [CLUSTER | LOCAL] QUERIES
`,
		//line sql.y: 2262
		SeeAlso: `CANCEL QUERY
`,
	},
	//line sql.y: 2278
	`SHOW JOBS`: {
		ShortDescription: `list background jobs`,
		//line sql.y: 2279
		Category: hMisc,
		//line sql.y: 2280
		Text: `SHOW JOBS
`,
		//line sql.y: 2281
		SeeAlso: `CANCEL JOB, PAUSE JOB, RESUME JOB
`,
	},
	//line sql.y: 2289
	`SHOW TRACE`: {
		ShortDescription: `display an execution trace`,
		//line sql.y: 2290
		Category: hMisc,
		//line sql.y: 2291
		Text: `
SHOW [KV] TRACE FOR SESSION
SHOW [KV] TRACE FOR <statement>
`,
		//line sql.y: 2294
		SeeAlso: `EXPLAIN
`,
	},
	//line sql.y: 2315
	`SHOW SESSIONS`: {
		ShortDescription: `list open client sessions`,
		//line sql.y: 2316
		Category: hMisc,
		//line sql.y: 2317
		Text: `SHOW [CLUSTER | LOCAL] SESSIONS
`,
	},
	//line sql.y: 2333
	`SHOW TABLES`: {
		ShortDescription: `list tables`,
		//line sql.y: 2334
		Category: hDDL,
		//line sql.y: 2335
		Text: `SHOW TABLES [FROM <databasename>]
`,
		//line sql.y: 2336
		SeeAlso: `WEBDOCS/show-tables.html
`,
	},
	//line sql.y: 2348
	`SHOW TRANSACTION`: {
		ShortDescription: `display current transaction properties`,
		//line sql.y: 2349
		Category: hCfg,
		//line sql.y: 2350
		Text: `SHOW TRANSACTION {ISOLATION LEVEL | PRIORITY | STATUS}
`,
		//line sql.y: 2351
		SeeAlso: `WEBDOCS/show-transaction.html
`,
	},
	//line sql.y: 2370
	`SHOW CREATE TABLE`: {
		ShortDescription: `display the CREATE TABLE statement for a table`,
		//line sql.y: 2371
		Category: hDDL,
		//line sql.y: 2372
		Text: `SHOW CREATE TABLE <tablename>
`,
		//line sql.y: 2373
		SeeAlso: `WEBDOCS/show-create-table.html
`,
	},
	//line sql.y: 2381
	`SHOW CREATE VIEW`: {
		ShortDescription: `display the CREATE VIEW statement for a view`,
		//line sql.y: 2382
		Category: hDDL,
		//line sql.y: 2383
		Text: `SHOW CREATE VIEW <viewname>
`,
		//line sql.y: 2384
		SeeAlso: `WEBDOCS/show-create-view.html
`,
	},
	//line sql.y: 2392
	`SHOW USERS`: {
		ShortDescription: `list defined users`,
		//line sql.y: 2393
		Category: hPriv,
		//line sql.y: 2394
		Text: `SHOW USERS
`,
		//line sql.y: 2395
		SeeAlso: `CREATE USER, DROP USER, WEBDOCS/show-users.html
`,
	},
	//line sql.y: 2441
	`PAUSE JOB`: {
		ShortDescription: `pause a background job`,
		//line sql.y: 2442
		Category: hMisc,
		//line sql.y: 2443
		Text: `PAUSE JOB <jobid>
`,
		//line sql.y: 2444
		SeeAlso: `SHOW JOBS, CANCEL JOB, RESUME JOB
`,
	},
	//line sql.y: 2452
	`CREATE TABLE`: {
		ShortDescription: `create a new table`,
		//line sql.y: 2453
		Category: hDDL,
		//line sql.y: 2454
		Text: `
CREATE TABLE [IF NOT EXISTS] <tablename> ( <elements...> ) [<interleave>]
CREATE TABLE [IF NOT EXISTS] <tablename> [( <colnames...> )] AS <source>

Table elements:
   <name> <type> [<qualifiers...>]
   [UNIQUE] INDEX [<name>] ( <colname> [ASC | DESC] [, ...] )
                           [STORING ( <colnames...> )] [<interleave>]
   FAMILY [<name>] ( <colnames...> )
   [CONSTRAINT <name>] <constraint>

Table constraints:
   PRIMARY KEY ( <colnames...> )
   FOREIGN KEY ( <colnames...> ) REFERENCES <tablename> [( <colnames...> )]
   UNIQUE ( <colnames... ) [STORING ( <colnames...> )] [<interleave>]
   CHECK ( <expr> )

Column qualifiers:
  [CONSTRAINT <constraintname>] {NULL | NOT NULL | UNIQUE | PRIMARY KEY | CHECK (<expr>) | DEFAULT <expr>}
  FAMILY <familyname>, CREATE [IF NOT EXISTS] FAMILY [<familyname>]
  REFERENCES <tablename> [( <colnames...> )]
  COLLATE <collationname>

Interleave clause:
   INTERLEAVE IN PARENT <tablename> ( <colnames...> ) [CASCADE | RESTRICT]

`,
		//line sql.y: 2480
		SeeAlso: `SHOW TABLES, CREATE VIEW, SHOW CREATE TABLE,
WEBDOCS/create-table.html
WEBDOCS/create-table-as.html
`,
	},
	//line sql.y: 2878
	`TRUNCATE`: {
		ShortDescription: `empty one or more tables`,
		//line sql.y: 2879
		Category: hDML,
		//line sql.y: 2880
		Text: `TRUNCATE [TABLE] <tablename> [, ...] [CASCADE | RESTRICT]
`,
		//line sql.y: 2881
		SeeAlso: `WEBDOCS/truncate.html
`,
	},
	//line sql.y: 2889
	`CREATE USER`: {
		ShortDescription: `define a new user`,
		//line sql.y: 2890
		Category: hPriv,
		//line sql.y: 2891
		Text: `CREATE USER <name> [ [WITH] PASSWORD <passwd> ]
`,
		//line sql.y: 2892
		SeeAlso: `DROP USER, SHOW USERS, WEBDOCS/create-user.html
`,
	},
	//line sql.y: 2910
	`CREATE VIEW`: {
		ShortDescription: `create a new view`,
		//line sql.y: 2911
		Category: hDDL,
		//line sql.y: 2912
		Text: `CREATE VIEW <viewname> [( <colnames...> )] AS <source>
`,
		//line sql.y: 2913
		SeeAlso: `CREATE TABLE, SHOW CREATE VIEW, WEBDOCS/create-view.html
`,
	},
	//line sql.y: 2927
	`CREATE INDEX`: {
		ShortDescription: `create a new index`,
		//line sql.y: 2928
		Category: hDDL,
		//line sql.y: 2929
		Text: `
CREATE [UNIQUE] INDEX [IF NOT EXISTS] [<idxname>]
       ON <tablename> ( <colname> [ASC | DESC] [, ...] )
       [STORING ( <colnames...> )] [<interleave>]

Interleave clause:
   INTERLEAVE IN PARENT <tablename> ( <colnames...> ) [CASCADE | RESTRICT]

`,
		//line sql.y: 2937
		SeeAlso: `CREATE TABLE, SHOW INDEXES, SHOW CREATE INDEX,
WEBDOCS/create-index.html
`,
	},
	//line sql.y: 3076
	`RELEASE`: {
		ShortDescription: `complete a retryable block`,
		//line sql.y: 3077
		Category: hTxn,
		//line sql.y: 3078
		Text: `RELEASE [SAVEPOINT] cockroach_restart
`,
		//line sql.y: 3079
		SeeAlso: `SAVEPOINT, WEBDOCS/savepoint.html
`,
	},
	//line sql.y: 3087
	`RESUME JOB`: {
		ShortDescription: `resume a background job`,
		//line sql.y: 3088
		Category: hMisc,
		//line sql.y: 3089
		Text: `RESUME JOB <jobid>
`,
		//line sql.y: 3090
		SeeAlso: `SHOW JOBS, CANCEL JOB, PAUSE JOB
`,
	},
	//line sql.y: 3098
	`SAVEPOINT`: {
		ShortDescription: `start a retryable block`,
		//line sql.y: 3099
		Category: hTxn,
		//line sql.y: 3100
		Text: `SAVEPOINT cockroach_restart
`,
		//line sql.y: 3101
		SeeAlso: `RELEASE, WEBDOCS/savepoint.html
`,
	},
	//line sql.y: 3115
	`BEGIN`: {
		ShortDescription: `start a transaction`,
		//line sql.y: 3116
		Category: hTxn,
		//line sql.y: 3117
		Text: `
BEGIN [TRANSACTION] [ <txnparameter> [[,] ...] ]
START TRANSACTION [ <txnparameter> [[,] ...] ]

Transaction parameters:
   ISOLATION LEVEL { SNAPSHOT | SERIALIZABLE }
   PRIORITY { LOW | NORMAL | HIGH }

`,
		//line sql.y: 3125
		SeeAlso: `COMMIT, ROLLBACK, WEBDOCS/begin-transaction.html
`,
	},
	//line sql.y: 3138
	`COMMIT`: {
		ShortDescription: `commit the current transaction`,
		//line sql.y: 3139
		Category: hTxn,
		//line sql.y: 3140
		Text: `
COMMIT [TRANSACTION]
END [TRANSACTION]
`,
		//line sql.y: 3143
		SeeAlso: `BEGIN, ROLLBACK, WEBDOCS/commit-transaction.html
`,
	},
	//line sql.y: 3156
	`ROLLBACK`: {
		ShortDescription: `abort the current transaction`,
		//line sql.y: 3157
		Category: hTxn,
		//line sql.y: 3158
		Text: `ROLLBACK [TRANSACTION] [TO [SAVEPOINT] cockroach_restart]
`,
		//line sql.y: 3159
		SeeAlso: `BEGIN, COMMIT, SAVEPOINT, WEBDOCS/rollback-transaction.html
`,
	},
	//line sql.y: 3273
	`CREATE DATABASE`: {
		ShortDescription: `create a new database`,
		//line sql.y: 3274
		Category: hDDL,
		//line sql.y: 3275
		Text: `CREATE DATABASE [IF NOT EXISTS] <name>
`,
		//line sql.y: 3276
		SeeAlso: `WEBDOCS/create-database.html
`,
	},
	//line sql.y: 3345
	`INSERT`: {
		ShortDescription: `create new rows in a table`,
		//line sql.y: 3346
		Category: hDML,
		//line sql.y: 3347
		Text: `
INSERT INTO <tablename> [[AS] <name>] [( <colnames...> )]
       <selectclause>
       [ON CONFLICT [( <colnames...> )] {DO UPDATE SET ... [WHERE <expr>] | DO NOTHING}]
       [RETURNING <exprs...>]
`,
		//line sql.y: 3352
		SeeAlso: `UPSERT, UPDATE, DELETE, WEBDOCS/insert.html
`,
	},
	//line sql.y: 3369
	`UPSERT`: {
		ShortDescription: `create or replace rows in a table`,
		//line sql.y: 3370
		Category: hDML,
		//line sql.y: 3371
		Text: `
UPSERT INTO <tablename> [AS <name>] [( <colnames...> )]
       <selectclause>
       [RETURNING <exprs...>]
`,
		//line sql.y: 3375
		SeeAlso: `INSERT, UPDATE, DELETE, WEBDOCS/upsert.html
`,
	},
	//line sql.y: 3451
	`UPDATE`: {
		ShortDescription: `update rows of a table`,
		//line sql.y: 3452
		Category: hDML,
		//line sql.y: 3453
		Text: `UPDATE <tablename> [[AS] <name>] SET ... [WHERE <expr>] [RETURNING <exprs...>]
`,
		//line sql.y: 3454
		SeeAlso: `INSERT, UPSERT, DELETE, WEBDOCS/update.html
`,
	},
	//line sql.y: 3622
	`<SELECTCLAUSE>`: {
		ShortDescription: `access tabular data`,
		//line sql.y: 3623
		Category: hDML,
		//line sql.y: 3624
		Text: `
Select clause:
  TABLE <tablename>
  VALUES ( <exprs...> ) [ , ... ]
  SELECT ... [ { INTERSECT | UNION | EXCEPT } [ ALL | DISTINCT ] <selectclause> ]
`,
	},
	//line sql.y: 3635
	`SELECT`: {
		ShortDescription: `retrieve rows from a data source and compute a result`,
		//line sql.y: 3636
		Category: hDML,
		//line sql.y: 3637
		Text: `
SELECT [DISTINCT]
       { <expr> [[AS] <name>] | [ [<dbname>.] <tablename>. ] * } [, ...]
       [ FROM <source> ]
       [ WHERE <expr> ]
       [ GROUP BY <expr> [ , ... ] ]
       [ HAVING <expr> ]
       [ WINDOW <name> AS ( <definition> ) ]
       [ { UNION | INTERSECT | EXCEPT } [ ALL | DISTINCT ] <selectclause> ]
       [ ORDER BY <expr> [ ASC | DESC ] [, ...] ]
       [ LIMIT { <expr> | ALL } ]
       [ OFFSET <expr> [ ROW | ROWS ] ]
`,
		//line sql.y: 3649
		SeeAlso: `WEBDOCS/select.html
`,
	},
	//line sql.y: 3709
	`TABLE`: {
		ShortDescription: `select an entire table`,
		//line sql.y: 3710
		Category: hDML,
		//line sql.y: 3711
		Text: `TABLE <tablename>
`,
		//line sql.y: 3712
		SeeAlso: `SELECT, VALUES, WEBDOCS/table-expressions.html
`,
	},
	//line sql.y: 3955
	`VALUES`: {
		ShortDescription: `select a given set of values`,
		//line sql.y: 3956
		Category: hDML,
		//line sql.y: 3957
		Text: `VALUES ( <exprs...> ) [, ...]
`,
		//line sql.y: 3958
		SeeAlso: `SELECT, TABLE, WEBDOCS/table-expressions.html
`,
	},
	//line sql.y: 4063
	`<SOURCE>`: {
		ShortDescription: `define a data source for SELECT`,
		//line sql.y: 4064
		Category: hDML,
		//line sql.y: 4065
		Text: `
Data sources:
  <tablename> [ @ { <idxname> | <indexhint> } ]
  <tablefunc> ( <exprs...> )
  ( { <selectclause> | <source> } )
  <source> [AS] <alias> [( <colnames...> )]
  <source> { [INNER] | { LEFT | RIGHT } [OUTER] } JOIN <source> ON <expr>
  <source> { [INNER] | { LEFT | RIGHT } [OUTER] } JOIN <source> USING ( <colnames...> )
  <source> NATURAL { [INNER] | { LEFT | RIGHT } [OUTER] } JOIN <source>
  <source> CROSS JOIN <source>
  <source> WITH ORDINALITY
  '[' EXPLAIN ... ']'
  '[' SHOW ... ']'

Index hints:
  '{' FORCE_INDEX = <idxname> [, ...] '}'
  '{' NO_INDEX_JOIN [, ...] '}'

`,
		//line sql.y: 4083
		SeeAlso: `WEBDOCS/table-expressions.html
`,
	},
}
