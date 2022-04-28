## KPHP issues

### Can't compile valid programs

1. False-typed value `$x` can't be used inside `if ($x)`.

```php
<?php

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
<?php

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
