This test demonstrates that the Markdown document can have additional content
surrounding the code blocks without affecting the structure of the tests.

In this test, we take some **unformatted** JSON as input ...

```
{ "foo": 1, "bar": 2, "baz": 3 }
```

And then we **pretty-print** it to make it more readable:

```au:assert
{
  "foo": 1,
  "bar": 2,
  "baz": 3
}
```
