<?php

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
        return $x / $y;
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
        return $x % $y;
    }
    echo "invalid argument in %\n";
    return 0;
}