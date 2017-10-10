// Code generated by help.awk. DO NOT EDIT.
// GENERATED FILE DO NOT EDIT

package parser

var helpMessages = map[string]HelpMessageBody{
	//line sql.y: 934
	`ALTER`: {
		//line sql.y: 935
		Category: hGroup,
		//line sql.y: 936
		Text: `ALTER TABLE, ALTER INDEX, ALTER VIEW, ALTER DATABASE
`,
	},
	//line sql.y: 944
	`ALTER TABLE`: {
		ShortDescription: `change the definition of a table`,
		//line sql.y: 945
		Category: hDDL,
		//line sql.y: 946
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
		//line sql.y: 968
		SeeAlso: `WEBDOCS/alter-table.html
`,
	},
	//line sql.y: 979
	`ALTER VIEW`: {
		ShortDescription: `change the definition of a view`,
		//line sql.y: 980
		Category: hDDL,
		//line sql.y: 981
		Text: `
ALTER VIEW [IF EXISTS] <name> RENAME TO <newname>
`,
		//line sql.y: 983
		SeeAlso: `WEBDOCS/alter-view.html
`,
	},
	//line sql.y: 990
	`ALTER DATABASE`: {
		ShortDescription: `change the definition of a database`,
		//line sql.y: 991
		Category: hDDL,
		//line sql.y: 992
		Text: `
ALTER DATABASE <name> RENAME TO <newname>
`,
		//line sql.y: 994
		SeeAlso: `WEBDOCS/alter-database.html
`,
	},
	//line sql.y: 1001
	`ALTER INDEX`: {
		ShortDescription: `change the definition of an index`,
		//line sql.y: 1002
		Category: hDDL,
		//line sql.y: 1003
		Text: `
ALTER INDEX [IF EXISTS] <idxname> <command>

Commands:
  ALTER INDEX ... RENAME TO <newname>
  ALTER INDEX ... SPLIT AT <selectclause>
  ALTER INDEX ... SCATTER [ FROM ( <exprs...> ) TO ( <exprs...> ) ]

`,
		//line sql.y: 1011
		SeeAlso: `WEBDOCS/alter-index.html
`,
	},
	//line sql.y: 1221
	`BACKUP`: {
		ShortDescription: `back up data to external storage`,
		//line sql.y: 1222
		Category: hCCL,
		//line sql.y: 1223
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
		//line sql.y: 1240
		SeeAlso: `RESTORE, WEBDOCS/backup.html
`,
	},
	//line sql.y: 1248
	`RESTORE`: {
		ShortDescription: `restore data from external storage`,
		//line sql.y: 1249
		Category: hCCL,
		//line sql.y: 1250
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
		//line sql.y: 1266
		SeeAlso: `BACKUP, WEBDOCS/restore.html
`,
	},
	//line sql.y: 1280
	`IMPORT`: {
		ShortDescription: `load data from file in a distributed manner`,
		//line sql.y: 1281
		Category: hCCL,
		//line sql.y: 1282
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
		//line sql.y: 1300
		SeeAlso: `CREATE TABLE
`,
	},
	//line sql.y: 1395
	`CANCEL`: {
		//line sql.y: 1396
		Category: hGroup,
		//line sql.y: 1397
		Text: `CANCEL JOB, CANCEL QUERY
`,
	},
	//line sql.y: 1403
	`CANCEL JOB`: {
		ShortDescription: `cancel a background job`,
		//line sql.y: 1404
		Category: hMisc,
		//line sql.y: 1405
		Text: `CANCEL JOB <jobid>
`,
		//line sql.y: 1406
		SeeAlso: `SHOW JOBS, PAUSE JOBS, RESUME JOB
`,
	},
	//line sql.y: 1414
	`CANCEL QUERY`: {
		ShortDescription: `cancel a running query`,
		//line sql.y: 1415
		Category: hMisc,
		//line sql.y: 1416
		Text: `CANCEL QUERY <queryid>
`,
		//line sql.y: 1417
		SeeAlso: `SHOW QUERIES
`,
	},
	//line sql.y: 1425
	`CREATE`: {
		//line sql.y: 1426
		Category: hGroup,
		//line sql.y: 1427
		Text: `
CREATE DATABASE, CREATE TABLE, CREATE INDEX, CREATE TABLE AS,
CREATE USER, CREATE VIEW
`,
	},
	//line sql.y: 1441
	`DELETE`: {
		ShortDescription: `delete rows from a table`,
		//line sql.y: 1442
		Category: hDML,
		//line sql.y: 1443
		Text: `DELETE FROM <tablename> [WHERE <expr>]
              [LIMIT <expr>]
              [RETURNING <exprs...>]
`,
		//line sql.y: 1446
		SeeAlso: `WEBDOCS/delete.html
`,
	},
	//line sql.y: 1459
	`DISCARD`: {
		ShortDescription: `reset the session to its initial state`,
		//line sql.y: 1460
		Category: hCfg,
		//line sql.y: 1461
		Text: `DISCARD ALL
`,
	},
	//line sql.y: 1473
	`DROP`: {
		//line sql.y: 1474
		Category: hGroup,
		//line sql.y: 1475
		Text: `DROP DATABASE, DROP INDEX, DROP TABLE, DROP VIEW, DROP USER
`,
	},
	//line sql.y: 1484
	`DROP VIEW`: {
		ShortDescription: `remove a view`,
		//line sql.y: 1485
		Category: hDDL,
		//line sql.y: 1486
		Text: `DROP VIEW [IF EXISTS] <tablename> [, ...] [CASCADE | RESTRICT]
`,
		//line sql.y: 1487
		SeeAlso: `WEBDOCS/drop-index.html
`,
	},
	//line sql.y: 1499
	`DROP TABLE`: {
		ShortDescription: `remove a table`,
		//line sql.y: 1500
		Category: hDDL,
		//line sql.y: 1501
		Text: `DROP TABLE [IF EXISTS] <tablename> [, ...] [CASCADE | RESTRICT]
`,
		//line sql.y: 1502
		SeeAlso: `WEBDOCS/drop-table.html
`,
	},
	//line sql.y: 1514
	`DROP INDEX`: {
		ShortDescription: `remove an index`,
		//line sql.y: 1515
		Category: hDDL,
		//line sql.y: 1516
		Text: `DROP INDEX [IF EXISTS] <idxname> [, ...] [CASCADE | RESTRICT]
`,
		//line sql.y: 1517
		SeeAlso: `WEBDOCS/drop-index.html
`,
	},
	//line sql.y: 1537
	`DROP DATABASE`: {
		ShortDescription: `remove a database`,
		//line sql.y: 1538
		Category: hDDL,
		//line sql.y: 1539
		Text: `DROP DATABASE [IF EXISTS] <databasename> [CASCADE | RESTRICT]
`,
		//line sql.y: 1540
		SeeAlso: `WEBDOCS/drop-database.html
`,
	},
	//line sql.y: 1560
	`DROP USER`: {
		ShortDescription: `remove a user`,
		//line sql.y: 1561
		Category: hPriv,
		//line sql.y: 1562
		Text: `DROP USER [IF EXISTS] <user> [, ...]
`,
		//line sql.y: 1563
		SeeAlso: `CREATE USER, SHOW USERS
`,
	},
	//line sql.y: 1605
	`EXPLAIN`: {
		ShortDescription: `show the logical plan of a query`,
		//line sql.y: 1606
		Category: hMisc,
		//line sql.y: 1607
		Text: `
EXPLAIN <statement>
EXPLAIN [( [PLAN ,] <planoptions...> )] <statement>

Explainable statements:
    SELECT, CREATE, DROP, ALTER, INSERT, UPSERT, UPDATE, DELETE,
    SHOW, EXPLAIN, EXECUTE

Plan options:
    TYPES, EXPRS, METADATA, QUALIFY, INDENT, VERBOSE, DIST_SQL

`,
		//line sql.y: 1618
		SeeAlso: `WEBDOCS/explain.html
`,
	},
	//line sql.y: 1676
	`PREPARE`: {
		ShortDescription: `prepare a statement for later execution`,
		//line sql.y: 1677
		Category: hMisc,
		//line sql.y: 1678
		Text: `PREPARE <name> [ ( <types...> ) ] AS <query>
`,
		//line sql.y: 1679
		SeeAlso: `EXECUTE, DEALLOCATE, DISCARD
`,
	},
	//line sql.y: 1701
	`EXECUTE`: {
		ShortDescription: `execute a statement prepared previously`,
		//line sql.y: 1702
		Category: hMisc,
		//line sql.y: 1703
		Text: `EXECUTE <name> [ ( <exprs...> ) ]
`,
		//line sql.y: 1704
		SeeAlso: `PREPARE, DEALLOCATE, DISCARD
`,
	},
	//line sql.y: 1727
	`DEALLOCATE`: {
		ShortDescription: `remove a prepared statement`,
		//line sql.y: 1728
		Category: hMisc,
		//line sql.y: 1729
		Text: `DEALLOCATE [PREPARE] { <name> | ALL }
`,
		//line sql.y: 1730
		SeeAlso: `PREPARE, EXECUTE, DISCARD
`,
	},
	//line sql.y: 1750
	`GRANT`: {
		ShortDescription: `define access privileges`,
		//line sql.y: 1751
		Category: hPriv,
		//line sql.y: 1752
		Text: `
GRANT {ALL | <privileges...> } ON <targets...> TO <grantees...>

Privileges:
  CREATE, DROP, GRANT, SELECT, INSERT, DELETE, UPDATE

Targets:
  DATABASE <databasename> [, ...]
  [TABLE] [<databasename> .] { <tablename> | * } [, ...]

`,
		//line sql.y: 1762
		SeeAlso: `REVOKE, WEBDOCS/grant.html
`,
	},
	//line sql.y: 1770
	`REVOKE`: {
		ShortDescription: `remove access privileges`,
		//line sql.y: 1771
		Category: hPriv,
		//line sql.y: 1772
		Text: `
REVOKE {ALL | <privileges...> } ON <targets...> FROM <grantees...>

Privileges:
  CREATE, DROP, GRANT, SELECT, INSERT, DELETE, UPDATE

Targets:
  DATABASE <databasename> [, <databasename>]...
  [TABLE] [<databasename> .] { <tablename> | * } [, ...]

`,
		//line sql.y: 1782
		SeeAlso: `GRANT, WEBDOCS/revoke.html
`,
	},
	//line sql.y: 1869
	`RESET`: {
		ShortDescription: `reset a session variable to its default value`,
		//line sql.y: 1870
		Category: hCfg,
		//line sql.y: 1871
		Text: `RESET [SESSION] <var>
`,
		//line sql.y: 1872
		SeeAlso: `RESET CLUSTER SETTING, WEBDOCS/set-vars.html
`,
	},
	//line sql.y: 1884
	`RESET CLUSTER SETTING`: {
		ShortDescription: `reset a cluster setting to its default value`,
		//line sql.y: 1885
		Category: hCfg,
		//line sql.y: 1886
		Text: `RESET CLUSTER SETTING <var>
`,
		//line sql.y: 1887
		SeeAlso: `SET CLUSTER SETTING, RESET
`,
	},
	//line sql.y: 1913
	`SET CLUSTER SETTING`: {
		ShortDescription: `change a cluster setting`,
		//line sql.y: 1914
		Category: hCfg,
		//line sql.y: 1915
		Text: `SET CLUSTER SETTING <var> { TO | = } <value>
`,
		//line sql.y: 1916
		SeeAlso: `SHOW CLUSTER SETTING, RESET CLUSTER SETTING, SET SESSION,
WEBDOCS/cluster-settings.html
`,
	},
	//line sql.y: 1937
	`SET SESSION`: {
		ShortDescription: `change a session variable`,
		//line sql.y: 1938
		Category: hCfg,
		//line sql.y: 1939
		Text: `
SET [SESSION] <var> { TO | = } <values...>
SET [SESSION] TIME ZONE <tz>
SET [SESSION] CHARACTERISTICS AS TRANSACTION ISOLATION LEVEL { SNAPSHOT | SERIALIZABLE }

`,
		//line sql.y: 1944
		SeeAlso: `SHOW SESSION, RESET, DISCARD, SHOW, SET CLUSTER SETTING, SET TRANSACTION,
WEBDOCS/set-vars.html
`,
	},
	//line sql.y: 1961
	`SET TRANSACTION`: {
		ShortDescription: `configure the transaction settings`,
		//line sql.y: 1962
		Category: hTxn,
		//line sql.y: 1963
		Text: `
SET [SESSION] TRANSACTION <txnparameters...>

Transaction parameters:
   ISOLATION LEVEL { SNAPSHOT | SERIALIZABLE }
   PRIORITY { LOW | NORMAL | HIGH }

`,
		//line sql.y: 1970
		SeeAlso: `SHOW TRANSACTION, SET SESSION,
WEBDOCS/set-transaction.html
`,
	},
	//line sql.y: 2109
	`SHOW`: {
		//line sql.y: 2110
		Category: hGroup,
		//line sql.y: 2111
		Text: `
SHOW SESSION, SHOW CLUSTER SETTING, SHOW DATABASES, SHOW TABLES, SHOW COLUMNS, SHOW INDEXES,
SHOW CONSTRAINTS, SHOW CREATE TABLE, SHOW CREATE VIEW, SHOW USERS, SHOW TRANSACTION, SHOW BACKUP,
SHOW JOBS, SHOW QUERIES, SHOW SESSIONS, SHOW TRACE
`,
	},
	//line sql.y: 2136
	`SHOW SESSION`: {
		ShortDescription: `display session variables`,
		//line sql.y: 2137
		Category: hCfg,
		//line sql.y: 2138
		Text: `SHOW [SESSION] { <var> | ALL }
`,
		//line sql.y: 2139
		SeeAlso: `WEBDOCS/show-vars.html
`,
	},
	//line sql.y: 2160
	`SHOW BACKUP`: {
		ShortDescription: `list backup contents`,
		//line sql.y: 2161
		Category: hCCL,
		//line sql.y: 2162
		Text: `SHOW BACKUP <location>
`,
		//line sql.y: 2163
		SeeAlso: `WEBDOCS/show-backup.html
`,
	},
	//line sql.y: 2171
	`SHOW CLUSTER SETTING`: {
		ShortDescription: `display cluster settings`,
		//line sql.y: 2172
		Category: hCfg,
		//line sql.y: 2173
		Text: `
SHOW CLUSTER SETTING <var>
SHOW ALL CLUSTER SETTINGS
`,
		//line sql.y: 2176
		SeeAlso: `WEBDOCS/cluster-settings.html
`,
	},
	//line sql.y: 2193
	`SHOW COLUMNS`: {
		ShortDescription: `list columns in relation`,
		//line sql.y: 2194
		Category: hDDL,
		//line sql.y: 2195
		Text: `SHOW COLUMNS FROM <tablename>
`,
		//line sql.y: 2196
		SeeAlso: `WEBDOCS/show-columns.html
`,
	},
	//line sql.y: 2204
	`SHOW DATABASES`: {
		ShortDescription: `list databases`,
		//line sql.y: 2205
		Category: hDDL,
		//line sql.y: 2206
		Text: `SHOW DATABASES
`,
		//line sql.y: 2207
		SeeAlso: `WEBDOCS/show-databases.html
`,
	},
	//line sql.y: 2215
	`SHOW GRANTS`: {
		ShortDescription: `list grants`,
		//line sql.y: 2216
		Category: hPriv,
		//line sql.y: 2217
		Text: `SHOW GRANTS [ON <targets...>] [FOR <users...>]
`,
		//line sql.y: 2218
		SeeAlso: `WEBDOCS/show-grants.html
`,
	},
	//line sql.y: 2226
	`SHOW INDEXES`: {
		ShortDescription: `list indexes`,
		//line sql.y: 2227
		Category: hDDL,
		//line sql.y: 2228
		Text: `SHOW INDEXES FROM <tablename>
`,
		//line sql.y: 2229
		SeeAlso: `WEBDOCS/show-index.html
`,
	},
	//line sql.y: 2247
	`SHOW CONSTRAINTS`: {
		ShortDescription: `list constraints`,
		//line sql.y: 2248
		Category: hDDL,
		//line sql.y: 2249
		Text: `SHOW CONSTRAINTS FROM <tablename>
`,
		//line sql.y: 2250
		SeeAlso: `WEBDOCS/show-constraints.html
`,
	},
	//line sql.y: 2263
	`SHOW QUERIES`: {
		ShortDescription: `list running queries`,
		//line sql.y: 2264
		Category: hMisc,
		//line sql.y: 2265
		Text: `SHOW [CLUSTER | LOCAL] QUERIES
`,
		//line sql.y: 2266
		SeeAlso: `CANCEL QUERY
`,
	},
	//line sql.y: 2282
	`SHOW JOBS`: {
		ShortDescription: `list background jobs`,
		//line sql.y: 2283
		Category: hMisc,
		//line sql.y: 2284
		Text: `SHOW JOBS
`,
		//line sql.y: 2285
		SeeAlso: `CANCEL JOB, PAUSE JOB, RESUME JOB
`,
	},
	//line sql.y: 2293
	`SHOW TRACE`: {
		ShortDescription: `display an execution trace`,
		//line sql.y: 2294
		Category: hMisc,
		//line sql.y: 2295
		Text: `
SHOW [KV] TRACE FOR SESSION
SHOW [KV] TRACE FOR <statement>
`,
		//line sql.y: 2298
		SeeAlso: `EXPLAIN
`,
	},
	//line sql.y: 2319
	`SHOW SESSIONS`: {
		ShortDescription: `list open client sessions`,
		//line sql.y: 2320
		Category: hMisc,
		//line sql.y: 2321
		Text: `SHOW [CLUSTER | LOCAL] SESSIONS
`,
	},
	//line sql.y: 2337
	`SHOW TABLES`: {
		ShortDescription: `list tables`,
		//line sql.y: 2338
		Category: hDDL,
		//line sql.y: 2339
		Text: `SHOW TABLES [FROM <databasename>]
`,
		//line sql.y: 2340
		SeeAlso: `WEBDOCS/show-tables.html
`,
	},
	//line sql.y: 2352
	`SHOW TRANSACTION`: {
		ShortDescription: `display current transaction properties`,
		//line sql.y: 2353
		Category: hCfg,
		//line sql.y: 2354
		Text: `SHOW TRANSACTION {ISOLATION LEVEL | PRIORITY | STATUS}
`,
		//line sql.y: 2355
		SeeAlso: `WEBDOCS/show-transaction.html
`,
	},
	//line sql.y: 2374
	`SHOW CREATE TABLE`: {
		ShortDescription: `display the CREATE TABLE statement for a table`,
		//line sql.y: 2375
		Category: hDDL,
		//line sql.y: 2376
		Text: `SHOW CREATE TABLE <tablename>
`,
		//line sql.y: 2377
		SeeAlso: `WEBDOCS/show-create-table.html
`,
	},
	//line sql.y: 2385
	`SHOW CREATE VIEW`: {
		ShortDescription: `display the CREATE VIEW statement for a view`,
		//line sql.y: 2386
		Category: hDDL,
		//line sql.y: 2387
		Text: `SHOW CREATE VIEW <viewname>
`,
		//line sql.y: 2388
		SeeAlso: `WEBDOCS/show-create-view.html
`,
	},
	//line sql.y: 2396
	`SHOW USERS`: {
		ShortDescription: `list defined users`,
		//line sql.y: 2397
		Category: hPriv,
		//line sql.y: 2398
		Text: `SHOW USERS
`,
		//line sql.y: 2399
		SeeAlso: `CREATE USER, DROP USER, WEBDOCS/show-users.html
`,
	},
	//line sql.y: 2445
	`PAUSE JOB`: {
		ShortDescription: `pause a background job`,
		//line sql.y: 2446
		Category: hMisc,
		//line sql.y: 2447
		Text: `PAUSE JOB <jobid>
`,
		//line sql.y: 2448
		SeeAlso: `SHOW JOBS, CANCEL JOB, RESUME JOB
`,
	},
	//line sql.y: 2456
	`CREATE TABLE`: {
		ShortDescription: `create a new table`,
		//line sql.y: 2457
		Category: hDDL,
		//line sql.y: 2458
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
		//line sql.y: 2484
		SeeAlso: `SHOW TABLES, CREATE VIEW, SHOW CREATE TABLE,
WEBDOCS/create-table.html
WEBDOCS/create-table-as.html
`,
	},
	//line sql.y: 2909
	`TRUNCATE`: {
		ShortDescription: `empty one or more tables`,
		//line sql.y: 2910
		Category: hDML,
		//line sql.y: 2911
		Text: `TRUNCATE [TABLE] <tablename> [, ...] [CASCADE | RESTRICT]
`,
		//line sql.y: 2912
		SeeAlso: `WEBDOCS/truncate.html
`,
	},
	//line sql.y: 2920
	`CREATE USER`: {
		ShortDescription: `define a new user`,
		//line sql.y: 2921
		Category: hPriv,
		//line sql.y: 2922
		Text: `CREATE USER <name> [ [WITH] PASSWORD <passwd> ]
`,
		//line sql.y: 2923
		SeeAlso: `DROP USER, SHOW USERS, WEBDOCS/create-user.html
`,
	},
	//line sql.y: 2941
	`CREATE VIEW`: {
		ShortDescription: `create a new view`,
		//line sql.y: 2942
		Category: hDDL,
		//line sql.y: 2943
		Text: `CREATE VIEW <viewname> [( <colnames...> )] AS <source>
`,
		//line sql.y: 2944
		SeeAlso: `CREATE TABLE, SHOW CREATE VIEW, WEBDOCS/create-view.html
`,
	},
	//line sql.y: 2958
	`CREATE INDEX`: {
		ShortDescription: `create a new index`,
		//line sql.y: 2959
		Category: hDDL,
		//line sql.y: 2960
		Text: `
CREATE [UNIQUE] INDEX [IF NOT EXISTS] [<idxname>]
       ON <tablename> ( <colname> [ASC | DESC] [, ...] )
       [STORING ( <colnames...> )] [<interleave>]

Interleave clause:
   INTERLEAVE IN PARENT <tablename> ( <colnames...> ) [CASCADE | RESTRICT]

`,
		//line sql.y: 2968
		SeeAlso: `CREATE TABLE, SHOW INDEXES, SHOW CREATE INDEX,
WEBDOCS/create-index.html
`,
	},
	//line sql.y: 3107
	`RELEASE`: {
		ShortDescription: `complete a retryable block`,
		//line sql.y: 3108
		Category: hTxn,
		//line sql.y: 3109
		Text: `RELEASE [SAVEPOINT] cockroach_restart
`,
		//line sql.y: 3110
		SeeAlso: `SAVEPOINT, WEBDOCS/savepoint.html
`,
	},
	//line sql.y: 3118
	`RESUME JOB`: {
		ShortDescription: `resume a background job`,
		//line sql.y: 3119
		Category: hMisc,
		//line sql.y: 3120
		Text: `RESUME JOB <jobid>
`,
		//line sql.y: 3121
		SeeAlso: `SHOW JOBS, CANCEL JOB, PAUSE JOB
`,
	},
	//line sql.y: 3129
	`SAVEPOINT`: {
		ShortDescription: `start a retryable block`,
		//line sql.y: 3130
		Category: hTxn,
		//line sql.y: 3131
		Text: `SAVEPOINT cockroach_restart
`,
		//line sql.y: 3132
		SeeAlso: `RELEASE, WEBDOCS/savepoint.html
`,
	},
	//line sql.y: 3146
	`BEGIN`: {
		ShortDescription: `start a transaction`,
		//line sql.y: 3147
		Category: hTxn,
		//line sql.y: 3148
		Text: `
BEGIN [TRANSACTION] [ <txnparameter> [[,] ...] ]
START TRANSACTION [ <txnparameter> [[,] ...] ]

Transaction parameters:
   ISOLATION LEVEL { SNAPSHOT | SERIALIZABLE }
   PRIORITY { LOW | NORMAL | HIGH }

`,
		//line sql.y: 3156
		SeeAlso: `COMMIT, ROLLBACK, WEBDOCS/begin-transaction.html
`,
	},
	//line sql.y: 3169
	`COMMIT`: {
		ShortDescription: `commit the current transaction`,
		//line sql.y: 3170
		Category: hTxn,
		//line sql.y: 3171
		Text: `
COMMIT [TRANSACTION]
END [TRANSACTION]
`,
		//line sql.y: 3174
		SeeAlso: `BEGIN, ROLLBACK, WEBDOCS/commit-transaction.html
`,
	},
	//line sql.y: 3187
	`ROLLBACK`: {
		ShortDescription: `abort the current transaction`,
		//line sql.y: 3188
		Category: hTxn,
		//line sql.y: 3189
		Text: `ROLLBACK [TRANSACTION] [TO [SAVEPOINT] cockroach_restart]
`,
		//line sql.y: 3190
		SeeAlso: `BEGIN, COMMIT, SAVEPOINT, WEBDOCS/rollback-transaction.html
`,
	},
	//line sql.y: 3304
	`CREATE DATABASE`: {
		ShortDescription: `create a new database`,
		//line sql.y: 3305
		Category: hDDL,
		//line sql.y: 3306
		Text: `CREATE DATABASE [IF NOT EXISTS] <name>
`,
		//line sql.y: 3307
		SeeAlso: `WEBDOCS/create-database.html
`,
	},
	//line sql.y: 3376
	`INSERT`: {
		ShortDescription: `create new rows in a table`,
		//line sql.y: 3377
		Category: hDML,
		//line sql.y: 3378
		Text: `
INSERT INTO <tablename> [[AS] <name>] [( <colnames...> )]
       <selectclause>
       [ON CONFLICT [( <colnames...> )] {DO UPDATE SET ... [WHERE <expr>] | DO NOTHING}]
       [RETURNING <exprs...>]
`,
		//line sql.y: 3383
		SeeAlso: `UPSERT, UPDATE, DELETE, WEBDOCS/insert.html
`,
	},
	//line sql.y: 3400
	`UPSERT`: {
		ShortDescription: `create or replace rows in a table`,
		//line sql.y: 3401
		Category: hDML,
		//line sql.y: 3402
		Text: `
UPSERT INTO <tablename> [AS <name>] [( <colnames...> )]
       <selectclause>
       [RETURNING <exprs...>]
`,
		//line sql.y: 3406
		SeeAlso: `INSERT, UPDATE, DELETE, WEBDOCS/upsert.html
`,
	},
	//line sql.y: 3482
	`UPDATE`: {
		ShortDescription: `update rows of a table`,
		//line sql.y: 3483
		Category: hDML,
		//line sql.y: 3484
		Text: `UPDATE <tablename> [[AS] <name>] SET ... [WHERE <expr>] [RETURNING <exprs...>]
`,
		//line sql.y: 3485
		SeeAlso: `INSERT, UPSERT, DELETE, WEBDOCS/update.html
`,
	},
	//line sql.y: 3653
	`<SELECTCLAUSE>`: {
		ShortDescription: `access tabular data`,
		//line sql.y: 3654
		Category: hDML,
		//line sql.y: 3655
		Text: `
Select clause:
  TABLE <tablename>
  VALUES ( <exprs...> ) [ , ... ]
  SELECT ... [ { INTERSECT | UNION | EXCEPT } [ ALL | DISTINCT ] <selectclause> ]
`,
	},
	//line sql.y: 3666
	`SELECT`: {
		ShortDescription: `retrieve rows from a data source and compute a result`,
		//line sql.y: 3667
		Category: hDML,
		//line sql.y: 3668
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
		//line sql.y: 3680
		SeeAlso: `WEBDOCS/select.html
`,
	},
	//line sql.y: 3740
	`TABLE`: {
		ShortDescription: `select an entire table`,
		//line sql.y: 3741
		Category: hDML,
		//line sql.y: 3742
		Text: `TABLE <tablename>
`,
		//line sql.y: 3743
		SeeAlso: `SELECT, VALUES, WEBDOCS/table-expressions.html
`,
	},
	//line sql.y: 3986
	`VALUES`: {
		ShortDescription: `select a given set of values`,
		//line sql.y: 3987
		Category: hDML,
		//line sql.y: 3988
		Text: `VALUES ( <exprs...> ) [, ...]
`,
		//line sql.y: 3989
		SeeAlso: `SELECT, TABLE, WEBDOCS/table-expressions.html
`,
	},
	//line sql.y: 4094
	`<SOURCE>`: {
		ShortDescription: `define a data source for SELECT`,
		//line sql.y: 4095
		Category: hDML,
		//line sql.y: 4096
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
		//line sql.y: 4114
		SeeAlso: `WEBDOCS/table-expressions.html
`,
	},
}
