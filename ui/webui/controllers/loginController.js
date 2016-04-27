
var godan_api = "http://localhost:8000/";

angular.module('Godan').controller('LoginCtrl', ["$scope", "$resource", "$uibModal", "$window", function QueriesCtrl($scope, $resource, $uibModal, $window) {
	$scope.login = function() {

		var url = godan_api + "login";

		var login_data = {
			"username": $scope.username,
			"password": $scope.password,
		}

		var taskResult = $resource(url, {}, {}).save(login_data);
		taskResult.$promise.then(
			function(result) {
				// TODO: Redirect to index?
				$window.sessionStorage.token = result.accesToken;
			},
			function(result) {
				delete $window.sessionStorage.token;
			}
		);
	}
}]);
