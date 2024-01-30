# JSON Pretty Printer

This document describes a test of a JSON pretty-printer, and proves that Aureus
can make assertions that are interleave with regular Markdown content.

The input to the pretty printer is shown below:

```json
{ "one": 1, "two": 2 }
```

If the pretty-printer is functioning correctly, it should produce JSON that is
indented and and separated onto new lines, as follows:

```json au:assertion
{
  "one": 1,
  "two": 2
}
```
