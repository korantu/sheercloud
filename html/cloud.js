function Controller($scope, $http) {
    $scope.checked = true
    $scope.result = "None"

    $scope.loggedIn = function() {
	return true;
    }

    $scope.login = function() {
	console.log("Attempting to login...")
	$http.post("/api/login", '{ "Username": "login_name", "Password": "login_password" }' ).success(function(data) {
	    if ( data.Success ) {
		$scope.user_name = "Code:" + data.Session;
	    } else {
		$scope.user_name = "Login failed.";
	    }
	} ).error( function(data) {
	    $scope.user_name = "Oops: " + data
	})

    }
}
