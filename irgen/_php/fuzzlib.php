<?php

/**
 * @param string $name
 * @return bool
 */
function _visit_function($name) {
    static $counters = [];
    if ($counters[$name] > 10) {
        echo "$name reached call limit\n";
        return false;
    }
    $counters[$name]++;
    return true;
}

function make_positive_inf(): float {
#ifndef KPHP
    return INF;
#endif
    return 1.7976931348623157E+308 * 10;
}

function make_negative_inf(): float {
#ifndef KPHP
    return -INF;
#endif
    return -(1.7976931348623157E+308 * 10);
}

function make_nan(): float {
#ifndef KPHP
    return NAN;
#endif
    return acos(-8);
}

function dump_with_pos($file, $line, $v) {
    var_dump(["$file:$line" => $v]);
}

/**
 * @param mixed $x
 * @param mixed $y
 */
function float_eq2($x, $y) {
    return $x == $y;
}

/**
 * @param mixed $x
 * @param mixed $y
 */
function float_eq3($x, $y) {
    return $x === $y;
}

/**
 * @param mixed $x
 * @param mixed $y
 */
function float_neq2($x, $y) {
    return $x != $y;
}

/**
 * @param mixed $x
 * @param mixed $y
 */
function float_neq3($x, $y) {
    return $x !== $y;
}

/**
 * @param float $x
 * @param float $y
 * @return float
 */
function _safe_float_div($x, $y) {
    if ($y > 0.0 || $y < 0.0) {
        try {
            return $x / $y;
        } catch (\Throwable $e) {
        }
    }
    echo "invalid argument in /\n";
    return 0.0;
}

/**
 * @param int $x
 * @param int $y
 */
function _safe_int_div($x, $y) {
    if ($y > 0 || $y < 0) {
        try {
            return $x / $y;
        } catch (\Throwable $e) {
        }
    }
    echo "invalid argument in /\n";
    return 0;
}

/**
 * @param float $x
 * @param float $y
 * @return float
 */
function _safe_float_mod($x, $y) {
    if ($y > 0.0 || $y < 0.0) {
        try {
            return $x % $y;
        } catch (\Throwable $e) {
        }
    }
    echo "invalid argument in %\n";
    return 0.0;
}

/**
 * @param int $x
 * @param int $y
 */
function _safe_int_mod($x, $y) {
    if ($y > 0 || $y < 0) {
        try {
            return $x % $y;
        } catch (\Throwable $e) {
        }
    }
    echo "invalid argument in %\n";
    return 0;
}

#ifndef KPHP
function tuple(...$args) {
    // turn off PhpStorm native inferring
    return ${'args'};
}
#endif
