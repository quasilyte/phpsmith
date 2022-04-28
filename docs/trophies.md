## KPHP issues

### Can't compile valid programs

1. False-typed value `$x` can't be used inside `if ($x)`.

```
<?php

function func3() {
  $v0 = false;
  if ($v0) {
    var_dump($v0);
  }
}

func3();
```
