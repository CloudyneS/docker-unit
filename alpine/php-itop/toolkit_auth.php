<?php

if (!isset($_SERVER['PHP_AUTH_USER']) || !isset($_SERVER['PHP_AUTH_PW']) || $_SERVER['PHP_AUTH_USER'] !== 'toolkit' || $_SERVER['PHP_AUTH_PW'] !== 'toolkit') {
    header('WWW-Authenticate: Basic realm="ITOP Toolkit"');
    header('HTTP/1.0 401 Unauthorized');
    echo 'Unauthorized';
    exit;
}
