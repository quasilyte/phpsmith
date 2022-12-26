<?php

function make_positive_inf(): float {
    $x = 1.0;
    $y = 0.0;
    return $x / $y;
}

function make_negative_inf(): float {
    $x = -1.0;
    $y = 0.0;
    return $x / $y;
}

function make_nan(): float {
    $x = 0.0;
    $y = 0.0;
    return $x / $y;
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
 * @kphp-template $x
 * @kphp-template $y
 */
function _safe_div($x, $y) {
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
 * @kphp-template $x
 * @kphp-template $y
 */
function _safe_mod($x, $y) {
    if ($y > 0 || $y < 0) {
        try {
            return $x % $y;
        } catch (\Throwable $e) {
        }
    }
    echo "invalid argument in %\n";
    return 0;
}