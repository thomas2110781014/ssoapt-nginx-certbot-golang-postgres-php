<html>
<body>

<p>This is a showcase of our app for the SSOA-PT exercise for FH Burgenland!</p>

  <form action="index.php" method="post">
    Title: <input type="text" name="title"><br>
    Description: <input type="text" name="description"><br>
    Due Date: <input type="text" name="duedate"><br>
    <input type="submit" name="buttonAdd" value="Add Task!">
    <input type="submit" name="buttonDelete" value="Delete Task!">
  </form>

<?php

$url_backend = "http://backend:8000";

function r($var){
    echo '<pre>';
    print_r($var);
    echo '</pre>';
}

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
    return $response;
}

function addTasks() {
    $ch = curl_init( $GLOBALS['url_backend'] );
    // r($GLOBALS['data']);
    # Setup request to send json via POST.
    $payload = json_encode($GLOBALS['data']);
    curl_setopt( $ch, CURLOPT_POSTFIELDS, $payload );
    curl_setopt( $ch, CURLOPT_HTTPHEADER, array('Content-Type:application/json'));
    # Return response instead of printing.
    curl_setopt( $ch, CURLOPT_RETURNTRANSFER, true );
    # Send request.
    $result = curl_exec($ch);
    curl_close($ch);
    # Print response.
    // echo "<p>Adding the following tasks: <p>";
    // echo "<pre>$result</pre>";
}

$tasks_array = getTasks();
$tasks_array_decoded = json_decode($tasks_array, true);

foreach ($tasks_array_decoded as $value) {
    echo $value['id'] . ": " . $value['title'] . " " . $value['description'] . " " . $value['duedate'] . "<br />";
}

if(isset($_POST['title']) and isset($_POST['buttonAdd'])) {
  $data[] = '';
  $data['title'] = $_POST['title'];
  $data['description'] = $_POST['description'];
  $data['duedate'] = $_POST['duedate'];

  addTasks($data);
  echo "<meta http-equiv='refresh' content='0'>";
}

?>

</body>
</html>
