<?php

$index=0;

while (true){

    if($index>=10){

        throw new Exception("error here!");
//        break;

    }

    echo 'echo success'.PHP_EOL;

    $index++;

    sleep(1);

}
