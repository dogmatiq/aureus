# JSON Pretty Printer (diff)

This document shows how to configure an assertion to expect a difference between
an input and output.

```json
{
  "value": 1
}
```

```json
{
  "value": 2
}
```

```diff au:assertion="diff"
 {
-  "value": 1
+  "value": 2
 }
```
