```sql au:fixture au:group=grp role=schema
CREATE TABLE example (value INTEGER);
```

```json au:fixture au:group=grp role=events
[{"value":1}]
```

```sql au:input au:group=grp
SELECT value FROM example;
```

```json au:output au:group=grp
[{"value":1}]
```
