var godan_api = "http://localhost:8000/";

angular.module('Godan').controller('QueriesCtrl', ["$scope", "$resource", "$uibModal", "$window", function QueriesCtrl($scope, $resource, $uibModal, $window) {

	$scope.animationsEnabled = true;

	$scope.submitQuery = function() {
		var query = godan_api + "ips";
		
		var params = {}
		if ($scope.inputPort) {
			params["port"] = $scope.inputPort;
		}
		if ($scope.inputService) {
			params["service"] = $scope.inputService;
		}
		if ($scope.inputRegexp) {
			params["regexp"] = $scope.inputRegexp;
		}

		if ($scope.inputIP) {
			query = query + "/" + $scope.inputIP
		}

		if (Object.keys(params).length > 0) {
			query += "?"+$.param(params);
		}

		$scope.resultsTable = $resource(query, {}, {
			query: {
				method: 'GET',
				isArray: true,
				headers: {'Authorization': 'Bearer ' + $window.sessionStorage.token}
			}
		}).query();
	}

	$scope.rowClicked = function(element) {
		var modalInstance = $uibModal.open({
			animation: $scope.animationsEnabled,
			templateUrl: 'modalTemplate',
			controller: 'ModalInstanceCtrl',
			size: 'lg',
			resolve: {
				element: function() {
					return element;
				}
			}
		});
		$scope.detailedContent = element.Content;
	};

}]);

angular.module('Godan').controller('ModalInstanceCtrl', function ($scope, $uibModalInstance, element) {
  $scope.closeit = function () {
    $uibModalInstance.close();
  };
	$scope.content = atob(element.Content); 
	$scope.elements = [element];
});
