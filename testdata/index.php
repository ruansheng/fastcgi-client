<?php
$data = file_get_contents("php://input");
echo $data;
/*
if(isset($_POST['name'])) {
	$name = $_POST['name'];
	echo "post:".$name;
}else if(isset($_GET['name'])) {
	$name = $_GET['name'];
	echo "get:".$name;
}
*/