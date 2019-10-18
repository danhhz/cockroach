colscan
=======

This prototype explores a what a columnar disk format could look like in
CockroachDB and its effect on end-to-end vectorized execution performance.

It's inspired by TiDB's [TiFlash]. Columnar disk formats are great for OLAP and
OLAP-light read queries, but row disk formats are better for pure OLTP style
point lookups and low-latency mutations. At a high level, a special "columnar
replica" which applies raft log entries to a [columnar data engine] is added as
a raft learner replica.

Being a learner replica, the raft leader sends it raft traffic in a timely
manner, but it doesn't vote or otherwise affect the OLTP part of the system.
These columnar replicas would likely be pinned to their own machines to further
isolate them (and because they want different machine specs).

The optimizer would be in charge of selecting which reads are appropriate for
columnar replicas. It already needs row count estimates to plan joins and other
things, so hopefully the heuristics are straightforward. Potentially this is as
easy and sending anything with a row count estimate greater than some threshold.

Our ScanRequest already supports two `ScanFormat`s for how the requested data
should be returned, the original protobuf based `KEY_VALUES` and the rocksdb
batch based `BATCH_RESPONSE`. A `COLUMNAR` option is added which directly
returns the arrow-based format used by our vectorized execution engine. The only
KV request supported by columnar replicas is `ScanRequest` with the `COLUMNAR`
format selected, all others return an error pointing the client to another
replica.

**Note** The above is not exactly what is implemented here. Right now, all
replicas have both a row and columnar engine, any AddSSTable requests affecting
user-keyspace ingest data into the columnar engine, and any `ScanRequest` with
the `COLUMNAR` format selected is served using the columnar engine. This will
serve inconsistent data when there is transactional traffic, but it works for
IMPORT followed by read-only traffic, which is enough for a prototype.

[TiFlash]: https://www.percona.com/live/19/sessions/making-htap-real-with-tiflash-a-tidb-native-columnar-extension
[columnar data engine]: ../colengine

# Challenges

- **Schema plumbing**: Both our replication and our engine abstractions are
  built around KVs, but for columnar data to make any sense, the KVs need to be
  decoded into the actual table columns, requiring the table schemas.

  This prototype adds a hook into SQL table leasing to request the columnar
  schema it should use. This decoding currently happens at the
  proposal-evaluation phase, which essentially doubles proposal sizes, doubling
  raft traffic and raft log sizes. An alternative is to do it below raft at
  application time, but this is terrifying.

- **Schema changes**: Our SQL layer knows how to run no-downtime online table
  schema migrations by carefully interpreting previously written kvs and
  backfilling when necessary. A columnar engine would have to be able to write
  data for a given table in multiple schemas and know how to merge them together
  at query time.

- **Streaming mutations**: Most columnar engines expect that mutations are done
  in batches to minimize write amplification. However, while we don't expect the
  low-latency mutation of an LSM, we'll want to apply raft log entries as they
  come in, so analytics queries can be served promptly. See [Mutations] below.

- **MVCC**: In our existing replicas, MVCC exists entirely on top of an LSM that
  presents a KV abstraction. For columnar scans to be performant, MVCC will need
  to be more tightly coupled into the columnar storage.

- **Abstractions**: Logic for encoding/decoding KVs, MVCC, and schema changes
  will all be somewhat necessary down in the columnar engine, but our current
  abstractions for these are not amenable for use there.

  Perhaps this will be a forcing function for developing first class
  abstractions for these. It would be interesting to remove them enough from
  CockroachDB SQL's use they they're appropriate dependencies for a columnar
  engine.

[mutations]: #mutations

# Columnar Engine Design

TODO

- Optimize for cold data
- Minimize data copies

## Mutations

TODO

# Future Work
- [ ] Rename pkg/col/colstore to pkg/col/coldb
- [ ] Make namespaces a []byte instead of int64
- [ ] Get secondary indexes working
- [ ] Get multiple namespaces plumbed
- [ ] Get all TPCH queries that are supported by the vectorized engine working
- [ ] Make coltypes.T a protobuf
- [ ] Allow leasing an offline table
- [ ] Better name for ColumnarSchemaer
- [ ] NULLs
- [ ] Speed up row comparisons. Probably something like [abbreviated keys]
- [ ] Compression

[abbreviated keys]: https://brandur.org/sortsupport#abbreviated-keys
