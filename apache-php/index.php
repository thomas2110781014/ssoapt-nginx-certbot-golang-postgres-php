<?php 

$url_backend = "http://backend:8000";


function getTasks() {
    // create curl resource
    $ch = curl_init();
    // set url
    curl_setopt($ch, CURLOPT_URL, $GLOBALS['url_backend']);
    // return the transfer as a string
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
    // $response contains the response string
    $response = curl_exec($ch);
    // close curl resource to free up system resources
    curl_close($ch);
    
    // return response
    echo $response;
}

echo '<h1>Hello World!</h1>'; 

getTasks();

?>


