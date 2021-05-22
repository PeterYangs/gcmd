<?php

$index=0;

while (true){

    if($index>=10){

        throw new Exception("error here!");

    }

    echo 'echo success'.PHP_EOL;

    $index++;

    sleep(1);

}
