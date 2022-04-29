## KPHP issues

### Critical error

1. Segfault in `modulo` with invalid args at run time

```php
function modxy($x, $y) {
  do {} while (false); // prevent inlining and const folding
  return $x % $y;
}

$v8 = modxy(-9223372036854775808, -1);
var_dump($v8);
```

### Different results from PHP

1. Overflowing numeric string conversions generate different values

```php
$result = (int)decbin(-1);
var_dump($result);
```

2. Mismatching results for `ucwords`

```php
$v3 = ucwords('204c');
var_dump($v3);
```

3. `~` is escaped in rawurlencode

```php
$v1 = "~";
var_dump(rawurlencode($v1));
```

### Can't compile valid programs

1. False-typed value `$x` can't be used inside `if ($x)`.

```php
function func3() {
  $v0 = false;
  if ($v0) {
    var_dump($v0);
  }
}

func3();
```

2. Concatenation breaks in `case` if it's after `continue` in a loop

```php
function func0() {
  $v2 = 'abc';
  $v4 = 0;
  while ($v4++ < 6) {
    continue;
    switch (($v2)) {
      case 'ab' . 'c':
        break;
    }
  }
}

func0();
```
