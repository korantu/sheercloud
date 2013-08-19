function Controller($scope, $http) {
    $scope.checked = true
    $scope.result = "None"

    $scope.loggedIn = function() {
	return true;
    }

    $scope.login = function() {
	console.log("Attempting to login...")
	
	$http.get("/api/login").success(function(data) {
	    $scope.user_name = data + ":" + $scope.login_name + " " + $scope.login_password
	} )

    }
}
