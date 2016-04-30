var godan_api = "http://localhost:8000/";

angular.module('Godan').controller('UsersCtrl', ["$scope", "$resource", "$uibModal", "$window", function UsersCtrl($scope, $resource, $uibModal, $window) {
	if (!$window.sessionStorage.token) $window.location.href = '#/login';

	var url = godan_api + "users";
	$scope.usersTable = $resource(url, {}, {
		query: {
			method: 'GET',
			isArray: true,
			headers: {'Authorization': 'Bearer ' + $window.sessionStorage.token}
		}
	}).query();
}]);
