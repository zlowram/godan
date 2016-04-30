var godan_api = "http://localhost:8000/";

angular.module('Godan').controller('NavCtrl', ['$scope', '$location', '$resource', '$window', '$uibModal', function NavCtrl($scope, $location, $resource, $window, $uibModal) {
	var url = godan_api + "users/";

	$scope.isAdmin = function() {
		if ($window.sessionStorage.role == "admin") {
			return true;
		} else {
			return false;
		}
	}

	$scope.getUsername = function() { return $window.sessionStorage.username; }

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

	$scope.profile = function() {
		var user = $resource(url + $scope.getUsername(), {}, {
			query: {
				method: 'GET',
				headers: {'Authorization': 'Bearer ' + $window.sessionStorage.token}
			}
		}).query();

		$uibModal.open({
			animation: $scope.animationsEnabled,
			templateUrl: 'profileModalTemplate',
			controller: 'ProfileModalInstanceCtrl',
			size: 'lg',
			resolve: {
				element: function() {
					return user;
				}
			}
		});
	};

	$scope.logout = function() {
		delete $window.sessionStorage.token;
		delete $window.sessionStorage.role;
		delete $window.sessionStorage.username;
		$window.location.href = '#/login';
	}
}]);

angular.module('Godan').controller('ProfileModalInstanceCtrl', function ($scope, $resource, $uibModalInstance, $window, element) {
  var url = godan_api + "users/";

  $scope.user = element;

  $scope.submitChanges = function () {
	var user = {}
	if ($scope.user.Email) user["Email"] = $scope.user.Email;
	if ($scope.user.NewPassword != $scope.user.NewPasswordConf) {
		$scope.passwordMismatch = true;
	} else if ($scope.user.NewPassword) user["Password"] = $scope.user.NewPassword;
	
	$resource(url + element.Username, {}, {
		update: {
			method: 'PUT',
			headers: {'Authorization': 'Bearer ' + $window.sessionStorage.token}
		}
	}).update(user);

	delete $scope.user.NewPassword;
	delete $scope.user.NewPasswordConf;
	$uibModalInstance.close();
  };

  $scope.closeit = function () {
    $uibModalInstance.close();
  };

});
