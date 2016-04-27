var godan_api = "http://localhost:8000/";

angular.module('Godan').controller('UsersCtrl', ["$scope", "$resource", "$uibModal", "$window", function UsersCtrl($scope, $resource, $uibModal, $window) {
	var url = godan_api + "users";
	$scope.usersTable = $resource(url, {}, {
		query: {
			method: 'GET',
			isArray: true,
			headers: {'Authorization': 'Bearer ' + $window.sessionStorage.token}
		}
	}).query();
}]);
