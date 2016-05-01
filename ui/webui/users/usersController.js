var godan_api = "http://localhost:8000/";
var url = godan_api + "users";

angular.module('Godan').controller('UsersCtrl', ["$scope", "$resource", "$uibModal", "$window", function UsersCtrl($scope, $resource, $uibModal, $window) {
	if (!$window.sessionStorage.token) $window.location.href = '#/login';

	$scope.animationsEnabled = true;

	getUsers = function() {
		return $resource(url, {}, {
			query: {
				method: 'GET',
				isArray: true,
				headers: {'Authorization': 'Bearer ' + $window.sessionStorage.token}
			}
		}).query();
	}

	$scope.usersTable = getUsers();

	$scope.editUser= function(user) {
		var modalInstance = $uibModal.open({
			animation: $scope.animationsEnabled,
			templateUrl: 'usersModalTemplate',
			controller: 'UsersModalInstanceCtrl',
			size: 'lg',
			resolve: {
				element: function() {
					return user;
				},
				update: function() {
					return true;
				}
			}
		});
	};

	$scope.addUser = function() {
		var modalInstance = $uibModal.open({
			animation: $scope.animationsEnabled,
			templateUrl: 'usersModalTemplate',
			controller: 'UsersModalInstanceCtrl',
			size: 'lg',
			resolve: {
				element: function() {
					return {};
				},
				update: function() {
					return false;
				}
			}
		}).result.then(
			function(result) {
				$scope.usersTable = getUsers();
				$scope.$applyAsync();
			}
		);
	}

	$scope.removeUser = function(user) {
		$uibModal.open({
			animation: $scope.animationsEnabled,
			templateUrl: 'confirmationModalTemplate',
			controller: 'ConfirmationModalCtrl',
			size: 'lg',
			resolve: {
				element: function() {
					return user;
				},
				title: function() { return "Confirmation"; },
				message: function() { return "Are you sure you want to delete the user " + user.Username + "?"; }
			}
		}).result.then(
			function(result) {
				$scope.usersTable = getUsers();
				$scope.$applyAsync();
			}
		);
	}
}]);

angular.module('Godan').controller('UsersModalInstanceCtrl', function ($scope, $resource, $uibModalInstance, $window, element, update) {
  $scope.user = element;

  $scope.roles = [
		{"name": "Admin", "value": "admin"},
		{"name": "User", "value": "user"},
  ];

  $scope.submitChanges = function () {
	if (!$scope.user.Username) {
		$scope.usernameEmpty = true;
	} else if (!$scope.user.Role) {
		$scope.roleEmpty = true;
	} else if (!$scope.user.NewPassword && !update) {
		$scope.passwordEmpty = true;
	} else if ($scope.user.NewPassword != $scope.user.NewPasswordConf) {
		$scope.passwordMismatch = true;
	} else {
		var user = {
			"Username": $scope.user.Username,
			"Email": $scope.user.Email,
			"Role": $scope.user.Role
		}
		if ($scope.user.NewPassword) user["Password"] = $scope.user.NewPassword;
		
		if (update) {
			$resource(url + "/" + element.Username, {}, {
				update: {
					method: 'PUT',
					headers: {'Authorization': 'Bearer ' + $window.sessionStorage.token}
				}
			}).update(user);
		} else {
			$resource(url, {}, {
				save: {
					method: 'POST',
					headers: {'Authorization': 'Bearer ' + $window.sessionStorage.token}
				}
			}).save(user);
		}

		delete $scope.user.NewPassword;
		delete $scope.user.NewPasswordConf;
		$uibModalInstance.close();
	}
  };

  $scope.closeit = function () {
    $uibModalInstance.close();
  };

});

angular.module('Godan').controller('ConfirmationModalCtrl', function ($scope, $resource, $uibModalInstance, $window, element, title, message) {
  $scope.title = title;
  $scope.message = message;

  $scope.ok = function () {
	$resource(url + "/" + element.Username, {}, {
		del: {
			method: 'DELETE',
			headers: {'Authorization': 'Bearer ' + $window.sessionStorage.token}
		}
	}).del();

    $uibModalInstance.close();
  };

  $scope.cancel = function () {
    $uibModalInstance.close();
  };
});
