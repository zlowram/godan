angular.module('Godan').controller('NavCtrl', ['$scope', '$location', '$window', function NavCtrl($scope, $location, $window) {
	$scope.isActive = function (viewLocation) { 
        return viewLocation === $location.path();
    };

	$scope.isLogged = function () {
		if ($window.sessionStorage.token) {
			return true;
		} else {
			return false;
		}
	}

	$scope.logout = function() {
		delete $window.sessionStorage.token
		$window.location.href = '#/login';
	}
}]);
