```text au:fixture=first
FIXTURE 1
```

```text au:fixture=third au:group=grp au:requires=first
FIXTURE 3
```

```text au:fixture=second au:group=grp au:requires=first
FIXTURE 2
```

```text au:input au:group=grp
INPUT
```

```text au:output au:group=grp
OUTPUT
```
