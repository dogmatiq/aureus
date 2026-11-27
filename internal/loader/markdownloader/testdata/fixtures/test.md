```text au:fixture=first foo=bar
FIXTURE 1
```

```text au:fixture=second foo=baz
FIXTURE 2
```

```text au:input au:group=grp au:requires=first,second
INPUT
```

```text au:output au:group=grp
OUTPUT
```

```text au:fixture=third
FIXTURE 3
```

```text au:input au:group=other au:requires=third
INPUT
```

```text au:output au:group=other
OUTPUT
```
