```text au:fixture=base
FIXTURE 1
```

```text au:fixture=first au:requires=base
FIXTURE 2
```

```text au:fixture=second au:requires=base
FIXTURE 3
```

```text au:input au:group=grp au:requires="first, second, first"
INPUT
```

```text au:output au:group=grp au:requires=base,second
OUTPUT
```
