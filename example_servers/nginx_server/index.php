<?php print("GX Port: " . $_SERVER['SERVER_PORT'] . "\n"); ?>
<?php usleep(50000); ?>
<?php print("50 ms, "); ob_flush(); flush(); ?>
<?php usleep(50000); ?>
<?php print("100 ms, "); ob_flush(); flush(); ?>
<?php usleep(50000); ?>
<?php print("150 ms, "); ob_flush(); flush(); ?>
<?php usleep(50000); ?>
<?php print("200 ms, "); ob_flush(); flush(); ?>
<?php usleep(50000); ?>
<?php print("250 ms\n"); ob_flush(); flush(); ?>
